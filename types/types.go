package types

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"gitlab.com/InfoBlogFriends/server/decoder"
)

const (
	FilePost = iota + 1
	FileAvatar
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

type NullTime struct {
	sql.NullTime
}

//MarshalJSON method is called by json.Marshal,
//whenever it is of type NullString
func (x *NullTime) MarshalJSON() ([]byte, error) {
	if !x.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(x.Time)
}

func (x *NullTime) UnmarshalJSON(data []byte) error {
	err := decoder.NewDecoder().Decode(bytes.NewBuffer(data), &x.Time)
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

	var source []byte
	switch t := value.(type) {
	case string:
		*x = NewNullUUID(t)
		return nil
	case []byte:
		if len(t) == 0 {
			source = nil
		} else {
			source = t
		}
	case nil:
		*x = NullUUID{}
	default:
		return errors.New("incompatible type for NullUUID")
	}

	uuidRaw, err := uuid.FromBytes(source)
	if err != nil {
		return err
	}
	x.Binary = source
	x.Valid = true

	x.String = uuidRaw.String()

	return nil
}

// Value implements the driver Valuer interface.
func (x NullUUID) Value() (driver.Value, error) {
	if !x.Valid {
		return nil, nil
	}

	return x.Binary, nil
}

type User struct{}
