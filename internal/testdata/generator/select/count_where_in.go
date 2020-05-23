package testdata

type SelectCountWithWhereBlockQuery struct {
	Count int `rel:"test_user:u"`
	Where struct {
		FullName string `sql:"u.full_name islike"`
	}
}
