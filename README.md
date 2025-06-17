# Queue Management System

Sistem manajemen antrean berbasis Go menggunakan:
- [Gin](https://gin-gonic.com/) sebagai HTTP framework
- [Asynq](https://github.com/hibiken/asynq) sebagai job/task queue
- [Redis](https://redis.io/) sebagai message broker

## 📦 Fitur

- Menambahkan job ke antrean (`POST /api/queue`)
- Menampilkan semua job yang tersimpan (`GET /api/queue`)
- Melakukan retry pada job yang gagal (`POST /api/queue/:id/retry`)

---

## 🚀 Cara Menjalankan

### 1. Jalankan Redis

Pastikan Redis sudah berjalan secara lokal di `localhost:6379`.  
Jika belum punya, kamu bisa install atau gunakan Docker:

```bash
docker run -p 6379:6379 redis
```

### 2. Jalankan Aplikasi

```bash
go run main.go
```

Server akan berjalan di `http://localhost:8080`.

---

## 📫 API Endpoint

### ➕ Tambah Job

```http
POST /api/queue
Content-Type: application/json

{
  "email": "user@example.com",
  "url": "https://example.com/task"
}
```

### 📄 Lihat Semua Job

```http
GET /api/queue
```

### 🔁 Retry Job Gagal

```http
POST /api/queue/:id/retry
Content-Type: application/json

{
  "email": "user@example.com",
  "url": "https://example.com/task"
}
```

## 🛠️ Struktur Proyek

```
.
├── internal/
│   ├── api/          # Handler HTTP (PostJob, GetJobs, RetryJob)
│   ├── queue/        # Handler asynq untuk background processing
│   ├── store/        # In-memory data store untuk menyimpan job
│   └── logstream/    # Handler WebSocket log streaming
├── model/            # Struct model data (Email payload)
├── main.go           # Entry point aplikasi
```

---

## 📋 Catatan Tambahan

- Gunakan [Asynqmon](https://github.com/hibiken/asynqmon) untuk memonitor job Asynq:

```bash
go install github.com/hibiken/asynqmon@latest
asynqmon --redis-addr=localhost:6379
```

Akses: `http://localhost:8080` (atau port yang digunakan oleh asynqmon)

---
