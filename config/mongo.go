package config

type Config struct {
	MongoDBHost     string
	MongoDBPort     string
	MongoDBUsername string
	MongoDBPassword string
}

func LoadConfig() *Config {
	return &Config{
		MongoDBHost:     "localhost",
		MongoDBPort:     "27017",
		MongoDBUsername: "username",
		MongoDBPassword: "password",
	}
}
