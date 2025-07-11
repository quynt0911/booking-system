package repository

import (
	"database/sql"
	"expert-service/internal/model"
	"time"

	"github.com/google/uuid"
)

type ExpertScheduleRepository interface {
	Create(schedule *model.ExpertSchedule) error
	GetByExpertID(expertID uuid.UUID) ([]*model.ExpertSchedule, error)
	Update(schedule *model.ExpertSchedule) error
	Delete(id uuid.UUID) error
}

type expertScheduleRepository struct {
	db *sql.DB
}

func NewExpertScheduleRepository(db *sql.DB) ExpertScheduleRepository {
	return &expertScheduleRepository{db: db}
}

func (r *expertScheduleRepository) Create(schedule *model.ExpertSchedule) error {
	query := `INSERT INTO expert_schedules (id, expert_id, day_of_week, start_time, end_time, is_active, created_at, updated_at)
              VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	schedule.ID = uuid.New()
	now := time.Now()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now
	_, err := r.db.Exec(query, schedule.ID, schedule.ExpertID, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime, schedule.IsActive, schedule.CreatedAt, schedule.UpdatedAt)
	return err
}

func (r *expertScheduleRepository) GetByExpertID(expertID uuid.UUID) ([]*model.ExpertSchedule, error) {
	query := `SELECT id, expert_id, day_of_week, start_time, end_time, is_active, created_at, updated_at FROM expert_schedules WHERE expert_id = $1 AND is_active = true`
	rows, err := r.db.Query(query, expertID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var schedules []*model.ExpertSchedule
	for rows.Next() {
		s := &model.ExpertSchedule{}
		err := rows.Scan(&s.ID, &s.ExpertID, &s.DayOfWeek, &s.StartTime, &s.EndTime, &s.IsActive, &s.CreatedAt, &s.UpdatedAt)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, s)
	}
	return schedules, nil
}

func (r *expertScheduleRepository) Update(schedule *model.ExpertSchedule) error {
	query := `UPDATE expert_schedules SET day_of_week=$1, start_time=$2, end_time=$3, is_active=$4, updated_at=$5 WHERE id=$6`
	schedule.UpdatedAt = time.Now()
	_, err := r.db.Exec(query, schedule.DayOfWeek, schedule.StartTime, schedule.EndTime, schedule.IsActive, schedule.UpdatedAt, schedule.ID)
	return err
}

func (r *expertScheduleRepository) Delete(id uuid.UUID) error {
	query := `UPDATE expert_schedules SET is_active=false, updated_at=$1 WHERE id=$2`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}
