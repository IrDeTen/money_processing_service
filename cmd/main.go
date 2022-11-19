package main

import (
	"log"

	"github.com/IrDeTen/money_processing_service.git/pkg/config"
	"github.com/IrDeTen/money_processing_service.git/pkg/logger"
	"github.com/IrDeTen/money_processing_service.git/server"
	"github.com/spf13/viper"
)

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
