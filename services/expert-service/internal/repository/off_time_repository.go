package repository

import (
	"database/sql"
	"expert-service/internal/model"
	"time"

	"github.com/google/uuid"
)

type OffTimeRepository interface {
	Create(offTime *model.OffTime) error
	GetByExpertID(expertID uuid.UUID) ([]*model.OffTime, error)
	GetByExpertIDAndDateRange(expertID uuid.UUID, date time.Time) ([]*model.OffTime, error)
	Delete(id uuid.UUID) error
}

type offTimeRepository struct {
	db *sql.DB
}

func NewOffTimeRepository(db *sql.DB) OffTimeRepository {
	return &offTimeRepository{db: db}
}

func (r *offTimeRepository) Create(offTime *model.OffTime) error {
	query := `
		INSERT INTO expert_off_times (id, expert_id, start_datetime, end_datetime, reason, is_recurring, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`

	offTime.ID = uuid.New()
	offTime.CreatedAt = time.Now()

	return r.db.QueryRow(query,
		offTime.ID, offTime.ExpertID, offTime.StartDateTime, offTime.EndDateTime,
		offTime.Reason, offTime.IsRecurring, offTime.CreatedAt).
		Scan(&offTime.ID, &offTime.CreatedAt)
}

func (r *offTimeRepository) GetByExpertID(expertID uuid.UUID) ([]*model.OffTime, error) {
	query := `
		SELECT id, expert_id, start_datetime, end_datetime, reason, is_recurring, created_at
		FROM expert_off_times WHERE expert_id = $1
		ORDER BY start_datetime DESC`
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
			&offTime.StartDateTime, &offTime.EndDateTime,
			&offTime.Reason, &offTime.IsRecurring,
			&offTime.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		offTimes = append(offTimes, offTime)
	}
	return offTimes, nil
}

func (r *offTimeRepository) GetByExpertIDAndDateRange(expertID uuid.UUID, date time.Time) ([]*model.OffTime, error) {
	query := `
		SELECT id, expert_id, start_datetime, end_datetime, reason, is_recurring, created_at
		FROM expert_off_times 
		WHERE expert_id = $1 AND start_datetime <= $2 AND end_datetime >= $2`
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
			&offTime.StartDateTime, &offTime.EndDateTime,
			&offTime.Reason, &offTime.IsRecurring,
			&offTime.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		offTimes = append(offTimes, offTime)
	}
	return offTimes, nil
}

func (r *offTimeRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM expert_off_times WHERE id = $1`
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
