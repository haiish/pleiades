package redis

import (
	"context"
	"fmt"
	"log"
	"pleiades/server/utils"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
)

var RedisClient *redis.Client

func InitRedis() {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file~ (T_T)")
	}

	// Read environment variables
	redisHost := utils.LoadEnvVariable("REDIS_HOST", "localhost")
	redisPort := utils.LoadEnvVariable("REDIS_PORT", "6379")
	redisPassword := utils.LoadEnvVariable("REDIS_PASSWORD", "")
	redisDB := utils.LoadEnvVariable("REDIS_DB", "0")

	fmt.Println("REDIS_HOST:", redisHost)
	fmt.Println("REDIS_PORT:", redisPort)
	fmt.Println("REDIS_DB:", redisDB)

	// Connect to Redis
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: redisPassword,
		DB:       stringToInt(redisDB, 0),
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Test connection
	_, err = RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Unable to connect to Redis~ (つ﹏<)･ﾟ｡ :", err)
	}

	fmt.Println("Connected to Redis~ ₍˶ˆ꒳ˆ˶₎✨")
}

func stringToInt(s string, defaultVal int) int {
	val, err := strconv.Atoi(s)
	if err != nil {
		return defaultVal
	}
	return val
}
