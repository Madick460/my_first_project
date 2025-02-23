package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/Rasikrr/my_project/configs"
	"github.com/Rasikrr/my_project/internal/cache"
	"github.com/Rasikrr/my_project/internal/cache/answers"
	"github.com/Rasikrr/my_project/internal/clients/gpt"
	"github.com/Rasikrr/my_project/internal/databases"
	http "github.com/Rasikrr/my_project/internal/ports/http"
	mcqS "github.com/Rasikrr/my_project/internal/services/mcq"
	"github.com/Rasikrr/my_project/internal/util"
	"github.com/Rasikrr/my_project/internal/workers"
	"github.com/hashicorp/go-multierror"
	"github.com/redis/go-redis/v9"
	"log"
	HTTP "net/http"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
)

type Starter interface {
	Start() error
}

type App struct {
	name     string
	config   configs.Config
	postgres *databases.Postgres

	workers []workers.Worker

	redisClient  *redis.Client
	cacheClient  cache.Cache
	answersCache answers.Cache
	hasher       util.Hasher
	gptClient    gpt.Client

	mcqService mcqS.Service
	httpServer *http.Server
}

// nolint: gocritic
func InitApp(ctx context.Context, name string) *App {
	log.Printf("starting app initialization: %s\n", name)
	log.Printf("go version: %s\n", runtime.Version())
	cfg, err := configs.Parse()
	if err != nil {
		panic(err)
	}
	app := &App{
		name:   name,
		config: cfg,
	}
	for _, init := range []func(ctx context.Context) error{
		app.InitPostgres,
		app.InitRedis,
		app.InitRepositories,
		app.InitUtil,
		app.InitCache,
		app.InitClients,
		app.InitServices,
		app.InitHTTPServer,
		//app.InitWorkers,
	} {
		if err := init(ctx); err != nil {
			log.Fatalf("init app error: %v", err)
		}
	}
	return app
}

func (a *App) InitRepositories(_ context.Context) error {
	return nil
}

func (a *App) InitUtil(_ context.Context) error {
	a.hasher = util.NewHasher()
	return nil
}

func (a *App) InitClients(_ context.Context) error {
	a.gptClient = gpt.NewClient()
	return nil
}

func (a *App) InitRedis(ctx context.Context) error {
	var err error
	a.redisClient, err = databases.NewRedis(ctx, &a.config)
	if err != nil {
		return fmt.Errorf("failed to init redis: %w", err)
	}
	log.Println("Redis connected")
	return nil
}

func (a *App) InitCache(_ context.Context) error {
	a.cacheClient = cache.NewRedisCache(a.redisClient)
	a.answersCache = answers.NewCache(a.cacheClient)
	return nil
}

func (a *App) InitServices(_ context.Context) error {
	a.mcqService = mcqS.NewService(
		a.gptClient,
		a.answersCache,
	)
	return nil
}

func (a *App) InitHTTPServer(_ context.Context) error {
	a.httpServer = http.NewServer(
		&a.config,
		a.mcqService,
	)
	return nil
}

// nolint
//func (a *App) InitWorkers(_ context.Context) error {
//	a.workers = []workers.Worker{
//		dbInfo.NewWorker(a.usersRepository),
//	}
//	return nil
//}

func (a *App) Start(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	stop := make(chan struct{})

	go a.handleShutdown(ctx, cancel, stop)

	fns := make([]any, 0, len(a.workers)+1)
	fns = append(fns, a.httpServer.Start)
	for _, w := range a.workers {
		fns = append(fns, w.Run)
	}

	if err := runParallel(ctx, fns...); err != nil {
		return err
	}
	<-stop
	return nil
}

func (a *App) handleShutdown(_ context.Context, cancel context.CancelFunc, s chan struct{}) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	log.Println("Received shutdown signal")

	if err := a.httpServer.Shutdown(context.Background()); err != nil {
		log.Println("Error while shutting down HTTP server:", err)
	}

	a.postgres.Close()

	if err := a.redisClient.Close(); err != nil {
		log.Println("Error while closing redis client:", err)
	}
	log.Println("Redis client closed gracefully")
	cancel()
	close(s)
}

func (a *App) InitPostgres(ctx context.Context) error {
	//var err error
	//a.postgres, err = databases.NewPostgres(ctx, &a.config)
	//if err != nil {
	//	return fmt.Errorf("failed to init postgres: %w", err)
	//}
	//log.Println("Postgres connected")
	return nil
}

func runParallel(ctx context.Context, fns ...any) error {
	var (
		wg     sync.WaitGroup
		errRes *multierror.Error
	)
	for _, f := range fns {
		wg.Add(1)
		go func(fn any) {
			defer wg.Done()
			switch fun := fn.(type) {
			case func() error:
				if err := fun(); err != nil {
					if errors.Is(err, HTTP.ErrServerClosed) {
						log.Println("HTTP server closed gracefully")
						return
					}
					errRes = multierror.Append(errRes, err)
				}
			case func(ctx context.Context) error:
				if err := fun(ctx); err != nil {
					errRes = multierror.Append(errRes, err)
				}
			default:
				log.Printf("unknown function type: %T", fn)
			}
		}(f)
	}
	wg.Wait()
	return errRes.ErrorOrNil()
}
