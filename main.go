package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"redir/database"
	"redir/datatypes"
	"redir/platform"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
	"github.com/oschwald/geoip2-golang"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Ughmain() {
	dbUser := "shqrladmin"           // Replace with your MySQL username
	dbPass := "TakiArte88!!Reeereee" // Replace with your MySQL password
	dbHost := "68.178.220.114"       // The shared IP address of your server (as shown in the screenshot)
	dbName := "shqrl"                // Replace with your MySQL database name

	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPass, dbHost, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&datatypes.Entry{})

	for i := 1; i <= 5; i++ {
		start := time.Now()

		var entry datatypes.Entry
		result := db.First(&entry, 1)
		if result.Error != nil {
			log.Fatal(result.Error)
		}

		duration := time.Since(start)
		fmt.Printf("Query %d took %v ms\n", i, duration.Milliseconds())

		// Optional: Print the retrieved entry details
		fmt.Printf("ID: %d\n", entry.ID)
		fmt.Printf("User: %s\n", entry.User)
		fmt.Printf("RealURL: %s\n", entry.RealURL)
		fmt.Printf("Archived: %t\n", entry.Archived)
		fmt.Printf("Date: %s\n", entry.Date.Format(time.RFC3339))

		// Wait for 1 second before the next query
		time.Sleep(1 * time.Second)
	}
}

func main() {

	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Fatalf("Failed to load the env vars: %v", err)
		}
	}

	ipDB, err := geoip2.Open("data/GeoLite2-City.mmdb")
	if err != nil {
		log.Fatal(err)
	}
	defer ipDB.Close()

	db, err := database.Connect()
	if err != nil {
		log.Fatal(err)
	}

	rtr := platform.New(db, ipDB)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := http.ListenAndServe(":"+port, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}

}
