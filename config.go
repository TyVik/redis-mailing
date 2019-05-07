package main

// config represents application configuration.
type config struct {
	ServicePort    string   `envconfig:"PORT" default:"8080"`
	LogLevel       string   `envconfig:"LOG_LEVEL" default:"debug"`
	RedisAddr      string   `envconfig:"REDIS_ADDR" required:"true"`
	PubSubChannels []string `envconfig:"REDIS_CHANNELS" required:"true"`
}
