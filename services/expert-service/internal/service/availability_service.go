package service

import (
	"encoding/json"
	"expert-service/internal/cache"
	"expert-service/internal/model"
	"expert-service/internal/repository"
	"fmt"
	"time"

	"expert-service/internal/client"
	"expert-service/internal/utils"

	"github.com/google/uuid"
)

// isTimeInRange checks if a time string falls within a start and end time range
func (s *expertAvailabilityService) isTimeInRange(timeStr, startTimeStr, endTimeStr string) bool {
	t, _ := time.ParseInLocation("15:04", timeStr, time.Local)
	start, _ := time.ParseInLocation("15:04", startTimeStr, time.Local)
	end, _ := time.ParseInLocation("15:04", endTimeStr, time.Local)
	return (t.Equal(start) || t.After(start)) && t.Before(end)
}

type ExpertAvailabilityService interface {
	CheckAvailability(req *model.CheckAvailabilityRequest) (bool, error)
	CreateOffTime(req *model.CreateOffTimeRequest) (*model.OffTime, error)
	GetExpertOffTimes(expertID string) ([]*model.OffTime, error)
	DeleteOffTime(id string) error
	CreateAvailability(req *model.CreateAvailabilityRequest) (*model.Availability, error)
	GetAvailabilityByID(id string) (*model.Availability, error)
	UpdateAvailability(id string, req *model.UpdateAvailabilityRequest) (*model.Availability, error)
	DeleteAvailability(id string) error
	GetAvailabilities(expertID string, startDate, endDate time.Time, isBooked *bool, token string) ([]*model.AvailabilitySlot, error)
	BookAvailability(id string) error
	CreateRecurringAvailability(req *model.CreateRecurringAvailabilityRequest) ([]*model.Availability, error)
	GetOffTimeByID(id string) (*model.OffTime, error)
	IsExpertProfileExists(userID uuid.UUID) (bool, error)
}

type expertAvailabilityService struct {
	expertRepo         repository.ExpertRepository
	scheduleRepo       repository.ScheduleRepository
	offTimeRepo        repository.OffTimeRepository
	expertScheduleRepo repository.ExpertScheduleRepository
	cache              cache.AvailabilityCache
}

func NewExpertAvailabilityService(
	expertRepo repository.ExpertRepository,
	scheduleRepo repository.ScheduleRepository,
	offTimeRepo repository.OffTimeRepository,
	expertScheduleRepo repository.ExpertScheduleRepository,
	cache cache.AvailabilityCache,
) ExpertAvailabilityService {
	return &expertAvailabilityService{
		expertRepo:         expertRepo,
		scheduleRepo:       scheduleRepo,
		offTimeRepo:        offTimeRepo,
		expertScheduleRepo: expertScheduleRepo,
		cache:              cache,
	}
}

// CheckAvailability kiểm tra chuyên gia có rảnh không tại thời điểm yêu cầu
func (s *expertAvailabilityService) CheckAvailability(req *model.CheckAvailabilityRequest) (bool, error) {
	expertUUID, err := uuid.Parse(req.ExpertID)
	if err != nil {
		return false, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return false, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return false, fmt.Errorf("không tìm thấy chuyên gia với ID %s", req.ExpertID)
	}

	// Parse date and time
	date, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return false, fmt.Errorf("định dạng ngày không hợp lệ")
	}

	// Check cache first
	cacheKey := fmt.Sprintf("%s:%s", req.ExpertID, req.Date)
	if cached, err := s.cache.GetAvailability(cacheKey); err == nil && cached != nil {
		var isAvailable bool
		if err := json.Unmarshal(cached, &isAvailable); err == nil {
			return isAvailable, nil
		}
	}

	// Check if expert is on off-time
	var offTimes []*model.OffTime
	for d := date; !d.After(date); d = d.AddDate(0, 0, 1) {
		dayOffTimes, err := s.offTimeRepo.GetByExpertIDAndDateRange(expertUUID, d)
		if err != nil {
			return false, fmt.Errorf("không thể kiểm tra thời gian nghỉ: %v", err)
		}
		offTimes = append(offTimes, dayOffTimes...)
	}
	if len(offTimes) > 0 {
		data, _ := json.Marshal(false)
		if err := s.cache.SetAvailability(cacheKey, data); err != nil {
			return false, fmt.Errorf("failed to cache availability: %w", err)
		}
		return false, nil
	}

	// Get expert's schedule for the day
	dayOfWeek := int(date.Weekday())
	schedules, err := s.scheduleRepo.GetByExpertIDAndDay(req.ExpertID, dayOfWeek)
	if err != nil {
		return false, fmt.Errorf("không thể lấy lịch làm việc: %v", err)
	}

	// Check if the requested time falls within any schedule
	isAvailable := false
	for _, schedule := range schedules {
		if s.isTimeInRange(req.Time, schedule.StartTime, schedule.EndTime) {
			isAvailable = true
			break
		}
	}

	// Cache the result
	data, _ := json.Marshal(isAvailable)
	if err := s.cache.SetAvailability(cacheKey, data); err != nil {
		return false, fmt.Errorf("failed to cache availability: %w", err)
	}
	return isAvailable, nil
}

// CreateOffTime tạo thời gian nghỉ cho chuyên gia
func (s *expertAvailabilityService) CreateOffTime(req *model.CreateOffTimeRequest) (*model.OffTime, error) {
	expertUUID, err := uuid.Parse(req.ExpertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %s", req.ExpertID)
	}

	startDateTime, err := time.Parse("2006-01-02T15:04:05Z", req.StartDateTime)
	if err != nil {
		return nil, fmt.Errorf("định dạng thời gian bắt đầu không hợp lệ")
	}

	endDateTime, err := time.Parse("2006-01-02T15:04:05Z", req.EndDateTime)
	if err != nil {
		return nil, fmt.Errorf("định dạng thời gian kết thúc không hợp lệ")
	}

	if endDateTime.Before(startDateTime) {
		return nil, fmt.Errorf("thời gian kết thúc phải sau thời gian bắt đầu")
	}

	// Kiểm tra trùng off-time
	existingOffTimes, err := s.offTimeRepo.GetByExpertID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra off-time hiện có: %v", err)
	}
	for _, off := range existingOffTimes {
		if startDateTime.Before(off.EndDateTime) && endDateTime.After(off.StartDateTime) {
			return nil, fmt.Errorf("Đã tồn tại off-time giao với khoảng thời gian này")
		}
	}

	offTime := &model.OffTime{
		ExpertID:      expertUUID,
		StartDateTime: startDateTime,
		EndDateTime:   endDateTime,
		Reason:        req.Reason,
		IsRecurring:   req.IsRecurring,
	}

	err = s.offTimeRepo.Create(offTime)
	if err != nil {
		return nil, fmt.Errorf("không thể tạo thời gian nghỉ: %v", err)
	}

	// Invalidate cache
	s.cache.InvalidateExpert(req.ExpertID)
	return offTime, nil
}

// GetExpertOffTimes lấy danh sách thời gian nghỉ của chuyên gia
func (s *expertAvailabilityService) GetExpertOffTimes(expertID string) ([]*model.OffTime, error) {
	expertUUID, err := uuid.Parse(expertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %s", expertID)
	}

	offTimes, err := s.offTimeRepo.GetByExpertID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy danh sách thời gian nghỉ: %v", err)
	}
	return offTimes, nil
}

// DeleteOffTime xóa thời gian nghỉ
func (s *expertAvailabilityService) DeleteOffTime(id string) error {
	offTimeID, err := uuid.Parse(id)
	if err != nil {
		return fmt.Errorf("invalid off time ID format: %v", err)
	}

	err = s.offTimeRepo.Delete(offTimeID)
	if err != nil {
		return fmt.Errorf("không thể xóa thời gian nghỉ: %v", err)
	}
	return nil
}

// CreateAvailability creates a new availability slot
func (s *expertAvailabilityService) CreateAvailability(req *model.CreateAvailabilityRequest) (*model.Availability, error) {
	expertUUID, err := uuid.Parse(req.ExpertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %s", req.ExpertID)
	}

	availability := &model.Availability{
		ID:        uuid.New(),
		ExpertID:  req.ExpertID,
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		IsBooked:  false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Lưu vào Redis với key pattern: availability:{expert_id}:{date}
	key := fmt.Sprintf("availability:%s:%s", req.ExpertID, req.Date)
	data, err := json.Marshal(availability)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal availability: %w", err)
	}

	if err := s.cache.SetAvailability(key, data); err != nil {
		return nil, fmt.Errorf("failed to cache availability: %w", err)
	}

	return availability, nil
}

// GetAvailabilityByID retrieves an availability slot by ID
func (s *expertAvailabilityService) GetAvailabilityByID(id string) (*model.Availability, error) {
	// Tìm trong Redis với pattern: availability:*:*
	// TODO: Implement search by ID in Redis
	// For now, return error as this might need a different Redis structure
	return nil, fmt.Errorf("get availability by ID not implemented with Redis")
}

// UpdateAvailability updates an existing availability slot
func (s *expertAvailabilityService) UpdateAvailability(id string, req *model.UpdateAvailabilityRequest) (*model.Availability, error) {
	// TODO: Implement update in Redis
	// This might need to fetch all availabilities for the expert and date,
	// update the specific one, and save back
	return nil, fmt.Errorf("update availability not implemented with Redis")
}

// DeleteAvailability deletes an availability slot
func (s *expertAvailabilityService) DeleteAvailability(id string) error {
	// TODO: Implement delete in Redis
	return fmt.Errorf("delete availability not implemented with Redis")
}

// Hàm cập nhật để lấy off-times với logic cải tiến
func (s *expertAvailabilityService) GetAvailabilities(expertID string, startDate, endDate time.Time, isBooked *bool, token string) ([]*model.AvailabilitySlot, error) {
	expertUUID, err := uuid.Parse(expertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}

	// 1. Lấy schedules định kỳ
	schedules, err := s.expertScheduleRepo.GetByExpertID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy schedules: %v", err)
	}

	// 2. Lấy tất cả off-times của expert (cả recurring và non-recurring)
	allOffTimes, err := s.offTimeRepo.GetByExpertID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy off-times: %v", err)
	}

	// Lọc off-times có liên quan đến khoảng thời gian yêu cầu
	var relevantOffTimes []*model.OffTime
	for _, offTime := range allOffTimes {
		if isOffTimeRelevant(offTime, startDate, endDate) {
			relevantOffTimes = append(relevantOffTimes, offTime)
		}
	}

	// 3. Tạo service token cho internal call
	serviceToken, err := utils.GenerateServiceToken("expert-service")
	if err != nil {
		return nil, fmt.Errorf("không thể tạo service token: %v", err)
	}

	// 4. Lấy bookings đã xác nhận từ API booking-service
	bookingServiceURL := "http://booking-service-dev:8082"
	bookings, err := client.GetBookingsByExpertAndDateRange(bookingServiceURL, expertID, serviceToken, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("không thể lấy bookings: %v", err)
	}

	// Debug log
	fmt.Printf("Debug - Retrieved %d bookings from booking service\n", len(bookings))
	for _, b := range bookings {
		fmt.Printf("Debug - Booking: time=%s-%s, status=%s\n",
			b.ScheduledTime.Format("2006-01-02 15:04:05"),
			b.EndTime.Format("2006-01-02 15:04:05"),
			b.Status)
	}

	// 5. Sinh ra các slot rảnh từ schedules
	var expertSchedules []*model.ExpertSchedule
	for _, s := range schedules {
		expertSchedules = append(expertSchedules, s)
	}
	slots := generateSlotsFromSchedules(expertSchedules, startDate, endDate)

	// Debug log
	fmt.Printf("Debug - Generated %d slots before filtering\n", len(slots))

	// 6. Loại bỏ các slot trùng với off-times và bookings
	slots = removeSlotsByOffTimes(slots, relevantOffTimes)
	fmt.Printf("Debug - %d slots after removing off-times\n", len(slots))

	slots = removeSlotsByBookings(slots, bookings)
	fmt.Printf("Debug - %d slots after removing booking\n", len(slots))

	// 7. Trả về kết quả
	var result []*model.AvailabilitySlot
	for _, slot := range slots {
		s := slot
		result = append(result, &s)
	}

	// 8. Cache kết quả
	if len(result) > 0 {
		slotsByDate := make(map[string][]*model.AvailabilitySlot)
		for _, slot := range result {
			slotsByDate[slot.Date] = append(slotsByDate[slot.Date], slot)
		}
		for date, slots := range slotsByDate {
			data, err := json.Marshal(slots)
			if err == nil {
				key := fmt.Sprintf("availability:%s:%s", expertID, date)
				s.cache.SetAvailability(key, data)
			}
		}
	}

	if result == nil {
		result = []*model.AvailabilitySlot{}
	}

	return result, nil
}

func isOffTimeRelevant(offTime *model.OffTime, startDate, endDate time.Time) bool {
	if !offTime.IsRecurring {
		// Off-time không lặp lại: kiểm tra xem có giao với khoảng thời gian yêu cầu không
		offStartDate := offTime.StartDateTime.Truncate(24 * time.Hour)
		offEndDate := offTime.EndDateTime.Truncate(24 * time.Hour)

		return !offEndDate.Before(startDate) && !offStartDate.After(endDate)
	} else {
		// Off-time lặp lại: luôn có thể liên quan nếu có cùng thứ trong tuần
		return true
	}
}

// generateSlotsFromSchedules: sinh slot rảnh từ schedules định kỳ, chia nhỏ mỗi slot 30 phút
func generateSlotsFromSchedules(schedules []*model.ExpertSchedule, startDate, endDate time.Time) []model.AvailabilitySlot {
	var slots []model.AvailabilitySlot
	const slotDuration = 30 // phút
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		weekday := int(d.Weekday())
		fmt.Printf("Debug - Checking date %s, weekday=%d\n", d.Format("2006-01-02"), weekday)
		for _, s := range schedules {
			if s.DayOfWeek == weekday && s.IsActive {
				fmt.Printf("Debug - Found matching schedule for weekday %d\n", weekday)
				// Parse start & end time từ timestamp
				startTime, err1 := time.Parse(time.RFC3339, s.StartTime)
				endTime, err2 := time.Parse(time.RFC3339, s.EndTime)
				if err1 != nil || err2 != nil {
					fmt.Printf("Debug - Error parsing time: %v, %v\n", err1, err2)
					continue
				}

				// Chỉ lấy giờ và phút
				startHour, startMin := startTime.Hour(), startTime.Minute()
				endHour, endMin := endTime.Hour(), endTime.Minute()

				// Tạo thời gian với ngày hiện tại
				startDateTime := time.Date(d.Year(), d.Month(), d.Day(), startHour, startMin, 0, 0, time.Local)
				endDateTime := time.Date(d.Year(), d.Month(), d.Day(), endHour, endMin, 0, 0, time.Local)

				fmt.Printf("Debug - Creating slots from %v to %v\n", startDateTime, endDateTime)

				for t := startDateTime; t.Add(time.Minute*slotDuration).Before(endDateTime) || t.Add(time.Minute*slotDuration).Equal(endDateTime); t = t.Add(time.Minute * slotDuration) {
					slotStart := t.Format("15:04")
					slotEnd := t.Add(time.Minute * slotDuration).Format("15:04")
					slots = append(slots, model.AvailabilitySlot{
						Date:      d.Format("2006-01-02"),
						StartTime: slotStart,
						EndTime:   slotEnd,
					})
					fmt.Printf("Debug - Created slot: %s %s-%s\n", d.Format("2006-01-02"), slotStart, slotEnd)
				}
			}
		}
	}
	return slots
}

// removeSlotsByOffTimes: loại bỏ slot trùng với off-time
func removeSlotsByOffTimes(slots []model.AvailabilitySlot, offTimes []*model.OffTime) []model.AvailabilitySlot {
	var result []model.AvailabilitySlot
	for _, slot := range slots {
		slotDate, _ := time.Parse("2006-01-02", slot.Date)
		slotStart, _ := time.ParseInLocation("2006-01-02 15:04", slot.Date+" "+slot.StartTime, time.Local)
		slotEnd, _ := time.ParseInLocation("2006-01-02 15:04", slot.Date+" "+slot.EndTime, time.Local)

		conflict := false
		for _, off := range offTimes {
			// Kiểm tra xem off-time có áp dụng cho ngày này không
			if shouldApplyOffTime(off, slotDate) {
				// Chuyển đổi off-time thành thời gian của ngày slot để so sánh
				offStartTime := off.StartDateTime.Format("15:04")
				offEndTime := off.EndDateTime.Format("15:04")

				// Tạo thời gian off-time cho ngày của slot
				offStart, _ := time.ParseInLocation("2006-01-02 15:04", slot.Date+" "+offStartTime, time.Local)
				offEnd, _ := time.ParseInLocation("2006-01-02 15:04", slot.Date+" "+offEndTime, time.Local)

				// Kiểm tra xung đột thời gian
				if slotStart.Before(offEnd) && slotEnd.After(offStart) {
					fmt.Printf("REMOVE slot %s %s-%s vì trùng off-time (recurring: %v)\n",
						slot.Date, slot.StartTime, slot.EndTime, off.IsRecurring)
					conflict = true
					break
				}
			}
		}
		if !conflict {
			result = append(result, slot)
		}
	}
	return result
}

// shouldApplyOffTime kiểm tra xem off-time có áp dụng cho ngày cụ thể không
func shouldApplyOffTime(offTime *model.OffTime, targetDate time.Time) bool {
	if !offTime.IsRecurring {
		// Off-time không lặp lại: chỉ áp dụng trong khoảng thời gian cụ thể
		offStartDate := offTime.StartDateTime.Truncate(24 * time.Hour)
		offEndDate := offTime.EndDateTime.Truncate(24 * time.Hour)
		targetDateTruncated := targetDate.Truncate(24 * time.Hour)

		return !targetDateTruncated.Before(offStartDate) && !targetDateTruncated.After(offEndDate)
	} else {
		// Off-time lặp lại: áp dụng cho cùng thứ trong tuần và cùng thời gian
		offWeekday := offTime.StartDateTime.Weekday()
		targetWeekday := targetDate.Weekday()

		return offWeekday == targetWeekday
	}
}

// removeSlotsByBookings: loại bỏ slot đã bị booking
func removeSlotsByBookings(slots []model.AvailabilitySlot, bookings []client.Booking) []model.AvailabilitySlot {
	var result []model.AvailabilitySlot
	for _, slot := range slots {
		booked := false
		slotStart, _ := time.ParseInLocation("2006-01-02 15:04", slot.Date+" "+slot.StartTime, time.Local)
		slotEnd, _ := time.ParseInLocation("2006-01-02 15:04", slot.Date+" "+slot.EndTime, time.Local)

		for _, b := range bookings {
			// Convert booking times to local timezone
			bookingStart := b.ScheduledTime.Local()
			bookingEnd := b.EndTime.Local()

			// Debug logs
			fmt.Printf("Comparing slot %s %s-%s with booking %s-%s\n",
				slot.Date, slot.StartTime, slot.EndTime,
				bookingStart.Format("15:04"), bookingEnd.Format("15:04"))

			// Kiểm tra xem slot có overlap với booking không
			if slotStart.Before(bookingEnd) && slotEnd.After(bookingStart) {
				fmt.Printf("REMOVE slot %s %s-%s vì trùng booking %s-%s\n",
					slot.Date, slot.StartTime, slot.EndTime,
					bookingStart.Format("15:04"), bookingEnd.Format("15:04"))
				booked = true
				break
			}
		}
		if !booked {
			result = append(result, slot)
		}
	}
	return result
}

// BookAvailability books an availability slot
func (s *expertAvailabilityService) BookAvailability(id string) error {
	// TODO: Implement booking in Redis
	// This might need to fetch all availabilities for the expert and date,
	// update the specific one's IsBooked status, and save back
	return fmt.Errorf("book availability not implemented with Redis")
}

// CreateRecurringAvailability creates multiple availability slots for recurring schedules
func (s *expertAvailabilityService) CreateRecurringAvailability(req *model.CreateRecurringAvailabilityRequest) ([]*model.Availability, error) {
	expertUUID, err := uuid.Parse(req.ExpertID)
	if err != nil {
		return nil, fmt.Errorf("invalid expert ID format: %v", err)
	}
	expert, err := s.expertRepo.GetByID(expertUUID)
	if err != nil {
		return nil, fmt.Errorf("không thể kiểm tra chuyên gia: %v", err)
	}
	if expert == nil {
		return nil, fmt.Errorf("không tìm thấy chuyên gia với ID %s", req.ExpertID)
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("định dạng ngày bắt đầu không hợp lệ")
	}

	endDate, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("định dạng ngày kết thúc không hợp lệ")
	}

	if endDate.Before(startDate) {
		return nil, fmt.Errorf("ngày kết thúc phải sau ngày bắt đầu")
	}

	var createdAvailabilities []*model.Availability
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		weekday := int(d.Weekday())
		for _, w := range req.DaysOfWeek {
			if weekday == w {
				availability := &model.Availability{
					ID:        uuid.New(),
					ExpertID:  req.ExpertID,
					Date:      d.Format("2006-01-02"),
					StartTime: req.StartTime,
					EndTime:   req.EndTime,
					IsBooked:  false,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}

				// Lưu vào Redis
				key := fmt.Sprintf("availability:%s:%s", req.ExpertID, d.Format("2006-01-02"))
				data, err := json.Marshal(availability)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal availability: %w", err)
				}

				if err := s.cache.SetAvailability(key, data); err != nil {
					return nil, fmt.Errorf("failed to cache availability: %w", err)
				}

				createdAvailabilities = append(createdAvailabilities, availability)
			}
		}
	}
	return createdAvailabilities, nil
}

// GetOffTimeByID lấy off-time theo id
func (s *expertAvailabilityService) GetOffTimeByID(id string) (*model.OffTime, error) {
	uuidID, err := uuid.Parse(id)
	if err != nil {
		return nil, fmt.Errorf("invalid off time ID format: %v", err)
	}
	offTime, err := s.offTimeRepo.GetByID(uuidID)
	if err != nil {
		return nil, err
	}
	return offTime, nil
}

func (s *expertAvailabilityService) IsExpertProfileExists(userID uuid.UUID) (bool, error) {
	return s.expertRepo.IsExpertProfileExists(userID)
}
