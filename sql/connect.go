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

// ConnectDB is a function that initializes a connection to a database
// The function configures a new logger to write on standard output
// and to colorize the console logs.
// The logger level is set to Warn and slow SQL queries are set to take at least 1 second
//
// The function retrieves a data source name (dsn) that is established by the GetDsn() function
// and then initiates a connection using the gorm Open method.
// If an error occurs during the opening of the database connection, the function will log a fatal error.
//
// Next, the function configures the database connection pool to have a maximum lifetime of 1 hour.
// Before the function ends, the function will log a message to terminal that the connection to the database has opened.
//
// GetDsn is a separate function to retrieve data source name.
// The DB connection string is stored in the global DB variable, of type *gorm.DB.
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

	//Config connection pool
	sqlDB, _ := DB.DB()
	sqlDB.SetConnMaxLifetime(time.Hour) // Maximum lifetime of a connection

	log.Println("Connection Opened to database")
	return nil
}

func GetDsn() string {
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_DATABASE")

	//Create dsn
	dsn := username + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbname + "?charset=utf8mb4&parseTime=True&loc=Local"

	return dsn
}
