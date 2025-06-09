package repository

import (
	"database/sql"
	"expert-service/internal/model"
	"time"
)

type OffTimeRepository interface {
	Create(offTime *model.OffTime) error
	GetByExpertID(expertID int) ([]*model.OffTime, error)
	GetByExpertIDAndDateRange(expertID int, date time.Time) ([]*model.OffTime, error)
	Delete(id int) error
}

type offTimeRepository struct {
	db *sql.DB
}

func NewOffTimeRepository(db *sql.DB) OffTimeRepository {
	return &offTimeRepository{db: db}
}

func (r *offTimeRepository) Create(offTime *model.OffTime) error {
	query := `
		INSERT INTO off_times (expert_id, start_date, end_date, reason)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	return r.db.QueryRow(query,
		offTime.ExpertID, offTime.StartDate, offTime.EndDate, offTime.Reason).
		Scan(&offTime.ID, &offTime.CreatedAt)
}

func (r *offTimeRepository) GetByExpertID(expertID int) ([]*model.OffTime, error) {
	query := `
		SELECT id, expert_id, start_date, end_date, reason, created_at
		FROM off_times WHERE expert_id = $1
		ORDER BY start_date DESC`
	rows, err := r.db.Query(query, expertID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offTimes []*model.OffTime
	for rows.Next() {
		offTime := &model.OffTime{}
		err := rows.Scan(
			&offTime.ID, &offTime.ExpertID,
			&offTime.StartDate, &offTime.EndDate,
			&offTime.Reason, &offTime.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		offTimes = append(offTimes, offTime)
	}
	return offTimes, nil
}

func (r *offTimeRepository) GetByExpertIDAndDateRange(expertID int, date time.Time) ([]*model.OffTime, error) {
	query := `
		SELECT id, expert_id, start_date, end_date, reason, created_at
		FROM off_times 
		WHERE expert_id = $1 AND start_date <= $2 AND end_date >= $2`
	rows, err := r.db.Query(query, expertID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var offTimes []*model.OffTime
	for rows.Next() {
		offTime := &model.OffTime{}
		err := rows.Scan(
			&offTime.ID, &offTime.ExpertID,
			&offTime.StartDate, &offTime.EndDate,
			&offTime.Reason, &offTime.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		offTimes = append(offTimes, offTime)
	}
	return offTimes, nil
}

func (r *offTimeRepository) Delete(id int) error {
	query := `DELETE FROM off_times WHERE id = $1`
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
