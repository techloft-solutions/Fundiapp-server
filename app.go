package app

var UserIDNil UserID

type UserID string

func (id UserID) String() string {
	return string(id)
}

type Config struct {
	DB struct {
		DSN string `toml:"dsn"`
	}
}

type App struct {
	// SQLite database used by SQLite service implementations.
	//DB *sqlite.DB
	// Configuration path and parsed config data.
	//config Config
	// HTTP server for handling HTTP communication.
	// SQLite services are attached to it before running.
	//HTTPServer *http.Server
}

func NewApp() *App {
	return &App{
		//Config:     DefaultConfig(),
		//ConfigPath: DefaultConfigPath,

		//DB:         sqlite.NewDB(""),
		//HTTPServer: http.NewServer(),
	}
}

/*
func (a *App) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if m.DB != nil {
		if err := a.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}
*/

type AuthUser struct {
	ID            string `json:"user_id"`
	Name          string `json:"name"`
	PictureUrl    string `json:"picture_url"`
	Provider      string `json:"provider"`
	Email         string `json:"email_address"`
	EmailVerified bool   `json:"email_verified"`
}
