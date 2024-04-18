package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"RateLmiter/config"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

// Types for maintaining Rate Limiter data.
type BucketDetails struct {
	LastChecked int64 `json:"last_checked"`
	Token   int64 `json:"token"`
}

var rdb *redis.Client
var configurations config.Configurations

func main() {
	
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetConfigType("yml")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("error reading config", err)
	}
	err := viper.Unmarshal(&configurations)
	if err != nil {
		fmt.Println("unable to get config struct-", err)
	}

	rdb = redis.NewClient(&redis.Options{
    	Addr: configurations.REDIS_HOST,
		Password: configurations.REDIS_PASS,
		DB: configurations.REDIS_DB,
	});
	ping, err := rdb.Ping(rdb.Context()).Result()
	fmt.Println("response on ping-", ping, err)

	router := gin.Default()
	router.GET("/requestLimit", requestLimit)
	router.Run("localhost:8080")
}

var CheckAccess = redis.NewScript(`
local key = KEYS[1]
local curTime = ARGV[1]
local token = tonumber(ARGV[2])
local rate = tonumber(ARGV[3])

local value = cjson.decode(redis.call("GET", key))
if not value then
  return False
end
local timepassed = curTime - value.last_checked
local bucket = value.token + timepassed * token / rate
value.last_checked = curTime
value.token = bucket
if bucket > token then
   value.token = token
end

if bucket < 1 then
  redis.call("SET", key, cjson.encode(value))
  return False
end

value.token = value.token - 1
redis.call("SET", key, cjson.encode(value))

return true
`)

func requestLimit(c *gin.Context) {
	ctx := context.Background()
	res, getErr := rdb.Get(ctx, c.ClientIP()).Result()
	fmt.Println("Response from cache for key-", c.ClientIP(), " is-", res)
	if getErr != nil {
       fmt.Println("Error on getting client from redis-", getErr)
	   value, err := json.Marshal(BucketDetails{
				LastChecked: time.Now().Unix(),
				Token: configurations.BUCKET_SIZE,
			})
		if err != nil{
			fmt.Println("error on intializing cache-", err)
		} else {
			fmt.Println("value to set in cache- ", value)
		}
	    setErr := rdb.Set(ctx, c.ClientIP(), value, 0).Err()
		if setErr != nil {
			fmt.Println("Error on seting value-", setErr)
			panic(setErr)
		}
	} else {
        allow, err := CheckAccess.Run(rdb.Context(), rdb, []string{c.ClientIP()}, time.Now().Unix(), configurations.REFILL_TOKEN, configurations.REFILL_RATE ).Bool()
		fmt.Println("Result-", allow)
		if err != nil {
			fmt.Println("Result Error-", err)
		}
		if (allow) {
			c.IndentedJSON(http.StatusAccepted, nil)
		} else {
			c.IndentedJSON(http.StatusTooManyRequests, nil)
		}
		c.IndentedJSON(http.StatusTooManyRequests, nil)
	}
}