// Code generated by protoc-gen-go. DO NOT EDIT.
// source: timestamp.proto

package timestamp

import (
	fmt "fmt"
	math "math"

	proto "github.com/golang/protobuf/proto"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// A Timestamp represents a point in time independent of any time zone
// or calendar, represented as seconds and fractions of seconds at
// nanosecond resolution in UTC Epoch time. It is encoded using the
// Proleptic Gregorian Calendar which extends the Gregorian calendar
// backwards to year one. It is encoded assuming all minutes are 60
// seconds long, i.e. leap seconds are "smeared" so that no leap second
// table is needed for interpretation. Range is from
// 0001-01-01T00:00:00Z to 9999-12-31T23:59:59.999999999Z.
// By restricting to that range, we ensure that we can convert to
// and from  RFC 3339 date strings.
// See [https://www.ietf.org/rfc/rfc3339.txt](https://www.ietf.org/rfc/rfc3339.txt).
//
// # Examples
//
// Example 1: Compute Timestamp from POSIX `time()`.
//
//     Timestamp timestamp;
//     timestamp.set_seconds(time(NULL));
//     timestamp.set_nanos(0);
//
// Example 2: Compute Timestamp from POSIX `gettimeofday()`.
//
//     struct timeval tv;
//     gettimeofday(&tv, NULL);
//
//     Timestamp timestamp;
//     timestamp.set_seconds(tv.tv_sec);
//     timestamp.set_nanos(tv.tv_usec * 1000);
//
// Example 3: Compute Timestamp from Win32 `GetSystemTimeAsFileTime()`.
//
//     FILETIME ft;
//     GetSystemTimeAsFileTime(&ft);
//     UINT64 ticks = (((UINT64)ft.dwHighDateTime) << 32) | ft.dwLowDateTime;
//
//     // A Windows tick is 100 nanoseconds. Windows epoch 1601-01-01T00:00:00Z
//     // is 11644473600 seconds before Unix epoch 1970-01-01T00:00:00Z.
//     Timestamp timestamp;
//     timestamp.set_seconds((INT64) ((ticks / 10000000) - 11644473600LL));
//     timestamp.set_nanos((INT32) ((ticks % 10000000) * 100));
//
// Example 4: Compute Timestamp from Java `System.currentTimeMillis()`.
//
//     long millis = System.currentTimeMillis();
//
//     Timestamp timestamp = Timestamp.newBuilder().setSeconds(millis / 1000)
//         .setNanos((int) ((millis % 1000) * 1000000)).build();
//
//
// Example 5: Compute Timestamp from current time in Python.
//
//     timestamp = Timestamp()
//     timestamp.GetCurrentTime()
//
// # JSON Mapping
//
// In JSON format, the Timestamp type is encoded as a string in the
// [RFC 3339](https://www.ietf.org/rfc/rfc3339.txt) format. That is, the
// format is "{year}-{month}-{day}T{hour}:{min}:{sec}[.{frac_sec}]Z"
// where {year} is always expressed using four digits while {month}, {day},
// {hour}, {min}, and {sec} are zero-padded to two digits each. The fractional
// seconds, which can go up to 9 digits (i.e. up to 1 nanosecond resolution),
// are optional. The "Z" suffix indicates the timezone ("UTC"); the timezone
// is required. A proto3 JSON serializer should always use UTC (as indicated by
// "Z") when printing the Timestamp type and a proto3 JSON parser should be
// able to accept both UTC and other timezones (as indicated by an offset).
//
// For example, "2017-01-15T01:30:15.01Z" encodes 15.01 seconds past
// 01:30 UTC on January 15, 2017.
//
// In JavaScript, one can convert a Date object to this format using the
// standard [toISOString()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Date/toISOString]
// method. In Python, a standard `datetime.datetime` object can be converted
// to this format using [`strftime`](https://docs.python.org/2/library/time.html#time.strftime)
// with the time format spec '%Y-%m-%dT%H:%M:%S.%fZ'. Likewise, in Java, one
// can use the Joda Time's [`ISODateTimeFormat.dateTime()`](
// http://www.joda.org/joda-time/apidocs/org/joda/time/format/ISODateTimeFormat.html#dateTime--
// ) to obtain a formatter capable of generating timestamps in this format.
//
//
type Timestamp struct {
	// Represents seconds of UTC time since Unix epoch
	// 1970-01-01T00:00:00Z. Must be from 0001-01-01T00:00:00Z to
	// 9999-12-31T23:59:59Z inclusive.
	Seconds int64 `protobuf:"varint,1,opt,name=seconds,proto3" json:"seconds,omitempty"`
	// Non-negative fractions of a second at nanosecond resolution. Negative
	// second values with fractions must still have non-negative nanos values
	// that count forward in time. Must be from 0 to 999,999,999
	// inclusive.
	Nanos                int32    `protobuf:"varint,2,opt,name=nanos,proto3" json:"nanos,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Timestamp) Reset()         { *m = Timestamp{} }
func (m *Timestamp) String() string { return proto.CompactTextString(m) }
func (*Timestamp) ProtoMessage()    {}
func (*Timestamp) Descriptor() ([]byte, []int) {
	return fileDescriptor_edac929d8ae1e24f, []int{0}
}

func (m *Timestamp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Timestamp.Unmarshal(m, b)
}
func (m *Timestamp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Timestamp.Marshal(b, m, deterministic)
}
func (m *Timestamp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Timestamp.Merge(m, src)
}
func (m *Timestamp) XXX_Size() int {
	return xxx_messageInfo_Timestamp.Size(m)
}
func (m *Timestamp) XXX_DiscardUnknown() {
	xxx_messageInfo_Timestamp.DiscardUnknown(m)
}

var xxx_messageInfo_Timestamp proto.InternalMessageInfo

func (m *Timestamp) GetSeconds() int64 {
	if m != nil {
		return m.Seconds
	}
	return 0
}

func (m *Timestamp) GetNanos() int32 {
	if m != nil {
		return m.Nanos
	}
	return 0
}

func init() {
	proto.RegisterType((*Timestamp)(nil), "sproto.Timestamp")
}

func init() { proto.RegisterFile("timestamp.proto", fileDescriptor_edac929d8ae1e24f) }

var fileDescriptor_edac929d8ae1e24f = []byte{
	// 191 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2f, 0xc9, 0xcc, 0x4d,
	0x2d, 0x2e, 0x49, 0xcc, 0x2d, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2b, 0x06, 0xd3,
	0x4a, 0xd6, 0x5c, 0x9c, 0x21, 0x30, 0x29, 0x21, 0x09, 0x2e, 0xf6, 0xe2, 0xd4, 0xe4, 0xfc, 0xbc,
	0x94, 0x62, 0x09, 0x46, 0x05, 0x46, 0x0d, 0xe6, 0x20, 0x18, 0x57, 0x48, 0x84, 0x8b, 0x35, 0x2f,
	0x31, 0x2f, 0xbf, 0x58, 0x82, 0x49, 0x81, 0x51, 0x83, 0x35, 0x08, 0xc2, 0x71, 0x6a, 0x60, 0xe4,
	0x12, 0x4d, 0xce, 0xcf, 0xd5, 0x4b, 0xcf, 0x2c, 0xc9, 0x28, 0x4d, 0xd2, 0x2b, 0xae, 0xcc, 0x4b,
	0x4e, 0x2b, 0x2d, 0x29, 0x2d, 0x4a, 0x75, 0xe2, 0x83, 0x1b, 0x1a, 0x00, 0xb2, 0x26, 0x80, 0x31,
	0x4a, 0x13, 0xaa, 0x28, 0x39, 0x3f, 0x57, 0x1f, 0xa1, 0x50, 0x3f, 0x3d, 0x5f, 0x1f, 0xe2, 0x12,
	0x7d, 0xb8, 0x0b, 0x7f, 0x30, 0x32, 0x2e, 0x62, 0x62, 0x0e, 0x0b, 0x70, 0x5a, 0xc5, 0xa4, 0x10,
	0x5c, 0x99, 0x97, 0x1c, 0x9c, 0x9f, 0x56, 0xa2, 0x07, 0x36, 0x26, 0xa9, 0x34, 0x4d, 0x2f, 0x3c,
	0x35, 0x27, 0xc7, 0x3b, 0x2f, 0xbf, 0x3c, 0x2f, 0xa4, 0xb2, 0x20, 0xb5, 0x38, 0x89, 0x0d, 0xac,
	0xd9, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x56, 0xd0, 0x49, 0x34, 0xe1, 0x00, 0x00, 0x00,
}
