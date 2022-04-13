package config

// Database is a struct that contains DB's configuration variables
type Database struct {
	Host     string
	Port     string
	User     string
	DB       string
	Password string
	Timezone string
}

type Redis struct {
	Host      string
	Port      string
	Password  string
	CacheTime string
}
