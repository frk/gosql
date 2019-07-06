package gosql

type AfterScanner interface {
	AfterScan()
}

type ErrorHandler interface {
	HandleError(err error) error
}
