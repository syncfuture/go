package timestamp

// XXX_WellKnownType allows this non-Google timestamp protobuf to interact properly with jsonpb
// and other libraries as it is was the standard WKT timestamp.
func (*Timestamp) XXX_WellKnownType() string { return "Timestamp" }
