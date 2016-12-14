package main

import(
    "github.com/yungtrizzle/groove/data"
    "github.com/pelletier/go-toml"
    "log"
)



func main(){
    
    config,ok:=toml.LoadFile("/config/groove.toml")
    
    if ok!=nil{
        log.Fatal(ok)
    }
    
    
    rcfg:= data.RedisConfig{
        config.Get("Redis.protocol").(string),
        config.Get("Redis.address").(string),
        config.Get("Redis.port").(int),
        config.Get("Redis.poolsize").(int),
    }
    
    pcfg:=data.PostgresConfig{
     config.Get("Postgres.host").(string),
      config.Get("Postgres.port").(int),
       config.Get("Postgres.user").(string),
       "",
        config.Get("Postgres.database").(string),
    }
    
    data.InitRedis(&rcfg)
    data.InitPostgres(&pcfg)
    
}
