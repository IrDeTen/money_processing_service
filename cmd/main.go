package main

import (
	"log"

	_ "github.com/IrDeTen/money_processing_service.git/docs"
	"github.com/IrDeTen/money_processing_service.git/pkg/config"
	"github.com/IrDeTen/money_processing_service.git/pkg/logger"
	"github.com/IrDeTen/money_processing_service.git/server"
	"github.com/spf13/viper"
)

// @title Money Processing Service API
// @version 1.0
// @description Test task
// @BasePath /processing
func main() {
	if err := config.InitConfig(); err != nil {
		log.Fatal(err)
	}
	logger.InitLogger(
		viper.GetString("app.log.dir"),
		viper.GetString("app.log.file"),
	)

	app := server.NewApp()
	app.Run(viper.GetString("app.http_port"))
}
