package gosql

type directive struct {
	// This field serves as an indicator that the type is actually a directive
	// as opposed to other types declared by this package. Used by helper
	// functions defined in the internal/ packages.
	_isdir struct{}
}

type (
	// The Column directive can be used in Where blocks to generate
	// single column predicates and column to column comparisons.
	Column directive

	// This is inteded to be used with Update and Delete commands to
	// explicitly indicate that, if a "where block" is missing, the
	// command should be executed against all the rows of the relation.
	//
	// This acts as a safeguard against unintentionally omitting the
	// "where block" and then generating a query that would delete/update
	// every single rows in a table.
	All directive

	Relation  directive
	LeftJoin  directive
	RightJoin directive
	FullJoin  directive
	CrossJoin directive

	// The Return directive produces a postgres RETURNING clause. The columns
	// to be returned have to be specified in the struct field's tag.
	//
	// The Return directive can be used in Insert, Update, and Delete commands.
	Return directive

	// The Force directive allows for specifying columns that are usually
	// omitted from a query to actually be included in the query by the
	// command in which the directive is used. The columns to be included
	// have to be specified in the struct field's tag.
	//
	// For example if a table has an "id" column whose value is auto-generated
	// by the database, the corresponding Go struct field's column can be
	// flagged as "auto" which will exclude the column from INSERT/UPDATE
	// queries. However, there may be scenarios in an app where the "id"
	// value is already pre-generated and it needs to be INSERTed together
	// with the record, this is where the Force directive can be used to
	// tell the command to include the "id" column in the query.
	Force directive

	// The Default directive produces the DEFAULT marker in place of values
	// of those columns that are listed in the directive's field tag. This can
	// be used to specify those columns of an INSERT/UPDATE command that should
	// have their value set to the default as defined by the database table.
	Default directive
)

type AfterScanner interface {
	AfterScan()
}

type ErrorHandler interface {
	HandleError(err error) error
}
