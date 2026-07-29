package repositories

import (
	"auth-service/core/db"
	"auth-service/core/enums"
	globalErrors "auth-service/core/errors"
	"context"
	"fmt"

	"github.com/google/uuid"
)

type AdminRepository struct {
	db db.Database
}

func NewAdminRepository(db db.Database) AdminRepository {
	return AdminRepository{
		db: db,
	}
}

func (a AdminRepository) UpdateRole(ctx context.Context, userID uuid.UUID, userRole enums.Role) error {
	query := `
	UPDATE users
	SET role=$1
	WHERE id=$2;
`
	tag, err := a.db.Exec(ctx, query, userRole, userID)
	if err != nil {
		return fmt.Errorf("role repository : update, err: %w", err)
	}

	if tag.RowsAffected() == 0 {
		return globalErrors.ErrUserNotFound
	}

	return nil
}

func (a AdminRepository) GetRole(ctx context.Context, userID uuid.UUID) (enums.Role, error) {
	query := `
	SELECT role
	FROM users
	WHERE id=$1;
`
	resultRow := a.db.QueryRow(ctx, query, userID)

	var role enums.Role
	if err := resultRow.Scan(&role); err != nil {
		return enums.User, fmt.Errorf("role repository : get, err: %w", err)
	}

	return role, nil
}
