# RealWorld API - Go + Gin

API backend cho RealWorld application sử dụng Go 1.22 và Gin framework.

## Cấu trúc Project

```
news/
├── config/          # Configuration management
├── controllers/     # HTTP handlers
├── services/        # Business logic
├── repositories/    # Data access layer
├── models/          # Database models
├── dto/             # Data transfer objects
├── middlewares/     # Middleware (auth, error handling)
├── utils/           # Utilities (JWT, password, slug)
├── database/        # Database setup và migrations
└── main.go          # Entry point
```

## Yêu cầu

- Go 1.22
- Docker và Docker Compose
- MySQL (chạy qua docker-compose)

## Setup

### 1. Setup Go PATH (nếu Go chưa có trong PATH)

Nếu Go chưa có trong PATH, thêm vào:

```bash
export PATH=$PATH:/home/pham.ngoc.kha@sun-asterisk.com/sdk/go1.22.0/bin
export GOPATH=$HOME/go
```

Để thêm vĩnh viễn, thêm vào `~/.bashrc` hoặc `~/.zshrc`:
```bash
echo 'export PATH=$PATH:/home/pham.ngoc.kha@sun-asterisk.com/sdk/go1.22.0/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
source ~/.bashrc
```

Hoặc chạy script setup:
```bash
./setup.sh
```

### 2. Kiểm tra Go version

```bash
go version
```

### 3. Cài đặt dependencies

```bash
go mod download
```

Hoặc:
```bash
go mod tidy
```

### 3. Setup environment variables

Tạo file `.env` (hoặc export environment variables):

```bash
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=news_user
export DB_PASSWORD=news_password
export DB_NAME=news_db
export JWT_SECRET=your-secret-key-change-this-in-production
export PORT=8080
```

### 4. Khởi động MySQL với Docker Compose

```bash
docker-compose up -d
```

Kiểm tra MySQL đã chạy:
```bash
docker-compose ps
```

### 5. Chạy server

```bash
go run main.go
```

Server sẽ chạy tại `http://localhost:8080`

## API Endpoints

### Authentication

- `POST /api/users` - Đăng ký user mới
- `POST /api/users/login` - Đăng nhập
- `GET /api/user` - Lấy thông tin user hiện tại (cần auth)
- `PUT /api/user` - Cập nhật thông tin user (cần auth)

### Profiles

- `GET /api/profiles/:username` - Lấy profile của user
- `POST /api/profiles/:username/follow` - Follow user (cần auth)
- `DELETE /api/profiles/:username/follow` - Unfollow user (cần auth)

## Testing API

### Đăng ký user

```bash
curl -X POST http://localhost:8080/api/users \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "username": "johndoe",
      "email": "john@example.com",
      "password": "password123"
    }
  }'
```

### Đăng nhập

```bash
curl -X POST http://localhost:8080/api/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "user": {
      "email": "john@example.com",
      "password": "password123"
    }
  }'
```

### Lấy thông tin user hiện tại

```bash
curl -X GET http://localhost:8080/api/user \
  -H "Authorization: Token YOUR_JWT_TOKEN"
```

### Lấy profile

```bash
curl -X GET http://localhost:8080/api/profiles/johndoe
```

### Follow user

```bash
curl -X POST http://localhost:8080/api/profiles/johndoe/follow \
  -H "Authorization: Token YOUR_JWT_TOKEN"
```

## Database Schema

Database schema được định nghĩa trong `database/migrations.sql` và sẽ tự động chạy khi khởi động MySQL container lần đầu.

Các bảng chính:
- `users` - Thông tin người dùng
- `follows` - Quan hệ follow giữa users
- `articles` - Bài viết (sẽ implement ở phase sau)
- `comments` - Comment trên bài viết (sẽ implement ở phase sau)
- `tags` - Tags (sẽ implement ở phase sau)
- `article_tags` - Quan hệ many-to-many giữa articles và tags
- `favorites` - User favorite article (sẽ implement ở phase sau)

## Development Notes

- Code được viết đơn giản, dễ hiểu cho người mới học Gin
- Sử dụng raw SQL queries (không dùng ORM)
- JWT token có thời hạn 24 giờ
- Error handling thống nhất theo RealWorld spec format

