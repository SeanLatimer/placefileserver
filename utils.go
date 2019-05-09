package main

import (
	"sync/atomic"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func GetIntSlice(key string) []int {
	return cast.ToIntSlice(viper.Get(key))
}

type ReloadCounter struct{ flag uint32 }

func (c *ReloadCounter) Current() uint32 {
	return atomic.LoadUint32(&c.flag)
}

func (c *ReloadCounter) Increment() uint32 {
	return atomic.AddUint32(&c.flag, 1)
}
