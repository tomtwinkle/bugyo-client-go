package main

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/tomtwinkle/bugyo-client-go"
	"log"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	var config bugyoclient.BugyoConfig
	if err := envconfig.Process("", &config); err != nil {
		log.Fatal(err)
	}
	client, err := bugyoclient.NewClient(&config, bugyoclient.WithDebug())
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Login(); err != nil {
		log.Fatal(err)
	}
	if err := client.Punchmark(bugyoclient.ClockTypeGoOut); err != nil {
		log.Fatal(err)
	}
}
