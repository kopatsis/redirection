package main

import (
	"log"
	"net/http"
	"os"
	"redir/database"
	"redir/platform"

	_ "github.com/go-sql-driver/mysql"

	"github.com/joho/godotenv"
	"github.com/oschwald/geoip2-golang"
)

func main() {

	if err := godotenv.Load(); err != nil {
		if os.Getenv("APP_ENV") != "production" {
			log.Fatalf("Failed to load the env vars: %v", err)
		}
	}

	ipDB, err := geoip2.Open("geolite/GeoLite2-City.mmdb")
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
