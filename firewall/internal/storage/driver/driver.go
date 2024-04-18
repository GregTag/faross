package driver

import (
	"firewall/internal/entity"
	"firewall/pkg/config"
	"log"
	"strings"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func migrateDB(connection gorm.Dialector) (*gorm.DB, error) {
	db, err := gorm.Open(connection, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&entity.Repository{},
		&entity.Package{},
	)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func NewSQLiteDB(path string) (*gorm.DB, error) {
	return migrateDB(sqlite.Open(path))
}

func NewPostgreSQLDB(dsn string) (*gorm.DB, error) {
	return migrateDB(postgres.Open(dsn))
}

func NewDB() (*gorm.DB, error) {
	driver := config.Koanf.MustString("dbDriver")
	creds := config.Koanf.MustString("dbCreds")
	driver = strings.ToLower(driver)
	var db *gorm.DB
	var err error

	switch driver {
	case "sqlite":
		db, err = NewSQLiteDB(creds)
	case "postgresql":
		db, err = NewPostgreSQLDB(creds)
	default:
		log.Panicln("No driver ", driver)
	}

	return db, err
}
