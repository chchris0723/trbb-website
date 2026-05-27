package models

import "time"

// Event status
const (
	EventDraft     = 0
	EventPublished = 1
	EventFull      = 2
	EventEnded     = 3
	EventCancelled = 4
)

// Registration status
const (
	RegPending   = 0
	RegConfirmed = 1
	RegCancelled = 2
	RegRefunded  = 3
)

type Event struct {
	ID              uint64    `json:"id"`
	UUID            string    `json:"uuid"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	EventType       int       `json:"event_type"`
	Location        string    `json:"location"`
	CoverURL        string    `json:"cover_url"`
	StartAt         time.Time `json:"start_at"`
	EndAt           time.Time `json:"end_at"`
	RegStartAt      time.Time `json:"reg_start_at"`
	RegEndAt        time.Time `json:"reg_end_at"`
	MaxParticipants *int      `json:"max_participants"`
	Fee             float64   `json:"fee"`
	Status          int       `json:"status"`
	CreatorID       uint64    `json:"creator_id"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// computed
	RegisteredCount int  `json:"registered_count,omitempty"`
	IsRegistered    bool `json:"is_registered,omitempty"`
}

func (e *Event) IsOpen() bool {
	now := time.Now()
	return e.Status == EventPublished &&
		now.After(e.RegStartAt) && now.Before(e.RegEndAt)
}

type EventRegistration struct {
	ID        uint64    `json:"id"`
	UUID      string    `json:"uuid"`
	EventID   uint64    `json:"event_id"`
	UserID    uint64    `json:"user_id"`
	Status    int       `json:"status"`
	Note      string    `json:"note"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Registration form snapshot (stored as individual columns)
	RegNameZh          string  `json:"reg_name_zh"`
	RegNameEn          string  `json:"reg_name_en"`
	RegIDNumber        string  `json:"reg_id_number"`
	RegPassportNumber  string  `json:"reg_passport_number"`
	RegGender          *int    `json:"reg_gender"`
	RegBirthday        string  `json:"reg_birthday"`
	RegPhone           string  `json:"reg_phone"`
	RegEmail           string  `json:"reg_email"`
	RegShirtSize       string  `json:"reg_shirt_size"`
	RegFoodType        *int    `json:"reg_food_type"`
	RegAddress         string  `json:"reg_address"`
	RegEmergencyContact  string `json:"reg_emergency_contact"`
	RegEmergencyPhone    string `json:"reg_emergency_phone"`
	RegEmergencyRelation string `json:"reg_emergency_relation"`

	// joined
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	Email       string `json:"email,omitempty"`
	EventTitle  string `json:"event_title,omitempty"`
}
