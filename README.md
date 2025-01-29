# Answer Backend 🚀

## Overview

Answer is a Q&A platform built with Go and PostgreSQL, inspired by Stack Overflow but tailored for Turkish users.

## 🔥 Latest Updates

- ✨ Enhanced error handling and logging system
- 🔒 Improved user authentication flow
- 👤 Admin seeding from environment variables
- 🎯 Optimized database transactions
- 📝 Better logging for user registration

## 🛠 Tech Stack

- Go
- PostgreSQL
- GORM
- Gin
- JWT

## 🚀 Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Make

### Environment Variables

Copy `.env.example` to `.env` and update the values:

```bash
cp .env.example .env
```

Required environment variables for admin seeding:

```env
ADMIN_USERNAME=your_admin_username
ADMIN_PASSWORD=your_admin_password
ADMIN_EMAIL=your_admin_email
ADMIN_ROLE=SUPER_ADMIN
ADMIN_STATUS=active
```

### Installation

```bash
# Clone the repository
git clone https://github.com/yourusername/answer-backend.git

# Navigate to the project directory
cd answer-backend

# Install dependencies
go mod download

# Run migrations
make migrate

# Start the server
make run
```

## 📚 Documentation

- API documentation can be found in [docs/API.md](docs/API.md)
- Database schema can be found in [docs/SCHEMA.md](docs/SCHEMA.md)

## 🧪 Testing

```bash
make test
```

## 🤝 Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
