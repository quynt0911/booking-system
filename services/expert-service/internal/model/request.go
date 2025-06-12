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
	ExpertID  int    `json:"expert_id" binding:"required"`
	DayOfWeek int    `json:"day_of_week" binding:"required,min=0,max=6"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type CreateOffTimeRequest struct {
	ExpertID  int    `json:"expert_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Reason    string `json:"reason"`
}

type CheckAvailabilityRequest struct {
	ExpertID int    `json:"expert_id" binding:"required"`
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
