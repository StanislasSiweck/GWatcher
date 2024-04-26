package sql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func ConnectDB() error {

	// Initialize new logger
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			Colorful:      true,        // Enable/Disable colors in consol
			LogLevel:      logger.Warn, // Log Level
		},
	)
	var err error
	dsn := GetDsn()

	//Setup db connection form env data
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		log.Println("Error opening database connection :", err)
		return err
	}

	sqlDB, _ := DB.DB()
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	log.Println("Connection Opened to database")
	return nil
}

func GetDsn() string {
	username := getEnv("DB_USERNAME", "root")
	password := getEnv("DB_PASSWORD", "")
	host := getEnv("DB_HOST", "127.0.0.1")
	port := getEnv("DB_PORT", "3306")
	dbname := getEnv("DB_DATABASE", "test")

	//Create dsn
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	return dsn
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
