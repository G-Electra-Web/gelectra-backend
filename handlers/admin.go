package handlers

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"time"

	"github.com/gelectra/gelectra-backend/internal/database"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func WhitelistEmailsFromCSV(c *fiber.Ctx) error {
	// Retrieve the file from the form data
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to retrieve the file",
		})
	}

	// Open the uploaded file
	f, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open the file",
		})
	}
	defer f.Close()

	// Create a CSV reader
	reader := csv.NewReader(f)
	// Read the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to read the CSV file",
		})
	}

	// Iterate over each record in the CSV file
	for _, record := range records {
		// Validate the CSV structure (assume [email, name, role1, role2, ...])
		if len(record) < 3 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid CSV format",
			})
		}

		email := record[0]
		name := record[1]
		roles := record[2:] // All remaining columns are roles

		// Check if the user already exists
		var user database.User
		err := database.DB.Table("users").Where("email = ?", email).First(&user).Error
		if err == gorm.ErrRecordNotFound {
			// Create a new user
			user = database.User{
				FullName: name,
				Email:    email,
				JoinDate: time.Now(),
			}

			err = database.DB.Table("users").Create(&user).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": err.Error(),
				})
			}
		} else if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Iterate through each role and dynamically assign it to the user
		for _, role := range roles {
			// Generic structure to hold role data
			roleData := map[string]interface{}{
				"user_id":   user.ID,
				"join_date": time.Now(),
			}

			// Attempt to insert the role into the corresponding table
			err = database.DB.Table(role).Create(roleData).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to assign role '%s': %s", role, err.Error()),
				})
			}
		}
	}

	return c.JSON(fiber.Map{
		"message": "CSV processed successfully, users whitelisted, and roles assigned",
	})
}

func WhitelistEmails(c *fiber.Ctx) error {
	// Define a struct to capture the incoming JSON payload
	type WhitelistInput struct {
		Email string `json:"whitelist_email"`
		Name  string `json:"name"`
		Role  string `json:"role"`
	}

	// Parse the incoming JSON into the WhitelistInput struct
	var input WhitelistInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON format",
		})
	}

	var user database.User

	err := database.DB.Table("users").Where("email = ?", input.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		// Create a new user
		user = database.User{
			FullName: input.Name,
			Email:    input.Email,
			Role:     input.Role,
			JoinDate: time.Now(),
		}

		err = database.DB.Create(&user).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		// Dynamically assign the role
		roleData := map[string]interface{}{
			"user_id":   user.ID,
			"join_date": time.Now(),
		}

		// Attempt to insert the role into the corresponding table
		err = database.DB.Table(input.Role).Create(roleData).Error
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to assign role '%s': %s", input.Role, err.Error()),
			})
		}

		return c.JSON(fiber.Map{
			"message": "Whitelisted user successfully and assigned role",
			"user":    user,
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error": "User already exists",
	})
}

func GetWhitelistedPending(c *fiber.Ctx) error {
	var users []database.User

	err := database.DB.Table("users").Find(&users).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"users": users,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	// Struct to hold the incoming request data
	type DeleteUserRequest struct {
		Email string `json:"email"`
		Role  string `json:"role"` // The role should be dynamically handled
	}

	var req DeleteUserRequest

	// Parse and validate the request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	// Fetch the user from the database
	var user database.User
	err := database.DB.Table("users").Where("email = ?", req.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Use dynamic role to handle the table name
	roleTable := req.Role // Assuming the table name corresponds to the role name

	// Check if the role table exists and delete the role-specific entry
	err = database.DB.Table(roleTable).Where("user_id = ?", user.ID).Delete(nil).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete role data for user",
		})
	}

	// Delete the user entry from the users table
	err = database.DB.Table("users").Where("email = ?", req.Email).Delete(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User and associated role data deleted successfully",
	})
}

func ChangeUserRole(c *fiber.Ctx) error {
	// Struct to hold the incoming request data
	type ChangeUserRoleRequest struct {
		Email   string `json:"email"`
		NewRole string `json:"new_role"` // New role will be dynamic
	}

	var req ChangeUserRoleRequest

	// Parse and validate the request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	// Fetch the user from the database
	var user database.User
	err := database.DB.Table("users").Where("email = ?", req.Email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// Delete the current role from the corresponding table dynamically
		currentRoleTable := user.Role
		if err := tx.Table(currentRoleTable).Where("user_id = ?", user.ID).Delete(nil).Error; err != nil {
			return err
		}

		// Add the user to the new role's table dynamically
		newRoleTable := req.NewRole
		roleData := map[string]interface{}{
			"user_id":    user.ID,
			"role_name":  req.NewRole,
			"created_at": time.Now(),
		}

		// Insert new role data into the dynamically determined table
		if err := tx.Table(newRoleTable).Create(roleData).Error; err != nil {
			return err
		}

		// Update the user role in the users table
		if err := tx.Table("users").Where("id = ?", user.ID).Update("role", req.NewRole).Error; err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to change user role",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User role changed successfully",
	})
}

func CreateRole(c *fiber.Ctx) error {
	// Define a struct to hold the incoming JSON data
	type CreateRoleRequest struct {
		RoleName string `json:"role_name"`
	}

	var req CreateRoleRequest

	// Parse the JSON body into the struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON input",
		})
	}

	// Validate role name
	if req.RoleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Role name is required",
		})
	}

	// Ensure the role name is valid and doesn't contain harmful characters
	// This can be enhanced with a stricter regex if needed
	validRoleName := regexp.MustCompile(`^[a-zA-Z_]+$`).MatchString
	if !validRoleName(req.RoleName) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid role name. Only letters and underscores are allowed.",
		})
	}

	// Check if the table already exists
	if database.DB.Migrator().HasTable(req.RoleName) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": fmt.Sprintf("Role table '%s' already exists", req.RoleName),
		})
	}

	// Define the structure of the new role table dynamically
	// For example, we can map to a Member structure, or different structures depending on the role
	type Role struct {
		ID        uint      `gorm:"primaryKey"`
		UserID    uint      `gorm:"not null"`
		RoleName  string    `gorm:"size:100;not null"`
		CreatedAt time.Time `gorm:"autoCreateTime"`
		UpdatedAt time.Time `gorm:"autoUpdateTime"`
	}

	// Attempt to create the new table
	if err := database.DB.Table(req.RoleName).AutoMigrate(&Role{}); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create role table",
		})
	}

	// Success response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": fmt.Sprintf("Role table '%s' created successfully", req.RoleName),
	})
}

func DeleteRole(c *fiber.Ctx) error {
	// Define a struct to hold the incoming JSON data
	type DeleteRoleRequest struct {
		RoleName string `json:"role_name"`
	}

	var req DeleteRoleRequest

	// Parse the JSON body into the struct
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid JSON input",
		})
	}

	if req.RoleName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Role name is required",
		})
	}

	// Check if the table exists
	isRole := database.DB.Migrator().HasTable(req.RoleName)

	if !isRole {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Role table does not exist",
		})
	}

	// Define a struct for the role table records
	type RoleRecord struct {
		UserID   uint
		JoinDate time.Time
		// Add any additional fields from the role table if necessary
	}

	// Retrieve all records from the role's table
	var roleRecords []RoleRecord
	err := database.DB.Table(req.RoleName).Find(&roleRecords).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve records from role table",
		})
	}

	// If there are records, move them to the members table and update the user table
	if len(roleRecords) > 0 {
		for _, record := range roleRecords {
			// Move to members table
			member := database.Member{
				UserID:          record.UserID,
				RegistrationNum: "Default RegNum", // You can modify or generate a default registration number
				MembershipLevel: "Member",         // Set default membership level
				JoinDate:        record.JoinDate,
			}

			err = database.DB.Table("members").Create(&member).Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to move records to members table",
				})
			}

			// Update the user's role in the user table (set to 'member' or a new role)
			err = database.DB.Table("users").
				Where("id = ?", record.UserID).
				Update("role", "members").Error
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to update user role",
				})
			}
		}
	}

	// Delete the role table
	err = database.DB.Migrator().DropTable(req.RoleName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete role table",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Role table deleted successfully, records moved to members, and user roles updated",
	})
}

// Notices

func CreateNotice(c *fiber.Ctx) error {
	// Create a new instance of Notice
	notice := new(database.Notice)

	// Parse the incoming JSON request body into the notice struct
	if err := c.BodyParser(notice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse request body",
		})
	}

	if notice.Title == "" || notice.Content == "" || notice.CreatedBy == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Title, content, and created_by are required fields",
		})
	}

	notice.CreatedAt = time.Now()

	if err := database.DB.Create(&notice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create notice",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Notice created successfully",
		"data":    notice,
	})
}

func DeleteNotice(c *fiber.Ctx) error {
	// Get the notice ID from the request parameters
	noticeID := c.Params("id")

	// Check if the notice exists
	var notice database.Notice
	if err := database.DB.First(&notice, noticeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Notice not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find notice",
		})
	}

	// Delete the notice
	if err := database.DB.Delete(&notice).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete notice",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Notice deleted successfully",
	})
}

func GetNotice(c *fiber.Ctx) error {
	// Get the notice ID from the request parameters
	noticeID := c.Params("id")

	// Retrieve the notice from the database
	var notice database.Notice
	if err := database.DB.First(&notice, noticeID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Notice not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve notice",
		})
	}

	// Return the retrieved notice
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   notice,
	})
}

func GetAllNotices(c *fiber.Ctx) error {
	// Create a slice to hold all notices
	var notices []database.Notice

	// Retrieve all notices from the database
	if err := database.DB.Find(&notices).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve notices",
		})
	}

	// Return the retrieved notices
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   notices,
	})
}

// Events
func CreateEvent(c *fiber.Ctx) error {
	// Create a new instance of Event
	event := new(database.Event)

	// Parse the incoming JSON request body into the event struct
	if err := c.BodyParser(event); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse request body",
		})
	}

	// Validate required fields
	if event.EventName == "" || event.Description == "" || event.Location == "" || event.EventDate.IsZero() || event.ImageURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Event name, description, location, event date, and image URL are required fields",
		})
	}

	// Set the CreatedAt field to the current time
	event.CreatedAt = time.Now()

	// Insert the new event into the database
	if err := database.DB.Create(&event).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create event",
		})
	}

	// Return a success response with the created event
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Event created successfully",
		"data":    event,
	})
}

func DeleteEvent(c *fiber.Ctx) error {
	// Get the event ID from the request parameters
	eventID := c.Params("id")

	// Check if the event exists
	var event database.Event
	if err := database.DB.First(&event, eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Event not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to find event",
		})
	}

	// Delete the event
	if err := database.DB.Delete(&event).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete event",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Event deleted successfully",
	})
}

func GetEvent(c *fiber.Ctx) error {
	// Get the event ID from the request parameters
	eventID := c.Params("id")

	// Retrieve the event from the database
	var event database.Event
	if err := database.DB.First(&event, eventID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Event not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve event",
		})
	}

	// Return the retrieved event
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   event,
	})
}

func GetAllEvents(c *fiber.Ctx) error {
	// Create a slice to hold all events
	var events []database.Event

	// Retrieve all events from the database
	if err := database.DB.Find(&events).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve events",
		})
	}

	// Return the retrieved events
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   events,
	})
}

func GetUpcomingEvents(c *fiber.Ctx) error {
	// Create a slice to hold upcoming events
	var upcomingEvents []database.Event

	// Get the current time
	currentTime := time.Now()

	// Retrieve events where the EventDate is greater than the current time
	if err := database.DB.Where("event_date > ?", currentTime).Find(&upcomingEvents).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve upcoming events",
		})
	}

	// Return the retrieved upcoming events
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   upcomingEvents,
	})
}

func GetPastEvents(c *fiber.Ctx) error {
	// Create a slice to hold past events
	var pastEvents []database.Event

	// Get the current time
	currentTime := time.Now()

	// Retrieve events where the EventDate is less than or equal to the current time
	if err := database.DB.Where("event_date <= ?", currentTime).Find(&pastEvents).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to retrieve past events",
		})
	}

	// Return the retrieved past events
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   pastEvents,
	})
}

// Edit Member
func EditMemberDetails(c *fiber.Ctx) error {
	memberID := c.Params("id")
	var member database.Member

	// Find the member by ID
	if err := database.DB.First(&member, memberID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Member not found",
		})
	}

	// Parse the JSON body to get the updated details
	if err := c.BodyParser(&member); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to parse request body",
		})
	}

	// Update the member details in the database
	if err := database.DB.Save(&member).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update member details",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Member details updated successfully",
		"data":    member,
	})
}

// Dynamic Admin
func MakeUserAdmin(c *fiber.Ctx) error {

	// Parse the user ID from the request parameters
	ID := c.Params("id")
	if ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// Convert ID to integer
	var uid uint
	if _, err := fmt.Sscanf(ID, "%d", &uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	// Check if user exists
	var user database.User
	if err := database.DB.First(&user, uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Check if user is already an admin
	var admin database.Admin
	if err := database.DB.First(&admin, "user_id = ?", uid).Error; err == nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is already an admin"})
	}

	// Create admin record
	newAdmin := database.Admin{
		ID:         (uid),
		Privileges: "full", // or based on input
		AssignedAt: time.Now(),
	}

	if err := database.DB.Create(&newAdmin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to promote user to admin"})
	}

	// Optionally update the user's Role field
	user.Role = "admin"
	if err := database.DB.Save(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user role"})
	}

	return c.JSON(fiber.Map{
		"message": "User promoted to admin successfully",
		"user":    user,
		"admin":   newAdmin,
	})
}

func RemoveUserAdmin(c *fiber.Ctx) error {
	// Parse the user ID from the request parameters
	ID := c.Params("id")
	if ID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID is required"})
	}

	// Convert ID to integer
	var uid int
	if _, err := fmt.Sscanf(ID, "%d", &uid); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid User ID"})
	}

	// Check if user is an admin
	var admin database.Admin
	if err := database.DB.First(&admin, "user_id = ?", uid).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User is not an admin"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Delete admin record
	if err := database.DB.Delete(&admin).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to remove admin role"})
	}

	return c.JSON(fiber.Map{
		"message": "User demoted from admin successfully",
		"admin":   admin,
	})
}

func GetAdmins(c *fiber.Ctx) error {
	var admins []database.Admin

	if err := database.DB.Find(&admins).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	return c.JSON(fiber.Map{
		"admins": admins,
	})
}

// Project
func ListPendingRequest(c *fiber.Ctx) error {
	var pendingRequests []database.ProjectRequest

	// Fetch all project requests with status 'pending'
	if err := database.DB.Where("status = ?", "pending").Find(&pendingRequests).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve pending requests"})
	}

	// If no pending requests found
	if len(pendingRequests) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "No pending requests found"})
	}

	return c.JSON(fiber.Map{
		"pending_requests": pendingRequests,
	})
}

func ApproveOrDeclineProject(c *fiber.Ctx) error {
	type RequestInput struct {
		ProjectID int    `json:"project_id" validate:"required"`
		Action    string `json:"action" validate:"required"` // "approve" or "decline"
	}

	var input RequestInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find the project request
	var project database.Project
	if err := database.DB.First(&project, input.ProjectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	}

	// Update project status based on action
	if input.Action == "approve" {
		project.Status = "approved"
	} else if input.Action == "decline" {
		project.Status = "declined"
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid action"})
	}

	if err := database.DB.Save(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update project status"})
	}

	return c.JSON(fiber.Map{
		"message": "Project status updated successfully",
		"project": project,
	})
}

func CompletedStatusUpdate(c *fiber.Ctx) error {
	type CompleteRequest struct {
		ProjectID int `json:"project_id" validate:"required"`
	}

	var input CompleteRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Find the project
	var project database.Project
	if err := database.DB.First(&project, input.ProjectID).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Project not found"})
	}

	// Update project status to "completed"
	project.Status = "completed"

	if err := database.DB.Save(&project).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update project status"})
	}

	return c.JSON(fiber.Map{
		"message": "Project marked as completed successfully",
		"project": project,
	})
}

// Home Page
func SetHomePageProject(c *fiber.Ctx) error {
	// Retrieve the project ID from the request parameter
	projectID, err := c.ParamsInt("project_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid project ID"})
	}

	// Unset any currently featured homepage project
	if err := database.DB.Model(&database.Project{}).Where("is_home_page_project = ?", true).
		Update("is_home_page_project", false).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to unset current homepage project"})
	}

	// Set the selected project as the homepage project
	if err := database.DB.Model(&database.Project{}).Where("project_id = ?", projectID).
		Update("is_home_page_project", true).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to set project as homepage project"})
	}

	return c.JSON(fiber.Map{
		"message": "Project set as homepage project successfully",
	})
}
func DeleteHomePageProject(c *fiber.Ctx) error {
	// Find the currently set homepage project
	if err := database.DB.Model(&database.Project{}).
		Where("is_home_page_project = ?", true).
		Update("is_home_page_project", false).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete homepage project"})
	}

	return c.JSON(fiber.Map{
		"message": "Homepage project removed successfully",
	})
}
