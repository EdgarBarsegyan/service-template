package user

import (
	"database/sql"
	"log/slog"
	"service-template/internal/app/config"
	"service-template/internal/app/events/publisher"
	userRepository "service-template/internal/persistence/repositories/user"

	"github.com/doug-martin/goqu/v9"
)

func Build(log *slog.Logger) (*UserCase, error) {
	cfg := config.GlobalInstance
	db, err := sql.Open("postgres", cfg.Db.Url)
	if err != nil {
		return nil, err
	}

	goquDb := goqu.New("postgres", db)
	memDispatcher := publisher.NewMemoryEventDispatcher(log)
	eventPublisher := publisher.New(memDispatcher)
	userRepo := userRepository.New(goquDb, eventPublisher)
	service := NewUserCase(userRepo, cfg)

	return service, nil
}
