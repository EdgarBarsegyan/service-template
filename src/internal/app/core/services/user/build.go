package user

import (
	"database/sql"
	"service-template/internal/config"
	userRepository "service-template/internal/persistence/repositories/user"

	"github.com/doug-martin/goqu/v9"
)

func Build() (*UserService, error) {
	cfg := config.GlobalInstance
	db, err := sql.Open("postgres", cfg.Db.Url)
	if err != nil {
		return nil, err
	}

	goquDb := goqu.New("postgres", db)
	userRepo := userRepository.New(goquDb)
	service := NewUserService(userRepo, cfg)

	return service, nil
}
