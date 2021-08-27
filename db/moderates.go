package db

import (
	"context"
	"gitlab.com/InfoBlogFriends/server/types"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type ModerateRepository struct {
	db *sqlx.DB
	crud
}

func (m *ModerateRepository) Create(ctx context.Context, moderate infoblog.Moderate) error {
	return m.create(ctx, &moderate)
}

func (m *ModerateRepository) Find(ctx context.Context, moderate infoblog.Moderate) (infoblog.Moderate, error) {
	err := m.find(ctx, &moderate, &moderate)

	return moderate, err
}

func (m *ModerateRepository) Update(ctx context.Context, moderate infoblog.Moderate) error {
	return m.update(ctx, &moderate)
}

func (m *ModerateRepository) Delete(ctx context.Context, moderate infoblog.Moderate) error {
	moderate.Active = types.NullBool{}
	return m.update(ctx, &moderate)
}

func (m *ModerateRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.Moderate, error) {
	var moderate []infoblog.Moderate
	err := m.list(ctx, &moderate, &infoblog.Moderate{}, limit, offset)

	return moderate, err
}

func (m *ModerateRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.Moderate, error) {
	var moderates []infoblog.Moderate
	err := m.crud.listx(ctx, &moderates, infoblog.Moderate{}, condition)
	if err != nil {
		return nil, err
	}

	return moderates, nil
}

func NewModerateRepository(db *sqlx.DB) infoblog.ModerateRepository {
	cr := crud{db: db}
	return &ModerateRepository{db: db, crud: cr}
}
