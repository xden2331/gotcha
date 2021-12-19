package gotcha

import "context"

type GetKeyFn func(ctx context.Context, keyer interface{}) string
type LoadFn func(ctx context.Context, missedItem interface{}, extraArgs interface{}) (interface{}, error)

type Load struct {
	GetKeyFn
	LoadFn
}

func (l *Load) getKey(ctx context.Context, item interface{}) string {
	return l.GetKeyFn(ctx, item)
}

func (l *Load) load(ctx context.Context, missedItem interface{}, extra interface{}) (interface{}, error) {
	return l.LoadFn(ctx, missedItem, extra)
}