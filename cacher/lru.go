package cacher

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

type MemStorage int64

const (
	DefaultLRUSize = 16
	Byte MemStorage = 1
	KiloByte = 1024 * Byte
	MegaByte = 1024 * KiloByte
	DefaultMinMemStorage MemStorage = 32 * MegaByte
	DefaultMaxMemStorage MemStorage = 128 * MegaByte
)

type lruCacher struct {
	cli *lru.Cache
	expiredKeyChan chan IPokeBall
	curMemSize, maxMemSize int64
	lock sync.RWMutex
}

func adjustSize(size int) int {
	if size <= 0 {
		size = DefaultLRUSize
	}
	return size
}

func NewLRUCacher(size int, evictedCallback func(key interface{}, val interface{})) Cacher {
	size = adjustSize(size)
	cli, _ := lru.NewWithEvict(size, evictedCallback)
	res := &lruCacher{cli: cli, expiredKeyChan: make(chan IPokeBall, 5000)}
	go res.asyncWash()
	return res
}

func (l *lruCacher) asyncWash() {
	defer func() {
		if r := recover(); r != nil {
			logger := log.New(os.Stdout, "", log.LstdFlags)
			logger.Print("[lruCacher.asyncWash] panic; r=%v; stack=%v", r, string(debug.Stack()))
		}
	}()
	t := time.NewTicker(time.Second)
	for {
		select {
		case  <- t.C:
			k, v, ok := l.cli.GetOldest()
			if ok {
				ball, ok := v.(interface{
					GetTTL() UnixTime
				})
				if ok && ball.GetTTL() >= UnixTime(time.Now().Unix()) {
					l.cli.Remove(k)
				}
			}
		}
	}
}

func (l *lruCacher) Get(ctx context.Context, key string) (IPokeBall, error) {
	val, ok := l.cli.Get(key)
	if !ok {
		return nil, ErrNotInCache
	}
	res, ok := val.(IPokeBall)
	if !ok {
		return nil, ErrCastError
	}
	return res, nil
}

func (l *lruCacher) Set(ctx context.Context, ball IPokeBall) error {
	key := ball.GetKey()
	l.cli.Add(key, ball)
	return nil
}

func (l *lruCacher) Contains(ctx context.Context, key string) bool {
	return l.cli.Contains(key)
}