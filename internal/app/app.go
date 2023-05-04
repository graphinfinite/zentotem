package app

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"
	"zentotem/config"
	"zentotem/internal/controllers"
	"zentotem/internal/migrations"
	"zentotem/internal/repository"

	"zentotem/pkg/postgres"
	redisclient "zentotem/pkg/redis"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
)

type App struct {
	httpServer *http.Server
}

func NewApp() *App {
	return &App{}
}

func (app *App) Run(cfg *config.Configuration) error {
	//
	//Logger
	//
	zlog := zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}).With().Timestamp().Logger()
	zlog.Info().Msgf("run with config: %#v", cfg)

	//
	// Postgres
	//
	zlog.Info().Msgf("connect to db...")
	db, err := postgres.NewPG(cfg.Psql)
	if err != nil {
		zlog.Err(err).Msg("App/NewPostgresPool ")
		return err
	}
	defer db.Db.Close()
	if err := migrations.Up(db.Db); err != nil {
		zlog.Err(err).Msg("Migrations ")
		return err
	}

	//
	//Redis
	//
	redclient := redisclient.InitDefault(cfg.RedisHost, cfg.RedisPwd, cfg.RedisDb, zlog)
	fmt.Print(redclient)

	// Layers
	user := controllers.NewUserController(repository.NewUserStorage(db), zlog)
	cash := controllers.NewMemoryCashController(redisclient.InitDefault(cfg.RedisHost, cfg.PgPwd, cfg.RedisDb, zlog), zlog)
	cript := controllers.NewCriptoController(zlog)

	//
	// Router
	//
	r := chi.NewRouter()
	r.Use(middleware.Timeout(5 * time.Second))
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("pong")) })
	r.Post("/redis/incr", cash.IncrementValue)
	r.Post("/sign/hmacsha512", cript.SingHmacsha512)
	r.Post("/postgres/users", user.AddUser)

	//
	// Http Server
	//
	app.httpServer = &http.Server{
		Addr:           net.JoinHostPort(cfg.ServerHost, cfg.ServerPort),
		Handler:        r,
		ReadTimeout:    20 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	zlog.Info().Msgf("server start... %s", app.httpServer.Addr)
	go func() {
		if err := app.httpServer.ListenAndServe(); err != nil {
			zlog.Err(err).Msg("💀")
		}
	}()

	// Shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return app.httpServer.Shutdown(ctx)

}
