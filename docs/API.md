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
  "status": "success",
  "message": "User registered successfully",
  "data": {
    "id": "integer",
    "username": "string",
    "email": "string",
    "avatar": "string",
    "status": "string",
    "role": "string",
    "created_at": "timestamp",
    "last_login_date": "timestamp"
  }
}
```

#### Error Responses

- **Code**: `400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "code": "username_taken",
    "message": "Username is already taken"
  }
}
```

- **Code**: `409 Conflict`

```json
{
  "status": "error",
  "error": {
    "code": "email_taken",
    "message": "Email is already taken"
  }
}
```

- **Code**: `403 Forbidden`

```json
{
  "status": "error",
  "error": {
    "code": "account_banned",
    "message": "This account is banned"
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
  "status": "success",
  "message": "Login successful",
  "data": {
    "token": "string",
    "user": {
      "id": "integer",
      "username": "string",
      "email": "string",
      "avatar": "string",
      "status": "string",
      "role": "string",
      "created_at": "timestamp",
      "last_login_date": "timestamp"
    }
  }
}
```

#### Error Responses

- **Code**: `400 Bad Request`

```json
{
  "status": "error",
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Identifier is required"
  }
}
```

- **Code**: `401 Unauthorized`

```json
{
  "status": "error",
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
    "status": "active" | "passive" | "banned" | "frozen"
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

### 🔐 Update Password

**Endpoint:** `PUT /api/v1/auth/password`

**Authentication Required:** Yes

**Request Body:**

```json
{
  "current_password": "string",
  "new_password": "string",
  "confirm_password": "string"
}
```

**Validation Rules:**

- `current_password`: Required
- `new_password`: Required, minimum 6 characters
- `confirm_password`: Required, must match new_password

**Success Response:**

```json
{
  "status": "success",
  "data": {
    "message": "Password updated successfully"
  }
}
```

**Error Responses:**

_Invalid Request Body (400 Bad Request)_

```json
{
  "status": "error",
  "error": {
    "code": "validation_error",
    "message": "Validation failed"
  }
}
```

_Current Password Incorrect (400 Bad Request)_

```json
{
  "status": "error",
  "error": {
    "code": "validation_error",
    "message": "Current password is incorrect"
  }
}
```

_Server Error (500 Internal Server Error)_

```json
{
  "status": "error",
  "error": {
    "code": "internal_error",
    "message": "Failed to update password"
  }
}
```

### 🚫 Ban User

**Endpoint:** `POST /api/v1/users/ban`

**Authentication Required:** Yes (Admin or Super Admin only)

**Request Body:**

```json
{
  "user_id": 123,
  "ban_reason": "Violation of community guidelines - Repeated spam",
  "ban_duration": "1_day" // Options: 1_day, 1_week, 1_month, permanent
}
```

**Validation Rules:**

- `user_id`: Required
- `ban_reason`: Required, minimum 10 characters, maximum 500 characters
- `ban_duration`: Required, must be one of: 1_day, 1_week, 1_month, permanent

**Success Response:**

```json
{
  "status": "success",
  "data": {
    "message": "User banned successfully"
  }
}
```

**Error Responses:**

_Invalid Request Body (400 Bad Request)_

```json
{
  "status": "error",
  "error": {
    "code": "validation_error",
    "message": "Validation failed"
  }
}
```

_Unauthorized (401 Unauthorized)_

```json
{
  "status": "error",
  "error": {
    "code": "unauthorized",
    "message": "You are not authorized to ban users"
  }
}
```

_Forbidden (403 Forbidden)_

```json
{
  "status": "error",
  "error": {
    "code": "forbidden",
    "message": "You cannot ban this user"
  }
}
```

_User Not Found (404 Not Found)_

```json
{
  "status": "error",
  "error": {
    "code": "not_found",
    "message": "User not found"
  }
}
```

**Notes:**

- Only Admin and Super Admin can ban users
- Admin cannot ban other admins or super admins
- Super admin cannot be banned
- Ban duration options:
  - 1_day: Ban for 24 hours
  - 1_week: Ban for 7 days
  - 1_month: Ban for 30 days
  - permanent: Permanent ban

### ❄️ Freeze Account

**Endpoint:** `POST /api/v1/users/freeze`

**Authentication Required:** Yes

**Request Body:**

```json
{
  "user_id": 123,
  "freeze_reason": "User requested account freeze"
}
```

**Validation Rules:**

- `user_id`: Required
- `freeze_reason`: Required

**Success Response:**

```json
{
  "status": "success",
  "data": {
    "message": "Account has been frozen successfully"
  }
}
```

**Error Responses:**

_Invalid Request Body (400 Bad Request)_

```json
{
  "status": "error",
  "error": {
    "code": "validation_error",
    "message": "Validation failed"
  }
}
```

_Unauthorized (401 Unauthorized)_

```json
{
  "status": "error",
  "error": {
    "code": "unauthorized",
    "message": "You are not authorized to perform this action"
  }
}
```

_Forbidden (403 Forbidden)_

```json
{
  "status": "error",
  "error": {
    "code": "forbidden",
    "message": "You don't have permission to freeze this account"
  }
}
```

_User Not Found (404 Not Found)_

```json
{
  "status": "error",
  "error": {
    "code": "not_found",
    "message": "User not found"
  }
}
```

**Notes:**

- Users can freeze their own accounts
- Admin cannot freeze other admin accounts
- Super admin accounts cannot be frozen
- Freeze reason is required for all users
- Frozen accounts cannot be accessed until unfrozen by an admin

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

## 👤 User Management Endpoints

### 🗑️ Delete Account

```http
DELETE /api/v1/users/:id
```

**Request Headers:**

```
Authorization: Bearer <token>
```

**Success Response (200):**

```json
{
  "status": "success",
  "data": {
    "message": "Account has been soft deleted successfully"
  }
}
```

**Error Responses:**

```json
{
  "status": "error",
  "error": {
    "code": "unauthorized",
    "message": "You can only delete your own account or be a SUPER_ADMIN to delete others"
  }
}
```

```json
{
  "status": "error",
  "error": {
    "code": "forbidden",
    "message": "SUPER_ADMIN accounts cannot be deleted"
  }
}
```

## 📝 Önemli Notlar

- Tüm endpoint'ler için başarılı yanıtlar `status: "success"` içerir
- Tüm hatalar `status: "error"` içerir
- Protected endpoint'ler için `Authorization` header'ı gereklidir
- Dondurulmuş veya silinmiş hesapların username ve email'leri yeni kayıtlar için kullanılabilir
- Banlanmış hesapların username ve email'leri korunur
- SUPER_ADMIN hesapları silinemez veya dondurulamaz
- Her kullanıcı kendi hesabını silebilir
- SUPER_ADMIN tüm hesapları yönetebilir

## 👑 Admin Endpoints

### 🔑 Admin Login

**Endpoint:** `POST /api/v1/admin/login`

**Request Body:**

```json
{
  "identifier": "string", // email or username
  "password": "string"
}
```

**Success Response (200):**

```json
{
  "status": "success",
  "data": {
    "token": "string",
    "user": {
      "id": "number",
      "username": "string",
      "email": "string",
      "status": "string",
      "role": "string",
      "avatar": "string",
      "created_at": "string"
    }
  }
}
```

**Error Responses:**

```json
{
  "status": "error",
  "error": {
    "code": "validation_error | unauthorized",
    "message": "Error message"
  }
}
```

### 👤 Get Admin Profile

**Endpoint:** `GET /api/v1/admin/me`

**Headers:**

```
Authorization: Bearer <token>
```

**Success Response (200):**

```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "number",
      "username": "string",
      "email": "string",
      "status": "string",
      "role": "string",
      "avatar": "string",
      "created_at": "string"
    }
  }
}
```

**Error Responses:**

```json
{
  "status": "error",
  "error": {
    "code": "unauthorized | forbidden | internal_error",
    "message": "Error message"
  }
}
```

**Notes:**

- Only users with `ADMIN` or `SUPER_ADMIN` role can access admin endpoints
- Admin account must be `active` to access admin endpoints
- `passive`, `frozen`, or soft deleted accounts cannot access admin endpoints

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

## 👤 User Management Endpoints

### 🗑️ Delete Account

```http
DELETE /api/v1/users/:id
```

**Request Headers:**

```
Authorization: Bearer <token>
```

**Success Response (200):**

```json
{
  "status": "success",
  "data": {
    "message": "Account has been soft deleted successfully"
  }
}
```

**Error Responses:**

```json
{
  "status": "error",
  "error": {
    "code": "unauthorized",
    "message": "You can only delete your own account or be a SUPER_ADMIN to delete others"
  }
}
```

```json
{
  "status": "error",
  "error": {
    "code": "forbidden",
    "message": "SUPER_ADMIN accounts cannot be deleted"
  }
}
```

## 📝 Önemli Notlar

- Tüm endpoint'ler için başarılı yanıtlar `status: "success"` içerir
- Tüm hatalar `status: "error"` içerir
- Protected endpoint'ler için `Authorization` header'ı gereklidir
- Dondurulmuş veya silinmiş hesapların username ve email'leri yeni kayıtlar için kullanılabilir
- Banlanmış hesapların username ve email'leri korunur
- SUPER_ADMIN hesapları silinemez veya dondurulamaz
- Her kullanıcı kendi hesabını silebilir
- SUPER_ADMIN tüm hesapları yönetebilir

## 👑 Admin Endpoints

### 🔑 Admin Login

**Endpoint:** `POST /api/v1/admin/login`

**Request Body:**

```json
{
  "identifier": "string", // email or username
  "password": "string"
}
```

**Success Response (200):**

```json
{
  "status": "success",
  "data": {
    "token": "string",
    "user": {
      "id": "number",
      "username": "string",
      "email": "string",
      "status": "string",
      "role": "string",
      "avatar": "string",
      "created_at": "string"
    }
  }
}
```

**Error Responses:**

```json
{
  "status": "error",
  "error": {
    "code": "validation_error | unauthorized",
    "message": "Error message"
  }
}
```

### 👤 Get Admin Profile

**Endpoint:** `GET /api/v1/admin/me`

**Headers:**

```
Authorization: Bearer <token>
```

**Success Response (200):**

```json
{
  "status": "success",
  "data": {
    "user": {
      "id": "number",
      "username": "string",
      "email": "string",
      "status": "string",
      "role": "string",
      "avatar": "string",
      "created_at": "string"
    }
  }
}
```

**Error Responses:**

```json
{
  "status": "error",
  "error": {
    "code": "unauthorized | forbidden | internal_error",
    "message": "Error message"
  }
}
```

**Notes:**

- Only users with `ADMIN` or `SUPER_ADMIN` role can access admin endpoints
- Admin account must be `active` to access admin endpoints
- `passive`, `frozen`, or soft deleted accounts cannot access admin endpoints

## 🎯 Rate Limiting 🚦

- Anonymous: 100 requests per minute
- Authenticated: 1000 requests per minute

## 📋 Data Types

### User Status

- `active`
- `passive`
- `banned`
- `frozen`

### User Roles

- `USER`
- `EDITOR`
- `ADMIN`
- `SUPER_ADMIN`

## 🔄 Recent Changes

- ✨ Enhanced error messages and logging
- 🔒 Improved authentication flow
- 👤 Admin seeding from environment variables
- 🎯 Optimized database operations
- 📝 Better registration process
