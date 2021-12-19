package gotcha

import (
	"context"
	"fmt"
	"github.com/xden2331/gotcha/cacher"
	"github.com/xden2331/gotcha/encoderDecoder"
	"os"
	"reflect"
	"testing"
	"time"
)

func toIntPtr(val int) *int {
	return &val
}

func toStrPtr(str string) *string {
	return &str
}

var expectedFirst = &dao{
	ID:      1,
	Email:   "first@mail.com",
	Address: "Rm first, second St., third State, fourth Country",
}

var expectedSecond =  &dao{
	ID:      2,
	Email:   "2_first@mail.com",
	Address: "2_Rm first, second St., third State, fourth Country",
}

var db = map[int64]*dao{
	1: expectedFirst,
	2: expectedSecond,
}

type dao struct {
	ID int64 `json:"id"`
	Email string `json:"email"`
	Address string `json:"address"`
}

var cacheRepo IGotcha
var ctx = context.Background()

func TestMain(m *testing.M) {
	resetRepo()
	os.Exit(m.Run())
}

func resetRepo() {
	cacheRepo = NewGotchaBuilder(func(ctx context.Context, i interface{}) string { // extract key func
		obj := i.(*dao)
		return fmt.Sprintf("%v", obj.ID)
	}, func(ctx context.Context, missedItem interface{}, extraArgs interface{}) (interface{}, error) { // Source db func
		obj := missedItem.(*dao)
		return db[obj.ID], nil
	}).WithTTL(time.Second * 3). // ttl
		WithCache(cacher.NewLRUCacher(1, func(k, v interface{}) {
			fmt.Printf("evicted: [key:%s, val:%s]\n", k, v)
	})). // cache
		WithEncoderDecoder(encoderDecoder.NewEncoderDecoder(encoderDecoder.EncoderDecoderType_OfficialJSON)). // encoder decoder
		WithStrategy(CacheThenSource).
		Build()
}

func TestGotcha_Get(t *testing.T) {
	defer resetRepo()
	firstItem := &dao{}
	info := cacheRepo.Get(ctx, &dao{ID: 1}, firstItem, nil)
	if info.Err != nil {
		t.Errorf("err: %s", info.Err)
		return
	}
	if info.Source != Source {
		t.Errorf("expected fetchForm:%s, got:%s", info.Source, Source)
		return
	}
	if !reflect.DeepEqual(firstItem, expectedFirst) {
		t.Errorf("expected %+v; got %+v", expectedFirst, firstItem)
		return
	}
}

func TestGotcha_GetFromCache(t *testing.T) {
	defer resetRepo()
	firstItem := &dao{}
	info := cacheRepo.Get(ctx, &dao{ID: 1}, firstItem, nil)
	if info.Err != nil {
		t.Errorf("err: %s", info.Err)
		return
	}
	if info.Source != Source {
		t.Errorf("expected fetchForm:%s, got:%s", info.Source, Source)
		return
	}
	if !reflect.DeepEqual(firstItem, expectedFirst) {
		t.Errorf("expected %+v; got %+v", expectedFirst, firstItem)
		return
	}

	// then should fetch from cache
	firstItem = &dao{}
	info = cacheRepo.Get(ctx, &dao{ID: 1}, firstItem, nil)
	if info.Err != nil {
		t.Errorf("err: %s", info.Err)
		return
	}
	if info.Source != Cache {
		t.Errorf("expected fetchForm:%s, got:%s", Cache, info.Source)
		return
	}
	if !reflect.DeepEqual(firstItem, expectedFirst) {
		t.Errorf("expected %+v; got %+v", expectedFirst, firstItem)
		return
	}
}

func TestGotcha_Expire(t *testing.T) {
	defer resetRepo()

	firstItem := &dao{}
	info := cacheRepo.Get(ctx, &dao{ID: 1}, firstItem, nil)
	if info.Err != nil {
		t.Errorf("err: %s", info.Err)
		return
	}
	if info.Source != Source {
		t.Errorf("expected fetchForm:%s, got:%s", info.Source, Source)
		return
	}
	if !reflect.DeepEqual(firstItem, expectedFirst) {
		t.Errorf("expected %+v; got %+v", expectedFirst, firstItem)
		return
	}

	time.Sleep(time.Second * 4) // sleep for 4 sec

	firstItem = &dao{}
	info = cacheRepo.Get(ctx, &dao{ID: 1}, firstItem, nil)
	if info.Err != nil {
		t.Errorf("err: %s", info.Err)
		return
	}
	if info.Source != Source {
		t.Errorf("expected fetchForm:%s, got:%s", info.Source, Source)
		return
	}
	if !reflect.DeepEqual(firstItem, expectedFirst) {
		t.Errorf("expected %+v; got %+v", expectedFirst, firstItem)
		return
	}
}

func TestGotcha_EvictedByTheSecondItem(t *testing.T) {
	defer resetRepo()

	// get the first
	firstItem := &dao{}
	info := cacheRepo.Get(ctx, &dao{ID: 1}, firstItem, nil)
	if info.Err != nil {
		t.Errorf("err: %s", info.Err)
		return
	}
	if info.Source != Source { // ensure from source
		t.Errorf("expected fetchForm:%s, got:%s", info.Source, Source)
		return
	}
	if !reflect.DeepEqual(firstItem, expectedFirst) {
		t.Errorf("expected %+v; got %+v", expectedFirst, firstItem)
		return
	}

	// get the second
	secondItem := &dao{ID: 2}
	info = cacheRepo.Get(ctx, secondItem, secondItem, nil)
	if info.Err != nil {
		t.Errorf("err: %s", info.Err)
		return
	}
	if info.Source != Source { // enusre from source
		t.Errorf("expected fetchForm:%s, got:%s", info.Source, Source)
		return
	}
	if !reflect.DeepEqual(secondItem, expectedSecond) {
		t.Errorf("expected %+v; got %+v", expectedSecond, secondItem)
		return
	}
	t.Logf("logs shows that the first item is evicted.\n")
}