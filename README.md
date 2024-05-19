# RateLmiter

RateLimter service implementing Token Bucket Algorithm using Redis Cluster for storing the usage data. 


Packages used - 
swaggo - for API documentation (https://github.com/swaggo/swag)
viper - for reading bucket configuration (https://github.com/spf13/viper)

Horizontal Scaling of Redis - https://redis.io/docs/latest/operate/oss_and_stack/management/scaling/#interact-with-the-cluster 


Setup redis cluster configuration - 
>> docker-compose up --build 

Initialize Swagger to generate the API's documentation - 
>> swag init - to redefine swagger 

Run the application using - 
>> go build 
>> go run .