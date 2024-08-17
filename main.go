package main

import (
	"log"
	"net/http"
	"os"
	"redir/platform"

	"github.com/joho/godotenv"
	"github.com/oschwald/geoip2-golang"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

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

	dsn := "user=postgres password=lab1 dbname=labsales sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn))
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
