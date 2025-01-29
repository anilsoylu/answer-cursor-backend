# 📚 Answer API Documentation

## 🔑 Authentication Endpoints

### 📝 Register User

```http
POST /api/v1/auth/register
```

**Request Body:**

```json
{
  "username": "string",
  "email": "string",
  "password": "string"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "id": "uuid",
    "username": "string",
    "email": "string",
    "avatar": "string",
    "status": "active",
    "role": "USER",
    "created_date": "timestamp"
  }
}
```

### 🔐 Login

```http
POST /api/v1/auth/login
```

**Request Body:**

```json
{
  "email": "string",
  "password": "string"
}
```

**Response:**

```json
{
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "string",
    "user": {
      "id": "uuid",
      "username": "string",
      "email": "string",
      "avatar": "string",
      "status": "active",
      "role": "USER",
      "created_date": "timestamp",
      "last_login_date": "timestamp"
    }
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
