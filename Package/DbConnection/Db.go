package DbConnection

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"
	"time"
)

var DB *pgxpool.Pool

func InitDB() {

	dbName := os.Getenv("DbName")
	dbPort := os.Getenv("DbPort")
	dbHost := os.Getenv("DbHost")
	dbUser := os.Getenv("DbUser")
	dbPassword := os.Getenv("DbPassword")
	maxConnsEnv := os.Getenv("MaxConns")
	minConnsEnv := os.Getenv("MinConns")
	sslMode := os.Getenv("DbSSLMode") // optional (disable, require, verify-full)

	if sslMode == "" {
		sslMode = "disable"
	}

	if dbName == "" || dbHost == "" || dbUser == "" {
		log.Fatal("Database environment variables missing")
	}

	var connStr string
	if dbPort == "" {
		connStr = "postgres://" + dbUser + ":" + dbPassword +
			"@" + dbHost + "/" + dbName + "?sslmode=" + sslMode
	} else {
		connStr = "postgres://" + dbUser + ":" + dbPassword +
			"@" + dbHost + ":" + dbPort + "/" + dbName + "?sslmode=" + sslMode
	}

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Invalid DB config:", err)
	}

	var maxConns int32 = 10
	var minConns int32 = 1

	if maxConnsEnv != "" {
		if v, err := strconv.Atoi(maxConnsEnv); err == nil && v > 0 {
			maxConns = int32(v)
		}
	}

	if minConnsEnv != "" {
		if v, err := strconv.Atoi(minConnsEnv); err == nil && v >= 0 {
			minConns = int32(v)
		}
	}

	cfg.MaxConns = maxConns
	cfg.MinConns = minConns
	cfg.MaxConnLifetime = time.Hour
	cfg.MaxConnIdleTime = 30 * time.Minute
	cfg.HealthCheckPeriod = time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	DB, err = pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatal("Unable to connect to DB:", err)
	}

	if err := DB.Ping(ctx); err != nil {
		log.Fatal("Database ping failed:", err)
	}

	log.Println("PostgreSQL connected")
}

func CloseDB() {
	if DB != nil {
		DB.Close()
		log.Println("PostgreSQL connection closed")
	}
}
