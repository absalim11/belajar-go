# Proyek API Kategori Go

## Ikhtisar
Proyek ini adalah latihan pembelajaran untuk membangun API RESTful menggunakan bahasa pemrograman Go. Fokus utamanya adalah menggunakan standar library Go tanpa dependensi eksternal untuk memahami mekanisme inti dari HTTP server, routing, dan penanganan data di Go.

## Struktur Proyek
- `main.go`: Berisi seluruh logika aplikasi, termasuk model, penyimpanan, dan handler HTTP.
- `go.mod`: Definisi modul Go.

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

### Cara Menjalankan
1. Pastikan Anda telah menginstal Go (versi 1.22 atau lebih tinggi direkomendasikan).
2. Jalankan aplikasi:
   ```bash
   go run main.go
   ```
3. Server akan berjalan di `http://localhost:8080`.

### Contoh Penggunaan (Tugas 1)
**Membuat Kategori:**
```bash
curl -X POST -H "Content-Type: application/json" -d '{"name":"Electronics", "description":"Gadgets"}' http://localhost:8080/categories
```

**Menampilkan Kategori:**
```bash
curl http://localhost:8080/categories
```
