package gosql

type directive struct {
	// This field serves as an indicator that the type is actually a directive
	// as opposed to other types declared by this package. Used by helper
	// functions defined in the internal/ packages.
	_isdir struct{}
}

type (
	// The Column directive can be used in Where blocks to generate simple,
	// single column predicates.
	Column    directive
	Relation  directive
	LeftJoin  directive
	RightJoin directive
	FullJoin  directive
	CrossJoin directive
)

type AfterScanner interface {
	AfterScan()
}

type ErrorHandler interface {
	HandleError(err error) error
}
