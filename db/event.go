package db

import (
	"context"

	"github.com/jmoiron/sqlx"
	infoblog "gitlab.com/InfoBlogFriends/server"
)

type EventRepository struct {
	db *sqlx.DB
	crud
}

func (c *EventRepository) Create(ctx context.Context, event infoblog.Event) error {
	return c.create(ctx, &event)
}

func (c *EventRepository) Find(ctx context.Context, event infoblog.Event) (infoblog.Event, error) {
	err := c.find(ctx, &event, &event)

	return event, err
}

func (c *EventRepository) Update(ctx context.Context, event infoblog.Event) error {
	return c.update(ctx, &event)
}

func (c *EventRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.Event, error) {
	var events []infoblog.Event
	err := c.list(ctx, &events, &infoblog.Event{}, limit, offset)

	return events, err
}

func (e *EventRepository) Listx(ctx context.Context, condition infoblog.Condition) ([]infoblog.Event, error) {
	var events []infoblog.Event
	err := e.crud.listx(ctx, &events, infoblog.Event{}, condition)
	if err != nil {
		return nil, err
	}

	return events, nil
}

func NewEventRepository(db *sqlx.DB) infoblog.EventRepository {
	cr := crud{db: db}
	return &EventRepository{db: db, crud: cr}
}
