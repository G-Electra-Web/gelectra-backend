package database

import (
	"time"
)

// User represents a user entity
type User struct {
	UserID       int       `gorm:"primaryKey;autoIncrement" db:"user_id"`
	Username     string    `gorm:"uniqueIndex" db:"username"`
	PasswordHash string    `db:"password_hash"`
	Email        string    `gorm:"uniqueIndex" db:"email"`
	FullName     string    `db:"full_name"`
	Role         string    `db:"role"`
	JoinDate     time.Time `json:"join_date" db:"join_date"`
	Facebook     string    `db:"facebook"` // Field for Facebook link
	Twitter      string    `db:"twitter"`  // Field for Twitter link
	LinkedIn     string    `db:"linkedin"` // Field for LinkedIn link
}

// Admin represents an admin entity
type Admin struct {
	AdminID    int       `gorm:"primaryKey;autoIncrement" db:"admin_id"`
	UserID     int       `gorm:"uniqueIndex;not null" db:"user_id"` // Foreign key referencing User
	Privileges string    `db:"privileges"`                          // e.g., "full", "limited"
	AssignedAt time.Time `db:"assigned_at"`
	User       User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Member represents a member entity
type Member struct {
	MemberID        int       `gorm:"primaryKey;autoIncrement" db:"member_id"`
	UserID          int       `gorm:"uniqueIndex;not null" db:"user_id"`
	RegistrationNum string    `db:"reg_num"`
	MembershipLevel string    `db:"membership_level"`
	JoinDate        time.Time `db:"join_date"`
	User            User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// CoreMember represents a core member entity
type CoreMember struct {
	CoreMemberID int       `gorm:"primaryKey;autoIncrement" db:"core_member_id"`
	UserID       int       `gorm:"uniqueIndex;not null" db:"user_id"`
	Role         string    `db:"role"`
	JoinDate     time.Time `db:"join_date"`
	User         User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Staff represents a staff entity
type Staff struct {
	StaffID    int    `gorm:"primaryKey;autoIncrement" db:"staff_id"`
	UserID     int    `gorm:"uniqueIndex;not null" db:"user_id"`
	Role       string `db:"role"`
	Department string `db:"department"`
	User       User   `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Event represents an event entity
type Event struct {
	EventID     int       `gorm:"primaryKey;autoIncrement" db:"event_id"`
	EventName   string    `db:"event_name"`
	EventDate   time.Time `db:"event_date"`
	Description string    `db:"description"`
	Location    string    `db:"location"`
	ImageURL    string    `db:"image_url"` // Path to the image
	CreatedAt   time.Time `db:"created_at"`
}

// Attendee represents an attendee entity
type Attendee struct {
	UserID  int   `gorm:"primaryKey" db:"user_id"`
	EventID int   `gorm:"primaryKey" db:"event_id"`
	User    User  `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Event   Event `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
}

// Project represents a project entity
type Project struct {
	ProjectID         int       `gorm:"primaryKey;autoIncrement" db:"project_id"`
	LeaderID          int       `gorm:"not null" db:"leader_id"` // Foreign key referencing User
	ProjectName       string    `db:"project_name"`
	Description       string    `db:"description"`
	ExtraLink         string    `db:"extra_link"`
	Status            string    `db:"status"`
	CreatedBy         int       `db:"created_by"`
	CreatedAt         time.Time `db:"created_at"`
	JoinCode          string    `db:"join_code"`
	IsHomePageProject bool      `gorm:"default:false" db:"is_home_page_project"`
	Users             []User    `gorm:"many2many:project_users;constraint:OnDelete:CASCADE"`
	Leader            User      `gorm:"foreignKey:LeaderID;constraint:OnDelete:CASCADE"`
}

// GalleryImage represents a gallery image entity
type GalleryImage struct {
	ImageID     int       `gorm:"primaryKey;autoIncrement" db:"image_id"`
	EventID     int       `gorm:"not null" db:"event_id"` // Foreign key referencing Event
	UserID      int       `gorm:"not null" db:"user_id"`  // Foreign key referencing User
	ImageURL    string    `db:"image_url"`
	Description string    `db:"description"`
	UploadedAt  time.Time `db:"uploaded_at"`
	Event       Event     `gorm:"foreignKey:EventID;constraint:OnDelete:CASCADE"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// ProjectRequest represents a project request entity
type ProjectRequest struct {
	RequestID   int       `gorm:"primaryKey;autoIncrement" db:"request_id"`
	ProjectID   int       `gorm:"not null" db:"project_id"` // Foreign key referencing Project
	UserID      int       `gorm:"not null" db:"user_id"`    // Foreign key referencing User
	Status      string    `db:"status"`
	RequestedAt time.Time `db:"requested_at"`
	Project     Project   `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
	User        User      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// ProjectMember represents a project member entity
type ProjectMember struct {
	ProjectID int     `gorm:"primaryKey" db:"project_id"`
	UserID    int     `gorm:"primaryKey" db:"user_id"`
	User      User    `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
	Project   Project `gorm:"foreignKey:ProjectID;constraint:OnDelete:CASCADE"`
}

// Notice represents a notice entity
type Notice struct {
	NoticeID  int       `gorm:"primaryKey;autoIncrement" db:"notice_id"`
	Title     string    `db:"title"`
	Content   string    `db:"content"`
	CreatedBy int       `db:"created_by"`
	CreatedAt time.Time `db:"created_at"`
}
