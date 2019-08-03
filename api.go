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
)

type AfterScanner interface {
	AfterScan()
}

type ErrorHandler interface {
	HandleError(err error) error
}
