package testdata

type SelectExistsWithWhereBlockQuery struct {
	Exists bool `rel:"test_user:u"`
	Where  struct {
		Email string `sql:"u.email,@lower"`
	}
}
