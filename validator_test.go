package validator_test

import (
	"testing"

	"github.com/goexl/validator"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	assert.NotNil(t, validator.New())
}
