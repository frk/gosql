package gosql

type directive struct{}

type (
	// The Column directive can be used in Where blocks to generate simple,
	// single column predicates.
	Column directive
)

type AfterScanner interface {
	AfterScan()
}

type ErrorHandler interface {
	HandleError(err error) error
}
