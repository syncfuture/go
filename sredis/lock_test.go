package sredis

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRedisLock_Lock(t *testing.T) {
	client := NewClient(&RedisConfig{Addrs: []string{"127.0.0.1:6379"}, Password: "Famous901"})
	wg := new(sync.WaitGroup)
	wg.Add(2)

	l := NewStrLock(client, "armos", "yingdao", time.Second*600)

	go func() {
		defer wg.Done()
		a, err := l.Lock()
		require.NoError(t, err)
		require.True(t, a)

		c, err := l.Locked()
		require.NoError(t, err)
		require.True(t, c)

		time.Sleep(time.Second * 3)
		err = l.Unlock()
		if err != ErrorUnlockFailed {
			require.NoError(t, err)
		}

	}()

	time.Sleep(time.Second * 1)

	go func() {
		defer func() {
			l.Unlock()
			wg.Done()
		}()
		b, err := l.Lock()
		require.NoError(t, err)
		require.False(t, b)

		time.Sleep(time.Second * 5)

		b, err = l.Lock()
		require.NoError(t, err)
		require.True(t, b)
	}()

	wg.Wait()

	d, err := l.Locked()
	require.NoError(t, err)
	require.False(t, d)
}
