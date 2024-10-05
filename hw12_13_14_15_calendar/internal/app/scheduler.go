package app

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/closer"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

// Scheduler структура планировщика с сервис-провайдером и кроном.
type Scheduler struct {
	serviceProvider *serviceProvider
}

// NewScheduler вернуть новый экземпляр шедулера с зависимостями.
func NewScheduler(ctx context.Context, cfg string) (*Scheduler, error) {
	configFileName = cfg
	a := &Scheduler{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run запуск шедулера.
func (s *Scheduler) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)

	err := s.serviceProvider.QueueService(context.Background()).CreateQueue(context.Background())
	if err != nil {
		return fmt.Errorf("failed to create queue: %w", err)
	}

	go func() {
		defer wg.Done()
		err = s.runCron()
	}()

	wg.Wait()

	return nil
}

func (s *Scheduler) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		s.initConfig,
		s.initServiceProvider,
		s.initLogger,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Scheduler) initConfig(_ context.Context) error {
	if _, err := config.New(configFileName); err != nil {
		return err
	}

	return nil
}

func (s *Scheduler) initLogger(_ context.Context) error {
	logger.InitByParams(
		s.serviceProvider.Config().Log.FileName,
		s.serviceProvider.Config().Log.Level,
		s.serviceProvider.Config().Log.MaxSize,
		s.serviceProvider.Config().Log.MaxBackups,
		s.serviceProvider.Config().Log.MaxAge,
		s.serviceProvider.Config().Log.Compress,
		s.serviceProvider.Config().Log.StdOut,
	)
	logger.Info("logger is enabled")
	return nil
}

func (s *Scheduler) initServiceProvider(_ context.Context) error {
	s.serviceProvider = newServiceProvider()
	return nil
}

func (s *Scheduler) runCron() error {
	log.Printf("Cron server is starting")
	ctx := context.Background()
	cron, err := gocron.NewScheduler()
	if err != nil {
		return fmt.Errorf("failed to create scheduler: %w", err)
	}

	s.notifyOnEventJob(ctx, cron)
	s.clearOldEvents(ctx, cron)

	cron.Start()

	time.Sleep(time.Minute)

	err = cron.Shutdown()
	if err != nil {
		logger.Error("cannot shutdown cron: %s", zap.Error(err))
	}

	return nil
}

func (s *Scheduler) notifyOnEventJob(ctx context.Context, cron gocron.Scheduler) {
	j, err := cron.NewJob(
		gocron.DurationJob(
			20*time.Second,
		),
		// gocron.DurationJob(
		//	5*time.Minute,
		// ),
		gocron.NewTask(
			func() {
				logger.Info("try")
				date := time.Now()
				err := s.serviceProvider.NotificationService(ctx).NotifyOnEvent(ctx)
				if err != nil {
					logger.Error("cannot notifiy events", zap.Time("date", date), zap.Error(err))
				}
			},
		),
	)
	if err != nil {
		log.Fatalf("cannot start job: %s", err)
	}

	logger.Info("job created", zap.String("ID", j.ID().String()))
}

func (s *Scheduler) clearOldEvents(ctx context.Context, cron gocron.Scheduler) {
	j, err := cron.NewJob(
		gocron.DailyJob(
			1,
			gocron.NewAtTimes(gocron.NewAtTime(23, 0, 0)),
		),
		gocron.NewTask(
			func() {
				err := s.serviceProvider.EventService(ctx).DeleteByDate(
					ctx,
					time.Now().AddDate(-1, 0, 0),
				)
				if err != nil {
					logger.Error("cannot clear old events", zap.Error(err))
				}
			},
		),
	)
	if err != nil {
		log.Fatalf("cannot start job: %s", err)
	}

	logger.Info("job created", zap.String("ID", j.ID().String()))
}
