package testdata

import (
	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

//BAD: missing relation field
type InsertTestBAD1 struct {
	// no record type ...
}

//BAD: missing `rel` tag
type InsertTestBAD2 struct {
	User *common.User
}

//BAD: invalid datatype kind
type InsertTestBAD3 struct {
	User string `rel:"users_table"`
}

//OK: user datatype
type InsertTestOK1 struct {
	UserRec *common.User `rel:"users_table"`
}

//OK: ignored datatype fields
type InsertTestOK2 struct {
	UserRec struct {
		_     string `sql:"name"` // ignore blank fields
		Name  string `sql:"-"`    // ignore "-" tags
		Name2 string ``           // ignore no `sql` tag
		Name3 string `sql:"name"` // all good
	} `rel:"users_table"`
}

//OK: unnamed iterator func
type SelectTestOK3 struct {
	User func(*common.User) error `rel:"users_table"`
}

type namedIteratorFunc func(*common.User) error

//OK: named iterator func
type SelectTestOK4 struct {
	User namedIteratorFunc `rel:"users_table"`
}

//OK: unnamed iterator interface
type SelectTestOK5 struct {
	User interface {
		Fn(*common.User) error
	} `rel:"users_table"`
}

type namedIterator interface {
	Fn(*common.User) error
}

//OK: named iterator interface
type SelectTestOK6 struct {
	User namedIterator `rel:"users_table"`
}

//OK: tag options with boolean operators
type SelectTestOK7 struct {
	Rel struct {
		a int `sql:"a,pk,auto"`
		b int `sql:"b,nullempty"`
		c int `sql:"c,ro,json"`
		d int `sql:"d,wo"`
		e int `sql:"e,+"`
		f int `sql:"f,coalesce"`
		g int `sql:"g,coalesce(-1)"`
	} `rel:"a_relation"`
}

//OK: nested fields
type InsertTestOK8 struct {
	Rel struct {
		Foobar common.Foo `sql:">foo_"`
	} `rel:"a_relation"`
}

//OK: where block
type DeleteTestOK9 struct {
	Rel   struct{} `rel:"a_relation"`
	Where struct {
		ID int `sql:"id"`
	}
}

//OK: where block with gosql.Column directive and all possible predicates
type DeleteTestOK10 struct {
	Rel   struct{} `rel:"a_relation"`
	Where struct {
		_ gosql.Column `sql:"column_a,notnull"`
		_ gosql.Column `sql:"column_b,isnull" bool:"and"`
		_ gosql.Column `sql:"column_c,nottrue" bool:"or"`
		_ gosql.Column `sql:"column_d,istrue"`
		_ gosql.Column `sql:"column_e,notfalse" bool:"or"`
		_ gosql.Column `sql:"column_f,isfalse" bool:"or"`
		_ gosql.Column `sql:"column_g,notunknown"`
		_ gosql.Column `sql:"column_h,isunknown"`
	}
}

//OK: nested where blocks
type DeleteTestOK11 struct {
	Rel   struct{} `rel:"a_relation"`
	Where struct {
		x struct {
			foo int          `sql:"column_foo"`
			_   gosql.Column `sql:"column_a,isnull"`
		} `sql:">"`
		y struct {
			_   gosql.Column `sql:"column_b,nottrue"`
			bar string       `sql:"column_bar" bool:"or"`
			z   struct {
				baz  bool         `sql:"column_baz"`
				quux string       `sql:"column_quux"`
				_    gosql.Column `sql:"column_c,istrue" bool:"or"`
			} `sql:">"`
		} `sql:">" bool:"or"`
		_   gosql.Column `sql:"column_d,notfalse" bool:"or"`
		_   gosql.Column `sql:"column_e,isfalse"`
		foo int          `sql:"column_foo"`
	}
}
