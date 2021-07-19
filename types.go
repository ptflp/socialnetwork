package infoblog

import (
	"database/sql"
	"encoding/json"
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
