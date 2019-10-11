package typetests

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"math/big"
	"net"
	"os"
	"testing"
	"time"
)

var _ = sql.DB{}
var _ = big.Int{}
var _ = net.IPNet{}
var _ = time.Time{}

type scanner struct {
	v interface{}
	s string
}

func (s *scanner) Scan(src interface{}) error {
	switch v := src.(type) {
	case []byte:
		s.s = string(v)
	case string:
		s.s = v
	}
	return nil
}

type valuer struct{ v string }

func (v valuer) Value() (driver.Value, error) {
	return []byte(v.v), nil
}

func Test_libpq_(t *testing.T) {
	const selattrs = `SELECT
	a.attname
	, t.typname
	FROM pg_attribute a
	LEFT JOIN pg_type t ON t.oid=a.atttypid
	WHERE a.attrelid='typetest_table'::regclass
	AND a.attnum > 0
	ORDER BY a.attnum ASC` //`

	rows, err := libpq.db.Query(selattrs)
	if err != nil {
		t.Fatal(err)
	}
	defer rows.Close()

	type column struct {
		name string
		typ  string
	}
	var cols []column
	for rows.Next() {
		var c column
		if err := rows.Scan(&c.name, &c.typ); err != nil {
			t.Fatal(err)
		}
		if c.name == "id" {
			continue
		}
		cols = append(cols, c)
	}
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	outfile, err := os.Create("./types.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer outfile.Close()
	defer outfile.Sync()

	tests := []struct {
		val interface{}
		dst interface{}
		dbg bool
	}{
		// prints a list of postgres types which can be used to hold
		// a value of a specific go type.
		{val: true, dst: new(bool)},
		{val: `foo`, dst: new(string)},
		{val: byte('a'), dst: new(byte)},
		{val: []byte(`foo`), dst: new([]byte)},
		{val: rune('魔'), dst: new(rune)},
		{val: int(123), dst: new(int)},
		{val: int8(123), dst: new(int8)},
		{val: int16(123), dst: new(int16)},
		{val: int32(123), dst: new(int32)},
		{val: int64(123), dst: new(int64)},
		{val: uint(123), dst: new(uint)},
		{val: uint8(123), dst: new(uint8)},
		{val: uint16(123), dst: new(uint16)},
		{val: uint32(123), dst: new(uint32)},
		{val: uint64(123), dst: new(uint64)},
		{val: float32(0.5), dst: new(float32)},
		{val: float64(0.5), dst: new(float64)},
		{val: time.Now(), dst: new(time.Time)},

		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]uint)}},
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]uint8)}},
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]uint16)}},
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]uint32)}},
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]uint64)}},

		// for go types that libpq does not translate directly to no
		// postgres type, the placeholders "valuer" and "scanner" are
		// used and as input a string representation of a literal is
		// provided. The actual valuer and scanner implementations,
		// once they've been implemented, can replace the placeholders.

		{val: valuer{`{t,f}`}, dst: &scanner{v: new([]bool)}},
		{val: valuer{`{'foo','bar'}`}, dst: &scanner{v: new([]string)}},
		{val: valuer{`{{'foo', 'bar'}, {'baz', 'quux'}}`}, dst: &scanner{v: new([][]string)}},
		{val: valuer{`"a"=>"1", "b"=>"2"`}, dst: &scanner{v: new(map[string]string)}},
		{val: valuer{`{"\"a\"=>\"1\"","\"b\"=>\"2\""}`}, dst: &scanner{v: new([]map[string]string)}},
		{val: valuer{`{"\\xdeadbeef","\\xdeadbeef"}`}, dst: &scanner{v: new([][]byte)}},
		{val: valuer{`c3548a67-e88a-4970-9f20-90ed62cbfd40`}, dst: &scanner{v: new([16]byte)}},
		{val: valuer{`{c3548a67-e88a-4970-9f20-90ed62cbfd40,a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11}`}, dst: &scanner{v: new([][16]byte)}},
		{val: valuer{`{'魔','馬'}`}, dst: &scanner{v: new([]rune)}},
		{val: valuer{`{{'魔','馬'},{'馬','魔'}}`}, dst: &scanner{v: new([][]rune)}},
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]int)}},
		{val: valuer{`[0,10)`}, dst: &scanner{v: new([2]int)}},
		{val: valuer{`{"[0,10)","[0,100)"}`}, dst: &scanner{v: new([][2]int)}},
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]int8)}},
		{val: valuer{`{"1 2 3","4 5 6"}`}, dst: &scanner{v: new([][]int8)}},
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]int16)}},
		{val: valuer{`{"1 2 3","4 5 6"}`}, dst: &scanner{v: new([][]int16)}},
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]int32)}},
		{val: valuer{`[0,10)`}, dst: &scanner{v: new([2]int32)}},
		{val: valuer{`{"[0,10)","[0,100)"}`}, dst: &scanner{v: new([][2]int32)}},
		{val: valuer{`$1.99`}, dst: &scanner{v: new(int64)}}, // money
		{val: valuer{`{1,2,3}`}, dst: &scanner{v: new([]int64)}},
		{val: valuer{`{$1.99,$0.99}`}, dst: &scanner{v: new([]int64)}}, // moneyarr
		{val: valuer{`[0,10)`}, dst: &scanner{v: new([2]int64)}},
		{val: valuer{`{"[0,10)","[0,100)"}`}, dst: &scanner{v: new([][2]int64)}},
		{val: valuer{`{0.5,1.5}`}, dst: &scanner{v: new([]float32)}},
		{val: valuer{`{0.5,1.5}`}, dst: &scanner{v: new([]float64)}},
		{val: valuer{`[0.5,10.5)`}, dst: &scanner{v: new([2]float64)}},                                                               // range
		{val: valuer{`(0.5,10.5)`}, dst: &scanner{v: new([2]float64)}},                                                               // point
		{val: valuer{`{"[0.5,10.5)","[0.5,100.5)"}`}, dst: &scanner{v: new([][2]float64)}},                                           // rangearr
		{val: valuer{`{"(0.5,10.5)","(0.5,100.5)"}`}, dst: &scanner{v: new([][2]float64)}},                                           // pointarr
		{val: valuer{`[(0.5,10.5),(0.5,100.5),(0.5,1000.5)]`}, dst: &scanner{v: new([][2]float64)}},                                  // path
		{val: valuer{`((0.5,10.5),(0.5,100.5),(0.5,1000.5))`}, dst: &scanner{v: new([][2]float64)}},                                  // polygon
		{val: valuer{`{"[(0.5,10.5),(0.5,100.5),(0.5,1000.5)]","[(0.5,10.5),(0.5,100.5)]"}`}, dst: &scanner{v: new([][][2]float64)}}, // patharr
		{val: valuer{`{"((0.5,10.5),(0.5,100.5),(0.5,1000.5))","((0.5,10.5),(0.5,100.5))"}`}, dst: &scanner{v: new([][][2]float64)}}, // polygonarr
		{val: valuer{`[(0.5,10.5),(0.5,100.5)]`}, dst: &scanner{v: new([2][2]float64)}},                                              // lseg
		{val: valuer{`(0.5,100.5),(0.5,10.5)`}, dst: &scanner{v: new([2][2]float64)}},                                                // box
		{val: valuer{`{"[(0.5,10.5),(0.5,100.5)]","[(0.5,10.5),(0.5,100.5)]"}`}, dst: &scanner{v: new([][2][2]float64)}},             // lsegarr
		{val: valuer{`{(0.5,100.5),(0.5,10.5);(0.5,100.5),(0.5,10.5)}`}, dst: &scanner{v: new([][2][2]float64)}},                     // boxarr
		{val: valuer{`{0.5,100.5,100.5}`}, dst: &scanner{v: new([3]float64)}},                                                        // line
		{val: valuer{`{"{0.5,100.5,100.5}","{0.5,100.5,100.5}"}`}, dst: &scanner{v: new([][3]float64)}},                              // linearr
		{val: valuer{`192.168.100.128/25`}, dst: &scanner{v: new(net.IPNet)}},
		{val: valuer{`{192.168.100.128/25,192.168.0.128/25}`}, dst: &scanner{v: new([]net.IPNet)}},
		{val: valuer{`{"1999-01-08 04:05:06","2004-01-08 05:08:17"}`}, dst: &scanner{v: new([]time.Time)}},                                                                            // timestamparr
		{val: valuer{`{"1999-01-08 13:05:06+01","2004-01-07 21:08:17+01"}`}, dst: &scanner{v: new([]time.Time)}},                                                                      // timestamptzarr
		{val: valuer{`{1999-01-08,2004-01-08}`}, dst: &scanner{v: new([]time.Time)}},                                                                                                  // datearr
		{val: valuer{`{04:05:06,05:08:17}`}, dst: &scanner{v: new([]time.Time)}},                                                                                                      // timearr
		{val: valuer{`{04:05:06-08,05:08:17+09}`}, dst: &scanner{v: new([]time.Time)}},                                                                                                // timetzarr
		{val: valuer{`[1999-01-08,2004-01-08)`}, dst: &scanner{v: new([2]time.Time)}},                                                                                                 // daterange
		{val: valuer{`["1999-01-08 04:05:06","2004-01-08 05:08:17")`}, dst: &scanner{v: new([2]time.Time)}},                                                                           // tsrange
		{val: valuer{`["1999-01-08 13:05:06+01","2004-01-07 21:08:17+01")`}, dst: &scanner{v: new([2]time.Time)}},                                                                     // tstzrange
		{val: valuer{`{"[1999-01-08,2004-01-08)","[1999-01-08,2004-01-08)"}`}, dst: &scanner{v: new([][2]time.Time)}},                                                                 // daterangearr
		{val: valuer{`{"[\"1999-01-08 04:05:06\",\"2004-01-08 05:08:17\")","[\"1999-01-08 04:05:06\",\"2004-01-08 05:08:17\")"}`}, dst: &scanner{v: new([][2]time.Time)}},             // tsrangearr
		{val: valuer{`{"[\"1999-01-08 13:05:06+01\",\"2004-01-07 21:08:17+01\")","[\"1999-01-08 13:05:06+01\",\"2004-01-07 21:08:17+01\")"}`}, dst: &scanner{v: new([][2]time.Time)}}, // tstzrangearr
		{val: valuer{`9999999999999999999999999999999999999999999999999`}, dst: &scanner{v: new(big.Int)}},                                                                            // numeric
		{val: valuer{`{9999999999999999999999999999999999999999999999999,99}`}, dst: &scanner{v: new([]big.Int)}},                                                                     // numericarr
		{val: valuer{`[99,9999999999999999999999999999999999999999999999999)`}, dst: &scanner{v: new([2]big.Int)}},                                                                    // numrange
		{val: valuer{`{"[99,9999999999999999999999999999999999999999999999999)","[888,888888888888888888888888)"}`}, dst: &scanner{v: new([][2]big.Int)}},                             // numrangearr
		{val: valuer{`99999999999999999999999999999999999999999999999.99`}, dst: &scanner{v: new(big.Float)}},                                                                         // numeric
		{val: valuer{`{99999999999999999999999999999999999999999999999.99,99}`}, dst: &scanner{v: new([]big.Float)}},                                                                  // numericarr
		{val: valuer{`[99,99999999999999999999999999999999999999999999999.99)`}, dst: &scanner{v: new([2]big.Float)}},                                                                 // numrange
		{val: valuer{`{"[99,99999999999999999999999999999999999999999999999.99)","[888,888888888888888888888888.99)"}`}, dst: &scanner{v: new([][2]big.Float)}},                       // numrangearr
		{val: valuer{`"a"=>"1", "b"=>"2"`}, dst: &scanner{v: new(map[string]sql.NullString)}},                                                                                         // hstore
		{val: valuer{`{"\"a\"=>\"1\"","\"b\"=>\"2\""}`}, dst: &scanner{v: new([]map[string]sql.NullString)}},                                                                          // hstorearr
	}

	fmt.Println(len(cols))

	debugcol := "col_uuidarr"

	for _, tt := range tests {
		var oktypes []string
		for _, cc := range cols {
			var id int
			row := libpq.db.QueryRow(`insert into typetest_table (`+cc.name+`) values ($1) returning id`, tt.val)
			if err := row.Scan(&id); err != nil {
				if tt.dbg && cc.name == debugcol {
					fmt.Printf("ERROR:: %v\n", err)
				}
				continue
			}

			row = libpq.db.QueryRow(`select `+cc.name+` from typetest_table where id = $1`, id)
			if err := row.Scan(tt.dst); err != nil {
				if tt.dbg && cc.name == debugcol {
					fmt.Printf("ERROR:: %v\n", err)
				}
				continue
			}

			if s, ok := tt.dst.(*scanner); ok {
				v, ok := tt.val.(valuer)
				if !ok || v.v != s.s {
					if tt.dbg && cc.name == debugcol {
						fmt.Printf("XXX %s <%T>: %s\n", cc.name, s.v, s.s)
					}
					continue
				}
			}
			oktypes = append(oktypes, cc.name[4:])
		}

		typ := tt.dst
		if s, ok := typ.(*scanner); ok {
			typ = s.v
		}

		// switch between true or false to either save the output
		// into a file or print it to stdout.
		if false {
			fmt.Fprintf(outfile, "%T: %+v\n", typ, oktypes)
		} else {
			fmt.Printf("%T: %+v\n", typ, oktypes)
		}
	}
}