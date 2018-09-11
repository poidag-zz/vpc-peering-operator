package wiring

type Config struct {
	Poller struct {
		Retries     int `envconfig:"default=5"`
		WaitSeconds int `envconfig:"default=60"`
	}
	ManageRoutes       bool `envconfig:"default=true"`
	WatchAllNamespaces bool `envconfig:"default=false"`
}
