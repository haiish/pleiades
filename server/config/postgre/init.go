package db

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"pleiades/server/utils"
	"time"

	"github.com/jackc/pgx/v4/pgxpool" // 🔁 use pgxpool instead of pgx
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB     *pgxpool.Pool // ✅ now uses a pool (auto-reconnects)
	GormDB *gorm.DB      // for ORM
)

func InitDB() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file, nyaaa (´；ω；`)!")
	}

	dbUser := utils.LoadEnvVariable("DB_USER", "postgres")
	dbPassword := utils.LoadEnvVariable("DB_PASSWORD", "password-default")
	dbHost := utils.LoadEnvVariable("DB_HOST", "localhost")
	dbPort := utils.LoadEnvVariable("DB_PORT", "5432")
	dbName := utils.LoadEnvVariable("DB_NAME", "postgres")
	dbSchema := utils.LoadEnvVariable("DB_SCHEMA", "public")
	env := utils.LoadEnvVariable("ENV", "development")
	sslMode := utils.LoadEnvVariable("SSL_STATUS_DB", "disable")

	fmt.Println("DB_USER:", dbUser)
	fmt.Println("DB_PASSWORD:", dbPassword)
	fmt.Println("DB_HOST:", dbHost)
	fmt.Println("DB_PORT:", dbPort)
	fmt.Println("DB_NAME:", dbName)
	fmt.Println("DB_SCHEMA:", dbSchema)
	fmt.Println("ENV:", env)

	// 🌐 Build DSN
	dsn := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   fmt.Sprintf("%s:%s", dbHost, dbPort),
		Path:   dbName,
	}
	q := dsn.Query()
	q.Add("sslmode", sslMode)
	if dbSchema != "" {
		q.Add("search_path", dbSchema)
	}
	dsn.RawQuery = q.Encode()
	databaseURL := dsn.String()

	// 🔁 Connect using pgxpool
	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatal("Failed to parse pgxpool config: ", err)
	}
	config.MaxConns = 20
	config.MaxConnLifetime = 5 * time.Minute

	DB, err = pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Unable to connect to PostgreSQL with pgxpool: ", err)
	}
	fmt.Println("Connected to PostgreSQL (pgxpool)~!! (๑˃̵ᴗ˂̵)و")

	// 🧋 Connect using GORM
	var gormLogger logger.Interface
	if env == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	GormDB, err = gorm.Open(postgres.New(postgres.Config{
		DSN:                  databaseURL,
		PreferSimpleProtocol: true,
	}), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		log.Fatal("Unable to connect to database with GORM: ", err)
	}

	sqlDB, err := GormDB.DB()
	if err != nil {
		log.Fatal("Failed to get generic DB from GORM: ", err)
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	err = sqlDB.Ping()
	if err != nil {
		log.Fatal("Failed to ping database from GORM: ", err)
	}
	fmt.Println("Connected to PostgreSQL (GORM)~!! (＾▽＾)")

	// 🔄 Optional background health check for reconnection logging
	go func() {
		for {
			time.Sleep(30 * time.Second)
			err := DB.Ping(context.Background())
			if err != nil {
				log.Println("⚠️ Lost DB connection, will auto-reconnect...")
			}
		}
	}()
}
