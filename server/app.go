package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/IrDeTen/money_processing_service.git/app"
	apphttp "github.com/IrDeTen/money_processing_service.git/app/delivery/http"
	"github.com/IrDeTen/money_processing_service.git/app/repo/postgres"
	"github.com/IrDeTen/money_processing_service.git/app/usecase"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type App struct {
	uc         app.IUsecase
	repo       app.IRepository
	httpServer *http.Server
}

func NewApp() *App {
	repo := postgres.NewRepo()
	uc := usecase.NewUsecase(repo)
	return &App{
		uc:   uc,
		repo: repo,
	}
}

func (a *App) Run(port string) error {
	defer a.repo.Close()
	//TODO change gin mod
	gin.SetMode(gin.ReleaseMode)
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

//TODO add init db

//TODO add run migrations
