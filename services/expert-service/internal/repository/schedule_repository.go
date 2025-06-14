package repository

import (
	"database/sql"
	"expert-service/internal/model"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

type ScheduleRepository interface {
	Create(schedule *model.Schedule) error
	GetByExpertID(expertID string) ([]*model.Schedule, error)
	GetByExpertIDAndDay(expertID string, dayOfWeek int) ([]*model.Schedule, error)
	Update(schedule *model.Schedule) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*model.Schedule, error)
	GetSchedules(req *model.GetSchedulesRequest) ([]*model.Schedule, error)
}

type scheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) ScheduleRepository {
	return &scheduleRepository{db: db}
}

func (r *scheduleRepository) Create(schedule *model.Schedule) error {
	query := `
		INSERT INTO schedules (id, expert_id, day_of_week, start_time, end_time, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at`

	schedule.ID = uuid.New()
	schedule.CreatedAt = time.Now()
	schedule.UpdatedAt = time.Now()

	return r.db.QueryRow(query,
		schedule.ID, schedule.ExpertID, schedule.DayOfWeek,
		schedule.StartTime, schedule.EndTime, schedule.IsActive,
		schedule.CreatedAt, schedule.UpdatedAt).
		Scan(&schedule.ID, &schedule.CreatedAt)
}

func (r *scheduleRepository) GetByExpertID(expertID string) ([]*model.Schedule, error) {
	query := `
		SELECT id, expert_id, day_of_week, start_time, end_time, is_active, created_at, updated_at
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
			&schedule.IsActive, &schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}

func (r *scheduleRepository) GetByExpertIDAndDay(expertID string, dayOfWeek int) ([]*model.Schedule, error) {
	query := `
		SELECT id, expert_id, day_of_week, start_time, end_time, is_active, created_at, updated_at
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
			&schedule.IsActive, &schedule.CreatedAt, &schedule.UpdatedAt,
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
		SET day_of_week = $1, start_time = $2, end_time = $3, 
			is_active = $4, updated_at = $5
		WHERE id = $6`

	schedule.UpdatedAt = time.Now()
	_, err := r.db.Exec(query,
		schedule.DayOfWeek, schedule.StartTime,
		schedule.EndTime, schedule.IsActive,
		schedule.UpdatedAt, schedule.ID)
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

func (r *scheduleRepository) GetByID(id uuid.UUID) (*model.Schedule, error) {
	query := `
		SELECT id, expert_id, user_id, availability_id, day_of_week, date, start_time, end_time, status, title, description, meeting_link, notes, is_active, created_at, updated_at
		FROM schedules WHERE id = $1`
	schedule := &model.Schedule{}
	err := r.db.QueryRow(query, id).Scan(
		&schedule.ID, &schedule.ExpertID, &schedule.UserID, &schedule.AvailabilityID,
		&schedule.DayOfWeek, &schedule.Date, &schedule.StartTime, &schedule.EndTime,
		&schedule.Status, &schedule.Title, &schedule.Description, &schedule.MeetingLink, &schedule.Notes,
		&schedule.IsActive, &schedule.CreatedAt, &schedule.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return schedule, err
}

func (r *scheduleRepository) GetSchedules(req *model.GetSchedulesRequest) ([]*model.Schedule, error) {
	baseQuery := `
		SELECT id, expert_id, user_id, availability_id, day_of_week, date, start_time, end_time, status, title, description, meeting_link, notes, is_active, created_at, updated_at
		FROM schedules
		WHERE is_active = true
		`

	var args []interface{}
	argCount := 1

	whereClauses := []string{}

	if req.ExpertID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("expert_id = $%d", argCount))
		args = append(args, req.ExpertID)
		argCount++
	}
	if req.UserID != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, req.UserID)
		argCount++
	}
	if req.Status != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("status = $%d", argCount))
		args = append(args, req.Status)
		argCount++
	}
	if req.Date != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("date = $%d", argCount))
		args = append(args, req.Date)
		argCount++
	}

	if len(whereClauses) > 0 {
		baseQuery += " AND " + strings.Join(whereClauses, " AND ")
	}

	baseQuery += " ORDER BY date, start_time "

	if req.Limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argCount)
		args = append(args, req.Limit)
		argCount++
	}
	if req.Offset > 0 {
		baseQuery += fmt.Sprintf(" OFFSET $%d", argCount)
		args = append(args, req.Offset)
		argCount++
	}

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var schedules []*model.Schedule
	for rows.Next() {
		schedule := &model.Schedule{}
		err := rows.Scan(
			&schedule.ID, &schedule.ExpertID, &schedule.UserID, &schedule.AvailabilityID,
			&schedule.DayOfWeek, &schedule.Date, &schedule.StartTime, &schedule.EndTime,
			&schedule.Status, &schedule.Title, &schedule.Description, &schedule.MeetingLink, &schedule.Notes,
			&schedule.IsActive, &schedule.CreatedAt, &schedule.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		schedules = append(schedules, schedule)
	}
	return schedules, nil
}
