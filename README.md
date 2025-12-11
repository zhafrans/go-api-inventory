# Inventory Management API

REST API untuk manajemen inventaris barang menggunakan Golang, Fiber, GORM, dan PostgreSQL.

## Fitur

1. **Authentication**

   - Register user baru
   - Login dengan JWT token

2. **Manajemen Item**

   - Create: Tambah item baru
   - Read: Lihat semua item atau detail item
   - Update: Edit informasi item
   - Update Stock: Tambah/kurangi stok barang
   - Delete: Hapus item

3. **Activity Log**
   - Melihat riwayat aktivitas

## Teknologi Stack

- **Golang** 1.20+
- **Fiber** - Web framework
- **GORM** - ORM untuk database
- **PostgreSQL** - Database
- **JWT** - Authentication

## Prerequisites

- Go 1.20 atau lebih baru
- PostgreSQL 15 atau lebih baru
- Git

## Instalasi dan Setup

### 1. Clone Repository

```bash
git clone <repository-url>
cd inventory-api
```

### 2. Jalankan Server
ini sudah termasuk menjalankan seeder

```bash
go run cmd/server/main.go
```

### (opsional : migrate fresh, jika ingin migrate ulang dan seeder ulang)

```bash
go run cmd/scripts/migrate_fresh.go
```
