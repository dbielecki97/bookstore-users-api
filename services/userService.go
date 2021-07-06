package services

import (
	"github.com/dbielecki97/bookstore-users-api/domain/users"
	"github.com/dbielecki97/bookstore-users-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(userId int64) (*users.User, *errors.RestErr) {
	result := users.User{ID: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}

	return &result, nil
}

func UpdateUser(u users.User) (*users.User, *errors.RestErr) {
	cur, err := GetUser(u.ID)
	if err != nil {
		return nil, err
	}

	if err := u.Validate(); err != nil {
		return nil, err
	}

	cur.FirstName = u.FirstName
	cur.LastName = u.LastName
	cur.Email = u.Email

	if err := cur.Update(); err != nil {
		return nil, err
	}

	return cur, nil
}

func PatchUser(u users.User) (*users.User, *errors.RestErr) {
	cur, err := GetUser(u.ID)
	if err != nil {
		return nil, err
	}

	if u.FirstName != "" {
		cur.FirstName = u.FirstName
	}

	if u.LastName != "" {
		cur.LastName = u.LastName
	}

	if u.Email != "" {
		cur.Email = u.Email
	}

	if err := cur.Update(); err != nil {
		return nil, err
	}

	return cur, nil
}
