package db

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type likesRepository struct {
	db *sqlx.DB
}

func (lr *likesRepository) Upsert(ctx context.Context, like infoblog.Like) error {
	createFields, err := infoblog.GetCreateFields("likes")
	if err != nil {
		return err
	}
	createFieldsPointers := infoblog.GetFieldsPointers(&like, "create")

	queryRaw := sq.Insert("likes").Columns(createFields...).Values(createFieldsPointers...)
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return err
	}
	args = append(args, &like.Active)
	query = strings.Join([]string{query, " ON DUPLICATE KEY UPDATE active = ?"}, "")
	_, err = lr.db.ExecContext(ctx, query, args...)

	return err
}

func (lr *likesRepository) Find(ctx context.Context, like *infoblog.Like) (infoblog.Like, error) {
	fields, err := infoblog.GetFields("likes")
	if err != nil {
		return infoblog.Like{}, err
	}

	queryRaw := sq.Select(fields...).From("likes").Where(sq.Eq{"type": like.Type, "foreign_uuid": like.ForeignUUID, "liker_uuid": like.LikerUUID})
	query, args, err := queryRaw.ToSql()

	likeFound := infoblog.Like{}
	err = lr.db.QueryRowxContext(ctx, query, args...).StructScan(&likeFound)

	return likeFound, err
}

func NewLikesRepository(db *sqlx.DB) infoblog.LikeRepository {
	return &likesRepository{db: db}
}
