package db

import (
	"context"
	"gitlab.com/InfoBlogFriends/server/types"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type HashTagRepository struct {
	db *sqlx.DB
	crud
}

func (h *HashTagRepository) Create(ctx context.Context, hashtag infoblog.HashTag) error {
	return h.create(ctx, &hashtag)
}

func (h *HashTagRepository) Find(ctx context.Context, hashtag infoblog.HashTag) (infoblog.HashTag, error) {
	err := h.find(ctx, &hashtag, &hashtag)

	return hashtag, err
}

func (h *HashTagRepository) Update(ctx context.Context, hashtag infoblog.HashTag) error {
	return h.update(ctx, &hashtag)
}

func (h *HashTagRepository) Delete(ctx context.Context, hashtag infoblog.HashTag) error {
	hashtag.Active = types.NullBool{}
	return h.update(ctx, &hashtag)
}

func (h *HashTagRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.HashTag, error) {
	var hashtag []infoblog.HashTag
	err := h.list(ctx, &hashtag, &infoblog.HashTag{}, limit, offset)

	return hashtag, err
}

func (h *HashTagRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.HashTag, error) {
	var hashtags []infoblog.HashTag
	err := h.crud.listx(ctx, &hashtags, infoblog.HashTag{}, condition)
	if err != nil {
		return nil, err
	}

	return hashtags, nil
}

func NewHashTagRepository(db *sqlx.DB) infoblog.HashTagRepository {
	cr := crud{db: db}
	return &HashTagRepository{db: db, crud: cr}
}
