package cacher

import (
	"context"
	lru "github.com/hashicorp/golang-lru"
	"reflect"
	"strings"
	"testing"
)

func Test_lruCacher_Get(t *testing.T) {
	cli, _ := lru.New(DefaultLRUSize)
	type fields struct {
		cli *lru.Cache
	}
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    IPokeBall
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "get nil",
			fields:  fields{cli: cli},
			args:    args{
				ctx: context.Background(),
				key: "emptyKey",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "set-get-1",
			fields:  fields{cli: cli},
			args:    args{
				ctx: context.Background(),
				key: "set_key_1",
			},
			want:    NewPokeBall("set_key_1",  UnixTime(100000), UnixTime(100000), []byte("set_key_1")),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lruCacher{
				cli: tt.fields.cli,
			}
			if strings.Contains(tt.name, "set") {
				l.Set(context.Background(), NewPokeBall(tt.args.key, UnixTime(100000), UnixTime(100000), []byte(tt.args.key)))
			}
			got, err := l.Get(tt.args.ctx, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_lruCacher_Set(t *testing.T) {
	cli, _ := lru.NewWithEvict(1, func(k, v interface{}) {
		t.Logf("evicted: key=%v, val=%v", k, v)
	})
	type fields struct {
		cli *lru.Cache
	}
	type args struct {
		ctx  context.Context
		ball IPokeBall
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "set_one",
			fields:  fields{cli: cli},
			args:    args{
				ctx:  context.Background(),
				ball: NewPokeBall("key_1", UnixTime(0), UnixTime(0), []byte("key_1")),
			},
			wantErr: false,
		},
		{
			name:    "set_two",
			fields:  fields{cli: cli},
			args:    args{
				ctx:  context.Background(),
				ball: NewPokeBall("key_2", UnixTime(0), UnixTime(0), []byte("key_2")),
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := &lruCacher{
				cli: tt.fields.cli,
			}
			if err := l.Set(tt.args.ctx, tt.args.ball); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}