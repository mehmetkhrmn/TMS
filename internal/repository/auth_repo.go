package repository

import (
	"TMS/internal/models"
	"database/sql"
)

func (r *Repository) GetAuthUserByUsername(username string) (*models.AuthUser, error) {
	var authUser models.AuthUser
	query := "SELECT id,username,password_hash,role,entity_id FROM auth_users WHERE username=$1"
	err := r.Db.QueryRow(query, username).Scan(&authUser.ID, &authUser.Username, &authUser.PasswordHash, &authUser.Role, &authUser.EntityId)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &authUser, nil
}
func (r *Repository) CreateAuthUser(tx *sql.Tx, user *models.AuthUser) error {
	query := "INSERT INTO auth_users (username, password_hash,role,entity_id) VALUES ($1, $2,$3,$4) RETURNING id"
	err := tx.QueryRow(query, user.Username, user.PasswordHash, user.Role, user.EntityId).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetAuthUser(id int) (*models.AuthUser, error) {
	var user models.AuthUser
	query := "SELECT * FROM auth_users WHERE id=$1"
	err := r.Db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
