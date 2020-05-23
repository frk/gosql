package testdata

type SelectNotExistsWithWhereBlockQuery struct {
	NotExists bool `rel:"test_user:u"`
	Where     struct {
		Email string `sql:"u.email,@lower"`
	}
}
