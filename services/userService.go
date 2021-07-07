package services

import (
	"github.com/dbielecki97/bookstore-users-api/domain/users"
	"github.com/dbielecki97/bookstore-utils-go/crypto"
	"github.com/dbielecki97/bookstore-utils-go/date"
	"github.com/dbielecki97/bookstore-utils-go/errors"
)

var (
	UserService userService = &defaultUserService{}
)

type userService interface {
	CreateUser(user users.User) (*users.User, *errors.RestErr)
	GetUser(userId int64) (*users.User, *errors.RestErr)
	UpdateUser(u users.User) (*users.User, *errors.RestErr)
	PatchUser(u users.User) (*users.User, *errors.RestErr)
	DeleteUser(userId int64) *errors.RestErr
	Search(status string) (users.Users, *errors.RestErr)
	FindByEmail(users.LoginRequest) (*users.User, *errors.RestErr)
}

type defaultUserService struct{}

func (s *defaultUserService) CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.Validate(); err != nil {
		return nil, err
	}

	user.DateCreated = date.GetNowDBString()
	user.Status = users.StatusActive
	enPass, err := crypto.Generate(user.Password)
	if err != nil {
		return nil, errors.NewInternalServerError("error processing request")
	}
	user.Password = enPass

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *defaultUserService) GetUser(userId int64) (*users.User, *errors.RestErr) {
	result := users.User{ID: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *defaultUserService) UpdateUser(u users.User) (*users.User, *errors.RestErr) {
	cur, err := s.GetUser(u.ID)
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

func (s *defaultUserService) PatchUser(u users.User) (*users.User, *errors.RestErr) {
	cur, err := s.GetUser(u.ID)
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

func (s *defaultUserService) DeleteUser(userId int64) *errors.RestErr {
	user := users.User{ID: userId}
	return user.Delete()
}

func (s *defaultUserService) Search(status string) (users.Users, *errors.RestErr) {
	dao := &users.User{}
	return dao.FindByStatus(status)
}

func (s *defaultUserService) FindByEmail(r users.LoginRequest) (*users.User, *errors.RestErr) {
	u := &users.User{Email: r.Email}

	if err := u.FindByEmail(); err != nil {
		return nil, err
	}

	err := crypto.Compare(u.Password, r.Password)
	if err != nil {
		return nil, errors.NewAuthenticationError("invalid credentials")
	}

	return u, nil
}
