package users

func (u *User) HideFields(isPublic bool) {
	u.Password = ""

	if isPublic {
		u.FirstName = ""
		u.LastName = ""
		u.Email = ""
	}
}

func (users Users) HideFields(isPublic bool) Users {
	results := make(Users, 0)
	for _, u := range users {
		u.HideFields(isPublic)
		results = append(results, u)
	}

	return results
}
