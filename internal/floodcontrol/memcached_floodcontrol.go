package floodcontrol

import (
	"context"
	"errors"
	"github.com/bradfitz/gomemcache/memcache"
	"strconv"
	"task/configs"
)

type MemcachedFloodControl struct {
	client *memcache.Client
	N      int64
	K      int64
}

func NewMemcachedFloodControl(config configs.Config) *MemcachedFloodControl {
	floodControl := &MemcachedFloodControl{
		client: memcache.New(config.Server),
		N:      int64(config.WindowSize.Seconds()),
		K:      int64(config.MaxRequests),
	}

	return floodControl
}

func (m *MemcachedFloodControl) getItem(key string) (*memcache.Item, error) {
	return m.client.Get(key)
}

func (m *MemcachedFloodControl) setItem(ctx context.Context, key string, count int64) (bool, error) {
	ttl, ok := ctx.Value("ttl").(int64)
	if !ok {
		ttl = 1
	}

	item := &memcache.Item{
		Key:        key,
		Value:      []byte(strconv.FormatInt(count, 10)),
		Expiration: int32(ttl),
	}

	err := m.client.Set(item)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (m *MemcachedFloodControl) getCount(key string) (int64, error) {
	item, err := m.client.Get(key)
	if err != nil {
		if errors.Is(err, memcache.ErrCacheMiss) {
			return 0, nil
		}
		return 0, err
	}

	count, err := strconv.ParseInt(string(item.Value), 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (m *MemcachedFloodControl) Check(ctx context.Context, userID int64) (bool, error) {
	key := "floodControl_user:" + strconv.FormatInt(userID, 10)
	count, err := m.getCount(key)
	if err != nil {
		return false, err
	}

	if count < m.K {
		if count == 0 {
			_, err := m.setItem(ctx, key, 1)
			if err != nil {
				return false, err
			}
		} else {
			_, err := m.client.Increment(key, 1)
			if err != nil {
				return false, err
			}
		}
	} else {
		return false, errors.New("flood")
	}

	return true, nil
}
