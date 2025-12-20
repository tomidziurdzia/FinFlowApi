package bootstrap

import (
	"log"

	"fin-flow-api/internal/infrastructure/config"
	"fin-flow-api/internal/infrastructure/db"
	"fin-flow-api/internal/infrastructure/hash"
	"fin-flow-api/internal/infrastructure/jwt"
	httptransport "fin-flow-api/internal/interfaces/http"
	userservices "fin-flow-api/internal/users/application/services"
	userpostgres "fin-flow-api/internal/users/infrastructure/persistence/postgres"
	usershttp "fin-flow-api/internal/users/interfaces/http"
)

type App struct {
	Server     *httptransport.Server
	DB         *db.DB
	Config     *config.Config
	UserService *userservices.UserService
}

func NewApp() (*App, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	database, err := db.NewDB(&cfg.Database)
	if err != nil {
		return nil, err
	}

	hashService := hash.NewService()
	jwtService := jwt.NewService()

	userRepo := userpostgres.NewRepository(database.Pool)

	userService := userservices.NewUserService(userRepo, hashService, cfg.App.SystemUser)

	userHandler := usershttp.NewHandler(userService)
	usershttp.SetHandler(userHandler)

	authHandler := usershttp.NewAuthHandler(userRepo, hashService, jwtService)
	usershttp.SetAuthHandler(authHandler)

	httpCfg := httptransport.Config{
		Addr:              cfg.Port,
		ReadTimeout:       cfg.Server.ReadTimeout,
		ReadHeaderTimeout: cfg.Server.ReadHeaderTimeout,
		WriteTimeout:      cfg.Server.WriteTimeout,
		IdleTimeout:        cfg.Server.IdleTimeout,
		ShutdownTimeout:   cfg.Server.ShutdownTimeout,
	}
	srv := httptransport.NewServer(httpCfg, jwtService)

	log.Println("Application initialized successfully")

	return &App{
		Server:      srv,
		DB:          database,
		Config:      cfg,
		UserService: userService,
	}, nil
}

func (a *App) Close() error {
	if a.DB != nil {
		a.DB.Close()
	}
	return nil
}