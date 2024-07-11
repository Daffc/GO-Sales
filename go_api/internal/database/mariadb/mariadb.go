package mariadb

import (
	"fmt"

	"github.com/Daffc/GO-Sales/internal/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewDatabaseConnection(c *config.Config) (*gorm.DB, error) {

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	return db, err
}
