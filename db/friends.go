package db

import (
	"context"
	"gitlab.com/InfoBlogFriends/server/types"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type FriendRepository struct {
	db *sqlx.DB
	crud
}

func (c *FriendRepository) Create(ctx context.Context, friend infoblog.Friend) error {
	return c.create(ctx, &friend)
}

func (c *FriendRepository) Find(ctx context.Context, friend infoblog.Friend) (infoblog.Friend, error) {
	err := c.find(ctx, &friend, &friend)

	return friend, err
}

func (c *FriendRepository) Update(ctx context.Context, friend infoblog.Friend) error {
	return c.update(ctx, &friend)
}

func (c *FriendRepository) Delete(ctx context.Context, friend infoblog.Friend) error {
	friend.Active = types.NullBool{}
	return c.update(ctx, &friend)
}

func (c *FriendRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.Friend, error) {
	var friend []infoblog.Friend
	err := c.list(ctx, &friend, &infoblog.Friend{}, limit, offset)

	return friend, err
}

func (f *FriendRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.Friend, error) {
	var friends []infoblog.Friend
	err := f.crud.listx(ctx, &friends, infoblog.Friend{}, condition)
	if err != nil {
		return nil, err
	}

	return friends, nil
}

func NewFriendRepository(db *sqlx.DB) infoblog.FriendRepository {
	cr := crud{db: db}
	return &FriendRepository{db: db, crud: cr}
}
