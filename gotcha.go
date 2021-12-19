package gotcha

import (
	"context"
	"errors"
	"github.com/xden2331/gotcha/cacher"
	"github.com/xden2331/gotcha/encoderDecoder"
	"time"
)

type IGotcha interface{
	Get(ctx context.Context, toFind, placeholder, extra interface{}) ProcessInfoer
}

type FetchFrom string

const (
	Cache FetchFrom = "cache"
	Source FetchFrom = "Source"
)

type Strategy int

const (
	SourceOnly Strategy = 0
	CacheThenSource Strategy = 1
)

type ProcessInfoer struct {
	Err    error
	Source FetchFrom
}

type gotcha struct {
	Load
	GotchaName string
	strategy Strategy
	//logger *log.Logger
	c cacher.Cacher
	ed encoderDecoder.EncoderDecoder
	ttl time.Duration

	afterGetFromCache func(ctx context.Context, k, v interface{})
	afterGetFromSource func(ctx context.Context, k, v interface{})
}

func (g *gotcha) Get(ctx context.Context, toFind, placeholder, extra interface{}) ProcessInfoer {
	info := ProcessInfoer{}
	key := g.getKey(ctx, toFind)
	if g.strategy != SourceOnly {
		ball, err := g.c.Get(ctx, key)
		if err != nil && !errors.Is(err, cacher.ErrNotInCache) {
			info.Err = err
			return info
		}
		if ball != nil && !ball.IsFaded() {
			info.Err = g.ed.Decode(ball.GetVal(), placeholder)
			if info.Err == nil {
				info.Source = Cache
			}
			if g.afterGetFromCache != nil {
				g.afterGetFromSource(ctx, key, placeholder)
			}
			return info
		}
	}

	srcObj, err := g.load(ctx, toFind, extra)
	if err != nil {
		info.Err = err
		return info
	}
	if g.afterGetFromSource != nil {
		g.afterGetFromSource(ctx, key, srcObj)
	}
	bs, err := g.ed.Encode(srcObj)
	if err != nil {
		info.Err = err
		return info
	}
	err = g.ed.Decode(bs, placeholder)
	if err != nil {
		info.Err = err
		return info
	}
	info.Source = Source
	err = g.c.Set(ctx, cacher.NewPokeBall(key, cacher.UnixTime(time.Now().Add(g.ttl).Unix()), bs))
	if err != nil {
		info.Err = err
		return info
	}
	return info
}

type DefaultBuilder gotcha

func NewDefaultGotcha(keyFunc getKeyFn, loadFunc loadFn) IGotcha {
	return &gotcha{
		Load:       Load{keyFunc, loadFunc},
		GotchaName: "",
		c:          cacher.NewLRUCacher(32, nil),
		ed:         encoderDecoder.NewEncoderDecoder(encoderDecoder.EncoderDecoderType_OfficialJSON),
		strategy: CacheThenSource,
	}
}

/**
NewGotchaBuilder
*/
func NewGotchaBuilder(keyFunc getKeyFn, loadFunc loadFn) DefaultBuilder {
	return DefaultBuilder(
		gotcha{
			Load:       Load{keyFunc, loadFunc},
			GotchaName: "",
			c:          nil,
			ed:         nil,
		},
	)
}

func (gb DefaultBuilder) WithCache(c cacher.Cacher) DefaultBuilder {
	gc := gotcha(gb)
	gc.c = c
	return DefaultBuilder(gc)
}

func (gb DefaultBuilder) WithEncoderDecoder(ed encoderDecoder.EncoderDecoder) DefaultBuilder {
	gc := gotcha(gb)
	gc.ed = ed
	return DefaultBuilder(gc)
}

func (gb DefaultBuilder) WithStrategy(s Strategy) DefaultBuilder {
	gc := gotcha(gb)
	gc.strategy = s
	return DefaultBuilder(gc)
}

func (gb DefaultBuilder) WithTTL(ttl time.Duration) DefaultBuilder {
	gc := gotcha(gb)
	gc.ttl = ttl
	return DefaultBuilder(gc)
}

func (gb DefaultBuilder) WithCacheCallback(afterGetFromCache func(context.Context, interface{}, interface{})) DefaultBuilder {
	gc := gotcha(gb)
	gc.afterGetFromCache = afterGetFromCache
	return DefaultBuilder(gc)
}

func (gb DefaultBuilder) WithSourceCallback(afterGetFromSrc func(context.Context, interface{}, interface{})) DefaultBuilder {
	gc := gotcha(gb)
	gc.afterGetFromSource = afterGetFromSrc
	return DefaultBuilder(gc)
}

func (gb DefaultBuilder) Build() IGotcha {
	gc := gotcha(gb)
	return &gc
}
