package config

type config struct {
	DB    DBConfig
	Redis RedisConfig
}

type DBConfig struct {
	DSN string
	S   int
}
type RedisConfig struct {
	Addr string
	A    int
}
