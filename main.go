package main

import (
	"log"
	"net/http"
	"os"
	"redir/database"
	"redir/platform"

	"github.com/go-redis/redis/v8"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/account"

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

	stripe.Key = os.Getenv("STRIPE_SECRET")

	acct, err := account.Get()
	if err != nil {
		log.Fatalf("Stripe API key test failed: %v", err)
	}
	log.Printf("Stripe API key test succeeded: Account ID = %s, Email = %s", acct.ID, acct.Email)

	rtr := platform.New(db, ipDB, rdb)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3500"
	}

	if err := http.ListenAndServe(":"+port, rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}

}
