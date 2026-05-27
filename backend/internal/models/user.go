package models

import "time"

const (
	RoleMember = 1
	RoleCoach  = 2
	RoleAdmin  = 8
	RoleSuper  = 9

	StatusPending   = 0
	StatusActive    = 1
	StatusSuspended = 2
	StatusRejected  = 3

	GenderMale   = 1
	GenderFemale = 2
	GenderOther  = 3

	FoodMeat  = 1
	FoodVeg   = 2
	FoodVegan = 3
)

type User struct {
	// ── 系統 ─────────────────────────────────────────────────
	ID            uint64     `json:"id"`
	UUID          string     `json:"uuid"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Role          int        `json:"role"`
	Status        int        `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	AvatarURL     string     `json:"avatar_url"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty"`

	// ── 個人資料（可修改）────────────────────────────────────
	DisplayName     string `json:"display_name"`
	NameZh          string `json:"name_zh"`
	NameEn          string `json:"name_en"`
	IDNumber        string `json:"id_number"`
	PassportNumber  string `json:"passport_number"`
	Gender          *int   `json:"gender"`
	Birthday        string `json:"birthday"`
	Phone           string `json:"phone"`
	ShirtSize       string `json:"shirt_size"`
	FoodType        *int   `json:"food_type"`
	Address         string `json:"address"`

	// ── 緊急聯絡人 ───────────────────────────────────────────
	EmergencyContact  string `json:"emergency_contact"`
	EmergencyPhone    string `json:"emergency_phone"`
	EmergencyRelation string `json:"emergency_relation"`
}

func (u *User) IsAdmin() bool  { return u.Role >= RoleAdmin }
func (u *User) IsSuper() bool  { return u.Role >= RoleSuper }
func (u *User) IsActive() bool { return u.Status == StatusActive }

func (u *User) FullProfile() map[string]any {
	return map[string]any{
		"id": u.ID, "uuid": u.UUID, "username": u.Username,
		"email": u.Email, "role": u.Role, "status": u.Status,
		"email_verified": u.EmailVerified, "avatar_url": u.AvatarURL,
		"created_at": u.CreatedAt, "updated_at": u.UpdatedAt,
		"display_name": u.DisplayName, "name_zh": u.NameZh, "name_en": u.NameEn,
		"id_number": u.IDNumber, "passport_number": u.PassportNumber,
		"gender": u.Gender, "birthday": u.Birthday, "phone": u.Phone,
		"shirt_size": u.ShirtSize, "food_type": u.FoodType, "address": u.Address,
		"emergency_contact": u.EmergencyContact,
		"emergency_phone":   u.EmergencyPhone,
		"emergency_relation": u.EmergencyRelation,
	}
}

func (u *User) PublicProfile() map[string]any {
	return map[string]any{
		"id": u.ID, "uuid": u.UUID, "username": u.Username,
		"email": u.Email, "display_name": u.DisplayName,
		"avatar_url": u.AvatarURL, "role": u.Role,
		"status": u.Status, "created_at": u.CreatedAt,
	}
}

func (u *User) RegistrationProfile() map[string]any {
	return map[string]any{
		"display_name": u.DisplayName, "name_zh": u.NameZh, "name_en": u.NameEn,
		"id_number": u.IDNumber, "passport_number": u.PassportNumber,
		"gender": u.Gender, "birthday": u.Birthday,
		"phone": u.Phone, "email": u.Email,
		"shirt_size": u.ShirtSize, "food_type": u.FoodType, "address": u.Address,
		"emergency_contact": u.EmergencyContact,
		"emergency_phone":   u.EmergencyPhone,
		"emergency_relation": u.EmergencyRelation,
	}
}
