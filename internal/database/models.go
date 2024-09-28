package database

import (
	"time"
)

// User represents a user entity
type User struct {
	ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string    `gorm:"uniqueIndex" json:"username"`
	PasswordHash string    `json:"-"`
	Email        string    `gorm:"uniqueIndex" json:"email"`
	FullName     string    `json:"full_name"`
	Role         string    `json:"role"`
	JoinDate     time.Time `json:"join_date"`
	Facebook     string    `json:"facebook"`
	Twitter      string    `json:"twitter"`
	LinkedIn     string    `json:"linkedin"`
}

// Admin represents an admin entity
type Admin struct {
	ID         uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint      `json:"user_id"`
	Privileges string    `json:"privileges"`
	AssignedAt time.Time `gorm:"autoCreateTime" json:"assigned_at"`
	User       User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

// Member represents a member entity
type Member struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          uint      `json:"user_id"`
	RegistrationNum string    `json:"registration_num"`
	MembershipLevel string    `json:"membership_level"`
	JoinDate        time.Time `json:"join_date"`
	User            User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

// CoreMember represents a core member entity
type CoreMember struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID   uint      `json:"user_id"`
	Role     string    `json:"role"`
	JoinDate time.Time `gorm:"autoCreateTime" json:"join_date"`
	User     User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

// Staff represents a staff entity
type Staff struct {
	ID         uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint   `json:"user_id"`
	Role       string `json:"role"`
	Department string `json:"department"`
	User       User   `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

// Event represents an event entity
type Event struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	EventName   string    `json:"event_name"`
	EventDate   time.Time `json:"event_date"`
	Description string    `json:"description"`
	Location    string    `json:"location"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
}

// Attendee represents an attendee entity
type Attendee struct {
	UserID  uint  `gorm:"primaryKey" json:"user_id"`
	EventID uint  `gorm:"primaryKey" json:"event_id"`
	User    User  `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Event   Event `gorm:"foreignKey:EventID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

// Project represents a project entity
type Project struct {
	ID                uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LeaderID          uint      `json:"leader_id"`
	ProjectName       string    `json:"project_name"`
	Description       string    `json:"description"`
	ExtraLink         string    `json:"extra_link"`
	Status            string    `json:"status"`
	CreatedBy         uint      `json:"created_by"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	JoinCode          string    `json:"join_code"`
	IsHomePageProject bool      `gorm:"default:false" json:"is_home_page_project"`
	Users             []User    `gorm:"many2many:project_users;constraint:OnDelete:CASCADE" json:"-"`
	Leader            User      `gorm:"foreignKey:LeaderID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

// ProjectRequest represents a project request entity
type ProjectRequest struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ProjectID   uint      `json:"project_id"`
	UserID      uint      `json:"user_id"`
	Status      string    `json:"status"`
	RequestedAt time.Time `gorm:"autoCreateTime" json:"requested_at"`
	Project     Project   `gorm:"foreignKey:ProjectID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	User        User      `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

// ProjectMember represents a project member entity
type ProjectMember struct {
	ProjectID uint    `gorm:"primaryKey" json:"project_id"`
	UserID    uint    `gorm:"primaryKey" json:"user_id"`
	User      User    `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
	Project   Project `gorm:"foreignKey:ProjectID;references:ID;constraint:OnDelete:CASCADE" json:"-"`
}

// Notice represents a notice entity
type Notice struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedBy uint      `json:"created_by"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
}
