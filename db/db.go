package db

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/charmbracelet/log"

	"github.com/alfredosa/go-youtube-reddit-automation/config"
	"github.com/jmoiron/sqlx"
)

func Connect(config config.Config) *sqlx.DB {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGHUP)

	dbConfig := fmt.Sprintf("host=localhost port=%d user=%s password=%s dbname=%s sslmode=disable", config.Database.Port, config.Database.User, config.Database.Password, config.Database.DB_Name)
	db := sqlx.MustConnect("postgres", dbConfig)

	go func() {
		sig := <-sigs
		log.Warn("\n Received an interrupt, stopping services... Signal: %v", sig)
		db.Close()
		log.Warn("Closed DB Connection")

		log.Info("Cleanup completed. Exiting...")

		os.Exit(0)
	}()

	return db
}
