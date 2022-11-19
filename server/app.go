package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/IrDeTen/money_processing_service.git/app"
	apphttp "github.com/IrDeTen/money_processing_service.git/app/delivery/http"
	appRepo "github.com/IrDeTen/money_processing_service.git/app/repo/postgres"
	appUC "github.com/IrDeTen/money_processing_service.git/app/usecase"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // required
	"github.com/spf13/viper"
)

type App struct {
	uc         app.IUsecase
	repo       app.IRepository
	httpServer *http.Server
}

func NewApp() *App {
	db := initDB()
	repo := appRepo.NewRepo(db)
	uc := appUC.NewUsecase(repo)
	return &App{
		uc:   uc,
		repo: repo,
	}
}

func (a *App) Run(port string) error {
	defer a.repo.Close()
	//TODO change gin mod
	gin.SetMode(gin.DebugMode)
	router := gin.New()
	router.Use(
		gin.RecoveryWithWriter( /*TODO add logger*/ nil),
	)

	if viper.GetBool("app.client.use") {
		router.Use(static.Serve("/", static.LocalFile(viper.GetString("app.client.dir"), true)))
	}

	apphttp.RegisterHTTPEndpoints(router, a.uc)

	a.httpServer = &http.Server{
		Addr:           ":" + port,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	l, err := net.Listen("tcp", a.httpServer.Addr)
	if err != nil {
		panic(err)
	}

	go func(l net.Listener) {
		if err := a.httpServer.Serve(l); err != nil {
			log.Fatalf("Failed to listen and serve: %+v", err)
		}
	}(l)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpServer.Shutdown(ctx)
}

// TODO add init db
func initDB() *sql.DB {
	dbString := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		viper.GetString("app.db.host"),
		viper.GetString("app.db.port"),
		viper.GetString("app.db.login"),
		viper.GetString("app.db.pass"),
		viper.GetString("app.db.name"),
	)
	db, err := sql.Open(
		"postgres",
		dbString,
	)
	if err != nil {
		panic(err)
	}
	runMigrations(db)
	return db
}

// TODO add run migrations
func runMigrations(db *sql.DB) {

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		panic(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		viper.GetString("app.db.name"),
		driver)
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion {
		fmt.Println(err)
	}
}
