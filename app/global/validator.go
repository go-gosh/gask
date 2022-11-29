package global

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

func init() {
	Validate.SetTagName("binding")
}
