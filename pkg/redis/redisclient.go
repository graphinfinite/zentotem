package redisclient

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
)

type Redis struct {
	redisClient *redis.Client
	log         zerolog.Logger
	TTL         time.Duration
}

// Init connect
func InitDefault(redisurl, redispwd string, redisdb int, logger zerolog.Logger) *Redis {
	red := Redis{
		redisClient: redis.NewClient(&redis.Options{
			Addr:     redisurl,
			Password: redispwd,
			DB:       redisdb,
		}),
		log: logger,
		TTL: 1 * time.Hour, //// some
	}
	return &red
}

func (r *Redis) SetParams(ctx context.Context, key string, value uint64) error {
	err := r.redisClient.Set(ctx, key, value, r.TTL).Err()
	if err != nil {
		r.log.Err(err).Msg("setparams")

	}

	return err
}

func (r *Redis) GetParams(ctx context.Context, key string) (uint64, error) {
	val, err := r.redisClient.Get(ctx, key).Uint64()
	if err != nil {
		r.log.Err(err).Msg("getparams")
	}

	return val, err
}

func (r *Redis) DeleteParams(ctx context.Context, key string) {
	res := r.redisClient.Del(ctx, key)
	r.log.Debug().Msg(res.String())
}
