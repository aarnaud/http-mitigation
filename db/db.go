package db

import (
	"github.com/go-redis/redis"
	"github.com/go-redis/redis_rate"
	"github.com/aarnaud/http-mitigation/config"
	log "github.com/Sirupsen/logrus"
)

var (
	Client *redis.Client
	Limiter *redis_rate.Limiter
)


func Connect(){
	Client = redis.NewClient(&redis.Options{
		Addr: config.Config.RedisAddr,
		Password: config.Config.RedisPassword,
		DB: config.Config.RedisDB,
	})

	pong, err := Client.Ping().Result()
	if err != nil {
		log.Panic(err)
	}

	if pong == "PONG" {
		log.Info("Ping redis success")
	}

	Limiter = redis_rate.NewLimiter(Client)
}