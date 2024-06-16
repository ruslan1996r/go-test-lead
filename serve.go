package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "leads/docs"
	"leads/handlers"
	"leads/storage"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

const DBPATH = "DB_PATH"

func health(ctx *gin.Context) {
	ctx.String(http.StatusOK, "ok")
}

func CreateApp() (*gin.Engine, error) {
	dbPath := os.Getenv(DBPATH)
	ctx := context.Background()

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("can't open database: %w", err)
	}

	sqlHelpers := storage.NewSQLHelper()
	sqlStorage, err := storage.New(db, sqlHelpers)
	if err != nil {
		log.Fatal("can't connect to storage: ", err)
	}

	if err := sqlStorage.Init(ctx); err != nil {
		log.Fatal("can't init storage", err)
	}

	if err := sqlStorage.Migrations(ctx); err != nil {
		log.Fatal("migrations failed", err)
	}

	clientsHandler := handlers.NewClientsHandlers(sqlStorage)

	r := gin.New()

	r.GET("/health", health)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	clientsHandler.InstallRoutes(r)

	return r, nil
}
