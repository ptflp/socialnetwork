package infoblog

import (
	"bytes"
	"database/sql"
	"encoding/json"

	"gitlab.com/InfoBlogFriends/server/decoder"
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
	return json.Marshal(x.String)
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
//whenever it is of type NullString
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
