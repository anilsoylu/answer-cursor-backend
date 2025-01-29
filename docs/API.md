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
