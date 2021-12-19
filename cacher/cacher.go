package cacher

import (
	"context"
	"fmt"
	"time"
)

var ErrNotInCache = fmt.Errorf("not in cache")
var ErrCastError = fmt.Errorf("cast to IPokeBall interface failed")
var ErrExceedStorage = fmt.Errorf("the new item cannot be fit into the member because of the size limitation")

type IPokeBall interface {
	GetKey() string
	GetVal() []byte
	GetTTL() UnixTime

	IsFaded() bool
}

type UnixTime int64

type PokeBall struct {
	key string
	ttl UnixTime
	val []byte
}

func NewPokeBall(key string, ttl UnixTime, val []byte) *PokeBall {
	return &PokeBall{
		key:    key,
		ttl:    ttl,
		val:    val,
	}
}

func (p *PokeBall) GetKey() string {
	return p.key
}

func (p *PokeBall) GetVal() []byte {
	return p.val
}

func (p *PokeBall) GetTTL() UnixTime {
	return p.ttl
}

func (p *PokeBall) IsFaded() bool {
	return time.Now().Unix() > int64(p.GetTTL())
}

type Cacher interface {
	Get(ctx context.Context, key string) (IPokeBall, error)
	Set(ctx context.Context, ball IPokeBall) error
	Contains(ctx context.Context, key string) bool
}
