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
	"testing"

	"github.com/observiq/bindplane-op/model"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var (
	updatesTestStore Store
	resourceMap      map[string]model.Resource
	testOptions      = Options{
		SessionsSecret:   "super-secret-key",
		MaxEventsToMerge: 1,
	}
)

func updatesTestSetup(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	updatesTestStore = NewMapStore(ctx, testOptions, zap.NewNop())
	resourceMap = map[string]model.Resource{}
	resources := []model.Resource{
		newTestProcessorType("pt1"),
		newTestProcessorType("pt2"),
		newTestProcessorType("pt3"),
		newTestProcessor("p1", "pt1"),
		newTestSourceType("st1"),
		newTestSourceType("st2"),
		newTestSourceType("st3"),
		newTestSourceType("st4"),
		newTestSourceType("st5"),
		newTestSource("s1", "st1"),
		newTestSource("s2", "st2"),
		newTestSource("s3", "st3"),
		newTestSourceWithProcessors("s4", "st5", []model.ResourceConfiguration{
			{ParameterizedSpec: model.ParameterizedSpec{Type: "pt2"}},
			{Name: "p1"},
		}),
		newTestDestinationType("dt1"),
		newTestDestinationType("dt2"),
		newTestDestinationType("dt3"),
		newTestDestinationType("dt4"),
		newTestDestination("d1", "dt1"),
		newTestDestination("d2", "dt2"),
		newTestDestination("d3", "dt3"),
		newTestConfiguration("c1", []string{"s1"}, []string{"st2"}, []string{"d1"}, []string{"dt2"}),
		newTestConfiguration("c2", []string{"s2"}, nil, []string{"d2"}, nil),
		newTestConfiguration("c3", []string{"s1", "s2", "s3"}, nil, []string{"d1", "d2", "d3"}, nil),
		newTestConfiguration("c4", nil, []string{"st4"}, nil, []string{"dt4"}),
		newTestConfiguration("c5", nil, nil, nil, nil),
		newTestConfiguration("c6", []string{"s4"}, nil, []string{"d3"}, nil),
		newTestConfiguration("c7", nil, []string{"st5"}, []string{"d3"}, nil),
	}
	for _, resource := range resources {
		resourceMap[resource.Name()] = resource
	}
	_, err := updatesTestStore.ApplyResources(context.Background(), resources)
	require.NoError(t, err)
}

func newTestSourceType(name string) *model.SourceType {
	return model.NewSourceType(name, []model.ParameterDefinition{})
}

func newTestSource(name string, sourceType string) *model.Source {
	return model.NewSource(name, sourceType, []model.Parameter{})
}

func newTestSourceWithProcessors(name string, sourceType string, processors []model.ResourceConfiguration) *model.Source {
	src := model.NewSource(name, sourceType, []model.Parameter{})
	src.Spec.Processors = processors
	return src
}

func newTestProcessorType(name string) *model.ProcessorType {
	return model.NewProcessorType(name, []model.ParameterDefinition{})
}

func newTestProcessor(name string, processorType string) *model.Processor {
	return model.NewProcessor(name, processorType, []model.Parameter{})
}

func newTestDestinationType(name string) *model.DestinationType {
	return model.NewDestinationType(name, []model.ParameterDefinition{})
}

func newTestDestination(name string, destinationType string) *model.Destination {
	return model.NewDestination(name, destinationType, []model.Parameter{})
}

func newTestConfiguration(name string, sources []string, sourceTypes []string, destinations []string, destinationTypes []string) *model.Configuration {
	c := &model.Configuration{
		ResourceMeta: model.ResourceMeta{
			APIVersion: model.V1,
			Kind:       model.KindDestinationType,
			Metadata: model.Metadata{
				Name: name,
			},
		},
		Spec: model.ConfigurationSpec{},
	}
	for _, source := range sources {
		c.Spec.Sources = append(c.Spec.Sources, model.ResourceConfiguration{Name: source})
	}
	for _, sourceType := range sourceTypes {
		c.Spec.Sources = append(c.Spec.Sources, model.ResourceConfiguration{ParameterizedSpec: model.ParameterizedSpec{Type: sourceType}})
	}
	for _, destination := range destinations {
		c.Spec.Destinations = append(c.Spec.Destinations, model.ResourceConfiguration{Name: destination})
	}
	for _, destinationType := range destinationTypes {
		c.Spec.Destinations = append(c.Spec.Destinations, model.ResourceConfiguration{ParameterizedSpec: model.ParameterizedSpec{Type: destinationType}})
	}
	return c
}

func addUpdates[T model.Resource](t *testing.T, names []string, events Events[T]) {
	for _, name := range names {
		resource, ok := resourceMap[name]
		require.True(t, ok)
		events.Include(resource.(T), EventTypeUpdate)
	}
}

func TestTransitiveUpdates(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	updatesTestSetup(t)
	tests := []struct {
		Name string

		Sources          []string
		SourceTypes      []string
		Processors       []string
		ProcessorTypes   []string
		Destinations     []string
		DestinationTypes []string
		Configurations   []string

		ExpectSources          []string
		ExpectSourceTypes      []string
		ExpectProcessors       []string
		ExpectProcessorTypes   []string
		ExpectDestinations     []string
		ExpectDestinationTypes []string
		ExpectConfigurations   []string
	}{
		{
			Name:                 "s1 source",
			Sources:              []string{"s1"},
			ExpectSources:        []string{"s1"},
			ExpectConfigurations: []string{"c1", "c3"},
		},
		{
			Name:                 "s2 source",
			Sources:              []string{"s2"},
			ExpectSources:        []string{"s2"},
			ExpectConfigurations: []string{"c2", "c3"},
		},
		{
			Name:                 "s1 sources, d2 destination",
			Sources:              []string{"s1"},
			Destinations:         []string{"d2"},
			ExpectSources:        []string{"s1"},
			ExpectDestinations:   []string{"d2"},
			ExpectConfigurations: []string{"c1", "c2", "c3"},
		},
		{
			Name:                 "st1-4 source types",
			SourceTypes:          []string{"st1", "st2", "st3", "st4"},
			ExpectSources:        []string{"s1", "s2", "s3"},
			ExpectSourceTypes:    []string{"st1", "st2", "st3", "st4"},
			ExpectConfigurations: []string{"c1", "c2", "c3", "c4"},
		},
		{
			Name:                   "dt2 destination type",
			DestinationTypes:       []string{"dt2"},
			ExpectDestinationTypes: []string{"dt2"},
			ExpectDestinations:     []string{"d2"},
			ExpectConfigurations:   []string{"c1", "c2", "c3"},
		},
		{
			Name:                 "s1 source, st1 sourceType, d2 destination",
			Sources:              []string{"s1"},
			SourceTypes:          []string{"st1"},
			Destinations:         []string{"d2"},
			Configurations:       []string{"c1"},
			ExpectSources:        []string{"s1"},
			ExpectSourceTypes:    []string{"st1"},
			ExpectDestinations:   []string{"d2"},
			ExpectConfigurations: []string{"c1", "c2", "c3"},
		},
		{
			Name:                 "s1 source, st1 sourceType, d2 destination",
			Sources:              []string{"s1"},
			SourceTypes:          []string{"st1"},
			Destinations:         []string{"d2"},
			Configurations:       []string{"c1"},
			ExpectSources:        []string{"s1"},
			ExpectSourceTypes:    []string{"st1"},
			ExpectDestinations:   []string{"d2"},
			ExpectConfigurations: []string{"c1", "c2", "c3"},
		},
		{
			Name:                 "p1",
			Processors:           []string{"p1"},
			ExpectProcessors:     []string{"p1"},
			ExpectSources:        []string{"s4"},
			ExpectConfigurations: []string{"c6"},
		},
		{
			Name:                 "pt2",
			ProcessorTypes:       []string{"pt2"},
			ExpectProcessorTypes: []string{"pt2"},
			ExpectSources:        []string{"s4"},
			ExpectConfigurations: []string{"c6"},
		},
		{
			Name:                 "pt1",
			ProcessorTypes:       []string{"pt1"},
			ExpectSources:        []string{"s4"},
			ExpectProcessors:     []string{"p1"},
			ExpectProcessorTypes: []string{"pt1"},
			ExpectConfigurations: []string{"c6"},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			updates := NewUpdates()

			// populate updates
			addUpdates(t, test.Sources, updates.Sources)
			addUpdates(t, test.SourceTypes, updates.SourceTypes)
			addUpdates(t, test.Processors, updates.Processors)
			addUpdates(t, test.ProcessorTypes, updates.ProcessorTypes)
			addUpdates(t, test.Destinations, updates.Destinations)
			addUpdates(t, test.DestinationTypes, updates.DestinationTypes)
			addUpdates(t, test.Configurations, updates.Configurations)

			// add transitive
			err := updates.addTransitiveUpdates(ctx, updatesTestStore)
			require.NoError(t, err)

			// compare results
			require.ElementsMatch(t, test.ExpectSources, updates.Sources.Keys(), "Sources")
			require.ElementsMatch(t, test.ExpectSourceTypes, updates.SourceTypes.Keys(), "SourceTypes")
			require.ElementsMatch(t, test.ExpectProcessors, updates.Processors.Keys(), "Processors")
			require.ElementsMatch(t, test.ExpectProcessorTypes, updates.ProcessorTypes.Keys(), "ProcessorTypes")
			require.ElementsMatch(t, test.ExpectDestinations, updates.Destinations.Keys(), "Destinations")
			require.ElementsMatch(t, test.ExpectDestinationTypes, updates.DestinationTypes.Keys(), "DestinationTypes")
			require.ElementsMatch(t, test.ExpectConfigurations, updates.Configurations.Keys(), "Configurations")
		})
	}
}
