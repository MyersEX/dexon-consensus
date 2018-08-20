// Copyright 2018 The dexon-consensus-core Authors
// This file is part of the dexon-consensus-core library.
//
// The dexon-consensus-core library is free software: you can redistribute it
// and/or modify it under the terms of the GNU Lesser General Public License as
// published by the Free Software Foundation, either version 3 of the License,
// or (at your option) any later version.
//
// The dexon-consensus-core library is distributed in the hope that it will be
// useful, but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Lesser
// General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the dexon-consensus-core library. If not, see
// <http://www.gnu.org/licenses/>.

package test

import (
	"time"

	"github.com/dexon-foundation/dexon-consensus-core/core/types"
)

// Event defines a scheduler event.
type Event struct {
	// ValidatorID is the ID of handler that this event deginated to.
	ValidatorID types.ValidatorID
	// Time is the expected execution time of this event.
	Time time.Time
	// ExecError record the error when handling this event.
	ExecError error
	// Payload is application specific data carried by this event.
	Payload interface{}
	// ParentTime is the time of parent event, this field is essential when
	// we need to calculate the latency the handler assigned.
	ParentTime time.Time
	// ExecInterval is the latency to execute this event
	ExecInterval time.Duration
}

// eventQueue implements heap.Interface.
type eventQueue []*Event

func (eq eventQueue) Len() int { return len(eq) }

func (eq eventQueue) Less(i, j int) bool {
	return eq[i].Time.Before(eq[j].Time)
}

func (eq eventQueue) Swap(i, j int) {
	eq[i], eq[j] = eq[j], eq[i]
}

func (eq *eventQueue) Push(x interface{}) {
	*eq = append(*eq, x.(*Event))
}

func (eq *eventQueue) Pop() interface{} {
	pos := len(*eq) - 1
	item := (*eq)[pos]
	*eq = (*eq)[0:pos]
	return item
}

// NewEvent is the constructor for Event.
func NewEvent(
	vID types.ValidatorID, when time.Time, payload interface{}) *Event {

	return &Event{
		ValidatorID: vID,
		Time:        when,
		Payload:     payload,
	}
}