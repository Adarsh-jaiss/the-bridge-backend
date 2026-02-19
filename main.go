package main

import (
	"fmt"

	db "github.com/adarsh-jaiss/the-bridge/internal/database"
	"github.com/adarsh-jaiss/the-bridge/pkg/config"
	"github.com/adarsh-jaiss/the-bridge/pkg/logger"
	"github.com/adarsh-jaiss/the-bridge/routes"
	_ "github.com/lib/pq"
)

// @title The Bridge API
// @version 1.0
// @description API for The Bridge - Professional network for seafarers
// @host localhost:8080
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg := config.Get()

	log, atomicLevel, err := logger.InitLogger(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer log.Sync()
	fmt.Println("DB USER:", cfg.DBConfig.User)
	fmt.Println("DB PASSWORD:", cfg.DBConfig.Password)
	fmt.Println("DBNAME",cfg.DBConfig.DBName)
	fmt.Println("HOST:",cfg.DBConfig.Host)
	fmt.Println("PORT",cfg.Port)

	dsn, err := db.CreateDSN(db.Config{
		Host:     cfg.DBConfig.Host,
		Port:     cfg.DBConfig.Port,
		DBName:   cfg.DBConfig.DBName,
		User:     cfg.DBConfig.User,
		Password: cfg.DBConfig.Password,
		SSLMode:  cfg.DBConfig.SSLMode,
	})
	if err != nil {
		panic(err)
	}

	conn, err := db.NewConnection("postgres", dsn)
	if err != nil {
		panic(err)
	}

	routes.Routes(conn, log, atomicLevel, cfg)
}
