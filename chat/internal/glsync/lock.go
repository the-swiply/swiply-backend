package glsync

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redsync/redsync/v4"
	"time"
)

type ChatLock struct {
	*redsync.Mutex
}

func (c *ChatLock) Lock(ctx context.Context) error {
	return c.LockContext(ctx)
}

func (c *ChatLock) Unlock(ctx context.Context) error {
	ok, err := c.UnlockContext(ctx)
	if err != nil {
		if c.Until().Before(time.Now()) {
			return fmt.Errorf("unlock called after mutex expiration")
		} else {
			return fmt.Errorf("can't release lock: %w", err)
		}
	}
	if !ok {
		return errors.New("can't release lock: status not ok")
	}

	return nil
}
