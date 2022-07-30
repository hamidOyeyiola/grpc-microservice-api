package interfaces

type Scanner interface {
	Scan(dest ...interface{}) error
}

type Iterator interface {
	Scanner
	Next() bool
}
