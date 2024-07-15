package config

import (
	"cassette/pkg/storage"
	"cassette/pkg/storage/filesystem"
	"os"

	"cassette/pkg/repository"
	"cassette/pkg/repository/gorm"
)

type Config struct {
	Storage    storage.Storage
	Repository repository.Repository

	Username string
	Password string
}

func FromEnvironment() (*Config, error) {
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")

	if username == "" {
		username = "admin"
	}

	if password == "" {
		password = "admin"
	}

	path := os.Getenv("DATA_PATH")

	if path == "" {
		path = "sessions"
	}

	s, err := filesystem.New(path)

	if err != nil {
		return nil, err
	}

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=" + os.Getenv("DB_HOST") + " user=" + os.Getenv("POSTGRES_USER") + " password=" + os.Getenv("POSTGRES_PASSWORD") + " dbname=" + os.Getenv("POSTGRES_DB") + " port=" + os.Getenv("DB_PORT") + " sslmode=disable"
	}

	r, err := gorm.NewPostgres(dsn)
	if err != nil {
		return nil, err
	}

	return &Config{
		Storage:    s,
		Repository: r,

		Username: username,
		Password: password,
	}, nil
}
