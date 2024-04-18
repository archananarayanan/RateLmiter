package config

type Configurations struct {
	REDIS_HOST string
	REDIS_PASS string
	REDIS_DB int
	REFILL_TOKEN int64
	REFILL_RATE int64
	BUCKET_SIZE int64
}