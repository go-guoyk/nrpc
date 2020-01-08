package nrpc

import (
	"context"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	ctx, ctxCancel := context.WithCancel(context.Background())
	go ctxCancel()
	err := do(ctx, func() (err error) {
		time.Sleep(time.Millisecond * 100)
		return
	})
	require.Equal(t, context.Canceled, err)

	ctx, ctxCancel = context.WithCancel(context.Background())
	err = do(ctx, func() (err error) {
		time.Sleep(time.Millisecond * 100)
		return
	})
	require.NoError(t, err)
}
