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

func (e *EventRepository) Create(ctx context.Context, event infoblog.Event) error {
	return e.create(ctx, &event)
}

func (e *EventRepository) Find(ctx context.Context, event infoblog.Event) (infoblog.Event, error) {
	err := e.find(ctx, &event, &event)

	return event, err
}

func (e *EventRepository) Update(ctx context.Context, event infoblog.Event) error {
	return e.update(ctx, &event)
}

func (e *EventRepository) List(ctx context.Context, limit, offset uint64) ([]infoblog.Event, error) {
	var events []infoblog.Event
	err := e.list(ctx, &events, &infoblog.Event{}, limit, offset)

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
