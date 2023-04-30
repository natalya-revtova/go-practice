package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("replacement", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Set("bbb", 200)
		c.Set("ccc", 300)
		c.Set("ddd", 400)

		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("replacement unused items", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		c.Set("bbb", 200)
		c.Set("ccc", 300)

		c.Set("bbb", 400)
		c.Get("aaa")
		c.Set("ccc", 500)
		c.Get("aaa")

		c.Set("ddd", 600)

		val, ok := c.Get("bbb")
		require.False(t, ok)
		require.Nil(t, val)

		c.Set("eee", 700)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("clear", func(t *testing.T) {
		c := NewCache(3)

		c.Set("aaa", 100)
		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		c.Clear()
		val, ok = c.Get("aaa")
		require.False(t, ok)
		require.Nil(t, val)
	})
}
