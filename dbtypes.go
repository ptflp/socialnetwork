package infoblog

import (
	"database/sql"

	"gitlab.com/InfoBlogFriends/server/decoder"
)

func NewNullString(s string) NullString {
	if len(s) == 0 {
		return NullString{}
	}
	return NullString{
		decoder.NewDecoder(),
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
		decoder.NewDecoder(),
		sql.NullInt64{
			Int64: n,
			Valid: true,
		},
	}
}

func NewNullBool(b bool) NullBool {
	if !b {
		return NullBool{}
	}
	return NullBool{
		decoder.NewDecoder(),
		sql.NullBool{
			Bool:  b,
			Valid: true,
		},
	}
}
