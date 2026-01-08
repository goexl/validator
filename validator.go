package validator

import (
	"github.com/goexl/validator/internal/core"
)

type Validator = core.Validator

func New() *Validator {
	return core.NewValidator()
}
