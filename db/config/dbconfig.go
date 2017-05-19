package dbconfig

// DBConfig is a basic interface that can provide a database connection Data Source Name (DSN)
type DBConfig interface {
	GetDSN() string
}
