### API Documentation

**Base URL**: `/user`

---

### Project Management

#### 1. Create Project
- **Endpoint**: `POST /projects`
- **Description**: Creates a new project.
- **Request Body**:
  ```json
  {
    "leader_id": 1,
    "project_name": "Project Title",
    "description": "Project description",
    "extra_link": "http://example.com",
    "user_ids": [2, 3, 4]
  }
  ```
- **Responses**:
  - `201 Created`: Project created successfully.
  - `400 Bad Request`: Invalid input.
  - `500 Internal Server Error`: Failed to create project.

---

#### 2. Join Project
- **Endpoint**: `POST /projects/join`
- **Description**: Allows a user to join a project using a join code.
- **Request Body**:
  ```json
  {
    "user_id": 1,
    "join_code": "unique-join-code"
  }
  ```
- **Responses**:
  - `200 OK`: User successfully joined the project.
  - `400 Bad Request`: Invalid input.
  - `404 Not Found`: Invalid join code or user not found.
  - `500 Internal Server Error`: Failed to join the project.

---

### User Profile Management

#### 3. Edit Profile
- **Endpoint**: `PUT /users/:user_id`
- **Description**: Edits a user’s profile.
- **Path Parameters**:
  - `user_id`: The ID of the user to edit.
- **Request Body**:
  ```json
  {
    "username": "new_username",
    "full_name": "New Full Name",
    "password": "newpassword123", // Optional
    "facebook": "http://facebook.com/user",
    "twitter": "http://twitter.com/user",
    "linkedin": "http://linkedin.com/in/user"
  }
  ```
- **Responses**:
  - `200 OK`: Profile updated successfully.
  - `400 Bad Request`: Invalid user ID or input.
  - `404 Not Found`: User not found.
  - `500 Internal Server Error`: Failed to update profile.

---