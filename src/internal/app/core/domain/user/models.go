package user

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

type Id struct {
	val uuid.UUID
}

func NewId(id uuid.UUID) (Id, error) {
	if uuid.UUID(id) == uuid.Nil {
		return Id{}, fmt.Errorf("domain validation error, user id is empty")
	}

	return Id{val: id}, nil
}

func (i Id) Value() uuid.UUID {
	return i.val
}

func (i Id) Equals(other Id) bool {
	return i.val == other.val
}

func (i Id) String() string {
	return i.val.String()
}

type UserName struct {
	val string
}

func NewUserName(str string) (UserName, error) {
	userNameCount := utf8.RuneCountInString(str)
	if userNameCount > 100 {
		return UserName{}, fmt.Errorf("domain validation error, user name more then 100 symbols")
	}

	return UserName{val: str}, nil
}

func (i UserName) Value() string {
	return i.val
}

func (un UserName) Compare(other UserName) int {
	return strings.Compare(un.val, other.val)
}

func (un UserName) Equals(other UserName) bool {
	return un.val == other.val
}

func (un UserName) String() string {
	return un.val
}
