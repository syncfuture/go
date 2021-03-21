// Go support for Protocol Buffers - Google's data interchange format
//
// Copyright 2016 The Go Authors.  All rights reserved.
// https://github.com/golang/protobuf
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//     * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//     * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//     * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package timestamp

import (
	"fmt"
	"time"

	google_tspb "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/pkg/errors"
)

const (
	// Seconds field of the earliest valid Timestamp.
	// This is time.Date(1, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	minValidSeconds = -62135596800
	// Seconds field just after the latest valid Timestamp.
	// This is time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC).Unix().
	maxValidSeconds = 253402300800
)

// validateTimestamp determines whether a Timestamp is valid.
// A valid timestamp represents a time in the range
// [0001-01-01, 10000-01-01) and has a Nanos field
// in the range [0, 1e9).
//
// If the Timestamp is valid, validateTimestamp returns nil.
// Otherwise, it returns an error that describes
// the problem.
//
// Every valid Timestamp can be represented by a time.Time, but the converse is not true.
func (m *Timestamp) ValidateTimestamp() error {
	if m == nil {
		return errors.New("timestamp: nil Timestamp")
	}
	if m.Seconds < minValidSeconds {
		return errors.Errorf("timestamp: %v before 0001-01-01", m)
	}
	if m.Seconds >= maxValidSeconds {
		return errors.Errorf("timestamp: %v after 10000-01-01", m)
	}
	if m.Nanos < 0 || m.Nanos >= 1e9 {
		return errors.Errorf("timestamp: %v: nanos not in range [0, 1e9)", m)
	}
	return nil
}

// Time converts the veqryn.protobuf.Timestamp proto to a time.Time.
// It returns an error if the argument is invalid.
//
// Unlike most Go functions, if Timestamp returns an error, the first return value
// is not the zero time.Time. Instead, it is the value obtained from the
// time.Unix function when passed the contents of the Timestamp, in the UTC
// locale. This may or may not be a meaningful time; many invalid Timestamps
// do map to valid time.Times.
//
// A nil Timestamp returns an error. The first return value in that case is
// undefined.
func (m *Timestamp) Time() (time.Time, error) {
	// Don't return the zero value on error, because corresponds to a valid
	// timestamp. Instead return whatever time.Unix gives us.
	var t time.Time
	if m == nil {
		t = time.Unix(0, 0).UTC() // treat nil like the empty Timestamp
	} else {
		t = time.Unix(m.Seconds, int64(m.Nanos)).UTC()
	}
	return t, m.ValidateTimestamp()
}

// SetToNow sets and returns the veqryn.protobuf.Timestamp for the current time.
func (m *Timestamp) SetToNow() *Timestamp {
	ts, err := m.SetToTime(time.Now())
	if err != nil {
		panic("timestamp: time.Now() out of Timestamp range")
	}
	return ts
}

// SetFromTime sets and returns the veqryn.protobuf.Timestamp proto to the time.Time.
// It returns an error if the resulting Timestamp is invalid.
func (m *Timestamp) SetToTime(t time.Time) (*Timestamp, error) {
	// initialize timestamp if pointer is nil
	if m == nil {
		*m = Timestamp{}
	}

	m.Seconds = t.Unix()
	m.Nanos = int32(t.Nanosecond())

	if err := m.ValidateTimestamp(); err != nil {
		return nil, err
	}
	return m, nil
}

// RFC3339 returns the RFC 3339 string for valid Timestamps. For invalid
// Timestamps, it returns an error message in parentheses.
func (m *Timestamp) RFC3339() string {
	t, err := m.Time()
	if err != nil {
		return fmt.Sprintf("(%v)", err)
	}
	return t.Format(time.RFC3339Nano)
}

// SetToString sets and returns the veqryn.protobuf.Timestamp proto to a time string with the given layout
func (m *Timestamp) SetToString(layout string, dateValue string) (*Timestamp, error) {
	t, err := time.Parse(layout, dateValue)
	if err != nil {
		return nil, err
	}
	return m.SetToTime(t)
}

// GoogleTimestamp converts the veqryn.protobuf.Timestamp proto into a google.protobuf.Timestamp proto
func (m *Timestamp) GoogleTimestamp() *google_tspb.Timestamp {
	return &google_tspb.Timestamp{Seconds: m.Seconds, Nanos: m.Nanos}
}

// SetToGoogleTimestamp sets and returns the veqryn.protobuf.Timestamp proto to a google.protobuf.Timestamp proto
func (m *Timestamp) SetToGoogleTimestamp(ts *google_tspb.Timestamp) *Timestamp {
	// initialize timestamp if pointer is nil
	if m == nil {
		*m = Timestamp{}
	}
	m.Seconds = ts.Seconds
	m.Nanos = ts.Nanos
	return m
}
