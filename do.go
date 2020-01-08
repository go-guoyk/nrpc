package nrpc

import "context"

func do(ctx context.Context, f func() (err error)) error {
	done := make(chan error, 1)
	go func() {
		done <- f()
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
