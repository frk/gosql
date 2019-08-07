// TODO(mkopriva):
// - given the gosql.Index directive inside an OnConflict block use pg_get_indexdef(index_oid)
//   to retrieve the index's definition, parse that to extract the index expression and then
//   use that expression when generating the ON CONFLICT clause.
//   (https://www.postgresql.org/message-id/204ADCAA-853B-4B5A-A080-4DFA0470B790%40justatheory.com)

package gosql

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

const (
	selectdbname = `SELECT current_database()` //`

	selectdbrelation = `SELECT
	c.oid
	, (c.relkind = 'v' OR c.relkind = 'm')
	FROM pg_class c
	WHERE c.relname = $1
	AND c.relnamespace = to_regnamespace($2)`  //`
)

func dbcheck(url string, cmds []*command) error {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return err
	} else if err := db.Ping(); err != nil {
		return err
	}
	defer db.Close()

	var dbname string // the name of the current database
	if err := db.QueryRow(selectdbname).Scan(&dbname); err != nil {
		return err
	}

	for _, cmd := range cmds {
		c := new(dbchecker)
		c.db = db
		c.dbname = dbname
		c.relid = cmd.rel.relid
		if err := c.load(); err != nil {
			return err
		}
	}
	return nil
}

type dbchecker struct {
	db     *sql.DB
	dbname string // name of the current database, used mainly for error reporting
	relid  relid
	rel    *dbrelation
}

func (c *dbchecker) load() error {
	rel := new(dbrelation)
	rel.name = c.relid.name
	rel.schema = c.relid.qual
	if len(rel.schema) == 0 {
		rel.schema = "public"
	}

	// retrieve relation info
	row := c.db.QueryRow(selectdbrelation, rel.name, rel.schema)
	if err := row.Scan(&rel.oid, &rel.isview); err != nil {
		if err == sql.ErrNoRows {
			return c.newerr(errNoDBRelation)
		}
		return err
	}

	c.rel = rel
	return nil
}

type dbrelation struct {
	oid    int64
	name   string // The name of the relation.
	schema string // The name of the schema to which the relation belongs.
	isview bool   // If set it indicates that the relation is a view.
}

// errors

type dbcheckError struct {
	code dbcheckErrCode
}

func (a *dbchecker) newerr(code dbcheckErrCode) error {
	return &dbcheckError{code: code}
}

func (e *dbcheckError) Error() string {
	return fmt.Sprintf(dbcheckErrCode2string[e.code])
}

type dbcheckErrCode uint

const (
	errNoDBRelation dbcheckErrCode = iota + 1
)

var dbcheckErrCode2string = map[dbcheckErrCode]string{
	errNoDBRelation: "no db relation found",
}
