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

package eventbus

import (
	"context"
	"sync"
	"time"

	"github.com/observiq/bindplane-op/internal/util"
)

const subscriberChannelBufferSize = 10

// UnsubscribeFunc is a function that allows a subscriber to unsubscribe
type UnsubscribeFunc func()

// Subscriber can be notified of events of type T. Instead of using this interface directly, use one of the eventbus.Subscribe
// functions to receive a channel of events.
type Subscriber[T any] interface {
	// Channel will be return to Subscribe calls to receive events
	Channel() <-chan T

	// Receive will be called when an event is available
	Receive(event T)

	// Close will be called when the subscriber is unsubscribed
	Close()
}

// Source is a source of events.
type Source[T any] interface {
	// Send the event to all of the subscribers
	Send(event T)

	// SubscribeUntilDone adds a subscriber to the source and automatically unsubscribes when the context is done. If the
	// context is nil, the unsubscribe function must be called to unsubscribe. Instead of using this method to subscribe
	// to a source, use one of the eventbus.Subscribe functions to receive a channel of events.
	SubscribeUntilDone(ctx context.Context, subscriber Subscriber[T], onUnsubscribe func()) UnsubscribeFunc

	// Subscribers returns the current number of subscribers
	Subscribers() int
}

// SubscriptionFilter can filter on events and map from an event to another type. It can also ignore events. If accept
// is false, the result is ignored and not sent to subscribers.
type SubscriptionFilter[T, R any] func(event T) (result R, accept bool)

var exists = struct{}{}

// implementation of EventBus
type source[T any] struct {
	// subscribers is the set of current subscribers, implemented as a map to an empty struct
	subscribers map[Subscriber[T]]struct{}
	mtx         sync.RWMutex
}

var _ Source[any] = (*source[any])(nil)

// ----------------------------------------------------------------------
// generic subscriptions

type subscription[T any] struct {
	channel chan T
	ctx     context.Context
	cancel  context.CancelFunc
}

func newSubscription[T any](options []SubscriptionOption[T]) Subscriber[T] {
	opts := makeSubscriptionOptions(options)
	if opts.unbounded {
		return newUnboundedSubscription[T](opts.unboundedInterval)
	}
	channel := opts.channel
	if channel == nil {
		channel = make(chan T, subscriberChannelBufferSize)
	}
	ctx, cancel := context.WithCancel(context.Background())
	return &subscription[T]{
		channel: channel,
		ctx:     ctx,
		cancel:  cancel,
	}
}

func (s *subscription[T]) Receive(event T) {
	select {
	case <-s.ctx.Done():
		close(s.channel)
	case s.channel <- event:
	}
}

func (s *subscription[T]) Channel() <-chan T {
	return s.channel
}

func (s *subscription[T]) Close() {
	s.cancel()
}

var _ Subscriber[int] = (*subscription[int])(nil)

// ----------------------------------------------------------------------
// filter subscriptions

type filterSubscription[T, R any] struct {
	subscription Subscriber[R]
	filter       SubscriptionFilter[T, R]
}

func newFilterSubscription[T, R any](filter SubscriptionFilter[T, R], options []SubscriptionOption[R]) *filterSubscription[T, R] {
	return &filterSubscription[T, R]{
		subscription: newSubscription(options),
		filter:       filter,
	}
}

func (s *filterSubscription[T, R]) Channel() <-chan T {
	panic("filterSubscription Channel() should not be called because it won't receive anything")
}
func (s *filterSubscription[T, R]) FilterChannel() <-chan R {
	return s.subscription.Channel()
}

func (s *filterSubscription[T, R]) Receive(event T) {
	filtered, accept := s.filter(event)
	if accept {
		s.subscription.Receive(filtered)
	}
}

func (s *filterSubscription[T, R]) Close() {
	s.subscription.Close()
}

var _ Subscriber[int] = (*filterSubscription[int, int])(nil)

// ----------------------------------------------------------------------
// unbounded

type unboundedSubscription[T any] struct {
	channel util.UnboundedChan[T]
}

const unboundedMinimumInterval = 50 * time.Millisecond

func newUnboundedSubscription[T any](interval time.Duration) *unboundedSubscription[T] {
	if interval < unboundedMinimumInterval {
		interval = unboundedMinimumInterval
	}
	return &unboundedSubscription[T]{
		channel: util.NewUnboundedChan[T](interval),
	}
}

func (s *unboundedSubscription[T]) Receive(event T) {
	s.channel.In() <- event
}

func (s *unboundedSubscription[T]) Channel() <-chan T {
	return s.channel.Out()
}

func (s *unboundedSubscription[T]) Close() {
	s.channel.Close()
}

var _ Subscriber[int] = (*unboundedSubscription[int])(nil)

// ----------------------------------------------------------------------
// package methods

// NewSource returns a new Source implementation for the specified event type T
func NewSource[T any]() Source[T] {
	return &source[T]{
		subscribers: make(map[Subscriber[T]]struct{}),
	}
}

// Subscribe subscribes to events on the bus and returns a channel to receive events and an unsubscribe function.
func Subscribe[T any](bus Source[T], options ...SubscriptionOption[T]) (<-chan T, UnsubscribeFunc) {
	return SubscribeUntilDone(nil, bus, options...)
}

// SubscribeUntilDone subscribes to events on the bus and returns a channel to receive events and an unsubscribe
// function. It automatically unsubscribes when the context is done.
func SubscribeUntilDone[T any](ctx context.Context, bus Source[T], options ...SubscriptionOption[T]) (<-chan T, UnsubscribeFunc) {
	subscription := newSubscription(options)
	opts := makeSubscriptionOptions(options)
	unsubscribe := bus.SubscribeUntilDone(ctx, subscription, opts.unsubscribeHook)
	return subscription.Channel(), unsubscribe
}

// SubscribeWithFilter TODO
func SubscribeWithFilter[T, R any](bus Source[T], filter SubscriptionFilter[T, R], options ...SubscriptionOption[R]) (<-chan R, UnsubscribeFunc) {
	return SubscribeWithFilterUntilDone(nil, bus, filter, options...)
}

// SubscribeWithFilterUntilDone TODO
func SubscribeWithFilterUntilDone[T, R any](ctx context.Context, source Source[T], filter SubscriptionFilter[T, R], options ...SubscriptionOption[R]) (<-chan R, UnsubscribeFunc) {
	subscription := newFilterSubscription(filter, options)
	opts := makeSubscriptionOptions(options)
	unsubscribe := source.SubscribeUntilDone(ctx, subscription, opts.unsubscribeHook)
	return subscription.FilterChannel(), unsubscribe
}

// SubscribeUntilDone adds the subscriber and returns a cancel function.
func (s *source[T]) SubscribeUntilDone(ctx context.Context, subscriber Subscriber[T], unsubscribeHook func()) UnsubscribeFunc {
	s.mtx.Lock()
	defer s.mtx.Unlock()

	// create a new subscription and add it
	s.subscribers[subscriber] = exists

	// create the unsubscribe function
	unsubscribe := func() {
		s.mtx.Lock()
		defer s.mtx.Unlock()
		delete(s.subscribers, subscriber)
		subscriber.Close()
		if unsubscribeHook != nil {
			unsubscribeHook()
		}
	}

	// wait for context done and auto-unsubscribe
	if ctx != nil {
		go func() {
			<-ctx.Done()
			unsubscribe()
		}()
	}

	return unsubscribe
}

// Send the event to all of the subscribers
func (s *source[T]) Send(event T) {
	for _, sub := range s.subscriberList() {
		sub.Receive(event)
	}
}

func (s *source[T]) subscriberList() []Subscriber[T] {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	subscribers := make([]Subscriber[T], 0, len(s.subscribers))
	for sub := range s.subscribers {
		subscribers = append(subscribers, sub)
	}
	return subscribers
}

// Subscribers returns the current number of subscribers
func (s *source[T]) Subscribers() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return len(s.subscribers)
}

// Relay will relay from source to destination. It runs a separate goroutine, consuming events from the source and
// sending events to the destination. When the supplied context is Done, the relay is automatically unsubscribed from
// the source and the destination will no longer receive events.
func Relay[T any](ctx context.Context, source Source[T], destination Source[T], options ...SubscriptionOption[T]) {
	channel, unsubscribe := Subscribe(source, options...)
	go relay(ctx, channel, unsubscribe, destination)
}

// RelayWithFilter will relay from source to destination with the specified filter. It runs a separate goroutine,
// consuming events from the source, running the filter, and sending events to the destination. When the supplied
// context is Done, the relay is automatically unsubscribed from the source and the destination will no longer receive
// events.
func RelayWithFilter[T, R any](ctx context.Context, source Source[T], filter SubscriptionFilter[T, R], destination Source[R], options ...SubscriptionOption[R]) {
	channel, unsubscribe := SubscribeWithFilter(source, filter, options...)
	go relay(ctx, channel, unsubscribe, destination)
}

func relay[T any](ctx context.Context, channel <-chan T, unsubscribe UnsubscribeFunc, destination Source[T]) {
	defer unsubscribe()
	for {
		select {
		case event, ok := <-channel:
			if ok {
				destination.Send(event)
			} else {
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

// SubscriptionMerger merges the second event into the first and returns true if the merge was successful. It should
// return false if a merge was not possible and the two individual events will be preserved and dispatched separately.
type SubscriptionMerger[S any, T *S] func(into, single T) bool

// RelayWithMerge will relay from source to destination, merging events before sending them to the destination. This can
// be used when there are lots of small individual events that can be more efficiently processed as a few larger events.
func RelayWithMerge[S any, T *S](ctx context.Context, source Source[T], merge SubscriptionMerger[S, T], destination Source[T], maxLatency time.Duration, maxEventsToMerge int, options ...SubscriptionOption[T]) {
	// constrain max events to at least 1
	if maxEventsToMerge < 1 {
		maxEventsToMerge = 1
	}
	channel, unsubscribe := Subscribe(source, options...)
	go func() {
		defer unsubscribe()

		maxLatencyTicker := time.NewTicker(maxLatency)
		defer maxLatencyTicker.Stop()

		var buffer []T
		mergeAndSend := func() {
			if len(buffer) == 0 {
				return
			}

			// merge: insert the first element merge into it
			prev := buffer[0]
			merged := []T{prev}

			for _, item := range buffer[1:] {
				if !merge(prev, item) {
					// unable to merge, append and start merging into this item
					merged = append(merged, item)
					prev = item
				}
			}

			// send the merged items
			for _, item := range merged {
				destination.Send(item)
			}

			// reset the buffer
			buffer = nil
		}

		// drain anything remaining when finished
		defer mergeAndSend()
		for {
			select {
			case event, ok := <-channel:
				if !ok {
					return
				}
				buffer = append(buffer, event)
				if len(buffer) >= maxEventsToMerge {
					mergeAndSend()
				}

			case <-maxLatencyTicker.C:
				// periodically drain the buffer to limit latency
				mergeAndSend()

			case <-ctx.Done():
				// send anything left in the buffer before stopping
				return
			}
		}
	}()
}
