package sql

import "bot-serveur-info/internal/pkg/sql/model"

// Migrate the schema
func Migrate() error {
	err := DB.AutoMigrate(&model.Guild{})
	if err != nil {
		return err
	}
	return DB.AutoMigrate(&model.Server{})
}
