package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Netflix/go-env"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

/*
CREATE INDEX idx_logs_created_at ON logs (created_at);

export DB_HOST=postgres
export DB_USER=postgres
export DB_PASS=p0stgr3s
export DB_NAME=logto
export DB_PORT=5432
*/
const (
	LOG_TABLE = "logs"
)

type Log struct {
	CreatedAt time.Time `gorm:"column:created_at"`
}

type Environment struct {
	DB_HOST string `env:"DB_HOST"`
	DB_USER string `env:"DB_USER"`
	DB_PASS string `env:"DB_PASS"`
	DB_NAME string `env:"DB_NAME"`
	DB_PORT int    `env:"DB_PORT"`
}

var cleanning bool = false

func clean(environment Environment) {
	if cleanning {
		return
	}

	cleanning = true
	defer func() {
		cleanning = false
	}()

	for {
		db, err := gorm.Open(postgres.New(postgres.Config{
			DSN:                  fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Shanghai", environment.DB_HOST, environment.DB_USER, environment.DB_PASS, environment.DB_NAME, environment.DB_PORT),
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		}), &gorm.Config{})

		if err != nil {
			log.Printf("failed to connect database: %s", err)
			time.Sleep(5 * time.Second)
			continue
		}

		log.Printf("connected to db")

		db.Table(LOG_TABLE).Exec("DELETE FROM logs WHERE created_at < now() - interval '1 day'")
		log.Printf("cleaned up logs")

		sqlDB, _ := db.DB()
		if sqlDB != nil {
			log.Printf("closing db connection")
			sqlDB.Close()
		}
		return
	}
}

func main() {

	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("DB ENV: %+v", environment)

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	ticker := time.NewTicker(24 * time.Hour)

	for {
		select {
		case <-ticker.C:
			go clean(environment)
		case <-ch:
			ticker.Stop()
			log.Printf("received signal, stopping")
			return
		}
	}
}
