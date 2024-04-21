package sql

// Migrate the schema
func Migrate() error {
	err := DB.AutoMigrate(&Guild{})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&Server{})
}
