package main

import (
	"log"
	"net/http"
	"os"
	"redir/database"
	"redir/platform"
	"time"

	"github.com/go-redis/redis/v8"
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

	redisAddr, redisPass := os.Getenv("REDIS_ADDR"), os.Getenv("REDIS_PASSWORD")
	if redisAddr == "" || redisPass == "" {
		log.Fatalf("cannot connect to redis as no addr and/or pass present in env")
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Username: "default",
		Password: redisPass,
		DB:       0,
	})

	var httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}

	rtr := platform.New(db, ipDB, rdb, httpClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3500"
	}

	if err := http.ListenAndServe(":"+port, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}

}
