package core

import (
	"context"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	"github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/goexl/exception"
	"github.com/goexl/gox"
	"github.com/goexl/gox/field"
	"github.com/goexl/validate"
	"github.com/goexl/validator/internal/internal/constant"
	"github.com/goexl/validator/internal/internal/converter"
	"github.com/goexl/validator/internal/internal/core"
)

var _ validate.Validator = (*Validator)(nil)

type Validator struct {
	validator  *validator.Validate
	translator *ut.UniversalTranslator
}

func NewValidator() *Validator {
	return &Validator{
		validator:  validator.New(),
		translator: ut.New(en.New(), zh.New()),
	}
}

func (v *Validator) Validate(ctx context.Context, target any) (err error) {
	if value := ctx.Value(constant.ContextTag); value != nil {
		err = v.validator.Var(target, value.(string))
	} else {
		err = v.validator.Struct(target)
	}

	if err != nil {
		err = v.localization(ctx, err.(validator.ValidationErrors))
	}

	return
}

func (v *Validator) localization(ctx context.Context, errors validator.ValidationErrors) (err error) {
	translations := v.translations(ctx, errors)
	// 国际化后，需要对每个字段名按规则进行转换
	fields := make([]gox.Field[any], 0, len(translations))
	for translation, message := range translations {
		realField := translation[strings.IndexRune(translation, constant.RuneDot)+1:]
		finalField := v.detectConverter(ctx).Convert(realField)
		fields = append(fields, field.New(finalField, message))
	}
	err = exception.New().Message(constant.MessageValidateError).Field(fields[0], fields[1:]...).Build()

	return
}

func (v *Validator) translations(
	ctx context.Context,
	errors validator.ValidationErrors,
) (translations validator.ValidationErrorsTranslations) {
	language := v.detectLanguage(ctx)
	if translate, found := v.translator.FindTranslator(language); found {
		translations = errors.Translate(translate)
	} else if zh, fz := v.translator.GetTranslator(constant.LanguageZh); fz {
		translations = errors.Translate(zh)
	}

	return
}

func (v *Validator) detectConverter(ctx context.Context) (language core.Converter) {
	if value := ctx.Value(constant.ContextConverter); value != nil {
		language = value.(core.Converter)
	} else {
		language = new(converter.Same)
	}

	return
}

func (v *Validator) detectLanguage(ctx context.Context) (language string) {
	if value := ctx.Value(constant.ContextAcceptLanguage); value != nil {
		language = value.(string)
	} else {
		language = "zh-CN"
	}

	return
}
