package dbinit

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/joho/godotenv"
)

const (
	maxDBConnections  = 10
	maxPingDBAttempts = 20
)

func GetRedis() (redis.Conn, error) {
	host := os.Getenv("hostRD")
	port := os.Getenv("portRD")
	c, err := redis.DialURL(fmt.Sprintf("redis://user:@%s:%s/0", host, port))
	if err != nil {
		return nil, err
	}
	return c, nil

}

func GetPostgres() (*sql.DB, error) {
	envFilePath := ".env"
	err := godotenv.Load(envFilePath)
	if err != nil {
		fmt.Println("err")
	}
	pass := os.Getenv("pass")
	user := os.Getenv("user")
	dbName := os.Getenv("dbName")
	host := os.Getenv("hostPG")
	port := os.Getenv("portPG")
	fmt.Println(port)
	sslMode := os.Getenv("sslMode")
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=%s",
		user, dbName, pass, host, port, sslMode)
	db, err := sql.Open("pgx", dsn)

	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(maxDBConnections)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	attemptsNumber := 0
	for range ticker.C {
		err = db.Ping()
		attemptsNumber++
		if err == nil {
			break
		}
		if attemptsNumber == maxPingDBAttempts {
			return nil, err
		}
	}
	return db, nil
}
