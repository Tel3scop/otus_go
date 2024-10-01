package app

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/closer"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/rs/cors"
	"go.uber.org/zap"
)

// App структура приложения с сервис-провайдером и GRPC-сервером.
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

var configFileName string

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// NewApp вернуть новый экземпляр приложения с зависимостями.
func NewApp(ctx context.Context, cfg string) (*App, error) {
	configFileName = cfg
	a := &App{}
	if err := a.initDeps(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

// Run запуск приложения.
func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()
	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		err := a.runHTTPServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initLogger,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	if _, err := config.New(configFileName); err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.InitByParams(
		a.serviceProvider.Config().Log.FileName,
		a.serviceProvider.Config().Log.Level,
		a.serviceProvider.Config().Log.MaxSize,
		a.serviceProvider.Config().Log.MaxBackups,
		a.serviceProvider.Config().Log.MaxAge,
		a.serviceProvider.Config().Log.Compress,
		a.serviceProvider.Config().Log.StdOut,
	)
	logger.Info("logger is enabled")
	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initHTTPServer(_ context.Context) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
		AllowCredentials: true,
	})

	// Создаем middleware для логирования
	loggerMiddleware := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			latency := time.Since(start)
			logger.Info("Request processed",
				zap.String("IP", r.RemoteAddr),
				zap.String("Time", time.Now().Format(time.RFC3339)),
				zap.String("Method", r.Method),
				zap.String("Path", r.URL.Path),
				zap.String("HTTP Version", r.Proto),
				zap.Int("Status", rw.statusCode),
				zap.Duration("Latency", latency),
				zap.String("User-Agent", r.UserAgent()),
			)
		})
	}

	// Создаем HTTP-сервер с применением middleware
	a.httpServer = &http.Server{
		Addr:              a.serviceProvider.Config().HTTP.Address,
		Handler:           loggerMiddleware(corsMiddleware.Handler(mux)),
		ReadHeaderTimeout: 3 * time.Second,
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.httpServer.Addr)

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
