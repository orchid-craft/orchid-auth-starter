package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

// ValidateStruct melakukan validasi struct dan menghasilkan pesan error berbahasa Indonesia yang jelas
func ValidateStruct(s any) error {
	if err := validate.Struct(s); err != nil {
		if valErrs, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, e := range valErrs {
				field := strings.ToLower(e.Field())
				switch e.Tag() {
				case "required":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' wajib diisi", field))
				case "email":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' harus berupa alamat email yang valid", field))
				case "min":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' minimal berukuran %s karakter", field, e.Param()))
				case "max":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' maksimal berukuran %s karakter", field, e.Param()))
				case "oneof":
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' harus salah satu dari [%s]", field, e.Param()))
				default:
					errMsgs = append(errMsgs, fmt.Sprintf("field '%s' tidak memenuhi aturan %s", field, e.Tag()))
				}
			}
			return fmt.Errorf("validasi gagal: %s", strings.Join(errMsgs, "; "))
		}
		return err
	}
	return nil
}
