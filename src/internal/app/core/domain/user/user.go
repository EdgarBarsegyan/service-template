package user

import (
	"fmt"
	"service-template/internal/app/core/domain/aggregate"
	"service-template/internal/app/core/domain/common"

	"github.com/google/uuid"
)

type User struct {
	*aggregate.Aggregate
	id       Id
	userName UserName
	email    common.Email
	deleted  bool
}

func New(userName string, email string) (*User, error) {
	dUseName, err := NewUserName(userName)
	if err != nil {
		return nil, err
	}

	dEmail, err := common.NewEmail(email)
	if err != nil {
		return nil, err
	}

	user := &User{
		Aggregate: aggregate.New(),
		id: Id{
			val: uuid.New(),
		},
		userName: dUseName,
		email:    dEmail,
	}
	user.TryPublishEvent(NewUserCreatedEvent(user.id))
	return user, nil
}

func Restore(id uuid.UUID, userName string, email string) (*User, error) {
	dId, err := NewId(id)
	if err != nil {
		return nil, fmt.Errorf("restore error, userId: `%s`, error: %v", id, err) 
	}

	dUserName, err := NewUserName(userName)
	if err != nil {
		return nil, fmt.Errorf("restore error, userId: `%s`, error: %v", dId, err) 
	}

	dEmail, err := common.NewEmail(email)
	if err != nil {
		return nil, fmt.Errorf("restore error, userId: `%s`, error: %v", dId, err) 
	}

	return &User{
		Aggregate: aggregate.New(),
		id:       dId,
		userName: dUserName,
		email:    dEmail,
	}, nil
}

func (u *User) Id() Id {
	return u.id
}

func (u *User) UserName() UserName {
	return u.userName
}

func (u *User) Email() common.Email {
	return u.email
}

func (u *User) SetEmail(email common.Email) common.Email {
	if u.email.Equals(email) {
		return u.email
	}

	u.email = email

	event, ok := u.TryGetEvent(string(UserUpdated))
	if ok {
		updatedEvent := event.(*UserUpdatedEvent)
		updatedEvent.AddUpdatedEventType(UserEmailUpdated)
	} else {
		updatedEvent := NewUserUpdatedEvent(u.id)
		updatedEvent.AddUpdatedEventType(UserEmailUpdated)
		u.TryPublishEvent(updatedEvent)
	}

	return u.email
}

func (u *User) Delete() error {
	if u.deleted {
		return fmt.Errorf("User is deleted")
	}

	u.deleted = true

	ok := u.TryPublishEvent(NewUserDeletedEvent(u.id))
	if !ok {
		return fmt.Errorf("contains user deleted event")
	}

	return nil
}
