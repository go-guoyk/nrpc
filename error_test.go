package nrpc

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUserError(t *testing.T) {
	err := errors.New("test error")
	ue := UserError(err)
	assert.True(t, IsUserError(ue))
	assert.Equal(t, err, errors.Unwrap(ue))
}
