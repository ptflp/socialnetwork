package db

import (
	"context"
	"database/sql"
	"strings"

	"gitlab.com/InfoBlogFriends/server/types"

	sq "github.com/Masterminds/squirrel"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type likesRepository struct {
	db *sqlx.DB
	crud
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
	fields, err := infoblog.GetFields(&infoblog.Like{})
	if err != nil {
		return infoblog.Like{}, err
	}

	queryRaw := sq.Select(fields...).From("likes").Where(sq.Eq{"type": like.Type, "foreign_uuid": like.ForeignUUID, "liker_uuid": like.LikerUUID})
	query, args, err := queryRaw.ToSql()
	if err != nil {
		return infoblog.Like{}, err
	}

	likeFound := infoblog.Like{}
	err = lr.db.QueryRowxContext(ctx, query, args...).StructScan(&likeFound)

	return likeFound, err
}

func (lr *likesRepository) CountByUser(ctx context.Context, user infoblog.User) (int64, error) {

	var count sql.NullInt64

	query, args, err := sq.Select("COUNT(user_uuid)").From("likes").Where(sq.Eq{"user_uuid": user.UUID, "active": 1}).ToSql()
	if err != nil {
		return count.Int64, err
	}

	err = lr.db.QueryRowContext(ctx, query, args...).Scan(&count)

	return count.Int64, err
}

func (lr *likesRepository) CountByPost(ctx context.Context, postUUID string) (uint64, error) {

	var count types.NullUint64

	query, args, err := sq.Select("COUNT(type)").From("likes").Where(sq.Eq{"type": 1, "foreign_uuid": types.NewNullUUID(postUUID), "active": 1}).ToSql()
	if err != nil {
		return count.Uint64.Uint64, err
	}

	err = lr.db.QueryRowContext(ctx, query, args...).Scan(&count)

	return count.Uint64.Uint64, err
}

func (l *likesRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.Like, error) {
	var likes []infoblog.Like
	err := l.crud.listx(ctx, &likes, infoblog.Like{}, condition)
	if err != nil {
		return nil, err
	}

	return likes, nil
}

func NewLikesRepository(db *sqlx.DB) infoblog.LikeRepository {
	return &likesRepository{db: db, crud: crud{db: db}}
}
