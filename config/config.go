package config

const (
	Host       string = "http://localhost:1323"
	LimitQuery uint64 = 100

	MariaDBUser     string = "dnsfilter"
	MariaDBPassword string = "rahasiadns"
	MariaDBDB       string = "dnsfilter"
	MariaDBHost     string = "127.0.0.1"
	MariaDBPort     string = "3306"

	RedisHost string = "localhost"
	RedisPort string = "6379"

	DnsLog string = "/var/www/html/dnslog/"
)
