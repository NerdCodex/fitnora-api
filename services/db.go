package services

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB

func ConnectToDB() {
	dbusername := os.Getenv("DB_USERNAME")
	dbpassword := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DATABASE")
	dbhost := os.Getenv("DB_HOST")
	dbport := os.Getenv("DB_PORT")

	// Validate env vars
	if dbusername == "" || dbpassword == "" || database == "" || dbhost == "" || dbport == "" {
		log.Fatal("DATABASE ENV VARIABLES ARE MISSING")
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbusername, dbpassword, dbhost, dbport, database,
	)

	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "fitnora.",
			SingularTable: true,
		},
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established successfully.")
}
