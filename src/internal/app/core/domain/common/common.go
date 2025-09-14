package common

import (
	"net/mail"
	"strings"
)

type Email struct {
	val string
}

func NewEmail(str string) (Email, error) {
	_, err := mail.ParseAddress(str)
	if err != nil {
		return Email{}, err
	}

	return Email{val: str}, nil
}

func (e Email) Value() string {
	return e.val
}

func (e Email) Compare(other Email) int {
	return strings.Compare(e.val, other.val)
}

func (e Email) Equals(other Email) bool {
	return e.val == other.val
}

func (e Email) String() string {
	return e.val
}




