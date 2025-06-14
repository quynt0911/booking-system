package repository

import (
	"database/sql"
	"expert-service/internal/model"

	"github.com/google/uuid"
)

type ExpertRepository interface {
	Create(expert *model.Expert) error
	GetByID(id uuid.UUID) (*model.Expert, error)
	GetByEmail(email string) (*model.Expert, error)
	GetAll() ([]*model.Expert, error)
	Update(expert *model.Expert) error
	Delete(id uuid.UUID) error
	GetByExpertise(expertise string) ([]*model.Expert, error)
}

type expertRepository struct {
	db *sql.DB
}

func NewExpertRepository(db *sql.DB) ExpertRepository {
	return &expertRepository{db: db}
}

func (r *expertRepository) Create(expert *model.Expert) error {
	query := `
        INSERT INTO experts (id, user_id, specialization, experience_years, hourly_rate, certifications, is_available, rating, total_reviews, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
        RETURNING id, created_at, updated_at`

	return r.db.QueryRow(query,
		expert.ID, expert.UserID, expert.Specialization,
		expert.ExperienceYears, expert.HourlyRate, expert.Certifications,
		expert.IsAvailable, expert.Rating, expert.TotalReviews,
		expert.CreatedAt, expert.UpdatedAt).
		Scan(&expert.ID, &expert.CreatedAt, &expert.UpdatedAt)
}

func (r *expertRepository) GetByID(id uuid.UUID) (*model.Expert, error) {
	expert := &model.Expert{}
	query := `
        SELECT id, user_id, specialization, experience_years, hourly_rate, certifications, is_available, rating, total_reviews, created_at, updated_at
        FROM experts WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&expert.ID, &expert.UserID, &expert.Specialization,
		&expert.ExperienceYears, &expert.HourlyRate, &expert.Certifications,
		&expert.IsAvailable, &expert.Rating, &expert.TotalReviews,
		&expert.CreatedAt, &expert.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return expert, err
}

func (r *expertRepository) GetByEmail(email string) (*model.Expert, error) {
	expert := &model.Expert{}
	query := `
        SELECT e.id, e.user_id, e.specialization, e.experience_years, e.hourly_rate, e.certifications, e.is_available, e.rating, e.total_reviews, e.created_at, e.updated_at
        FROM experts e
        JOIN users u ON e.user_id = u.id
        WHERE u.email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&expert.ID, &expert.UserID, &expert.Specialization,
		&expert.ExperienceYears, &expert.HourlyRate, &expert.Certifications,
		&expert.IsAvailable, &expert.Rating, &expert.TotalReviews,
		&expert.CreatedAt, &expert.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	return expert, err
}

func (r *expertRepository) GetAll() ([]*model.Expert, error) {
	query := `
        SELECT id, user_id, specialization, experience_years, hourly_rate, certifications, is_available, rating, total_reviews, created_at, updated_at
        FROM experts ORDER BY created_at DESC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var experts []*model.Expert
	for rows.Next() {
		expert := &model.Expert{}
		err := rows.Scan(
			&expert.ID, &expert.UserID, &expert.Specialization,
			&expert.ExperienceYears, &expert.HourlyRate, &expert.Certifications,
			&expert.IsAvailable, &expert.Rating, &expert.TotalReviews,
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
        SET specialization = $1, experience_years = $2, hourly_rate = $3, 
            certifications = $4, is_available = $5, updated_at = CURRENT_TIMESTAMP
        WHERE id = $6`

	_, err := r.db.Exec(query,
		expert.Specialization, expert.ExperienceYears, expert.HourlyRate,
		expert.Certifications, expert.IsAvailable, expert.ID)
	return err
}

func (r *expertRepository) Delete(id uuid.UUID) error {
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

func (r *expertRepository) GetByExpertise(expertise string) ([]*model.Expert, error) {
	query := `
        SELECT id, user_id, specialization, experience_years, hourly_rate, certifications, is_available, rating, total_reviews, created_at, updated_at
        FROM experts WHERE specialization = $1 AND is_available = true
        ORDER BY created_at DESC`

	rows, err := r.db.Query(query, expertise)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var experts []*model.Expert
	for rows.Next() {
		expert := &model.Expert{}
		err := rows.Scan(
			&expert.ID, &expert.UserID, &expert.Specialization,
			&expert.ExperienceYears, &expert.HourlyRate, &expert.Certifications,
			&expert.IsAvailable, &expert.Rating, &expert.TotalReviews,
			&expert.CreatedAt, &expert.UpdatedAt)
		if err != nil {
			return nil, err
		}
		experts = append(experts, expert)
	}
	return experts, nil
}
