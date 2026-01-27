# Panduan Deployment ke VPS

Dokumen ini menjelaskan langkah-langkah untuk mendeploy aplikasi Go Category API ke server VPS (Ubuntu/Debian).

## 1. Persiapan di VPS
Pastikan Go sudah terinstal di VPS. Jika belum, instal menggunakan perintah:
```bash
sudo apt update
sudo apt install golang-go -y
```

## 2. Clone Repositori
Masuk ke VPS dan clone repositori Anda:
```bash
git clone https://github.com/absalim11/belajar-go.git
cd belajar-go
```

## 3. Build Aplikasi
Build aplikasi menjadi file binary agar dapat dijalankan tanpa source code:
```bash
go build -o pos-api main.go
```

## 4. Konfigurasi Systemd (Rekomendasi)
Agar aplikasi tetap berjalan di background dan otomatis restart jika server reboot, gunakan Systemd.

Buat file service baru:
```bash
sudo nano /etc/systemd/system/pos-api.service
```

Masukkan konfigurasi berikut (sesuaikan `User` dan `WorkingDirectory`):
```ini
[Unit]
Description=Go Category API Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/root/belajar-go
ExecStart=/root/belajar-go/pos-api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

## 5. Menjalankan Service
Aktifkan dan jalankan service yang baru dibuat:
```bash
sudo systemctl daemon-reload
sudo systemctl enable pos-api
sudo systemctl start pos-api
```

## 6. Verifikasi
Cek status service:
```bash
sudo systemctl status pos-api
```

Aplikasi sekarang berjalan di `http://IP_VPS_ANDA:8080`.

---

## 7. Deployment Menggunakan Docker (Alternatif)
Jika Anda lebih suka menggunakan container, Anda dapat menggunakan Docker.

### Build Image
```bash
docker build -t pos-api .
```

### Jalankan Container
```bash
docker run -d -p 8080:8080 --name pos-api-container pos-api
```

### Cek Status Container
```bash
docker ps
docker logs -f category-api-container
```

---

## Tips: Menggunakan Nginx sebagai Reverse Proxy (Opsional)
Jika ingin menggunakan domain dan SSL (HTTPS), gunakan Nginx:
1. Instal Nginx: `sudo apt install nginx`
2. Konfigurasi server block untuk meneruskan traffic dari port 80 ke 8080.
