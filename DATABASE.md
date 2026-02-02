# Database Documentation

## Overview
Proyek ini menggunakan PostgreSQL 16 Alpine yang dijalankan melalui Docker Compose untuk memudahkan setup dan development.

## Quick Start

### 1. Menjalankan Database
```bash
docker-compose up -d
```

### 2. Cek Status Database
```bash
docker-compose ps
```

Expected output:
```
Name                      Command                 State                        Ports
belajar-go-postgres   docker-entrypoint.sh postgres   Up (healthy)   0.0.0.0:5433->5432/tcp
```

### 3. Koneksi ke Database
```bash
# Interactive psql shell (untuk manual exploration)
docker exec -it belajar-go-postgres psql -U postgres -d belajar_go

# Query langsung (non-interactive - cocok untuk automation/script)
docker exec belajar-go-postgres psql -U postgres -d belajar_go -c "SELECT * FROM categories;"
```

## Database Schema

### Table: categories
Menyimpan data kategori produk.

| Column | Type | Constraint | Description |
|--------|------|------------|-------------|
| id | SERIAL | PRIMARY KEY | ID kategori (auto increment) |
| name | VARCHAR(100) | NOT NULL | Nama kategori |
| description | TEXT | - | Deskripsi kategori |
| created_at | TIMESTAMP | DEFAULT NOW() | Waktu pembuatan record |
| updated_at | TIMESTAMP | DEFAULT NOW() | Waktu update record (auto update via trigger) |

**Sample Data:**
```sql
SELECT * FROM categories;

 id |   name   |        description         |         created_at         |         updated_at
----+----------+---------------------------+----------------------------+----------------------------
  1 | Makanan  | Kategori produk makanan   | 2026-02-01 12:13:24.308214 | 2026-02-01 12:13:24.308214
  2 | Minuman  | Kategori produk minuman   | 2026-02-01 12:13:24.308214 | 2026-02-01 12:13:24.308214
  3 | Bumbu    | Kategori produk bumbu dapur| 2026-02-01 12:13:24.308214 | 2026-02-01 12:13:24.308214
```

### Table: products
Menyimpan data produk dengan relasi ke kategori.

| Column | Type | Constraint | Description |
|--------|------|------------|-------------|
| id | SERIAL | PRIMARY KEY | ID produk (auto increment) |
| nama | VARCHAR(100) | NOT NULL | Nama produk |
| harga | INTEGER | NOT NULL | Harga produk (dalam rupiah) |
| stok | INTEGER | NOT NULL, DEFAULT 0 | Jumlah stok produk |
| category_id | INTEGER | FK to categories(id) | ID kategori (nullable) |
| created_at | TIMESTAMP | DEFAULT NOW() | Waktu pembuatan record |
| updated_at | TIMESTAMP | DEFAULT NOW() | Waktu update record (auto update via trigger) |

**Foreign Key:**
- `category_id` references `categories(id)` ON DELETE SET NULL
  - Jika kategori dihapus, `category_id` di produk akan di-set NULL

**Indexes:**
- `idx_products_category_id` pada kolom `category_id` untuk performa JOIN

**Sample Data:**
```sql
SELECT p.*, c.name as category_name
FROM products p
LEFT JOIN categories c ON p.category_id = c.id;

 id |      nama      | harga | stok | category_id | category_name
----+----------------+-------+------+-------------+---------------
  1 | Indomie Godog  |  3500 |   10 |           1 | Makanan
  2 | Vit 1000ml     |  3000 |   40 |           2 | Minuman
  3 | Kecap          | 12000 |   20 |           3 | Bumbu
```

## Database Features

### Auto-Update Timestamp
Database menggunakan trigger untuk otomatis update kolom `updated_at`:

```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_categories_updated_at BEFORE UPDATE ON categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_products_updated_at BEFORE UPDATE ON products
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
```

Setiap kali melakukan UPDATE, kolom `updated_at` akan otomatis diperbarui.

## Database Configuration

### Environment Variables
Konfigurasi database tersimpan di file `.env`:

```env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=belajar_go
DB_SSLMODE=disable
```

### Docker Compose Configuration
File `docker-compose.yml` menggunakan environment variables:

```yaml
environment:
  POSTGRES_USER: ${DB_USER:-postgres}
  POSTGRES_PASSWORD: ${DB_PASSWORD:-postgres}
  POSTGRES_DB: ${DB_NAME:-belajar_go}
ports:
  - "${DB_PORT:-5433}:5432"
```

## Useful Commands

### Docker Commands

**View Logs:**
```bash
# Follow logs
docker-compose logs -f postgres

# View last 100 lines
docker-compose logs --tail=100 postgres
```

**Execute SQL Commands:**
```bash
# Interactive psql shell
docker exec -it belajar-go-postgres psql -U postgres -d belajar_go

# Single command (non-interactive - lebih cocok untuk automation/script)
docker exec belajar-go-postgres psql -U postgres -d belajar_go -c "SELECT * FROM categories;"

# Multiple commands
docker exec belajar-go-postgres psql -U postgres -d belajar_go -c "
  SELECT COUNT(*) as total_products FROM products;
  SELECT COUNT(*) as total_categories FROM categories;
"
```

**Database Backup:**
```bash
# Backup to file
docker exec belajar-go-postgres pg_dump -U postgres belajar_go > backup.sql

# Restore from file
docker exec -i belajar-go-postgres psql -U postgres -d belajar_go < backup.sql
```

**Reset Database:**
```bash
# Stop and remove container + volume
docker-compose down -v

# Start fresh
docker-compose up -d
```

### SQL Commands

**View all tables:**
```sql
\dt
```

**Describe table structure:**
```sql
\d categories
\d products
```

**Count records:**
```sql
SELECT COUNT(*) FROM categories;
SELECT COUNT(*) FROM products;
```

**View products with category:**
```sql
SELECT
    p.id,
    p.nama,
    p.harga,
    p.stok,
    c.name as category_name
FROM products p
LEFT JOIN categories c ON p.category_id = c.id
ORDER BY p.id;
```

## Troubleshooting

### Port 5433 sudah digunakan
Ubah port di file `.env`:
```env
DB_PORT=5434  # atau port lain yang tersedia
```

Kemudian restart:
```bash
docker-compose down
docker-compose up -d
```

### Container tidak bisa start
Check logs untuk error:
```bash
docker-compose logs postgres
```

### Connection refused
Pastikan:
1. Container berjalan: `docker-compose ps`
2. Port mapping benar di `docker-compose.yml`
3. Konfigurasi `.env` sesuai dengan docker-compose

### Data hilang setelah restart
Jika menggunakan `docker-compose down -v`, volume akan terhapus.
Gunakan `docker-compose down` (tanpa -v) untuk keep data.

## Database Management Tools

### pgAdmin
Bisa menggunakan pgAdmin untuk GUI management:

Connection settings:
- Host: localhost
- Port: 5433
- Database: belajar_go
- Username: postgres
- Password: postgres

### DBeaver / TablePlus
Atau tools lain dengan connection string:
```
postgresql://postgres:postgres@localhost:5433/belajar_go?sslmode=disable
```

## Init Script

File `init.sql` akan otomatis dijalankan saat container pertama kali dibuat.
Script ini berisi:
1. CREATE TABLE untuk categories dan products
2. CREATE INDEX untuk performa
3. INSERT sample data
4. CREATE FUNCTION dan TRIGGER untuk auto-update timestamp

Jika ingin update `init.sql`:
1. Edit file `init.sql`
2. Reset database: `docker-compose down -v`
3. Start ulang: `docker-compose up -d`
