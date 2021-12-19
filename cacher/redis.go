package cacher

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

type redisCacher struct {
	cli *redis.Client
}

func NewRedisCacher(cli *redis.Client) Cacher {
	return &redisCacher{cli: cli}
}

func (r *redisCacher) Get(ctx context.Context, key string) (IPokeBall, error) {
	val, err := r.cli.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	ttl, err := r.cli.TTL(ctx, key).Result()
	if err != nil {
		log.Fatalf("query %v ttl failed; err=%v", key, err.Error())
	}
	return  &PokeBall{
		key:    key,
		ttl:    UnixTime(ttl),
		val:    []byte(val),
	}, nil
}

func (r *redisCacher) Set(ctx context.Context, ball IPokeBall) error {
	key, val, ttl := ball.GetKey(), ball.GetVal(), ball.GetTTL()
	_, err := r.cli.Set(ctx, key, string(val), time.Duration(ttl)).Result()
	return err
}

func (r *redisCacher) Contains(ctx context.Context, key string) bool {
	return false
}