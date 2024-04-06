package main

import (
	"log"
	"urlshortner/database"
	"urlshortner/router"

	"github.com/go-redis/redis"
	"github.com/joho/godotenv"
)

func init() {
	var RedisClient *redis.Client
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	log.Println("//")
	database.ConnectDb()
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	log.Println("///")
	pingres, err := RedisClient.Ping().Result()
	log.Printf("The err is %v", err)
	if err != nil {
		log.Println("////")
		panic(err)
	}
	log.Printf("The ping is %v", pingres)
}

func main() {
	router.ClientRoutes()
}
