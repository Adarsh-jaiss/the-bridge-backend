package main

import (
	"github.com/adarsh-jaiss/the-bridge/db"
	"github.com/adarsh-jaiss/the-bridge/pkg/config"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/adarsh-jaiss/the-bridge/routes"
)

func main() {
	cfg := config.Load()

	log, atomicLevel, err := logger.InitLogger(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	dsn, err := db.CreateDSN(db.Config{
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		DBName:   cfg.DBConfig.DBName,
		User:     cfg.DBConfig.User,
		Password: cfg.DBConfig.Password,
	})
	if err != nil {
		panic(err)
	}

	conn, err := db.NewConnection("mysql", dsn)
	if err != nil {
		panic(err)
	}

	routes.Routes(conn, log, atomicLevel)
}
