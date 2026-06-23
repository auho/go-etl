package buildtable

type Config struct {
	Recreate bool
	Truncate bool
}

type TableOption func(Tabler)

func WithConfig(config Config) func(Tabler) {
	return func(t Tabler) {
		t.withConfig(config)
	}
}
