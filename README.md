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

- Docker và Docker Compose
- Không cần cài đặt Go ở local (chạy trong Docker)

## Setup với Docker

### Development Mode (với Hot Reload)

**Khuyến nghị cho development** - Tự động rebuild khi có thay đổi code:

```bash
# Khởi động với development mode
docker compose -f docker-compose.dev.yml up

# Hoặc chạy background
docker compose -f docker-compose.dev.yml up -d
```

Tính năng:
- ✅ Hot reload tự động khi bạn sửa code
- ✅ Source code được mount vào container (không cần rebuild)
- ✅ Sử dụng Air để watch file changes
- ✅ Logs hiển thị real-time

### Production Mode

**Cho production hoặc test nhanh** - Build binary và chạy:

```bash
# Khởi động production mode
docker compose up -d
```

Lệnh này sẽ:
- Build Docker image cho backend (Go latest)
- Khởi động MySQL container
- Khởi động Backend container
- Tự động chạy database migrations

### Kiểm tra services đã chạy

```bash
docker compose ps
# hoặc
docker compose -f docker-compose.dev.yml ps
```

Bạn sẽ thấy containers:
- `news_mysql` / `news_mysql_dev` - MySQL database
- `news_backend` / `news_backend_dev` - Go API server

### Xem logs

```bash
# Development mode
docker compose -f docker-compose.dev.yml logs -f backend

# Production mode
docker compose logs -f backend

# MySQL logs
docker compose logs -f mysql
```

### Server đã sẵn sàng

Backend server sẽ chạy tại `http://localhost:8080`

## Environment Variables

Environment variables được định nghĩa trong `docker-compose.yml`. Bạn có thể override bằng cách tạo file `.env` hoặc sửa trực tiếp trong `docker-compose.yml`:

```yaml
environment:
  DB_HOST: mysql
  DB_PORT: 3306
  DB_USER: news_user
  DB_PASSWORD: news_password
  DB_NAME: news_db
  JWT_SECRET: your-secret-key-change-this-in-production
  PORT: 8080
```

## Testing

### Chạy Unit Tests trong Docker

**Development mode:**
```bash
# Chạy tất cả tests (với -mod=mod để ignore vendor)
docker compose -f docker-compose.dev.yml run --rm backend go test -mod=mod ./... -v

# Chạy tests cho một package cụ thể
docker compose -f docker-compose.dev.yml run --rm backend go test -mod=mod ./utils/... -v
docker compose -f docker-compose.dev.yml run --rm backend go test -mod=mod ./middlewares/... -v
docker compose -f docker-compose.dev.yml run --rm backend go test -mod=mod ./services/... -v

# Chạy tests với coverage
docker compose -f docker-compose.dev.yml run --rm backend go test -mod=mod ./... -cover
```

**Production mode:**
```bash
# Chạy tất cả tests
docker compose run --rm backend go test -mod=mod ./... -v
```

### Chạy tests với output đẹp hơn

```bash
# Development mode
docker compose -f docker-compose.dev.yml run --rm backend go test -mod=mod ./... -v -coverprofile=coverage.out

# Production mode
docker compose run --rm backend go test -mod=mod ./... -v -coverprofile=coverage.out
```

### Test Results Summary

Kết quả test hiện tại:
- ✅ **utils**: PASS - Coverage: 86.8%
- ✅ **middlewares**: PASS - Coverage: 93.9%
- ✅ **services**: PASS - Coverage: 0.0% (chỉ có validation tests)

## Docker Commands

### Development Mode

```bash
# Khởi động development (với hot reload)
docker compose -f docker-compose.dev.yml up

# Dừng development
docker compose -f docker-compose.dev.yml down

# Rebuild development image
docker compose -f docker-compose.dev.yml build --no-cache
docker compose -f docker-compose.dev.yml up -d
```

### Production Mode

```bash
# Dừng services
docker compose down

# Dừng và xóa volumes (xóa database data)
docker compose down -v

# Rebuild backend image
docker compose build backend
docker compose up -d

# Rebuild tất cả
docker compose build --no-cache
docker compose up -d
```

### Vào container để debug

```bash
# Development mode
docker compose -f docker-compose.dev.yml exec backend sh

# Production mode
docker compose exec backend sh

# MySQL container
docker compose exec mysql bash
```

## Development với Docker

### Hot Reload trong Development Mode

Khi chạy `docker compose -f docker-compose.dev.yml up`, Air sẽ tự động:
- Watch các file `.go` trong project
- Tự động rebuild khi có thay đổi
- Restart server tự động
- Hiển thị logs real-time

Chỉ cần sửa code và save, server sẽ tự động reload!

### Chạy Go commands trong container

**Development mode:**
```bash
# Chạy go mod tidy
docker compose -f docker-compose.dev.yml run --rm backend go mod tidy

# Chạy go build
docker compose -f docker-compose.dev.yml run --rm backend go build -o news main.go
```

**Production mode:**
```bash
# Chạy go mod tidy
docker compose run --rm backend go mod tidy
```

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

### Articles

- `GET /api/articles` - Lấy danh sách articles (query params: tag, author, favorited, limit, offset)
- `GET /api/articles/feed` - Lấy articles từ users đang follow (cần auth)
- `GET /api/articles/:slug` - Lấy article theo slug
- `POST /api/articles` - Tạo article mới (cần auth)
- `PUT /api/articles/:slug` - Cập nhật article (cần auth)
- `DELETE /api/articles/:slug` - Xóa article (cần auth)
- `POST /api/articles/:slug/favorite` - Favorite article (cần auth)
- `DELETE /api/articles/:slug/favorite` - Unfavorite article (cần auth)

### Comments

- `POST /api/articles/:slug/comments` - Thêm comment vào article (cần auth)
- `GET /api/articles/:slug/comments` - Lấy danh sách comments của article
- `DELETE /api/articles/:slug/comments/:id` - Xóa comment (cần auth, chỉ author mới xóa được)

### Tags

- `GET /api/tags` - Lấy tất cả tags (không cần auth)

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

### Tạo article

```bash
curl -X POST http://localhost:8080/api/articles \
  -H "Authorization: Token YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "article": {
      "title": "How to build REST APIs",
      "description": "Learn how to build REST APIs with Go and Gin",
      "body": "This is a comprehensive guide...",
      "tagList": ["go", "gin", "rest-api"]
    }
  }'
```

### Lấy danh sách articles

```bash
# Lấy tất cả articles
curl -X GET http://localhost:8080/api/articles

# Lọc theo tag
curl -X GET "http://localhost:8080/api/articles?tag=go"

# Lọc theo author
curl -X GET "http://localhost:8080/api/articles?author=johndoe"

# Pagination
curl -X GET "http://localhost:8080/api/articles?limit=10&offset=0"
```

### Lấy article theo slug

```bash
curl -X GET http://localhost:8080/api/articles/how-to-build-rest-apis
```

### Favorite article

```bash
curl -X POST http://localhost:8080/api/articles/how-to-build-rest-apis/favorite \
  -H "Authorization: Token YOUR_JWT_TOKEN"
```

### Feed articles (từ users đang follow)

```bash
curl -X GET http://localhost:8080/api/articles/feed \
  -H "Authorization: Token YOUR_JWT_TOKEN"
```

### Thêm comment vào article

```bash
curl -X POST http://localhost:8080/api/articles/how-to-build-rest-apis/comments \
  -H "Authorization: Token YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "comment": {
      "body": "Great article! Thanks for sharing."
    }
  }'
```

### Lấy danh sách comments của article

```bash
curl -X GET http://localhost:8080/api/articles/how-to-build-rest-apis/comments
```

### Xóa comment

```bash
curl -X DELETE http://localhost:8080/api/articles/how-to-build-rest-apis/comments/1 \
  -H "Authorization: Token YOUR_JWT_TOKEN"
```

### Lấy tất cả tags

```bash
curl -X GET http://localhost:8080/api/tags
```

## Database Schema

Database schema được định nghĩa trong `database/migrations.sql` và sẽ tự động chạy khi khởi động MySQL container lần đầu.

Các bảng chính:
- `users` - Thông tin người dùng
- `follows` - Quan hệ follow giữa users
- `articles` - Bài viết
- `comments` - Comment trên bài viết
- `tags` - Tags
- `article_tags` - Quan hệ many-to-many giữa articles và tags
- `favorites` - User favorite article

## Development Notes

- Code được viết đơn giản, dễ hiểu cho người mới học Gin
- Sử dụng raw SQL queries (không dùng ORM)
- JWT token có thời hạn 24 giờ
- Error handling thống nhất theo RealWorld spec format
- Slug tự động generate từ title và tự động update khi title thay đổi
- Pagination mặc định: limit=20, max limit=100
- Tất cả được chạy trong Docker, không cần cài Go ở local

## Dockerfile

### Dockerfile (Production)
- Multi-stage build:
  - Stage 1: Build Go application với golang:latest (Go >= 1.24)
  - Stage 2: Runtime image với alpine:latest (nhỏ gọn, ~10MB)
- Build binary và chạy trực tiếp

### Dockerfile.dev (Development)
- Sử dụng golang:latest image
- Cài đặt Air cho hot reload
- Mount source code vào container
- Tự động rebuild khi có thay đổi

## Troubleshooting

### Backend không kết nối được MySQL

Kiểm tra:
1. MySQL container đã chạy: `docker compose ps`
2. Backend đợi MySQL healthy: `docker compose logs backend`
3. DB_HOST trong docker-compose.yml phải là `mysql` (tên service)

### Port 8080 đã được sử dụng

Sửa port trong `docker-compose.yml`:
```yaml
ports:
  - "8081:8080"  # Thay 8081 bằng port bạn muốn
```

### Rebuild sau khi thay đổi code

**Development mode:** Không cần rebuild, chỉ cần save file và Air sẽ tự động reload!

**Production mode:**
```bash
docker compose build backend
docker compose up -d
```

### Hot reload không hoạt động

Kiểm tra:
1. Đang chạy development mode: `docker compose -f docker-compose.dev.yml up`
2. Source code đã được mount: `docker compose -f docker-compose.dev.yml exec backend ls -la /app`
3. Air đang chạy: `docker compose -f docker-compose.dev.yml logs backend`

