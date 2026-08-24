package datasource

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/fx"
)

func connString() (string, error) {
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	user := os.Getenv("PG_USER")
	pswd := os.Getenv("PG_PASSWORD")
	db := os.Getenv("PG_DB")

	if host == "" || port == "" || user == "" || pswd == "" || db == "" {
		return "", ErrBuildConnString
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(user, pswd),
		Host:   net.JoinHostPort(host, port),
		Path:   "/" + db,
	}

	return u.String(), nil
}

func InitDBPool(lc fx.Lifecycle) (*pgxpool.Pool, error) {
	connString, err := connString()
	if err != nil {
		return nil, fmt.Errorf("build connection string: %w", err)
	}

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		log.Printf("unable to parse connection string: %v\n", err)
		return nil, err
	}

	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	config.HealthCheckPeriod = time.Minute

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctxTimeout, config)
	if err != nil {
		log.Printf("unable to create connection pool: %v\n", err)
		return nil, err
	}

	if err := pingDB(pool); err != nil {
		log.Println("unable to ping db")
		pool.Close()
		return nil, fmt.Errorf("db ping: %w", err)
	}

	log.Println("db connection pool established")

	lc.Append(fx.Hook{
		OnStop: func(context.Context) error {
			pool.Close()
			return nil
		},
	})

	return pool, nil
}

func pingDB(pool *pgxpool.Pool) error {
	const (
		attempts    = 3
		pingTimeout = 5 * time.Second
		retryDelay  = 500 * time.Millisecond
	)

	var pingErr error

	for i := 1; i <= attempts; i++ {
		ctxTimeout, cancel := context.WithTimeout(context.Background(), pingTimeout)

		pingErr = pool.Ping(ctxTimeout)
		cancel()

		if pingErr == nil {
			return nil
		}

		log.Printf("db ping attempt %d failed: %v", i, pingErr)

		if i < attempts {
			time.Sleep(retryDelay)
		}
	}

	return pingErr
}
