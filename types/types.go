package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/volatiletech/null/v8"

	"github.com/google/uuid"

	"gitlab.com/InfoBlogFriends/server/decoder"
)

const (
	TypePost = iota + 1
	TypeAvatar
	TypeUserModerate
)

//NullString is a wrapper around sql.NullString
type NullString struct {
	sql.NullString
}

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullString
func (x *NullString) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", x.String)), nil
}

func (x *NullString) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.String)
	if err != nil {
		return err
	}
	x.Valid = true
	if len(x.String) == 0 {
		x.Valid = false
	}

	return nil
}

type NullBool struct {
	sql.NullBool
}

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullString
func (x *NullBool) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Bool)
}

func (x *NullBool) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.Bool)
	if err != nil {
		return err
	}
	x.Valid = true

	return nil
}

type NullInt64 struct {
	sql.NullInt64
}

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullString
func (x *NullInt64) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Int64)
}

func (x *NullInt64) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.Int64)
	if err != nil {
		return err
	}
	x.Valid = true

	return nil
}

type NullFloat64 struct {
	sql.NullFloat64
}

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullFloat64
func (x *NullFloat64) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Float64)
}

func (x *NullFloat64) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.Float64)
	if err != nil {
		return err
	}
	x.Valid = true

	return nil
}

type NullUUID struct {
	Binary []byte
	Valid  bool
	String string
}

func (x *NullUUID) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", x.String)), nil
}

func (x *NullUUID) UnmarshalJSON(data []byte) error {
	if data[0] == 110 {
		*x = NullUUID{}
		return nil
	}
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.String)
	if err != nil {
		return err
	}
	x.Valid = true
	uuidRaw, err := uuid.Parse(x.String)
	if err != nil {
		return err
	}
	x.Binary, err = uuidRaw.MarshalBinary()
	if err != nil {
		return err
	}

	return nil
}

// Scan implements the Scanner interface.
func (x *NullUUID) Scan(value interface{}) error {
	if value == nil {
		*x = NullUUID{}
		return nil
	}

	var dest []byte
	switch source := value.(type) {
	case string:
		*x = NewNullUUID(source)
		return nil
	case []byte:
		if len(source) == 0 {
			dest = nil
		} else {
			dest = make([]byte, len(source))
			copy(dest, source)
		}
	case nil:
		*x = NullUUID{}
	default:
		return errors.New("incompatible type for NullUUID")
	}

	uuidRaw, err := uuid.FromBytes(dest)
	if err != nil {
		return err
	}
	x.Binary = dest
	x.Valid = true

	x.String = uuidRaw.String()

	return nil
}

// Value implements the driver Valuer interface.
func (x NullUUID) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}
	b := make([]byte, len(x.Binary))
	copy(b, x.Binary)

	return b, nil
}

type NullUint64 struct {
	null.Uint64
}

type NullTime struct {
	null.Time
}

type User struct{}

type Chat struct{}
