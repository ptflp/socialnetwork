package infoblog

import "database/sql"

func NewNullString(s string) sql.NullString {
	if len(s) == 0 {
		return sql.NullString{}
	}
	return sql.NullString{
		String: s,
		Valid:  true,
	}
}

func NewNullInt64(n int64) sql.NullInt64 {
	if n == 0 {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: n,
		Valid: true,
	}
}

func NewNullBool(b bool) sql.NullBool {
	if !b {
		return sql.NullBool{}
	}
	return sql.NullBool{
		Bool:  b,
		Valid: true,
	}
}
