package repository

import (
	"database/sql"
	"expert-service/internal2/model"
)

type ExpertRepository interface {
	Create(expert *model.Expert) error
	GetByID(id int) (*model.Expert, error)
	GetAll() ([]*model.Expert, error)
	Update(expert *model.Expert) error
	Delete(id int) error
}

type expertRepository struct {
	db *sql.DB
}

func NewExpertRepository(db *sql.DB) ExpertRepository {
	return &expertRepository{db: db}
}

func (r *expertRepository) Create(expert *model.Expert) error {
	query := `
        INSERT INTO experts (name, email, specialization, status)
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query, expert.Name, expert.Email,
		expert.Specialization, expert.Status).
		Scan(&expert.ID, &expert.CreatedAt, &expert.UpdatedAt)
}

func (r *expertRepository) GetByID(id int) (*model.Expert, error) {
	expert := &model.Expert{}
	query := `
        SELECT id, name, email, specialization, status, created_at, updated_at
        FROM experts WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&expert.ID, &expert.Name, &expert.Email,
		&expert.Specialization, &expert.Status,
		&expert.CreatedAt, &expert.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return expert, err
}

func (r *expertRepository) GetAll() ([]*model.Expert, error) {
	query := `
        SELECT id, name, email, specialization, status, created_at, updated_at
        FROM experts ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var experts []*model.Expert
	for rows.Next() {
		expert := &model.Expert{}
		err := rows.Scan(&expert.ID, &expert.Name, &expert.Email,
			&expert.Specialization, &expert.Status,
			&expert.CreatedAt, &expert.UpdatedAt)
		if err != nil {
			return nil, err
		}
		experts = append(experts, expert)
	}
	return experts, nil
}

func (r *expertRepository) Update(expert *model.Expert) error {
	query := `
        UPDATE experts 
        SET name = $1, email = $2, specialization = $3, status = $4, updated_at = CURRENT_TIMESTAMP
        WHERE id = $5`

	_, err := r.db.Exec(query, expert.Name, expert.Email,
		expert.Specialization, expert.Status, expert.ID)
	return err
}

func (r *expertRepository) Delete(id int) error {
	query := `DELETE FROM experts WHERE id = $1`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return sql.ErrNoRows
	}
	return nil

}
