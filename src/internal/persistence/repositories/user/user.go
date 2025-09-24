package user

import (
	"context"
	"database/sql"
	"fmt"

	domainUser "service-template/internal/domain/user"
	"service-template/internal/persistence/repositories/base"

	"github.com/doug-martin/goqu/v9"
	"github.com/google/uuid"
)

type userTableSchema struct {
	tableName, id, username, email, created_at, updated_at string
}

var userSchema = userTableSchema{
	tableName:  "app.users",
	id:         "id",
	username:   "username",
	email:      "email",
	created_at: "created_at",
	updated_at: "updated_at",
}

type UserEntity struct {
	Id       uuid.UUID
	Username string
	Email    string
}

type UserRepository struct {
	*base.BaseRepository
	db *goqu.Database
}

func New(db *goqu.Database, publisher base.EventPublisher) *UserRepository {
	return &UserRepository{
		BaseRepository: base.New(publisher),
		db:             db,
	}
}

func (r *UserRepository) Create(ctx context.Context, user *domainUser.User) error {
	query, args, err := r.db.
		Insert(userSchema.tableName).
		Prepared(true).
		Rows(goqu.Record{
			userSchema.id:       user.Id().Value(),
			userSchema.username: user.UserName().Value(),
			userSchema.email:    user.Email().Value(),
		}).
		ToSQL()
	if err != nil {
		return fmt.Errorf("can not create sql query, error %v", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}

	r.BaseRepository.FlashEvents(ctx, user.Aggregate)

	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query, args, err := r.db.
		Delete(userSchema.tableName).
		Prepared(true).
		Where(
			goqu.C(userSchema.id).Eq(id),
		).
		ToSQL()
	if err != nil {
		return fmt.Errorf("can not create sql query, error %v", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to delete user: %v", err)
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return fmt.Errorf("can not delete user by id %s", id)
	}

	// r.BaseRepository.FlashEvents(user.Aggregate)

	return nil
}

func (r *UserRepository) Update(ctx context.Context, user *domainUser.User) error {
	query, args, err := r.db.
		Update(userSchema.tableName).
		Set(goqu.Record{
			userSchema.email: user.Email().Value(),
		}).
		Where(
			goqu.C(userSchema.id).Eq(user.Id().Value()),
		).
		// TODO Ужас, это не дефолтная настройка оказывается, нужно подумать что с этим можно сделать
		Prepared(true).
		ToSQL()
	if err != nil {
		return fmt.Errorf("can not create sql query, error %v", err)
	}

	_, err = r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("failed to update user: %v", err)
	}

	r.BaseRepository.FlashEvents(ctx, user.Aggregate)

	return nil
}

func (r *UserRepository) GetUser(ctx context.Context, id uuid.UUID) (*domainUser.User, error) {
	query, args, err := r.db.
		From(userSchema.tableName).
		Select(
			userSchema.id, userSchema.username, userSchema.email,
		).
		Where(
			goqu.C(userSchema.id).Eq(id),
		).
		Prepared(true).
		ToSQL()

	if err != nil {
		return nil, fmt.Errorf("can not create sql query, error %v", err)
	}

	var userEntity UserEntity
	err = r.db.QueryRowContext(ctx, query, args...).
		Scan(
			&userEntity.Id,
			&userEntity.Username,
			&userEntity.Email,
		)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user: %v", err)
	}

	user, err := domainUser.Restore(
		userEntity.Id,
		userEntity.Username,
		userEntity.Email,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUsers(ctx context.Context, limit int, page int) ([]*domainUser.User, int, error) {
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if page <= 0 {
		page = 1
	}

	offset := (page - 1) * limit
	query, args, err := r.db.
		From(userSchema.tableName).
		Select(
			userSchema.id, userSchema.username, userSchema.email,
		).
		Limit(uint(limit)).
		Offset(uint(offset)).
		Prepared(true).
		ToSQL()

	if err != nil {
		return nil, 0, fmt.Errorf("can not create sql query, error %v", err)
	}

	result, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to create user: %v", err)
	}
	defer result.Close()

	var userEntities []*UserEntity
	for result.Next() {
		user := &UserEntity{}
		err := result.Scan(
			&user.Id,
			&user.Username,
			&user.Email,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %v", err)
		}
		userEntities = append(userEntities, user)
	}

	if err = result.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %v", err)
	}

	users, err := mapToDomainSlice(userEntities)
	if err != nil {
		return nil, 0, err
	}

	totalQuery, totalArgs, err := r.db.
		From(userSchema.tableName).
		Select(goqu.COUNT(goqu.Star())).
		ToSQL()

	if err != nil {
		return nil, 0, fmt.Errorf("can not create count query: %v", err)
	}

	var total int
	err = r.db.QueryRowContext(ctx, totalQuery, totalArgs...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %v", err)
	}

	return users, total, nil
}

func mapToDomainSlice(userEntitySlice []*UserEntity) ([]*domainUser.User, error) {
	result := make([]*domainUser.User, 0, len(userEntitySlice))
	for _, v := range userEntitySlice {
		user, err := domainUser.Restore(
			v.Id,
			v.Username,
			v.Email,
		)
		if err != nil {
			return nil, err
		}

		result = append(result, user)
	}

	return result, nil
}
