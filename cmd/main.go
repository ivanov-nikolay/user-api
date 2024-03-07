package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ivanov-nikolay/user-api/dbinit"
	"github.com/ivanov-nikolay/user-api/internal/delivery"
	"github.com/ivanov-nikolay/user-api/internal/middleware"
	"github.com/ivanov-nikolay/user-api/internal/storage"
	"github.com/ivanov-nikolay/user-api/internal/usecase"
	_ "github.com/jackc/pgx/stdlib"
	"go.uber.org/zap"
)

func main() {
	zapLogger, err := zap.NewProduction()
	if err != nil {
		log.Printf("error in logger start")
		return
	}
	logger := zapLogger.Sugar()
	defer func() {
		err = logger.Sync()
		if err != nil {
			log.Printf("error in logger sync")
		}
	}()
	pgxDB, err := dbinit.GetPostgres()
	fmt.Println("eeee")

	if err != nil {
		logger.Errorf("error in connection to postgres: %s", err)
		return
	}
	logger.Infof("connected to postgres")
	defer func() {
		err = pgxDB.Close()
		if err != nil {
			logger.Errorf("error in close connection to mysql: %s", err)
		}
	}()

	redisConn, err := dbinit.GetRedis()
	if err != nil {
		logger.Infof("error on connection to redis: %s", err.Error())
	}
	defer func() {
		err = redisConn.Close()
		if err != nil {
			logger.Infof("error on redis close: %s", err.Error())
		}
	}()
	logger.Infof("connected to redis")

	s := storage.New(pgxDB, redisConn)
	u := usecase.New(s)
	h := delivery.New(u, logger)

	router := mux.NewRouter()

	router.HandleFunc("/user", h.CreateUserHandler).Methods(http.MethodPost)
	router.HandleFunc("/user/{USER_ID}", h.GetUserByIDHandlerID).Methods(http.MethodGet)
	router.HandleFunc("/user", h.UpdateUserHandler).Methods(http.MethodPut)
	router.HandleFunc("/user/{USER_ID}", h.DeleteUserHandler).Methods(http.MethodDelete)
	router.HandleFunc("/users", h.SearchUsersHandler).Methods(http.MethodGet)

	aclRouter := middleware.AccessLog(router, logger)

	port := os.Getenv("appPort")

	logger.Infow("starting server",
		"type", "START",
		"addr", port,
	)
	err = http.ListenAndServe(port, aclRouter)
	if err != nil {
		logger.Fatalf("errror in server start")
	}
}
