package api

import (
	"github.com/gelectra/gelectra-backend/handlers"
	"github.com/gofiber/fiber/v2"
)

func InitRoutes(app *fiber.App) {

	//Grouping
	user := app.Group("/user")
	admin := app.Group("/admin")

	//Ping

	// Admin
	admin.Use(handlers.JWTAuthMiddleware())
	admin.Post("/whitelist", handlers.WhitelistEmails)
	admin.Post("/whitelist/csv", handlers.WhitelistEmailsFromCSV)
	admin.Get("/whitelist", handlers.GetWhitelistedPending)
	admin.Delete("/whitelist", handlers.DeleteUser)

	admin.Get("/users", handlers.GetWhitelistedPending)
	admin.Delete("/users/delete", handlers.DeleteUser)
	admin.Put("/users/change-role", handlers.ChangeUserRole)
	admin.Post("/roles/create", handlers.CreateRole)
	admin.Delete("/roles/delete", handlers.DeleteRole)

	admin.Post("/notices", handlers.CreateNotice)       // Create a new notice
	admin.Delete("/notices/:id", handlers.DeleteNotice) // Delete a notice by ID
	admin.Get("/notices/:id", handlers.GetNotice)       // Get a specific notice by ID
	admin.Get("/notices", handlers.GetAllNotices)       // Get all notices

	admin.Post("/events", handlers.CreateEvent)               // Create a new event
	admin.Delete("/events/:id", handlers.DeleteEvent)         // Delete an event by ID
	admin.Get("/events/:id", handlers.GetEvent)               // Get a specific event by ID
	admin.Get("/events", handlers.GetAllEvents)               // Get all events
	admin.Get("/events/upcoming", handlers.GetUpcomingEvents) // Get all upcoming events
	admin.Get("/events/past", handlers.GetPastEvents)         // Get all past events

	// Member Management
	admin.Put("/members/:id", handlers.EditMemberDetails) // Edit member details

	// Admin Management
	admin.Post("/users/:id/admin", handlers.MakeUserAdmin)     // Promote a user to admin
	admin.Delete("/users/:id/admin", handlers.RemoveUserAdmin) // Demote a user from admin
	admin.Get("/admins", handlers.GetAdmins)                   // Get all admins

	// Project Management
	admin.Get("/projects/pending", handlers.ListPendingRequest)                  // List all pending project requests
	admin.Post("/projects/approve-or-decline", handlers.ApproveOrDeclineProject) // Approve or decline a project request
	admin.Post("/projects/complete", handlers.CompletedStatusUpdate)             // Mark a project as completed
	admin.Post("/projects/:project_id/home", handlers.SetHomePageProject)        // Set a project as homepage project
	admin.Delete("/projects/home", handlers.DeleteHomePageProject)               // Remove homepage project

	//User
	user.Use(handlers.JWTAuthMiddleware())
	user.Post("/projects", handlers.CreateProject)   // Create a new project
	app.Post("/projects/join", handlers.JoinProject) // Join a project

	// User Profile Management Routes
	app.Put("/users/:user_id", handlers.EditProfile) // Edit user profile

	//HomePage

	//Common
	app.Post("/api/signup", handlers.SignUp)                    // User Sign Up
	app.Post("/api/login", handlers.Login)                      // User Login
	app.Get("/api/projects", handlers.ListProjects)             // List all approved projects
	app.Get("/api/projects/:id", handlers.GetProject)           // Get details of a specific project
	app.Get("/api/projects/home", handlers.GetHomePageProjects) // Get homepage projects

	// Member Management Routes
	app.Get("/api/members", handlers.GetSortedMembers) // Get sorted members
}
