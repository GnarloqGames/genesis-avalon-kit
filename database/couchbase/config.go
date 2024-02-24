package couchbase

var (
	url      string = "couchbase://localhost"
	bucket   string
	username string
	password string
	scope    string
)

func URL() string {
	return url
}

func SetURL(value string) {
	url = value
}

func Bucket() string {
	return bucket
}

func SetBucket(value string) {
	bucket = value
}

func Username() string {
	return username
}

func SetUsername(value string) {
	username = value
}

func Password() string {
	return password
}

func SetPassword(value string) {
	password = value
}

func Scope() string {
	return scope
}

func SetScope(value string) {
	scope = value
}
