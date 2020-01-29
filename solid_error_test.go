package nrpc

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSolidError(t *testing.T) {
	err := errors.New("test error")
	ue := Solid(err)
	assert.True(t, IsSolid(ue))
	assert.Equal(t, err, errors.Unwrap(ue))
}
