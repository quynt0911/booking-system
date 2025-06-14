package model

type CreateExpertRequest struct {
	UserID          string   `json:"user_id" binding:"required"`
	Specialization  string   `json:"specialization" binding:"required"`
	ExperienceYears int      `json:"experience_years" binding:"required"`
	HourlyRate      float64  `json:"hourly_rate" binding:"required"`
	Certifications  []string `json:"certifications"`
	IsAvailable     bool     `json:"is_available"`
}

type CreateScheduleRequest struct {
	ExpertID       string `json:"expert_id" binding:"required"`
	UserID         string `json:"user_id" binding:"required"`
	AvailabilityID string `json:"availability_id" binding:"required"`
	DayOfWeek      int    `json:"day_of_week" binding:"required,min=0,max=6"`
	Date           string `json:"date" binding:"required"` // YYYY-MM-DD
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time" binding:"required"`
	Title          string `json:"title" binding:"required"`
	Description    string `json:"description,omitempty"`
	MeetingLink    string `json:"meeting_link,omitempty"`
	Notes          string `json:"notes,omitempty"`
}

type CreateOffTimeRequest struct {
	ExpertID      string `json:"expert_id" binding:"required"`
	StartDateTime string `json:"start_datetime" binding:"required"`
	EndDateTime   string `json:"end_datetime" binding:"required"`
	Reason        string `json:"reason"`
	IsRecurring   bool   `json:"is_recurring"`
}

type CheckAvailabilityRequest struct {
	ExpertID string `json:"expert_id" binding:"required"`
	Date     string `json:"date" binding:"required"`
	Time     string `json:"time" binding:"required"`
}

type UpdateExpertRequest struct {
	Specialization  *string  `json:"specialization,omitempty"`
	ExperienceYears *int     `json:"experience_years,omitempty"`
	HourlyRate      *float64 `json:"hourly_rate,omitempty"`
	Certifications  []string `json:"certifications,omitempty"`
	IsAvailable     *bool    `json:"is_available,omitempty"`
}

type CreateRecurringAvailabilityRequest struct {
	ExpertID   string `json:"expert_id" binding:"required"`
	DaysOfWeek []int  `json:"days_of_week" binding:"required"` // 0=Sunday, 6=Saturday
	StartTime  string `json:"start_time" binding:"required"`
	EndTime    string `json:"end_time" binding:"required"`
	StartDate  string `json:"start_date" binding:"required"` // YYYY-MM-DD
	EndDate    string `json:"end_date" binding:"required"`   // YYYY-MM-DD
}

type GetSchedulesRequest struct {
	ExpertID string `json:"expert_id,omitempty"`
	UserID   string `json:"user_id,omitempty"`
	Status   string `json:"status,omitempty"`
	Date     string `json:"date,omitempty"` // YYYY-MM-DD
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
}

type UpdateScheduleRequest struct {
	Status         *string `json:"status,omitempty"`
	Title          *string `json:"title,omitempty"`
	Description    *string `json:"description,omitempty"`
	MeetingLink    *string `json:"meeting_link,omitempty"`
	Notes          *string `json:"notes,omitempty"`
	ExpertID       *string `json:"expert_id,omitempty"`
	UserID         *string `json:"user_id,omitempty"`
	AvailabilityID *string `json:"availability_id,omitempty"`
	DayOfWeek      *int    `json:"day_of_week,omitempty"`
	Date           *string `json:"date,omitempty"`
	StartTime      *string `json:"start_time,omitempty"`
	EndTime        *string `json:"end_time,omitempty"`
	IsActive       *bool   `json:"is_active,omitempty"`
}

type CreateAvailabilityRequest struct {
	ExpertID  string `json:"expert_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type UpdateAvailabilityRequest struct {
	Date      *string `json:"date,omitempty"`
	StartTime *string `json:"start_time,omitempty"`
	EndTime   *string `json:"end_time,omitempty"`
	IsBooked  *bool   `json:"is_booked,omitempty"`
}
