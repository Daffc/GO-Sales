package mariadb

import (
	"fmt"

	"github.com/Daffc/GO-Sales/internal/config"

	"github.com/pressly/goose/v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabaseConnection(dbc *config.Database) (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbc.User, dbc.Password, dbc.Host, dbc.Port, dbc.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return db, err
}

func RunMigrations(db *gorm.DB, migrationsFolderPath string) error {

	// Get the underlying sql.DB from gorm.DB
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// Set the dialect to use for goose
	err = goose.SetDialect("mysql")
	if err != nil {
		return err
	}

	// Run the migrations
	err = goose.Up(sqlDB, migrationsFolderPath)
	if err != nil {
		return err
	}

	return nil
}
