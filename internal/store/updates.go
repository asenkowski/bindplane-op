// Copyright  observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package store

import (
	"context"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/observiq/bindplane-op/internal/eventbus"
	"github.com/observiq/bindplane-op/model"
)

// Updates are sent on the channel available from Store.Updates().
type Updates struct {
	Agents           Events[*model.Agent]
	AgentVersions    Events[*model.AgentVersion]
	Sources          Events[*model.Source]
	SourceTypes      Events[*model.SourceType]
	Processors       Events[*model.Processor]
	ProcessorTypes   Events[*model.ProcessorType]
	Destinations     Events[*model.Destination]
	DestinationTypes Events[*model.DestinationType]
	Configurations   Events[*model.Configuration]
}

// NewUpdates returns a New Updates struct
func NewUpdates() *Updates {
	// TODO: optimize allocate as needed
	return &Updates{
		Agents:           NewEvents[*model.Agent](),
		AgentVersions:    NewEvents[*model.AgentVersion](),
		Sources:          NewEvents[*model.Source](),
		SourceTypes:      NewEvents[*model.SourceType](),
		Processors:       NewEvents[*model.Processor](),
		ProcessorTypes:   NewEvents[*model.ProcessorType](),
		Destinations:     NewEvents[*model.Destination](),
		DestinationTypes: NewEvents[*model.DestinationType](),
		Configurations:   NewEvents[*model.Configuration](),
	}
}

// IncludeAgent will include the agent in the updates. While updates.Agents.Include can also be used directly, this
// matches the pattern of IncludeResource.
func (updates *Updates) IncludeAgent(agent *model.Agent, eventType EventType) {
	updates.Agents.Include(agent, eventType)
}

// IncludeResource will include the resource in the updates for the appropriate type. If the specified Resource is not
// supported by Updates, this will do nothing.
func (updates *Updates) IncludeResource(r model.Resource, eventType EventType) {
	switch r := r.(type) {
	case *model.AgentVersion:
		updates.AgentVersions.Include(r, eventType)
	case *model.Source:
		updates.Sources.Include(r, eventType)
	case *model.SourceType:
		updates.SourceTypes.Include(r, eventType)
	case *model.Processor:
		updates.Processors.Include(r, eventType)
	case *model.ProcessorType:
		updates.ProcessorTypes.Include(r, eventType)
	case *model.Destination:
		updates.Destinations.Include(r, eventType)
	case *model.DestinationType:
		updates.DestinationTypes.Include(r, eventType)
	case *model.Configuration:
		updates.Configurations.Include(r, eventType)
	}
}

// Empty returns true if all individual updates are empty
func (updates *Updates) Empty() bool {
	return updates.Size() == 0
}

// Size returns the sum of all updates of all types
func (updates *Updates) Size() int {
	return len(updates.Agents) +
		len(updates.AgentVersions) +
		len(updates.Sources) +
		len(updates.SourceTypes) +
		len(updates.Processors) +
		len(updates.ProcessorTypes) +
		len(updates.Destinations) +
		len(updates.DestinationTypes) +
		len(updates.Configurations)
}

// ----------------------------------------------------------------------
//
// add transitive updates based on updates that already exist. this could be optimized for a specific store and may not
// be used by all stores.

// TODO: how does this work in a distributed environment?
// pub/sub individual event => pub/sub include dependencies => subscribers
func (updates *Updates) addTransitiveUpdates(ctx context.Context, s Store) error {
	// for sourceTypes, add sources
	// for processorTypes, add sources and processors
	// for destinationTypes, add destinations
	// for sources and sourceTypes, add configurations
	// for processors and processorTypes, add configurations
	// for destinations and destinationTypes, add configurations

	var errs error

	err := updates.addProcessorUpdates(ctx, s)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	err = updates.addSourceUpdates(ctx, s)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	err = updates.addDestinationUpdates(ctx, s)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	err = updates.addConfigurationUpdates(ctx, s)
	if err != nil {
		errs = multierror.Append(errs, err)
	}

	return errs
}

func (updates *Updates) addSourceUpdates(ctx context.Context, s Store) error {
	if updates.SourceTypes.Empty() && updates.Processors.Empty() && updates.ProcessorTypes.Empty() {
		return nil
	}

	// get all of the sources
	sources, err := s.Sources(ctx)
	if err != nil {
		return err
	}

sourceLoop:
	for _, source := range sources {
		// updates to a SourceType will trigger updates of all of the Sources that use that SourceType.
		for _, sourceTypeEvent := range updates.SourceTypes.Updates() {
			sourceTypeName := sourceTypeEvent.Item.Name()

			if source.Spec.Type == sourceTypeName {
				updates.Sources.Include(source, EventTypeUpdate)
				continue sourceLoop
			}
		}

		// updates to a ProcessorType will trigger updates of all of the Sources that use that ProcessorType.
		for _, processorTypeEvent := range updates.ProcessorTypes.Updates() {
			processorTypeName := processorTypeEvent.Item.Name()
			for _, processor := range source.Spec.Processors {
				if processor.Type == processorTypeName {
					updates.Sources.Include(source, EventTypeUpdate)
					continue sourceLoop
				}
			}
		}

		// updates to a Processor will trigger updates of all of the Sources that use that Processor.
		for _, processorEvent := range updates.Processors.Updates() {
			processorName := processorEvent.Item.Name()
			for _, processor := range source.Spec.Processors {
				if processor.Name == processorName {
					updates.Sources.Include(source, EventTypeUpdate)
					continue sourceLoop
				}
			}
		}
	}

	return nil
}

func (updates *Updates) addProcessorUpdates(ctx context.Context, s Store) error {
	if updates.ProcessorTypes.Empty() {
		return nil
	}

	processors, err := s.Processors(ctx)
	if err != nil {
		return err
	}

	for _, processorTypeEvent := range updates.ProcessorTypes {
		if processorTypeEvent.Type == EventTypeUpdate {
			processorTypeName := processorTypeEvent.Item.Name()

			for _, processor := range processors {
				if processor.Spec.Type == processorTypeName {
					updates.Processors.Include(processor, EventTypeUpdate)
				}
			}
		}
	}

	return nil
}

func (updates *Updates) addDestinationUpdates(ctx context.Context, s Store) error {
	if updates.DestinationTypes.Empty() {
		return nil
	}

	// get all of the destinations
	destinations, err := s.Destinations(ctx)
	if err != nil {
		return err
	}

	// updates to a DestinationType will trigger updates of all of the Destinations that use that DestinationType.
	for _, destinationTypeEvent := range updates.DestinationTypes {
		if destinationTypeEvent.Type == EventTypeUpdate {
			destinationTypeName := destinationTypeEvent.Item.Name()

			for _, destination := range destinations {
				if destination.Spec.Type == destinationTypeName {
					updates.Destinations.Include(destination, EventTypeUpdate)
				}
			}
		}
	}

	return nil
}

func (updates *Updates) addConfigurationUpdates(ctx context.Context, s Store) error {
	configurations, err := s.Configurations(ctx)
	if err != nil {
		return err
	}

	for _, configuration := range configurations {
		// as a small optimization, before checking all of the sources and destinations for changes, check to see if we're
		// already updating this configuration.
		if updates.Configurations.Contains(configuration.Name(), EventTypeUpdate) {
			continue
		}
		updates.addConfigurationUpdatesFromComponents(configuration, s)
	}
	return nil
}

func (updates *Updates) addConfigurationUpdatesFromComponents(configuration *model.Configuration, _ Store) {
	for _, source := range configuration.Spec.Sources {
		if _, ok := updates.Sources[source.Name]; ok {
			updates.Configurations.Include(configuration, EventTypeUpdate)
			return
		}
		if _, ok := updates.SourceTypes[source.Type]; ok {
			updates.Configurations.Include(configuration, EventTypeUpdate)
			return
		}
	}
	for _, destination := range configuration.Spec.Destinations {
		if _, ok := updates.Destinations[destination.Name]; ok {
			updates.Configurations.Include(configuration, EventTypeUpdate)
			return
		}
		if _, ok := updates.DestinationTypes[destination.Type]; ok {
			updates.Configurations.Include(configuration, EventTypeUpdate)
			return
		}
	}
}

// ----------------------------------------------------------------------
// merge for use with RelayWithMerge

func mergeUpdates(into, single *Updates) bool {
	// first make sure we can safely merge
	safe := into.Agents.CanSafelyMerge(single.Agents) &&
		into.AgentVersions.CanSafelyMerge(single.AgentVersions) &&
		into.Sources.CanSafelyMerge(single.Sources) &&
		into.SourceTypes.CanSafelyMerge(single.SourceTypes) &&
		into.Processors.CanSafelyMerge(single.Processors) &&
		into.ProcessorTypes.CanSafelyMerge(single.ProcessorTypes) &&
		into.Destinations.CanSafelyMerge(single.Destinations) &&
		into.DestinationTypes.CanSafelyMerge(single.DestinationTypes) &&
		into.Configurations.CanSafelyMerge(single.Configurations)

	if !safe {
		return false
	}

	// merge individual events
	into.Agents.Merge(single.Agents)
	into.AgentVersions.Merge(single.AgentVersions)
	into.Sources.Merge(single.Sources)
	into.SourceTypes.Merge(single.SourceTypes)
	into.Processors.Merge(single.Processors)
	into.ProcessorTypes.Merge(single.ProcessorTypes)
	into.Destinations.Merge(single.Destinations)
	into.DestinationTypes.Merge(single.DestinationTypes)
	into.Configurations.Merge(single.Configurations)

	return true
}

// ----------------------------------------------------------------------

type storeUpdates struct {
	updates eventbus.Source[*Updates]
	// updatesInternal is an internal source used for notification. It will relay to the updates available to clients of
	// the store.
	updatesInternal eventbus.Source[*Updates]
}

func newStoreUpdates(ctx context.Context, maxEventsToMerge int) *storeUpdates {
	updates := eventbus.NewSource[*Updates]()
	updatesInternal := eventbus.NewSource[*Updates]()

	if maxEventsToMerge == 0 {
		maxEventsToMerge = 100
	}

	// introduce a separate relay with a large buffer to avoid blocking on changes
	eventbus.RelayWithMerge(
		ctx,
		updatesInternal,
		mergeUpdates,
		updates,
		200*time.Millisecond,
		maxEventsToMerge,
		eventbus.WithUnboundedChannel[*Updates](100*time.Millisecond),
	)

	return &storeUpdates{
		updates:         updates,
		updatesInternal: updatesInternal,
	}
}

// Updates returns the external channel that can be provided to external clients.
func (s *storeUpdates) Updates() eventbus.Source[*Updates] {
	return s.updates
}

// Send adds an Updates event to the internal channel where it can be merged and relayed to the external channel.
func (s *storeUpdates) Send(updates *Updates) {
	s.updatesInternal.Send(updates)
}
