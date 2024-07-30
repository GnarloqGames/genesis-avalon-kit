package cockroach

var (
	username string = ""
	password string = ""
	hostname string = ""
	port     uint16 = 26257
	database string = ""
)

func SetUsername(val string) {
	username = val
}

func SetPassword(val string) {
	password = val
}

func SetHostname(val string) {
	hostname = val
}

func SetPort(val uint16) {
	port = val
}

func SetDatabase(val string) {
	database = val
}
