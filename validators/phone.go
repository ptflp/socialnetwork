package validators

import (
	"strconv"

	"github.com/nyaruka/phonenumbers"
)

func CheckPhoneFormat(phone string) (string, error) {
	num, err := phonenumbers.Parse(phone, "RU")
	if err != nil {
		return "", err
	}
	phone = "7" + strconv.FormatUint(*num.NationalNumber, 10)

	return phone, nil
}
