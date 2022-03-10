// Created by Hisen at 2022/3/4.
package v1

// fixme
import (
	"github.com/marmotedu/component-base/pkg/validation"
	"github.com/marmotedu/component-base/pkg/validation/field"
)

func (u *User) Validate() field.ErrorList {
	val := validation.NewValidator(u)
	allErrs := val.Validate()
	if err := validation.IsValidPassword(u.Password); err != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("password"), err.Error(), ""))
	}
	return allErrs
}
