package gotcha

import "context"

type getKeyFn func(ctx context.Context, keyer interface{}) string
type loadFn func(ctx context.Context, missedItem interface{}, extraArgs interface{}) (interface{}, error)

type Load struct {
	getKeyFn
	loadFn
}

func (l *Load) getKey(ctx context.Context, item interface{}) string {
	return l.getKeyFn(ctx, item)
}

func (l *Load) load(ctx context.Context, missedItem interface{}, extra interface{}) (interface{}, error) {
	return l.loadFn(ctx, missedItem, extra)
}