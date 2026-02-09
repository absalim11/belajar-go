# Proyek Go POS API

## Ikhtisar
Proyek ini adalah latihan pembelajaran untuk membangun API RESTful menggunakan bahasa pemrograman Go. Fokus utamanya adalah menggunakan standar library Go tanpa dependensi eksternal untuk memahami mekanisme inti dari HTTP server, routing, dan penanganan data di Go.

**📖 [Quick Start Guide](QUICK_START.md)** - Panduan cepat untuk langsung mulai coding!

## Struktur Proyek
- `main.go`: Entry point aplikasi, setup router dan dependency injection
- `internal/`: Package internal untuk business logic
  - `category/`: Module untuk kategori (entity, handler, service, repository)
  - `product/`: Module untuk produk (entity, handler, service, repository)
  - `transaction/`: Module untuk transaksi dan checkout (entity, handler, service, repository) ⭐ NEW
- `pkg/`: Package yang bisa digunakan ulang
  - `database/`: Database connection configuration
  - `response/`: Standard API response format
- `go.mod`: Definisi modul Go
- `README.md`: Dokumentasi utama proyek
- `QUICK_START.md`: Panduan cepat untuk memulai
- `DATABASE.md`: Dokumentasi lengkap database schema dan commands
- `DEPLOYMENT.md`: Panduan instalasi dan deployment ke VPS
- `BOOTCAMP_SESSION_3.md`: Dokumentasi implementasi Session 3 (Search, Transaction, Report) ⭐ NEW
- `Dockerfile`: Konfigurasi container Docker
- `docker-compose.yml`: Konfigurasi Docker Compose untuk PostgreSQL
- `init.sql`: Script inisialisasi database (tabel, index, data awal)
- `.env`: File konfigurasi environment variables (jangan di-commit ke git)
- `.env.example`: Template file environment variables
- `POS_API_Collection.postman_collection.json`: Postman collection untuk testing API (Updated with Session 3)

---

## Tugas 1: API Kategori Dasar (CRUD)
Pada tugas ini, kita telah mengimplementasikan API CRUD (Create, Read, Update, Delete) dasar untuk entitas "Category".

### Fitur
- **Penyimpanan In-Memory**: Data disimpan dalam slice dengan akses thread-safe menggunakan `sync.Mutex`.
- **Hanya Standard Library**: Tidak menggunakan framework eksternal (seperti Gin atau Echo).
- **Routing Go 1.22**: Menggunakan `http.ServeMux` yang telah ditingkatkan untuk pola rute yang lebih bersih.

### Endpoint
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/categories` | Menampilkan semua kategori |
| `POST` | `/categories` | Membuat kategori baru |
| `GET` | `/categories/{id}` | Mendapatkan detail kategori berdasarkan ID |
| `PUT` | `/categories/{id}` | Memperbarui kategori berdasarkan ID |
| `DELETE` | `/categories/{id}` | Menghapus kategori berdasarkan ID |

---

## Tugas 2: API Produk (CRUD)
Pada tugas ini, kita mengimplementasikan API CRUD untuk entitas "Product" dengan data awal.

### Endpoint Produk
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/products` | Menampilkan semua produk |
| `GET` | `/products?name={keyword}` | Mencari produk berdasarkan nama (case-insensitive) ⭐ NEW |
| `POST` | `/products` | Membuat produk baru |
| `GET` | `/products/{id}` | Mendapatkan detail produk berdasarkan ID |
| `PUT` | `/products/{id}` | Memperbarui produk berdasarkan ID |
| `DELETE` | `/products/{id}` | Menghapus produk berdasarkan ID |

---

## Tugas 3: Sistem Transaksi & Reporting

### Fitur Baru
- **Search by Name**: Pencarian produk dengan filter nama (ILIKE)
- **Transaction System**: Checkout dengan database transaction untuk atomicity
- **Stock Management**: Automatic stock reduction & validation
- **Sales Report**: Laporan penjualan harian dengan produk terlaris

### Endpoint Transaksi
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `POST` | `/api/checkout` | Checkout transaksi dengan multiple items |

### Endpoint Report
| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/api/report/hari-ini` | Laporan penjualan hari ini (revenue, transaksi, produk terlaris) |

---

## Health Check
Endpoint untuk mengecek apakah server berjalan dengan baik.

| Method | Endpoint | Deskripsi |
|--------|----------|-----------|
| `GET` | `/health` | Mengecek status server |

---

## Setup Database dengan Docker

### Prasyarat
- Docker dan Docker Compose telah terinstal
- Port 5433 tidak digunakan oleh aplikasi lain

### Langkah-langkah Setup Database

1. **Jalankan PostgreSQL dengan Docker Compose**
   ```bash
   docker-compose up -d
   ```

   Perintah ini akan:
   - Membuat container PostgreSQL 16 Alpine
   - Menjalankan script `init.sql` untuk membuat tabel dan data awal
   - Expose port 5433 (mapped ke 5432 di container)
   - Membuat volume persistent untuk data PostgreSQL

2. **Cek Status Container**
   ```bash
   docker-compose ps
   ```

   Output yang diharapkan:
   ```
   Name                      Command                 State                        Ports
   belajar-go-postgres   docker-entrypoint.sh postgres   Up (healthy)   0.0.0.0:5433->5432/tcp
   ```

3. **Lihat Logs Database (Optional)**
   ```bash
   docker-compose logs -f postgres
   ```

4. **Koneksi ke Database (Optional)**
   ```bash
   docker exec -it belajar-go-postgres psql -U postgres -d belajar_go
   ```

### Konfigurasi Database

File `.env` sudah dikonfigurasi untuk koneksi ke database Docker:
```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=belajar_go
DB_SSLMODE=disable
SERVER_PORT=8080
```

### Perintah Docker Compose Berguna

**Stop Database:**
```bash
docker-compose stop
```

**Start Database:**
```bash
docker-compose start
```

**Restart Database:**
```bash
docker-compose restart
```

**Stop dan Hapus Container (Data tetap tersimpan di volume):**
```bash
docker-compose down
```

**Hapus Container dan Volume (Menghapus semua data):**
```bash
docker-compose down -v
```

**Rebuild dan Restart:**
```bash
docker-compose down
docker-compose up -d
```

### Struktur Database

Database `belajar_go` memiliki tabel:
- **categories**: Menyimpan data kategori produk
  - id (PRIMARY KEY)
  - name
  - description
  - created_at
  - updated_at

- **products**: Menyimpan data produk
  - id (PRIMARY KEY)
  - nama
  - harga
  - stok
  - category_id (FOREIGN KEY ke categories)
  - created_at
  - updated_at

- **transactions**: Menyimpan data transaksi checkout ⭐ NEW
  - id (PRIMARY KEY)
  - total_amount
  - created_at

- **transaction_details**: Menyimpan detail item per transaksi ⭐ NEW
  - id (PRIMARY KEY)
  - transaction_id (FOREIGN KEY ke transactions)
  - product_id (FOREIGN KEY ke products)
  - quantity
  - subtotal

Data awal akan otomatis dibuat melalui file `init.sql`.

**Untuk informasi lebih detail tentang database, lihat [DATABASE.md](DATABASE.md)**

---

## Cara Menjalankan Aplikasi

### Langkah Lengkap

1. **Pastikan Prasyarat Terpenuhi:**
   - Go versi 1.22 atau lebih tinggi telah terinstal
   - Docker dan Docker Compose telah terinstal

2. **Clone Repository dan Masuk ke Direktori:**
   ```bash
   cd belajar-go
   ```

3. **Jalankan Database PostgreSQL:**
   ```bash
   docker-compose up -d
   ```

4. **Verifikasi Database Berjalan:**
   ```bash
   docker-compose ps
   ```

5. **Install Dependencies Go:**
   ```bash
   go mod tidy
   ```

6. **Jalankan Aplikasi:**
   ```bash
   go run main.go
   ```

7. **Verifikasi Server Berjalan:**
   ```bash
   curl http://localhost:8080/health
   ```

   Output yang diharapkan:
   ```json
   {
     "status": "OK",
     "message": "Server is running smoothly"
   }
   ```

Server akan berjalan di `http://localhost:8080`.

### Testing dengan Postman
1. Import file `POS_API_Collection.postman_collection.json` ke Postman
2. Collection sudah berisi semua endpoint dengan contoh request
3. Variable `{{base_url}}` sudah diset ke `http://localhost:8080`

### Contoh Penggunaan (Tugas 1)
**Membuat Kategori:**
```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"Electronics", "description":"Gadgets"}' http://localhost:8080/categories
```

**Menampilkan Kategori:**
```bash
curl http://localhost:8080/categories
```

---

### Contoh Penggunaan (Tugas 2 - Product API)

**Membuat Produk:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"nama":"Indomie Goreng","harga":3500,"stok":100,"category_id":1}' \
  http://localhost:8080/products
```

**Menampilkan Semua Produk:**
```bash
curl http://localhost:8080/products
```

**Mendapatkan Detail Produk:**
```bash
curl http://localhost:8080/products/1
```

**Memperbarui Produk:**
```bash
curl -X PUT -H "Content-Type: application/json" \
  -d '{"nama":"Indomie Goreng Spesial","harga":4000,"stok":80,"category_id":1}' \
  http://localhost:8080/products/1
```

**Menghapus Produk:**
```bash
curl -X DELETE http://localhost:8080/products/1
```

---

### Contoh Penggunaan (Session 3 - Search, Transaction, Report) ⭐ NEW

**Search Produk by Name:**
```bash
# Search dengan keyword "indo"
curl "http://localhost:8080/products?name=indo"

# Search dengan keyword "vit"
curl "http://localhost:8080/products?name=vit"
```

**Checkout Transaksi - Single Item:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2}
    ]
  }' \
  http://localhost:8080/api/checkout
```

**Checkout Transaksi - Multiple Items:**
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{
    "items": [
      {"product_id": 1, "quantity": 2},
      {"product_id": 2, "quantity": 3},
      {"product_id": 3, "quantity": 1}
    ]
  }' \
  http://localhost:8080/api/checkout
```

**Daily Sales Report:**
```bash
curl http://localhost:8080/api/report/hari-ini
```

**Response Example (Checkout):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "total_amount": 16000,
    "created_at": "2026-02-09T06:03:53.412228Z",
    "details": [
      {
        "id": 1,
        "transaction_id": 1,
        "product_id": 1,
        "product_name": "Indomie Godog",
        "quantity": 2,
        "subtotal": 7000
      },
      {
        "id": 2,
        "transaction_id": 1,
        "product_id": 2,
        "product_name": "Vit 1000ml",
        "quantity": 3,
        "subtotal": 9000
      }
    ]
  }
}
```

**Response Example (Daily Report):**
```json
{
  "success": true,
  "data": {
    "total_revenue": 45500,
    "total_transaksi": 2,
    "produk_terlaris": {
      "nama": "Indomie Godog",
      "qty_terjual": 7
    }
  }
}
```

---

### Contoh Response

**Success Response (GET /products):**
```json
{
  "success": true,
  "data": [
    {
      "id": 1,
      "nama": "Indomie Godog",
      "harga": 3500,
      "stok": 10,
      "category_id": 1,
      "category_name": "Makanan",
      "created_at": "2026-02-01T12:13:24.308214Z",
      "updated_at": "2026-02-01T12:13:24.308214Z"
    }
  ]
}
```

**Error Response:**
```json
{
  "success": false,
  "message": "product not found"
}
```
