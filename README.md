# Answer Backend

🚀 Backend project for a Q&A platform - Turkish version of Stack Overflow.

## 🛠 Technologies

- Go
- PostgreSQL
- GORM
- Gin Framework
- JWT Authentication

## 🔥 Features

### 👤 User Management

- ✨ Registration and Login
- 🔒 JWT-based Authentication
- 👑 Role-based Authorization (USER, EDITOR, ADMIN, SUPER_ADMIN)
- 🚫 Account freezing and banning system
- 🗑️ Soft delete support
- 🔄 Username and email reuse system

### 🔐 Security

- 🔒 Password hashing (bcrypt)
- 🛡️ CORS protection
- 🔑 JWT token-based authentication
- 👮 Role-based authorization

### 💾 Database

- 📊 PostgreSQL
- 🔄 GORM ORM
- 📈 Migration system
- 🏷️ Custom indexes and constraints

## 🚀 Installation

1. Clone the repository:

```bash
git clone https://github.com/anilsoylu/answer-backend.git
```

2. Install required packages:

```bash
go mod download
```

3. Create `.env` file:

```bash
cp .env.example .env
```

4. Run migrations:

```bash
migrate -path internal/database/migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up
```

5. Start the application:

```bash
go run cmd/api/main.go
```

## 📝 Important Notes

- Usernames and emails from frozen or deleted accounts can be used for new registrations
- Usernames and emails from banned accounts are protected
- SUPER_ADMIN accounts cannot be deleted or frozen
- Users can delete their own accounts
- SUPER_ADMIN can manage all accounts

## 🤝 Contributing

1. Fork the project
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) for more details.

## Contact 📧

Anil Soylu - [@anilsoylu](https://github.com/anilsoylu)

Project Link: [https://github.com/anilsoylu/answer-backend](https://github.com/anilsoylu/answer-backend)
