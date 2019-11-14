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

package stime

// This file implements operations on veqryn.protobuf.Timestamp.

import (
	"time"

	google_tspb "github.com/golang/protobuf/ptypes/timestamp"
	tspb "github.com/syncfuture/go/sproto/timestamp"
)

// Timestamp converts a veqryn.protobuf.Timestamp proto to a time.Time.
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
func Timestamp(ts *tspb.Timestamp) (time.Time, error) {
	return ts.Time()
}

// TimestampNow returns a veqryn.protobuf.Timestamp for the current time.
func TimestampNow() *tspb.Timestamp {
	return (&tspb.Timestamp{}).SetToNow()
}

// TimestampProto converts the time.Time to a veqryn.protobuf.Timestamp proto.
// It returns an error if the resulting Timestamp is invalid.
func TimestampProto(t time.Time) (*tspb.Timestamp, error) {
	return (&tspb.Timestamp{}).SetToTime(t)
}

// TimestampString returns the RFC 3339 string for valid Timestamps. For invalid
// Timestamps, it returns an error message in parentheses.
func TimestampString(ts *tspb.Timestamp) string {
	return ts.RFC3339()
}

// StringTimestamp creates a veqryn.protobuf.Timestamp proto from a string with the given layout
func StringTimestamp(layout string, dateValue string) (*tspb.Timestamp, error) {
	return (&tspb.Timestamp{}).SetToString(layout, dateValue)
}

// ToGoogleTimestamp converts a veqryn.protobuf.Timestamp proto into a google.protobuf.Timestamp proto
func ToGoogleTimestamp(ts *tspb.Timestamp) *google_tspb.Timestamp {
	return ts.GoogleTimestamp()
}

// FromGoogleTimestamp converts a google.protobuf.Timestamp proto into a veqryn.protobuf.Timestamp proto
func FromGoogleTimestamp(ts *google_tspb.Timestamp) *tspb.Timestamp {
	return (&tspb.Timestamp{}).SetToGoogleTimestamp(ts)
}
