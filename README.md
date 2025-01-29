# Answer Backend

🚀 Stack Overflow'un Türkçe versiyonu olan soru-cevap platformunun backend projesi.

## 🛠 Teknolojiler

- Go
- PostgreSQL
- GORM
- Gin Framework
- JWT Authentication

## 🔥 Özellikler

### 👤 Kullanıcı Yönetimi

- ✨ Kayıt ve Giriş
- 🔒 JWT bazlı kimlik doğrulama
- 👑 Rol bazlı yetkilendirme (USER, EDITOR, ADMIN, SUPER_ADMIN)
- 🚫 Hesap dondurma ve banlama sistemi
- 🗑️ Soft delete desteği
- 🔄 Username ve email yeniden kullanım sistemi

### 🔐 Güvenlik

- 🔒 Şifre hashleme (bcrypt)
- 🛡️ CORS koruması
- 🔑 JWT token bazlı kimlik doğrulama
- 👮 Rol bazlı yetkilendirme

### 💾 Veritabanı

- 📊 PostgreSQL
- 🔄 GORM ORM
- 📈 Migration sistemi
- 🏷️ Özel index'ler ve constraint'ler

## 🚀 Kurulum

1. Repoyu klonlayın:

```bash
git clone https://github.com/anilsoylu/answer-backend.git
```

2. Gerekli paketleri yükleyin:

```bash
go mod download
```

3. `.env` dosyasını oluşturun:

```bash
cp .env.example .env
```

4. Migration'ları çalıştırın:

```bash
migrate -path internal/database/migrations -database "postgresql://user:password@localhost:5432/dbname?sslmode=disable" up
```

5. Uygulamayı başlatın:

```bash
go run cmd/api/main.go
```

## 📝 Önemli Notlar

- Dondurulmuş veya silinmiş hesapların username ve email'leri yeni kayıtlar için kullanılabilir
- Banlanmış hesapların username ve email'leri korunur
- SUPER_ADMIN hesapları silinemez veya dondurulamaz
- Her kullanıcı kendi hesabını silebilir
- SUPER_ADMIN tüm hesapları yönetebilir

## 🤝 Katkıda Bulunma

1. Fork yapın
2. Feature branch oluşturun (`git checkout -b feature/amazing-feature`)
3. Değişikliklerinizi commit edin (`git commit -m 'feat: add amazing feature'`)
4. Branch'inizi push edin (`git push origin feature/amazing-feature`)
5. Pull Request oluşturun

## 📄 Lisans

MIT License - daha fazla detay için [LICENSE](LICENSE) dosyasına bakın.

## Contact 📧

Anil Soylu - [@anilsoylu](https://github.com/anilsoylu)

Project Link: [https://github.com/anilsoylu/answer-backend](https://github.com/anilsoylu/answer-backend)
