package main

import (
	"cajueiro/code/transactions/models"
	"cajueiro/code/transactions/routers"
	"cajueiro/pkg/app"
	"cajueiro/pkg/exit"
	"cajueiro/pkg/logger"
	"cajueiro/pkg/server"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var api *app.App

func initenv() error {
	// capturando variáveis de ambiente
	viper.SetConfigFile(".env")
	err := viper.ReadInConfig()
	if err != nil {
		logrus.Fatal("Falha ao carregar: ", viper.ConfigFileUsed())
	}
	return err
}

func initapp() error {
	logrus.Info("Arquivo de configuração: ", viper.ConfigFileUsed())
	// armazenando configurações em um struct app
	var err error
	api, err = app.GetApp()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	return err
}

func initdb() error {
	// migrando os schemas do DB
	err := api.DB.Client.AutoMigrate(&models.Account{}, &models.Transaction{})
	if err != nil {
		logrus.Fatal("Erro na migração dos dados")
	}
	return err
}

func init() {
	if initenv() == nil {
		if initapp() == nil {
			if initdb() == nil {
				if api.Cfg.GetDebugMode() == "true" {
					logrus.Warn("Transactions Control rodando em modo Debug")
				} else {
					logrus.Warn("Transactions Control rodando")
				}
			}
		}
	}
}

func main() {

	defer api.DB.CloseDB()

	srv := server.
		GetServer().
		WithAddr(api.Cfg.GetAPIPort()).
		WithRouter(routers.GetRouter(api)).
		WithLogger(logger.Error)

	go func() {
		api.Log.Info("Iniciando servidor na porta ", api.Cfg.GetAPIPort())
		if err := srv.StartServer(); err != nil {
			api.Log.Fatal(err.Error())
		}
	}()

	exit.Init(func() {
		if err := srv.CloseServer(); err != nil {
			api.Log.Error(err.Error())
		}

		if err := api.DB.CloseDB(); err != nil {
			api.Log.Error(err.Error())
		}
	})
}
