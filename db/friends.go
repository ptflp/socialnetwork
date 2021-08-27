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

func (f *FriendRepository) Create(ctx context.Context, friend infoblog.Friend) error {
	return f.create(ctx, &friend)
}

func (f *FriendRepository) Find(ctx context.Context, friend infoblog.Friend) (infoblog.Friend, error) {
	err := f.find(ctx, &friend, &friend)

	return friend, err
}

func (f *FriendRepository) Update(ctx context.Context, friend infoblog.Friend) error {
	return f.update(ctx, &friend)
}

func (f *FriendRepository) Delete(ctx context.Context, friend infoblog.Friend) error {
	friend.Active = types.NullBool{}
	return f.update(ctx, &friend)
}

func (f *FriendRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.Friend, error) {
	var friend []infoblog.Friend
	err := f.list(ctx, &friend, &infoblog.Friend{}, limit, offset)

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
