# 🚀 Answer Backend

A powerful Q&A platform backend built with Go and PostgreSQL.

## 🛠 Tech Stack

- **Go** - Programming Language
- **Gin** - Web Framework
- **PostgreSQL** - Database
- **JWT** - Authentication
- **Air** - Live Reload
- **Validator** - Request Validation

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher
- Air (for development)

## 🔧 Installation

1. Clone the repository:

```bash
git clone https://github.com/anilsoylu/answer-backend.git
cd answer-backend
```

2. Install dependencies:

```bash
go mod download
```

3. Set up environment variables:

```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Install Air (for development):

```bash
go install github.com/cosmtrek/air@latest
```

## 🚀 Running the Application

### Development

```bash
air
```

### Production

```bash
go run cmd/api/main.go
```

## 📚 API Documentation

API documentation can be found in the [docs/API.md](docs/API.md) file.

## 🔒 Security

- Password hashing using bcrypt
- JWT-based authentication
- Input validation
- CORS protection
- Rate limiting

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
