package database

import (
	"errors"
	"fmt"
	"log"
	"tenant-Dynamin-DB/internals/config"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config, tenant string) (*gorm.DB, error) {
	dbname, driver, err := ResolveTenant(tenant)
	if err != nil {
		return nil, err
	}
	switch driver {
	case "postgres":
		dsn := fmt.Sprintf(
			"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			cfg.PG_HOST,
			cfg.PG_USER,
			cfg.PG_PASS,
			dbname,
			cfg.PG_PORT,
		)
		log.Println("postgress database connected sucessfully DB name :", dbname)

		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("postgres connect failed: %w", err)
		}

		return db, nil

	case "mysql":

		dsn := fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			cfg.MYSQL_USER,
			cfg.MYSQL_PASS,
			cfg.MYSQL_HOST,
			cfg.MYSQL_PORT,
			dbname,
		)
		log.Println("Mysql database connected sucessfully DB name :", dbname)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("mysql connect failed: %w", err)
		}

		return db, nil
	}

	return nil, errors.New("unsupported database driver")
}
