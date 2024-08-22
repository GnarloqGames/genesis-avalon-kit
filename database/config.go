package database

const (
	DriverCockroach = "cockroach"
	DriverMock      = "mock"
)

var (
	kind string = DriverCockroach
)

func SetKind(driver string) {
	switch driver {
	case DriverCockroach, DriverMock:
		kind = driver
	}
}
