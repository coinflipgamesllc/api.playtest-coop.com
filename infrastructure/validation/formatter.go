package validation

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// JSONFormatter formats binding/validation errors for JSON
type JSONFormatter struct{}

// NewJSONFormatter registers the formatter with gin validation and returns it
func NewJSONFormatter() *JSONFormatter {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}

			return name
		})
	}

	return &JSONFormatter{}
}

// Format returns a "field": "message" map of errors encountered during validation
func (JSONFormatter) Format(verr validator.ValidationErrors) map[string]string {
	errs := make(map[string]string)

	for _, e := range verr {
		err := e.ActualTag()
		if e.Param() != "" {
			err = translateToHumanReadable(err, e.Param())
		}

		errs[e.Field()] = err
	}

	return errs
}

func translateToHumanReadable(tag, param string) string {
	switch tag {
	case "email":
		return "invalid email format"
	case "gtefield":
		return fmt.Sprintf("%s must be greater than or equal to %s", tag, param)
	case "ltefield":
		return fmt.Sprintf("%s must be less than or equal to %s", tag, param)
	case "min":
		return fmt.Sprintf("minimum length allowed is %s", param)
	case "max":
		return fmt.Sprintf("maxiumum length allowed is %s", param)
	case "nefield":
		return fmt.Sprintf("%s must not equal %s", tag, param)
	case "required":
		return "this field is required"
	}

	return fmt.Sprintf("%s = %s", tag, param)
}
