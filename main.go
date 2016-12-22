package main

import (
	"github.com/pelletier/go-toml"
	"github.com/yungtrizzle/groove/app"
	"github.com/yungtrizzle/groove/data"
	"github.com/yungtrizzle/groove/web"
	"log"
)

func main() {

	config, ok := toml.LoadFile("config/groove.toml")

	if ok != nil {
		log.Fatal(ok)
	}

	rcfg := data.RedisConfig{
		config.Get("Redis.protocol").(string),
		config.Get("Redis.address").(string),
		int(config.Get("Redis.port").(int64)),
		int(config.Get("Redis.poolsize").(int64)),
	}

	pcfg := data.PostgresConfig{
		config.Get("Postgres.host").(string),
		int(config.Get("Postgres.port").(int64)),
		config.Get("Postgres.user").(string),
		config.Get("Postgres.key").(string),
		config.Get("Postgres.database").(string),
	}

	data.InitRedis(&rcfg)
	data.InitPostgres(&pcfg)
        
        web.Router()
	app.StartPool()
	

}
