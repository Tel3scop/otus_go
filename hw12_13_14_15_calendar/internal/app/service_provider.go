package app

import (
	"context"
	"log"

	eventApi "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/api/event"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/client/db"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/client/db/pg"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/client/db/transaction"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/closer"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/service"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/service/event"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/storage/sql"
)

type serviceProvider struct {
	config          *config.Config
	eventRepository storage.EventStorage
	eventService    service.EventService
	eventImpl       *eventApi.Implementation
	dbClient        db.Client
	txManager       db.TxManager
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) Config() *config.Config {
	if s.config == nil {
		cfg, err := config.New(configFileName)
		if err != nil {
			log.Fatalf("failed to get config: %s", err.Error())
		}
		s.config = cfg
	}
	return s.config
}

func (s *serviceProvider) DBClient(ctx context.Context) db.Client {
	if s.Config().Database != config.PostgresDatabaseType {
		return s.dbClient
	}

	if s.dbClient == nil {
		cl, err := pg.New(ctx, s.Config().Postgres.DSN)
		if err != nil {
			log.Fatalf("failed to create db client: %v", err)
		}

		err = cl.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}
		closer.Add(cl.Close)

		s.dbClient = cl
	}

	return s.dbClient
}

func (s *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if s.txManager != nil {
		return s.txManager
	}
	switch s.Config().Database {
	case config.MemoryDatabaseType:
		s.txManager = &memorystorage.NoOpTxManager{}
	case config.PostgresDatabaseType:
		s.txManager = transaction.NewTransactionManager(s.DBClient(ctx).DB())
	default:
		log.Fatalf("unknown db type: %s", s.Config().Database)
	}

	return s.txManager
}

func (s *serviceProvider) EventRepository(ctx context.Context) storage.EventStorage {
	if s.eventRepository != nil {
		return s.eventRepository
	}
	cfg := s.Config()
	switch cfg.Database {
	case config.MemoryDatabaseType:
		s.eventRepository = memorystorage.NewInMemoryEventStorage()
	case config.PostgresDatabaseType:
		s.eventRepository = sqlstorage.NewRepository(s.DBClient(ctx))
	default:
		log.Fatalf("unknown db type: %s", s.Config().Database)
	}

	return s.eventRepository
}

func (s *serviceProvider) EventService(ctx context.Context) service.EventService {
	if s.eventService == nil {
		s.eventService = event.NewService(
			s.EventRepository(ctx),
			s.TxManager(ctx),
		)
	}

	return s.eventService
}

func (s *serviceProvider) EventImpl(ctx context.Context) *eventApi.Implementation {
	if s.eventImpl == nil {
		s.eventImpl = eventApi.NewImplementation(s.EventService(ctx))
	}

	return s.eventImpl
}
