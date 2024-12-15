package internal

type App struct {
	db     *MongoDB
	config *Config
}

func NewApp(config *Config, db *MongoDB) *App {
	return &App{
		config: config,
		db:     db,
	}
}
