package sqlstorage

import (
	"context"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/client/db"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/entity"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage"
	"github.com/google/uuid"
)

const (
	tableName = "events"

	columnID           = "id"
	columnTitle        = "title"
	columnDateTime     = "datetime"
	columnDuration     = "duration"
	columnDescription  = "description"
	columnUserID       = "user_id"
	columnNotifyBefore = "notify_before"
)

type repo struct {
	db db.Client
}

// NewRepository создание репозитория для событий.
func NewRepository(db db.Client) storage.EventStorage {
	return &repo{db: db}
}

// Create создание события.
func (r *repo) Create(ctx context.Context, event entity.Event) (string, error) {
	event.ID = uuid.New().String()

	builder := sq.Insert(tableName).
		Columns(columnID, columnTitle, columnDateTime, columnDuration, columnDescription, columnUserID, columnNotifyBefore).
		Values(event.ID, event.Title, event.DateTime, event.Duration, event.Description, event.UserID, event.NotifyBefore).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	q := db.Query{
		Name:     "event.Create",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Error(err.Error())
		return "", err
	}

	return event.ID, nil
}

// Update обновление события.
func (r *repo) Update(ctx context.Context, eventID string, event entity.Event) error {
	builder := sq.Update(tableName).
		Set(columnTitle, event.Title).
		Set(columnDateTime, event.DateTime).
		Set(columnDuration, event.Duration).
		Set(columnDescription, event.Description).
		Set(columnUserID, event.UserID).
		Set(columnNotifyBefore, event.NotifyBefore).
		Where(sq.Eq{columnID: eventID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	q := db.Query{
		Name:     "event.Update",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

// Delete удаление события.
func (r *repo) Delete(ctx context.Context, eventID string) error {
	builder := sq.Delete(tableName).
		Where(sq.Eq{columnID: eventID}).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error(err.Error())
		return err
	}

	q := db.Query{
		Name:     "event.Delete",
		QueryRaw: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		logger.Error(err.Error())
	}

	return nil
}

// List список событий на определенный период.
func (r *repo) List(ctx context.Context, date time.Time, period entity.PeriodType) ([]entity.Event, error) {
	var builder sq.SelectBuilder

	switch period {
	case entity.PeriodDay:
		builder = sq.Select(
			columnID,
			columnTitle,
			columnDateTime,
			columnDuration,
			columnDescription,
			columnUserID,
			columnNotifyBefore,
		).
			From(tableName).
			Where(sq.Expr("DATE_TRUNC('day', "+columnDateTime+") = ?", date)).
			PlaceholderFormat(sq.Dollar)
	case entity.PeriodWeek:
		startOfWeek := date.AddDate(0, 0, -int(date.Weekday()))
		endOfWeek := startOfWeek.AddDate(0, 0, 7)
		builder = sq.Select(
			columnID,
			columnTitle,
			columnDateTime,
			columnDuration,
			columnDescription,
			columnUserID,
			columnNotifyBefore,
		).
			From(tableName).
			Where(sq.Expr(columnDateTime+" >= ? AND "+columnDateTime+" < ?", startOfWeek, endOfWeek)).
			PlaceholderFormat(sq.Dollar)
	case entity.PeriodMonth:
		startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		endOfMonth := startOfMonth.AddDate(0, 1, 0)
		builder = sq.Select(
			columnID,
			columnTitle,
			columnDateTime,
			columnDuration,
			columnDescription,
			columnUserID,
			columnNotifyBefore,
		).
			From(tableName).
			Where(sq.Expr(columnDateTime+" >= ? AND "+columnDateTime+" < ?", startOfMonth, endOfMonth)).
			PlaceholderFormat(sq.Dollar)
	default:
		return nil, storage.ErrInvalidPeriod
	}

	query, args, err := builder.ToSql()
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	q := db.Query{
		Name:     "event.ListEvents",
		QueryRaw: query,
	}
	var events []entity.Event
	err = r.db.DB().ScanAllContext(ctx, &events, q, args...)
	if err != nil {
		logger.Error(err.Error())
		return nil, err
	}

	return events, nil
}
