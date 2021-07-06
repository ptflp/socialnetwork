package cache

import (
	"errors"
	"time"
)

var (
	ErrCacheMiss = errors.New("key not found")
)

type Cache interface {
	Get(key string, ptrValue interface{}) error
	Set(key string, ptrValue interface{}, expires time.Duration)
	Del(key string) error
}
