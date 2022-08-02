package timestamp

import (
	"database/sql/driver"
	"time"
)

// Scan implements the Scanner interface of the database driver
func (m *Timestamp) Scan(value interface{}) error {

	// convert the interface to a time type
	dbTime, ok := value.(time.Time)

	if ok {
		_, err := m.SetToTime(dbTime)
		return err
	}

	return nil
}

// Value implements the db driver Valuer interface
func (m Timestamp) Value() (driver.Value, error) {
	return m.Time()
}

// XXX_WellKnownType allows this non-Google timestamp protobuf to interact properly with jsonpb
// and other libraries as it is was the standard WKT timestamp.
func (*Timestamp) XXX_WellKnownType() string { return "Timestamp" }
