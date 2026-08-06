package repository

import (
	"TMS/internal/models"
	"database/sql"
)

func (r *Repository) CreateRepresentativeTx(tx *sql.Tx, rep *models.Representative) error {
	query := "INSERT INTO representatives(name) values($1) RETURNING id, created_at"
	err := tx.QueryRow(query, rep.Name).Scan(&rep.ID, &rep.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) CreateRepresentative(rep *models.Representative) error {
	query := "INSERT INTO representatives(name) values($1) RETURNING id, created_at"
	err := r.Db.QueryRow(query, rep.Name).Scan(&rep.ID, &rep.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}
func (r *Repository) GetAllRepresentatives() ([]models.Representative, error) {
	var representatives []models.Representative
	query := "SELECT id,name,created_at,updated_at FROM representatives"
	rows, err := r.Db.Query(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var representative models.Representative
		rows.Scan(&representative.ID,
			&representative.Name,
			&representative.CreatedAt,
			&representative.UpdatedAt)
		representatives = append(representatives, representative)
	}
	return representatives, nil
}
func (r *Repository) GetRepresentative(id int) (*models.Representative, error) {
	var representative models.Representative
	query := "SELECT id,name,created_at,updated_at FROM representatives where id=$1"
	err := r.Db.QueryRow(query, id).Scan(&representative.ID, &representative.Name, &representative.CreatedAt, &representative.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &representative, nil
}
func (r *Repository) UpdateRepresentative(id int, representative *models.Representative) error {
	query := "UPDATE representatives SET name=$1 ,updated_at=NOW() WHERE id=$2"
	_, err := r.Db.Exec(query, representative.Name, id)
	if err != nil {
		return err
	}
	return nil
}
