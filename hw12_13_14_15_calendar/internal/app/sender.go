package app

import (
	"context"
	"fmt"
	"sync"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/closer"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/config"
)

// Sender структура рассыльщика с сервис-провайдером.
type Sender struct {
	serviceProvider *serviceProvider
}

// NewSender вернуть новый экземпляр рассыльщика с зависимостями.
func NewSender(ctx context.Context, cfg string) (*Sender, error) {
	configFileName = cfg
	a := &Sender{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run запуск рассыльщика.
func (s *Sender) Run() error {
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
		_ = s.serviceProvider.NotificationService(context.Background()).Send(context.Background())
	}()

	wg.Wait()

	return nil
}

func (s *Sender) initDeps(ctx context.Context) error {
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

func (s *Sender) initConfig(_ context.Context) error {
	if _, err := config.New(configFileName); err != nil {
		return err
	}

	return nil
}

func (s *Sender) initLogger(_ context.Context) error {
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

func (s *Sender) initServiceProvider(_ context.Context) error {
	s.serviceProvider = newServiceProvider()
	return nil
}
