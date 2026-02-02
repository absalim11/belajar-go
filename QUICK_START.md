# Quick Start Guide - POS API Belajar Go

Panduan cepat untuk menjalankan aplikasi POS API.

## Prerequisites
- Go 1.22 atau lebih tinggi
- Docker & Docker Compose
- curl atau Postman (untuk testing)

## Langkah 1: Setup Database

```bash
# Jalankan PostgreSQL dengan Docker
docker-compose up -d

# Verifikasi database running
docker-compose ps
# Output: belajar-go-postgres   Up (healthy)
```

## Langkah 2: Jalankan Aplikasi

```bash
# Install dependencies (jika belum)
go mod tidy

# Jalankan aplikasi
go run main.go

# Output: Server starting on http://localhost:8080
```

## Langkah 3: Test API

### Health Check
```bash
curl http://localhost:8080/health
```

### Get All Products
```bash
curl http://localhost:8080/products
```

### Get All Categories
```bash
curl http://localhost:8080/categories
```

### Create Product
```bash
curl -X POST http://localhost:8080/products \
  -H "Content-Type: application/json" \
  -d '{
    "nama": "Aqua 600ml",
    "harga": 3000,
    "stok": 100,
    "category_id": 2
  }'
```

## Testing dengan Postman

1. Buka Postman
2. Import file: `POS_API_Collection.postman_collection.json`
3. Collection berisi semua endpoint siap pakai
4. Variable `{{base_url}}` = `http://localhost:8080`

## Endpoint Summary

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check |
| GET | `/categories` | List semua kategori |
| GET | `/categories/{id}` | Detail kategori |
| POST | `/categories` | Buat kategori baru |
| PUT | `/categories/{id}` | Update kategori |
| DELETE | `/categories/{id}` | Hapus kategori |
| GET | `/products` | List semua produk |
| GET | `/products/{id}` | Detail produk |
| POST | `/products` | Buat produk baru |
| PUT | `/products/{id}` | Update produk |
| DELETE | `/products/{id}` | Hapus produk |

## Troubleshooting

### Port 8080 sudah digunakan
Ubah port di `.env`:
```env
SERVER_PORT=8081
```

### Database connection error
1. Pastikan Docker running: `docker-compose ps`
2. Check logs: `docker-compose logs postgres`
3. Restart: `docker-compose restart`

### Reset database ke kondisi awal
```bash
docker-compose down -v
docker-compose up -d
# Database akan di-reinitialize dengan data awal dari init.sql
```

## Dokumentasi Lengkap

- **README.md** - Overview dan cara menjalankan
- **DATABASE.md** - Dokumentasi database lengkap (schema, commands, troubleshooting)
- **DEPLOYMENT.md** - Panduan deployment ke VPS
- **POS_API_Collection.postman_collection.json** - Postman collection

## Stop Aplikasi

```bash
# Stop aplikasi Go (Ctrl+C di terminal yang menjalankan go run)

# Stop database
docker-compose stop

# Stop dan hapus container (data tetap ada)
docker-compose down

# Stop dan hapus semua termasuk data
docker-compose down -v
```

## Happy Coding! 🚀
