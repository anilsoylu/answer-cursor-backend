# 📚 Answer API Documentation

## 🔑 Authentication Endpoints

### 📝 Register User

Register a new user account.

- **URL**: `/api/v1/auth/register`
- **Method**: `POST`
- **Content-Type**: `application/json`

#### Request Body

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

#### Request Parameters

| Parameter | Type   | Required | Description                    |
| --------- | ------ | -------- | ------------------------------ |
| username  | string | Yes      | Unique username (min: 3 chars) |
| email     | string | Yes      | Valid email address            |
| password  | string | Yes      | Password (min: 6 chars)        |

#### Success Response

- **Code**: `201 Created`
- **Content**:

```json
{
  "success": true,
  "data": {
    "message": "User registered successfully"
  }
}
```

#### Error Responses

- **Code**: `400 Bad Request`

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Username must be at least 3 characters long"
  }
}
```

- **Code**: `409 Conflict`

```json
{
  "success": false,
  "error": {
    "code": "USERNAME_EXISTS",
    "message": "Username already exists"
  }
}
```

### 🔐 Login

Authenticate a user and receive a JWT token.

- **URL**: `/api/v1/auth/login`
- **Method**: `POST`
- **Content-Type**: `application/json`

#### Request Body

```json
{
  "identifier": "string",
  "password": "string"
}
```

#### Request Parameters

| Parameter  | Type   | Required | Description               |
| ---------- | ------ | -------- | ------------------------- |
| identifier | string | Yes      | Email address or username |
| password   | string | Yes      | User's password           |

#### Success Response

- **Code**: `200 OK`
- **Content**:

```json
{
  "success": true,
  "data": {
    "token": "string",
    "token_type": "Bearer",
    "expires_in": 86400
  }
}
```

#### Error Responses

- **Code**: `400 Bad Request`

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Identifier is required"
  }
}
```

- **Code**: `401 Unauthorized`

```json
{
  "success": false,
  "error": {
    "code": "INVALID_CREDENTIALS",
    "message": "Invalid identifier or password"
  }
}
```

### Update User Status

```http
PATCH /api/v1/users/status
```

#### Request Headers

| Name            | Type     | Description                                  |
| --------------- | -------- | -------------------------------------------- |
| `Authorization` | `string` | **Required**. JWT token with `Bearer` prefix |

#### Query Parameters

| Name      | Type     | Description                                                   |
| --------- | -------- | ------------------------------------------------------------- |
| `user_id` | `number` | **Optional**. Target user ID (Only for ADMIN and SUPER_ADMIN) |

#### Request Body

```json
{
    "status": "active" | "passive" | "banned"
}
```

#### Success Response

```json
{
  "status": "success",
  "data": {
    "message": "User status updated successfully"
  }
}
```

#### Error Response

```json
{
    "status": "error",
    "error": {
        "code": "validation_error" | "unauthorized" | "forbidden" | "not_found" | "internal_error",
        "message": "Error message"
    }
}
```

#### Notes

- Users can update their own status (only `active` and `passive`)
- `ADMIN` and `SUPER_ADMIN` can update other users' status using `user_id` query parameter
- `ADMIN` cannot change other admin's status
- `SUPER_ADMIN` status cannot be changed
- Valid status values are: `active`, `passive`, `banned`

### Update User Profile

```http
PUT /api/v1/users/profile
```

#### Request Headers

| Name            | Type     | Description                                  |
| --------------- | -------- | -------------------------------------------- |
| `Authorization` | `string` | **Required**. JWT token with `Bearer` prefix |

#### Request Body

```json
{
  "username": "string", // Optional, min: 3 chars
  "email": "string", // Optional, valid email
  "avatar": "string" // Optional
}
```

#### Success Response

```json
{
  "status": "success",
  "data": {
    "user": {
      "id": 1,
      "username": "string",
      "email": "string",
      "avatar": "string",
      "status": "active",
      "role": "USER"
    },
    "message": "Profile updated successfully"
  }
}
```

#### Error Response

```json
{
    "status": "error",
    "error": {
        "code": "validation_error" | "not_found" | "conflict" | "internal_error",
        "message": "Error message"
    }
}
```

#### Notes

- All fields in request body are optional
- Username must be at least 3 characters long
- Email must be valid
- Username and email must be unique
- Password update is handled by a separate endpoint

## 🔄 Response Codes

| Status Code | Description           |
| ----------- | --------------------- |
| 200         | Success               |
| 201         | Created               |
| 400         | Bad Request           |
| 401         | Unauthorized          |
| 403         | Forbidden             |
| 404         | Not Found             |
| 409         | Conflict              |
| 422         | Unprocessable Entity  |
| 429         | Too Many Requests     |
| 500         | Internal Server Error |

## 🔒 Authentication

All protected endpoints require a JWT token in the Authorization header:

```http
Authorization: Bearer <token>
```

## ⚠️ Error Response Format

```json
{
  "status": "error",
  "message": "Error message here",
  "errors": [
    {
      "field": "field_name",
      "message": "Error message for this field"
    }
  ]
}
```

## Error Codes 🚨

| Code                  | Description                    |
| --------------------- | ------------------------------ |
| VALIDATION_ERROR      | Request validation failed      |
| INVALID_INPUT         | Invalid input parameters       |
| UNAUTHORIZED          | Authentication required        |
| FORBIDDEN             | Access forbidden               |
| NOT_FOUND             | Resource not found             |
| USERNAME_EXISTS       | Username is already taken      |
| EMAIL_EXISTS          | Email is already registered    |
| INVALID_CREDENTIALS   | Invalid identifier or password |
| ACCOUNT_INACTIVE      | User account is not active     |
| INTERNAL_SERVER_ERROR | Unexpected server error        |
