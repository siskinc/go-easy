package str

import "regexp"

func CheckPhoneNumber(phoneNumber string) (matched bool, err error) {
	pattern := "^(13[0-9]|14[5-9]|15[0-3,5-9]|16[2,5,6,7]|17[0-8]|18[0-9]|19[1,3,5,8,9])\\d{8}$"
	matched, err = regexp.MatchString(pattern, phoneNumber)
	return
}
