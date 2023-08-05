package sql

// Migrate the schema
func Migrate() error {
	return DB.AutoMigrate(&Server{})
}
