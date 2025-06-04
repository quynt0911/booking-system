package model

type CreateExpertRequest struct {
    Name           string `json:"name" binding:"required"`
    Email          string `json:"email" binding:"required,email"`
    Specialization string `json:"specialization" binding:"required"`
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