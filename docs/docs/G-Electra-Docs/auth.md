
#### 1. Sign Up

**Endpoint:** `/api/signup`  
**Method:** `POST`  
**Description:** Create a new user account after validating email and checking if the user is whitelisted.

**Request Body:**
```json
{
  "username": "string (required)",
  "password": "string (required)",
  "full_name": "string (required)",
  "email": "string (required, valid email format)"
}
```

**Responses:**
- **200 OK**
  ```json
  {
    "message": "User signed up successfully",
    "user": {
      "user_id": 1,
      "username": "string",
      "email": "string",
      "full_name": "string",
      "role": "string",
      "join_date": "2023-09-22T00:00:00Z"
    }
  }
  ```
- **400 Bad Request**
  - If the email is not whitelisted:
  ```json
  {
    "error": "Email not whitelisted"
  }
  ```
  - If the user already exists:
  ```json
  {
    "error": "Email already exists"
  }
  ```
  - If other fields are not empty:
  ```json
  {
    "error": "Other fields must be empty"
  }
  ```
- **500 Internal Server Error**
  ```json
  {
    "error": "Failed to create user"
  }
  ```

---

#### 2. Login

**Endpoint:** `/api/login`  
**Method:** `POST`  
**Description:** Authenticate a user and return a JWT token if the credentials are valid.

**Request Body:**
```json
{
  "username": "string",
  "password": "string"
}
```

**Responses:**
- **200 OK**
  ```json
  {
    "token": "jwt_token_string"
  }
  ```
- **401 Unauthorized**
  ```json
  {
    "error": "Invalid credentials"
  }
  ```
- **400 Bad Request**
  ```json
  {
    "error": "Invalid input"
  }
  ```
- **500 Internal Server Error**
  ```json
  {
    "error": "Could not generate token"
  }
  ```