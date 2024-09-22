### Admin API Documentation

**Base URL**: `/admin`

---

### Authentication

All requests to the admin routes require a valid authorization token. The token must be included in the `Authorization` header of the request.

- **Header:**
  ```
  Authorization: Bearer YOUR_TOKEN_HERE
  ```

- **Response on Missing/Invalid Token:**
  ```json
  {
    "status": "error",
    "message": "Missing or invalid token"
  }
  ```

---

### User Management

#### 1. Whitelist a User
- **Endpoint**: `POST /whitelist`
- **Description**: Whitelists a user by their email, name, and role.
- **Request Body**:
  ```json
  {
    "whitelist_email": "user@example.com",
    "name": "John Doe",
    "role": "member"
  }
  ```
- **Responses**:
  - `200 OK`: User successfully whitelisted.
  - `400 Bad Request`: Invalid input or user already exists.
  - `500 Internal Server Error`: Server error.

---

#### 2. Whitelist Users from CSV
- **Endpoint**: `POST /whitelist/csv`
- **Description**: Upload a CSV file to whitelist multiple users.
- **Request**: Multipart form with a file field named `file`.
- **Responses**:
  - `200 OK`: Successfully processed CSV.
  - `400 Bad Request`: Invalid CSV format.
  - `500 Internal Server Error`: Server error.

---

#### 3. Get Whitelisted Users
- **Endpoint**: `GET /users`
- **Description**: Retrieves all whitelisted users.
- **Responses**:
  - `200 OK`: Returns an array of users.
  - `500 Internal Server Error`: Server error.

---

#### 4. Delete a User
- **Endpoint**: `DELETE /users/delete`
- **Description**: Deletes a user by their email.
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "role": "member"
  }
  ```
- **Responses**:
  - `200 OK`: User and associated role data deleted.
  - `404 Not Found`: User not found.
  - `500 Internal Server Error`: Server error.

---

#### 5. Change User Role
- **Endpoint**: `PUT /users/change-role`
- **Description**: Changes a user's role.
- **Request Body**:
  ```json
  {
    "email": "user@example.com",
    "new_role": "staff"
  }
  ```
- **Responses**:
  - `200 OK`: User role changed successfully.
  - `404 Not Found`: User not found.
  - `500 Internal Server Error`: Server error.

---

#### 6. Create Role Table
- **Endpoint**: `POST /roles/create`
- **Description**: Creates a new role table.
- **Request Body**:
  ```json
  {
    "role_name": "staff"
  }
  ```
- **Responses**:
  - `200 OK`: Role table created.
  - `400 Bad Request`: Role name is required.
  - `500 Internal Server Error`: Server error.

---

#### 7. Delete Role Table
- **Endpoint**: `DELETE /roles/delete`
- **Description**: Deletes a role table.
- **Request Body**:
  ```json
  {
    "role_name": "staff"
  }
  ```
- **Responses**:
  - `200 OK`: Role table deleted.
  - `404 Not Found`: Role table does not exist.
  - `500 Internal Server Error`: Server error.

---

### Notices

#### 1. Create a Notice
- **Endpoint**: `/admin/notices`
- **Method**: `POST`
- **Request Body**:
    ```json
    {
      "title": "Notice Title",
      "content": "Notice content",
      "created_by": 1
    }
    ```
- **Response**:
  - **Success**:
    ```json
    {
      "status": "success",
      "message": "Notice created successfully",
      "data": {
        "id": 1,
        "title": "Notice Title",
        "content": "Notice content",
        "created_by": 1,
        "created_at": "2024-09-22T14:30:00Z"
      }
    }
    ```
  - **Error**: `400 Bad Request` / `500 Internal Server Error`

---

#### 2. Delete a Notice
- **Endpoint**: `/admin/notices/:id`
- **Method**: `DELETE`
- **Response**:
  - **Success**:
    ```json
    {
      "status": "success",
      "message": "Notice deleted successfully"
    }
    ```
  - **Error**: `404 Not Found` / `500 Internal Server Error`

---

#### 3. Get a Notice
- **Endpoint**: `/admin/notices/:id`
- **Method**: `GET`
- **Response**:
  - **Success**:
    ```json
    {
      "status": "success",
      "data": {
        "id": 1,
        "title": "Notice Title",
        "content": "Notice content",
        "created_by": 1,
        "created_at": "2024-09-22T14:30:00Z"
      }
    }
    ```
  - **Error**: `404 Not Found` / `500 Internal Server Error`

---

#### 4. Get All Notices
- **Endpoint**: `/admin/notices`
- **Method**: `GET`
- **Response**:
  - **Success**:
    ```json
    {
      "status": "success",
      "data": [
        {
          "id": 1,
          "title": "Notice Title 1",
          "content": "Content 1",
          "created_by": 1,
          "created_at": "2024-09-22T14:30:00Z"
        },
        {
          "id": 2,
          "title": "Notice Title 2",
          "content": "Content 2",
          "created_by": 2,
          "created_at": "2024-09-22T15:00:00Z"
        }
      ]
    }
    ```
  - **Error**: `500 Internal Server Error`

---

### Event Endpoints

#### 1. Create an Event
- **Endpoint**: `POST /events`
- **Description**: Creates a new event.
- **Request Body**:
  ```json
  {
    "event_name": "Tech Talk",
    "description": "A talk on emerging tech trends.",
    "location": "Auditorium",
    "event_date": "2024-10-01T15:00:00Z",
    "image_url": "https://example.com/image.jpg"
  }
  ```
- **Response**:
  - **Success**:
    ```json
    {
      "status": "success",
      "message": "Event created successfully",
      "data": {
        "event_id": 1,
        "event_name": "Tech Talk",
        "description": "A talk on emerging tech trends.",
        "location": "Auditorium",
        "event_date": "2024-10-01T15:00:00Z",
        "created_at": "2024-09-22T12:00:00Z",
        "image_url": "https://example.com/image.jpg"
      }
    }
    ```
  - **Error**:
    ```json
    {
      "status": "error",
      "message": "Event name, description, location, event date, and image URL are required fields"
    }
    ```

---

#### 2. Delete an Event
- **Endpoint**: `DELETE /events/:id`
- **Description**: Deletes an event by ID.
- **Response**:
  - **Success**:
    ```json
    {
      "status": "success",
      "message": "Event deleted successfully"
    }
    ```
  - **Error**:
    ```json
    {
      "status": "error",
      "message": "Event not found"
    }
    ```

---

#### 3. Get an Event
- **Endpoint**: `GET /events/:id`
- **Description**: Retrieves an event by ID.
- **Response**:
  - **Success**:
    ```json
    {
      "status": "success",
      "data": {
        "event_id": 1,
        "event_name": "Tech Talk",
        "description": "A talk on emerging tech trends.",
        "location": "Auditorium",
        "event_date": "2024-10-01T15:00:00Z",
        "created_at": "2024-09-22T12:00:00Z",
        "image_url": "https://example.com/image.jpg"
      }
    }
    ```
  - **Error**:
    ```json
    {
      "status": "error",
      "message": "Event not found"
    }
    ```

---

#### 4. Get All Events
- **Endpoint**: `GET /events`
- **Description**: Retrieves all events.
- **Response**:
  ```json
  {
    "status": "success",
    "data": [
      {
        "event_id": 1,
        "event_name": "Tech Talk",
        "description": "A talk on emerging tech trends.",
        "location": "

Auditorium",
        "event_date": "2024-10-01T15:00:00Z",
        "created_at": "2024-09-22T12:00:00Z",
        "image_url": "https://example.com/image.jpg"
      }
    ]
  }
  ```

---

### Member Management

#### 1. Edit Member Details
- **Endpoint**: `PUT /members/:id`
- **Description**: Edits the details of a member.
- **Path Parameters**:
  - `id`: The ID of the member to edit.
- **Request Body**:
  ```json
  {
    "name": "New Name",
    "email": "new_email@example.com",
    "role": "member"
  }
  ```
- **Responses**:
  - `200 OK`: Member details updated successfully.
  - `404 Not Found`: Member not found.
  - `500 Internal Server Error`: Server error.

---

### Admin Management

#### 2. Promote a User to Admin
- **Endpoint**: `POST /users/:id/admin`
- **Description**: Promotes a user to admin status.
- **Path Parameters**:
  - `id`: The ID of the user to promote.
- **Responses**:
  - `200 OK`: User promoted to admin successfully.
  - `404 Not Found`: User not found.
  - `500 Internal Server Error`: Server error.

---

#### 3. Demote a User from Admin
- **Endpoint**: `DELETE /users/:id/admin`
- **Description**: Demotes a user from admin status.
- **Path Parameters**:
  - `id`: The ID of the user to demote.
- **Responses**:
  - `200 OK`: User demoted from admin successfully.
  - `404 Not Found`: User not found.
  - `500 Internal Server Error`: Server error.

---

#### 4. Get All Admins
- **Endpoint**: `GET /admins`
- **Description**: Retrieves a list of all admin users.
- **Responses**:
  - `200 OK`: Returns an array of admin users.
  - `500 Internal Server Error`: Server error.

---

### Project Management

#### 5. List Pending Project Requests
- **Endpoint**: `GET /projects/pending`
- **Description**: Lists all pending project requests.
- **Responses**:
  - `200 OK`: Returns an array of pending project requests.
  - `500 Internal Server Error`: Server error.

---

#### 6. Approve or Decline a Project Request
- **Endpoint**: `POST /projects/approve-or-decline`
- **Description**: Approves or declines a project request.
- **Request Body**:
  ```json
  {
    "project_id": "12345",
    "action": "approve" // or "decline"
  }
  ```
- **Responses**:
  - `200 OK`: Project request processed successfully.
  - `400 Bad Request`: Invalid input.
  - `404 Not Found`: Project not found.
  - `500 Internal Server Error`: Server error.

---

#### 7. Mark a Project as Completed
- **Endpoint**: `POST /projects/complete`
- **Description**: Marks a project as completed.
- **Request Body**:
  ```json
  {
    "project_id": "12345"
  }
  ```
- **Responses**:
  - `200 OK`: Project marked as completed successfully.
  - `404 Not Found`: Project not found.
  - `500 Internal Server Error`: Server error.

---

#### 8. Set a Project as Homepage Project
- **Endpoint**: `POST /projects/:project_id/home`
- **Description**: Sets a specific project as the homepage project.
- **Path Parameters**:
  - `project_id`: The ID of the project to set as homepage.
- **Responses**:
  - `200 OK`: Project set as homepage successfully.
  - `404 Not Found`: Project not found.
  - `500 Internal Server Error`: Server error.

---

#### 9. Delete Homepage Project
- **Endpoint**: `DELETE /projects/home`
- **Description**: Deletes the current homepage project.
- **Responses**:
  - `200 OK`: Homepage project deleted successfully.
  - `500 Internal Server Error`: Server error.

---