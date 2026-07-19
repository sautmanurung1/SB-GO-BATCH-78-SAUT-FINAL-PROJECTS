# Sistem Manajemen Portofolio Saham & Dividen

Backend RESTful API untuk Sistem Manajemen Portofolio Saham, dibangun dengan arsitektur Clean Architecture menggunakan Golang, framework Gin, dan PostgreSQL.

## Fitur Utama

- **Autentikasi & Otorisasi:** JWT (JSON Web Token) dengan Role-Based Access Control (RBAC) untuk `admin` dan `investor`.
- **Manajemen Data Master (Admin):** CRUD untuk Sektor dan Saham.
- **Pencatatan Transaksi (Investor):** 
  - Input transaksi (BUY/SELL) manual.
  - **Import PDF Trade Confirmation:** Parsing transaksi otomatis dari file PDF (mendukung deteksi Ticker, Tipe, Lot, Harga, dan Tanggal).
- **Kalkulasi Portofolio Otomatis:** Agregasi portofolio dengan validasi saldo lot sebelum transaksi SELL (mencegah *short selling* yang tidak valid).
- **Auto-Migration:** Sinkronisasi skema database secara otomatis pada saat runtime menggunakan `sql-migrate`.

## Prasyarat Lingkungan (Environment)

Pastikan dependensi berikut telah terinstal dan berjalan:
- **Golang** (versi 1.18+)
- **PostgreSQL** (versi 12+)

Konfigurasi database (Default, dapat di-override dengan *Environment Variables*):
- `DB_HOST`: localhost
- `DB_PORT`: 5432
- `DB_USER`: postgres
- `DB_PASSWORD`: postgres
- `DB_NAME`: portofolio_saham (Buat database ini sebelum menjalankan aplikasi)

## Cara Menjalankan (Penggunaan)

1. **Buat Database**
   Pastikan PostgreSQL berjalan, lalu buat database dengan nama `portofolio_saham`.

2. **Download Dependencies**
   ```bash
   go mod tidy
   ```

3. **Jalankan Aplikasi**
   Aplikasi akan secara otomatis menjalankan proses *auto-migration* pada saat pertama kali berjalan dan menyinkronkan skema pada tabel database.
   ```bash
   go run main.go
   ```
   Aplikasi berjalan secara default di port `:8080`.

## Daftar Path API Tersedia (Routes)

Semua endpoint berbasis pada base URL: `http://localhost:8080/api`

### 1. Autentikasi (Public)
| Method | Endpoint | Deskripsi |
| --- | --- | --- |
| `POST` | `/api/auth/register` | Pendaftaran user baru. |
| `POST` | `/api/auth/login` | Login untuk mendapatkan token JWT. |

### 2. Master Data (Role: `admin`)
Memerlukan header: `Authorization: Bearer <token_jwt>`

| Method | Endpoint | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/sectors` | Mendapatkan daftar sektor. |
| `POST` | `/api/sectors` | Menambah sektor baru. |
| `PUT` | `/api/sectors/:id` | Mengubah data sektor. |
| `DELETE` | `/api/sectors/:id` | Menghapus data sektor. |
| `GET` | `/api/stocks` | Mendapatkan daftar saham. |
| `POST` | `/api/stocks` | Menambah saham baru. |
| `PUT` | `/api/stocks/:id` | Mengubah data saham. |
| `DELETE` | `/api/stocks/:id` | Menghapus data saham. |

### 3. Transaksi & Portofolio (Role: `investor`)
Memerlukan header: `Authorization: Bearer <token_jwt>`

| Method | Endpoint | Deskripsi |
| --- | --- | --- |
| `GET` | `/api/transactions` | Mendapatkan riwayat transaksi pengguna login. |
| `POST` | `/api/transactions` | Mencatat transaksi manual (Payload JSON). |
| `POST` | `/api/transactions/import-pdf` | Import transaksi dari PDF (Multipart Form-Data, key: `file`). |
| `GET` | `/api/portfolio` | Mendapatkan agregasi portofolio pengguna saat ini. |

