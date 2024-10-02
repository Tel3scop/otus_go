package app

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/Tel3scop/helpers/logger"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/closer"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/config"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/interceptor"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/pkg/event_v1"
	// Register statik for swagger UI.
	_ "github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/statik"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

// App структура приложения с сервис-провайдером и GRPC-сервером.
type App struct {
	serviceProvider *serviceProvider
	httpServer      *http.Server
	grpcServer      *grpc.Server
	swaggerServer   *http.Server
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
	wg.Add(3)

	go func() {
		defer wg.Done()
		err := a.runHTTPServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runSwaggerServer()
		if err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := a.runGRPCServer()
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
		a.initGRPCServer,
		a.initHTTPServer,
		a.initSwaggerServer,
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

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(grpcMiddleware.ChainUnaryServer(
			interceptor.LogInterceptor,
		)),
	)

	reflection.Register(a.grpcServer)
	event_v1.RegisterEventServiceServer(a.grpcServer, a.serviceProvider.EventImpl(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on %s", a.serviceProvider.Config().GRPC.Address)

	list, err := net.Listen("tcp", a.serviceProvider.Config().GRPC.Address)
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := event_v1.RegisterEventServiceHandlerFromEndpoint(ctx, mux, a.serviceProvider.Config().HTTP.Address, opts)
	if err != nil {
		return err
	}

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

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api_event.swagger.json", serveSwaggerFile("/api_event.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:              a.serviceProvider.Config().Swagger.Address,
		Handler:           mux,
		ReadHeaderTimeout: time.Duration(a.serviceProvider.Config().Swagger.Timeout) * time.Second,
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", a.swaggerServer.Addr)

	err := a.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func(file http.File) {
			err = file.Close()
			if err != nil {
				log.Fatal(err)
			}
		}(file)

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
