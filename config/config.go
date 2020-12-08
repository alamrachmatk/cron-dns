package config

const (
	Host       string = "qdnsfilter.quantum.net.id"
	LimitQuery uint64 = 100

	MariaDBUser     string = "root"
	MariaDBPassword string = ""
	MariaDBDB       string = "kumparan_test"
	MariaDBHost     string = "localhost"
	MariaDBPort     string = "3306"

	RedisHost string = "localhost"
	RedisPort string = "6379"

	DnsLog string = "/var/www/html/dnslog/"
)
