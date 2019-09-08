package testdata

import (
	"database/sql"
	"encoding/json"
	"math/big"
	"net"
	"time"

	"github.com/frk/gosql"
	"github.com/frk/gosql/testdata/common"
)

type T struct{} // stub type

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
	} `rel:"relation_a"`
}

//OK: nested fields
type InsertTestOK8 struct {
	Rel struct {
		Foobar common.Foo `sql:">foo_"`
	} `rel:"relation_a"`
}

//OK: where block
type DeleteTestOK9 struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		ID int `sql:"id"`
	}
}

//OK: where block with gosql.Column directive and all possible predicates
type DeleteTestOK10 struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		_ gosql.Column `sql:"column_a notnull"`
		_ gosql.Column `sql:"column_b isnull" bool:"and"`
		_ gosql.Column `sql:"column_c nottrue" bool:"or"`
		_ gosql.Column `sql:"column_d istrue"`
		_ gosql.Column `sql:"column_e notfalse" bool:"or"`
		_ gosql.Column `sql:"column_f isfalse" bool:"or"`
		_ gosql.Column `sql:"column_g notunknown"`
		_ gosql.Column `sql:"column_h isunknown"`
	}
}

//OK: nested where blocks
type DeleteTestOK11 struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		x struct {
			foo int          `sql:"column_foo"`
			_   gosql.Column `sql:"column_a isnull"`
		} `sql:">"`
		y struct {
			_   gosql.Column `sql:"column_b nottrue"`
			bar string       `sql:"column_bar" bool:"or"`
			z   struct {
				baz  bool         `sql:"column_baz"`
				quux string       `sql:"column_quux"`
				_    gosql.Column `sql:"column_c istrue" bool:"or"`
			} `sql:">"`
		} `sql:">" bool:"or"`
		_   gosql.Column `sql:"column_d notfalse" bool:"or"`
		_   gosql.Column `sql:"column_e isfalse"`
		foo int          `sql:"column_foo"`
	}
}

//OK: where block with field items and specific comparison operators
type DeleteTestOK12 struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		a int `sql:"column_a <"`
		b int `sql:"column_b >"`
		c int `sql:"column_c <="`
		d int `sql:"column_d >="`
		e int `sql:"column_e ="`
		f int `sql:"column_f <>"`
		g int `sql:"column_g"` // defaults to "="
	}
}

//OK: where block with gosql.Column directive and comparison expressions
type DeleteTestOK13 struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		_ gosql.Column `sql:"column_a <> column_b"`
		_ gosql.Column `sql:"t.column_c=u.column_d"`
		_ gosql.Column `sql:"t.column_e>123"`
		_ gosql.Column `sql:"t.column_f = 'active'"`
		_ gosql.Column `sql:"t.column_g <> true"`
	}
}

//OK: where block with "between" predicates
type DeleteTestOK14 struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		a struct {
			x int `sql:"x"`
			y int `sql:"y"`
		} `sql:"column_a isbetween"`
		b struct {
			_ gosql.Column `sql:"column_x,x"`
			_ gosql.Column `sql:"column_y,y"`
		} `sql:"column_b isbetweensym"`
		c struct {
			_ gosql.Column `sql:"column_z,x"`
			z int          `sql:"y"`
		} `sql:"column_c notbetweensym"`
		d struct {
			z int          `sql:"x"`
			_ gosql.Column `sql:"column_z,y"`
		} `sql:"column_d notbetween"`
	}
}

//OK: where block with "distinct from" predicates
type DeleteTestOK_DistinctFrom struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		a int `sql:"column_a isdistinct"`
		b int `sql:"column_b notdistinct"`

		_ gosql.Column `sql:"column_c isdistinct column_x"`
		_ gosql.Column `sql:"column_d notdistinct column_y"`
	}
}

//OK: where block with array comparisons
type DeleteTestOK_ArrayComparisons struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		a []int   `sql:"column_a isin"`
		b [5]int  `sql:"column_b notin"`
		c []int   `sql:"column_c=any"`
		d [10]int `sql:"column_d >some"`
		e []int   `sql:"column_e <= all"`
	}
}

//OK: where block with pattern matching
type DeleteTestOK_PatternMatching struct {
	Rel   struct{} `rel:"relation_a"`
	Where struct {
		a string `sql:"column_a islike"`
		b string `sql:"column_b notlike"`
		c string `sql:"column_c issimilar"`
		d string `sql:"column_d notsimilar"`
		e string `sql:"column_e ~"`
		f string `sql:"column_f~*"`
		g string `sql:"column_g!~"`
		h string `sql:"column_h !~*"`
	}
}

//OK: DELETE with Using joinblock
type DeleteTestOK_Using struct {
	Rel   struct{} `rel:"relation_a:a"`
	Using struct {
		_ gosql.Relation  `sql:"relation_b:b"`
		_ gosql.LeftJoin  `sql:"relation_c:c,c.b_id = b.id"`
		_ gosql.RightJoin `sql:"relation_d:d,d.c_id = c.id;d.num > b.num"`
		_ gosql.FullJoin  `sql:"relation_e:e,e.d_id = d.id,e.is_foo isfalse"`
		_ gosql.CrossJoin `sql:"relation_f:f"`
	}
	Where struct {
		_ gosql.Column `sql:"a.id = d.a_id"`
	}
}

//OK: UPDATE with From joinblock
type UpdateTestOK_From struct {
	Rel  struct{} `rel:"relation_a:a"`
	From struct {
		_ gosql.Relation  `sql:"relation_b:b"`
		_ gosql.LeftJoin  `sql:"relation_c:c,c.b_id = b.id"`
		_ gosql.RightJoin `sql:"relation_d:d,d.c_id = c.id;d.num > b.num"`
		_ gosql.FullJoin  `sql:"relation_e:e,e.d_id = d.id,e.is_foo isfalse"`
		_ gosql.CrossJoin `sql:"relation_f:f"`
	}
	Where struct {
		_ gosql.Column `sql:"a.id = d.a_id"`
	}
}

//OK: SELECT with Join joinblock
type SelectTestOK_Join struct {
	Rel  struct{} `rel:"relation_a:a"`
	Join struct {
		_ gosql.LeftJoin  `sql:"relation_b:b,b.a_id = a.id"`
		_ gosql.LeftJoin  `sql:"relation_c:c,c.b_id = b.id"`
		_ gosql.RightJoin `sql:"relation_d:d,d.c_id = c.id;d.num > b.num"`
		_ gosql.FullJoin  `sql:"relation_e:e,e.d_id = d.id,e.is_foo isfalse"`
		_ gosql.CrossJoin `sql:"relation_f:f"`
	}
	Where struct {
		_ gosql.Column `sql:"a.id = d.a_id"`
	}
}

//OK: Update with All directive
type UpdateTestOK_All struct {
	Rel struct{} `rel:"relation_a:a"`
	_   gosql.All
}

//OK: Delete with All directive
type DeleteTestOK_All struct {
	Rel struct{} `rel:"relation_a:a"`
	_   gosql.All
}

//OK: Delete with Return directive
type DeleteTestOK_Return struct {
	Rel struct{}     `rel:"relation_a:a"`
	_   gosql.Return `sql:"*"`
}

//OK: Insert with Return directive
type InsertTestOK_Return struct {
	Rel struct{}     `rel:"relation_a:a"`
	_   gosql.Return `sql:"a.foo,a.bar,a.baz"`
}

//OK: Update with Return directive
type UpdateTestOK_Return struct {
	Rel struct{}     `rel:"relation_a:a"`
	_   gosql.Return `sql:"a.foo,a.bar,a.baz"`
}

//OK: Insert with Default directive
type InsertTestOK_Default struct {
	Rel struct{}      `rel:"relation_a:a"`
	_   gosql.Default `sql:"*"`
}

//OK: Update with Default directive
type UpdateTestOK_Default struct {
	Rel struct{}      `rel:"relation_a:a"`
	_   gosql.Default `sql:"a.foo,a.bar,a.baz"`
}

//OK: Insert with Force directive
type InsertTestOK_Force struct {
	Rel struct{}    `rel:"relation_a:a"`
	_   gosql.Force `sql:"*"`
}

//OK: Update with Force directive
type UpdateTestOK_Force struct {
	Rel struct{}    `rel:"relation_a:a"`
	_   gosql.Force `sql:"a.foo,a.bar,a.baz"`
}

type myerrorhandler struct{}

func (myerrorhandler) HandleError(e error) error { return e }

//OK: Select with ErrorHandler field
type SelectTestOK_ErrorHandler struct {
	Rel struct{} `rel:"relation_a:a"`
	eh  myerrorhandler
}

//OK: Insert with embedded ErrorHandler field
type InsertTestOK_ErrorHandler struct {
	Rel struct{} `rel:"relation_a:a"`
	myerrorhandler
}

//OK: Select with Count field
type SelectTestOK_Count struct {
	Count int `rel:"relation_a:a"`
}

//OK: Select with Exists field
type SelectTestOK_Exists struct {
	Exists bool `rel:"relation_a:a"`
}

//OK: Select with NotExists field
type SelectTestOK_NotExists struct {
	NotExists bool `rel:"relation_a:a"`
}

//OK: Delete with Relation directive
type DeleteTestOK_Relation struct {
	_ gosql.Relation `rel:"relation_a:a"`
}

//OK: Select with Limit directive
type SelectTestOK_LimitDirective struct {
	Rel []T         `rel:"relation_a:a"`
	_   gosql.Limit `sql:"25"`
}

//OK: Select with Limit field
type SelectTestOK_LimitField struct {
	Rel   []T `rel:"relation_a:a"`
	Limit int `sql:"10"`
}

//OK: Select with Offset directive
type SelectTestOK_OffsetDirective struct {
	Rel []T          `rel:"relation_a:a"`
	_   gosql.Offset `sql:"25"`
}

//OK: Select with Offset field
type SelectTestOK_OffsetField struct {
	Rel    []T `rel:"relation_a:a"`
	Offset int `sql:"10"`
}

//OK: Select with OrderBy directive
type SelectTestOK_OrderByDirective struct {
	Rel []T           `rel:"relation_a:a"`
	_   gosql.OrderBy `sql:"a.foo:nullsfirst,-a.bar:nullsfirst,-a.baz,a.quux:nullslast"`
}

//OK: Insert with Override directive
type InsertTestOK_OverrideDirective struct {
	Rel []T            `rel:"relation_a:a"`
	_   gosql.Override `sql:"system"`
}

//OK: Filter with TextSearch directive
type FilterTestOK_TextSearchDirective struct {
	_ T                `rel:"relation_a:a"`
	_ gosql.TextSearch `sql:"a.ts_document"`
}

//OK: Insert with onconflict block
type InsertTestOK_OnConflict struct {
	Rel        []T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Ignore
	}
}

//OK: Insert with onconflict block with column target
type InsertTestOK_OnConflictColumn struct {
	Rel        []T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Column `sql:"a.id"`
		_ gosql.Ignore
	}
}

//OK: Insert with onconflict block with constraint target and update action
type InsertTestOK_OnConflictConstraint struct {
	Rel        []T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Constraint `sql:"relation_constraint_xyz"`
		_ gosql.Update     `sql:"a.foo,a.bar,a.baz"`
	}
}

//OK: Insert with onconflict block with index target and update action
type InsertTestOK_OnConflictIndex struct {
	Rel        []T `rel:"relation_a:a"`
	OnConflict struct {
		_ gosql.Index  `sql:"relation_index_xyz"`
		_ gosql.Update `sql:"*"`
	}
}

//OK: Delete with Result field
type DeleteTestOK_ResultField struct {
	_     gosql.Relation `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"a.is_inactive istrue"`
	}
	Result []T
}

//OK: Delete with RowsAffected field
type DeleteTestOK_RowsAffected struct {
	_     gosql.Relation `rel:"relation_a:a"`
	Where struct {
		_ gosql.Column `sql:"a.is_inactive istrue"`
	}
	RowsAffected int
}

//OK: Select with Filter field
type SelectTestOK_FilterField struct {
	Rel    []T `rel:"relation_a:a"`
	Filter gosql.Filter
}

//OK: test field types basic
type SelectTestOK_FieldTypesBasic struct {
	Rel struct {
		f1  bool       `sql:"c1"`
		f2  byte       `sql:"c2"`
		f3  rune       `sql:"c3"`
		f4  int8       `sql:"c4"`
		f5  int16      `sql:"c5"`
		f6  int32      `sql:"c6"`
		f7  int64      `sql:"c7"`
		f8  int        `sql:"c8"`
		f9  uint8      `sql:"c9"`
		f10 uint16     `sql:"c10"`
		f11 uint32     `sql:"c11"`
		f12 uint64     `sql:"c12"`
		f13 uint       `sql:"c13"`
		f14 uintptr    `sql:"c14"`
		f15 float32    `sql:"c15"`
		f16 float64    `sql:"c16"`
		f17 complex64  `sql:"c17"`
		f18 complex128 `sql:"c18"`
		f19 string     `sql:"c19"`
	} `rel:"relation_a:a"`
}

//OK: test field types slices, arrays, maps, and pointers
type SelectTestOK_FieldTypesSlices struct {
	Rel struct {
		f1  []bool                    `sql:"c1"`
		f2  []byte                    `sql:"c2"`
		f3  []rune                    `sql:"c3"`
		f4  net.HardwareAddr          `sql:"c4"`
		f5  json.RawMessage           `sql:"c5"`
		f6  []json.Marshaler          `sql:"c6"`
		f7  []json.RawMessage         `sql:"c7"`
		f8  [][]byte                  `sql:"c8"`
		f9  [][2][2]float64           `sql:"c9"`
		f10 [][][2]float64            `sql:"c10"`
		f11 map[string]sql.NullString `sql:"c11"`
		f12 []map[string]*string      `sql:"c12"`
		f13 [][2]*big.Int             `sql:"c13"`
	} `rel:"relation_a:a"`
}

//OK: test field types interfaces
type SelectTestOK_FieldTypesInterfaces struct {
	Rel struct {
		f1 json.Marshaler   `sql:"c1"`
		f2 json.Unmarshaler `sql:"c2"`
		f3 interface {
			json.Marshaler
			json.Unmarshaler
		} `sql:"c3"`
	} `rel:"relation_a:a"`
}

//OK: test typeinfo.string()
type SelectTestOK_typeinfo_string struct {
	Rel struct {
		f01 bool                        `sql:"c01"`
		f02 *bool                       `sql:"c02"`
		f03 []bool                      `sql:"c03"`
		f04 string                      `sql:"c04"`
		f05 *string                     `sql:"c05"`
		f06 []string                    `sql:"c06"`
		f07 [][]string                  `sql:"c07"`
		f08 map[string]string           `sql:"c08"`
		f09 map[string]*string          `sql:"c09"`
		f10 []map[string]string         `sql:"c10"`
		f11 []map[string]*string        `sql:"c11"`
		f12 byte                        `sql:"c12"`
		f13 *byte                       `sql:"c13"`
		f14 []byte                      `sql:"c14"`
		f15 [][]byte                    `sql:"c15"`
		f16 [16]byte                    `sql:"c16"`
		f17 [][16]byte                  `sql:"c17"`
		f18 rune                        `sql:"c18"`
		f19 *rune                       `sql:"c19"`
		f20 []rune                      `sql:"c20"`
		f21 [][]rune                    `sql:"c21"`
		f22 int8                        `sql:"c22"`
		f23 *int8                       `sql:"c23"`
		f24 []int8                      `sql:"c24"`
		f25 [][]int8                    `sql:"c25"`
		f26 int16                       `sql:"c26"`
		f27 *int16                      `sql:"c27"`
		f28 []int16                     `sql:"c28"`
		f29 [][]int16                   `sql:"c29"`
		f30 int32                       `sql:"c30"`
		f31 *int32                      `sql:"c31"`
		f32 []int32                     `sql:"c32"`
		f33 [2]int32                    `sql:"c33"`
		f34 [][2]int32                  `sql:"c34"`
		f35 int64                       `sql:"c35"`
		f36 *int64                      `sql:"c36"`
		f37 []int64                     `sql:"c37"`
		f38 [2]int64                    `sql:"c38"`
		f39 [][2]int64                  `sql:"c39"`
		f40 float32                     `sql:"c40"`
		f41 *float32                    `sql:"c41"`
		f42 []float32                   `sql:"c42"`
		f43 float64                     `sql:"c43"`
		f44 *float64                    `sql:"c44"`
		f45 []float64                   `sql:"c45"`
		f46 [2]float64                  `sql:"c46"`
		f47 [][2]float64                `sql:"c47"`
		f48 [][][2]float64              `sql:"c48"`
		f49 [2][2]float64               `sql:"c49"`
		f50 [][2][2]float64             `sql:"c50"`
		f51 [3]float64                  `sql:"c51"`
		f52 [][3]float64                `sql:"c52"`
		f53 *net.IPNet                  `sql:"c53"`
		f54 []*net.IPNet                `sql:"c54"`
		f55 time.Time                   `sql:"c55"`
		f56 *time.Time                  `sql:"c56"`
		f57 []time.Time                 `sql:"c57"`
		f58 []*time.Time                `sql:"c58"`
		f59 [2]time.Time                `sql:"c59"`
		f60 [][2]time.Time              `sql:"c60"`
		f61 net.HardwareAddr            `sql:"c61"`
		f62 []net.HardwareAddr          `sql:"c62"`
		f63 big.Int                     `sql:"c63"`
		f64 *big.Int                    `sql:"c64"`
		f65 []big.Int                   `sql:"c65"`
		f66 []*big.Int                  `sql:"c66"`
		f67 [2]big.Int                  `sql:"c67"`
		f68 [2]*big.Int                 `sql:"c68"`
		f69 [][2]*big.Int               `sql:"c69"`
		f70 map[string]sql.NullString   `sql:"c70"`
		f71 []map[string]sql.NullString `sql:"c71"`
		f72 json.RawMessage             `sql:"c72"`
		f73 []json.RawMessage           `sql:"c73"`
	} `rel:"relation_a:a"`
}
