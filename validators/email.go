package validators

import "net/mail"

func CheckEmailFormat(email string) error {
	_, err := mail.ParseAddress(email)

	return err
}
