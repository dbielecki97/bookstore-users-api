package users

func (u *User) Marshall(isPublic bool) *User {
	u.Password = ""

	if isPublic {
		u.FirstName = ""
		u.LastName = ""
		u.Email = ""
	}

	return u
}

func (users Users) Marshall(isPublic bool) Users {
	results := make(Users, 0)
	for _, u := range users {
		results = append(results, *u.Marshall(isPublic))
	}

	return results
}
