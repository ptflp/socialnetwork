package infoblog

import (
	"database/sql"
)

func NewNullString(s string) NullString {
	if len(s) == 0 {
		return NullString{}
	}
	return NullString{
		sql.NullString{
			String: s,
			Valid:  true,
		},
	}
}

func NewNullInt64(n int64) NullInt64 {
	if n == 0 {
		return NullInt64{}
	}
	return NullInt64{
		sql.NullInt64{
			Int64: n,
			Valid: true,
		},
	}
}

func NewNullFloat64(n float64) NullFloat64 {
	if n == 0 {
		return NullFloat64{}
	}
	return NullFloat64{
		sql.NullFloat64{
			Float64: n,
			Valid:   true,
		},
	}
}

func NewNullBool(b bool) NullBool {
	if !b {
		return NullBool{}
	}
	return NullBool{
		sql.NullBool{
			Bool:  b,
			Valid: true,
		},
	}
}
