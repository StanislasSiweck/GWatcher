package sql

import (
	"fmt"
	"github.com/lmittmann/tint"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"log/slog"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectDB() error {

	level := getLogLevel(logger.Silent)

	// Initialize new logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			Colorful:      true,        // Enable/Disable colors in consol
			LogLevel:      level,       // Log Level
		},
	)
	var err error
	dsn := getDsn()

	//Setup db connection form env data
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		slog.Error("Can't opening database connection", tint.Err(err))
		return err
	}

	slog.Info("Connection Opened to database")
	return nil
}

func getLogLevel(defaultLevel logger.LogLevel) (level logger.LogLevel) {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG":
		level = logger.Error
	case "INFO":
		level = logger.Info
	case "WARN":
		level = logger.Warn
	case "ERROR":
		level = logger.Error
	default:
		level = defaultLevel
	}
	return
}

func getDsn() string {
	username := getEnv("DB_USERNAME", "root")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_DATABASE", "test")

	//Create dsn
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)

	return dsn
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
