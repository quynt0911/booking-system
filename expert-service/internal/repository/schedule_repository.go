package repository

import (
	"database/sql"
	"expert-service/internal/model"
	"github.com/google/uuid"
)

type ScheduleRepository interface {
	Create(schedule *model.Schedule) error
	GetByExpertID(expertID int) ([]*model.Schedule, error)
	GetByExpertIDAndDay(expertID, dayOfWeek int) ([]*model.Schedule, error)
	Update(schedule *model.Schedule) error
	Delete(id uuid.UUID) error
}

type scheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(schedule *model.Schedule) error {
	query := `
		INSERT INTO schedules (expert_id, day_of_week, start_time, end_time, is_active)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`
	return r.db.QueryRow(query,
		schedule.ExpertID, schedule.DayOfWeek,
		schedule.StartTime, schedule.EndTime,
		schedule.IsActive).Scan(&schedule.ID, &schedule.CreatedAt)
}

func (r *scheduleRepository) GetByExpertID(expertID int) ([]*model.Schedule, error) {
	query := `
		SELECT id, expert_id, day_of_week, start_time, end_time, is_active, created_at
		FROM schedules WHERE expert_id = $1 AND is_active = true
		ORDER BY day_of_week, start_time`
	rows, err := r.db.Query(query, expertID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.Schedule
	for rows.Next() {
		schedule := &model.Schedule{}
		err := rows.Scan(
			&schedule.ID, &schedule.ExpertID, &schedule.DayOfWeek,
			&schedule.StartTime, &schedule.EndTime,
			&schedule.IsActive, &schedule.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func (r *scheduleRepository) GetByExpertIDAndDay(expertID, dayOfWeek int) ([]*model.Schedule, error) {
	query := `
		SELECT id, expert_id, day_of_week, start_time, end_time, is_active, created_at
		FROM schedules 
		WHERE expert_id = $1 AND day_of_week = $2 AND is_active = true
		ORDER BY start_time`
	rows, err := r.db.Query(query, expertID, dayOfWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.Schedule
	for rows.Next() {
		schedule := &model.Schedule{}
		err := rows.Scan(
			&schedule.ID, &schedule.ExpertID, &schedule.DayOfWeek,
			&schedule.StartTime, &schedule.EndTime,
			&schedule.IsActive, &schedule.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func (r *scheduleRepository) Update(schedule *model.Schedule) error {
	query := `
		UPDATE schedules 
		SET day_of_week = $1, start_time = $2, end_time = $3, is_active = $4
		WHERE id = $5`
	_, err := r.db.Exec(query,
		schedule.DayOfWeek, schedule.StartTime,
		schedule.EndTime, schedule.IsActive, schedule.ID)
	return err
}

func (r *scheduleRepository) Delete(id uuid.UUID) error {
	query := `UPDATE schedules SET is_active = false WHERE id = $1`
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
