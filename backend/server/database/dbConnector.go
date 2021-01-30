package database

import (
	"context"
	"fmt"

	// Database driver for db interactions
	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/jmoiron/sqlx"
)

const (
	sqlDriver = "pgx"
)

// Database abstracts sqlx connection
type Database struct {
	*sqlx.DB
}

// DBConfig abstracts generation of a database configuration string
// Configuration can rely on networks and passing secrets so a method 
// is required instead of string primitive to allow retries in case of failure
type DBConfig interface {
	ConfigString(context.Context) (string, error)
}

// DBConfigFromValues is the default DBConfig type using set values
type DBConfigFromValues struct {
	Database string
	Host     string
	Port     string
	User     string
	Password string
}

// ConfigString returns DBConfigValues formatted into a configuartion string
func (dbConfig DBConfigFromValues) ConfigString(ctx context.Context) (string, error) {
	configString := fmt.Sprintf(
		"database=%s host=%s port=%s user=%s password=%s",
		dbConfig.Database,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.Password,
	)

	return configString, nil
}

// ConnectToDB creates a db connection with any predefined timeout
func ConnectToDB(ctx context.Context, dbConfig DBConfig) (*Database, error) {
	config, err := dbConfig.ConfigString(ctx)
	if err != nil {
		return &Database{}, err
	}

	conn, err := sqlx.ConnectContext(ctx, sqlDriver, config)
	if err != nil {
		return &Database{}, err
	}

	return &Database{conn}, nil
}

// Connected pings server and returns bool response status
func (db *Database) Connected(ctx context.Context) bool {
	if *db == (Database{}) {
		return false
	}

	err := db.PingContext(ctx)

	if err != nil {
		return false
	}

	return true
}

