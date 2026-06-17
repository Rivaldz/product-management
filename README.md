## Cara Menjalankan Menggunakan Docker (Rekomendasi)

### Persiapan

Sudah terinstall **Docker** dan **Docker Compose**.

### Langkah-langkah Menjalankan:

1. Perintah berikut ini dijalankan di root project:
   ```bash
   docker compose up --build
   ```
2. Docker Compose akan secara otomatis melakukan:

   - Membuat container database PostgreSQL (`db`) dan menunggu sampai siap menerima koneksi.
   - Menjalankan migrasi database (`migrate`) secara otomatis dari folder `./migrations` ke database lokal tersebut.
   - Membangun image aplikasi (`app`) dan menjalankannya di port `8080`.

3. Aplikasi siap diakses di:

   - Endpoint API: `http://localhost:8080` (contoh: `http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items`)
   - Swagger Docs: `http://localhost:8080/swagger/index.html`

4. Untuk mematikan aplikasi dan membersihkan container beserta volume database:
   ```bash
   docker compose down -v
   ```

---

## Penjelasan Konsep: Antigravity Persona vs Traditional Multi-Agent

Pada pengerjaan _live test_ saya ini, terdapat penjelasan konsep terkait arsitektur AI _agent_ yang digunakan:

1. **Bukan Multi-Agent Terpisah (Tidak Menggunakan OpenRouter, crewAI, dsb)**:
   Proyek ini dikerjakan di bawah ekosistem **Antigravity IDE**.
2. **Menggunakan Konsep Persona Terintegrasi**:
   Sebagai gantinya, Antigravity menggunakan konsep **Persona** yang didefinisikan secara deklaratif di dalam berkas konfigurasi `.antigravity/rules.json`.
   - **Tech Lead (`tech_lead`)**: Berperan melakukan analisis PRD, merancang arsitektur, dan mengidentifikasi risiko teknis.
   - **Golang Coder (`golang_coder`)**: Berperan melakukan implementasi kode Go yang _idiomatic_ sesuai prinsip _Clean Architecture_.
   - **QA Automation (`qa_test`)**: Berperan melakukan testing API via terminal, membandingkan _expected vs actual_, serta membuat laporan hasil uji (_gap analysis_).
3. **Kelebihan Pendekatan Ini**:
   AI bertindak sebagai satu agen terpadu (_unified agent_) yang memuat aturan perilaku (_workspace rules_), batasan kemampuan (_capabilities_), dan format output yang berbeda secara dinamis. Namun, AI tetap membagi riwayat konteks obrolan (_context sharing_) yang sama. Hal ini menghindari hilangnya informasi (_information loss_) yang sering terjadi pada komunikasi antar _agent_ tradisional, sekaligus menjamin hasil pengerjaan (analisis, kode, pengujian) terstruktur rapi sesuai spesialisasi masing-masing peran.

---

# Panduan Migrasi Database (Untuk Golang Dev)

Dokumen ini berisi panduan untuk menjalankan migrasi database pada project ini, serta cara mengatasi beberapa error umum yang sering terjadi saat melakukan migrasi.

## Persiapan Awal (.env)

Pastikan file `.env` sudah dibuat di root project dan memiliki konfigurasi untuk URL database (`PG_URL`).

```env
PG_URL="postgres://<user>:<password>@<host>:<port>/<dbname>?sslmode=disable"
```

**Penting:** Jika menggunakan server database yang tidak memiliki sertifikat SSL yang valid atau tidak menggunakan koneksi SSL, pastikan untuk menambahkan `?sslmode=disable` di akhir URL.

## Menjalankan Migrasi

Project ini menggunakan [`golang-migrate/migrate`](https://github.com/golang-migrate/migrate) dan `Makefile` untuk mempermudah eksekusi perintah.

- **Menerapkan semua migrasi (Up):**
  ```bash
  make migrate-up
  ```
- **Membatalkan semua migrasi (Down):**
  ```bash
  make migrate-down
  ```

---

## Troubleshooting Error Migrasi

Berikut adalah beberapa error yang mungkin kamu temui dan cara mengatasinya agar bisa lebih mandiri dalam melakukan _debugging_.

### 1. Error: `URL cannot be empty`

**Pesan Error:**

```text
migrate -path migrations -database "" up
error: failed to parse scheme from database URL: URL cannot be empty
```

**Penyebab:**
Nilai `${PG_URL}` di Makefile terbaca kosong. Ini berarti file `.env` tidak ada, atau variabel `PG_URL` di dalam file `.env` belum dideklarasikan/salah ketik.

**Solusi:**
Pastikan file `.env` benar-benar ada di root folder project. Makefile di project ini sudah di-setting untuk memuat file `.env` secara otomatis (menggunakan `include .env`). Periksa juga apakah nama variabelnya sudah benar `PG_URL`.

---

### 2. Error: `SSL is not enabled on the server`

**Pesan Error:**

```text
error: failed to open database: pq: SSL is not enabled on the server
```

**Penyebab:**
Driver `pq` (PostgreSQL) di Golang mencoba menggunakan mode koneksi aman (SSL) secara _default_. Namun, server database menolak koneksi tersebut karena server tidak mengaktifkan koneksi SSL.

**Solusi:**
Tambahkan parameter `?sslmode=disable` di paling belakang koneksi `PG_URL` pada file `.env` kamu.
Contoh:

```env
PG_URL="postgres://rivaldo:root1212@43.133.144.153:5432/temuh_dev?sslmode=disable"
```

---

### 3. Error: `no migration found for version...` (Database State = Dirty)

**Pesan Error:**

```text
error: no migration found for version 20260615100000: read down for version 20260615100000 .: file does not exist
```

_Atau error "Dirty database version ... Fix and force version."_

**Penyebab:**
Hal ini biasa disebut sebagai _Dirty Migration_. Tabel `schema_migrations` di database mencatat versi migrasi terakhir (misal `20260615100000`), tetapi file `.sql` untuk versi tersebut **sudah tidak ada** atau **dihapus** di folder `migrations/` lokal kamu. _golang-migrate_ tidak akan mau jalan jika state database tidak sinkron dengan file lokal.

**Solusi:**
Kamu perlu "memaksa" (force) database agar kembali ke _state_ versi terakhir yang file-nya tersedia di folder lokal.

1. **Cari versi terakhir yang valid.** Lihat folder `migrations/` dan cek angka depan dari file migrasi. Misal jika ada `000001_init_schema.up.sql`, maka versinya adalah `1`.
2. **Jalankan perintah force** sesuai versi yang ingin di-set:
   ```bash
   make migrate-force version=1
   ```
   _(Ganti `1` dengan versi target. Jika ingin mengosongkan log versi, bisa di-set ke versi `0` asalkan sesuai)_
3. Setelah di-force, database sudah tidak berstatus _dirty_. Kamu bisa kembali menjalankan:
   ```bash
   make migrate-up
   ```

> **Tips Best Practice:**
>
> - Jangan pernah mengubah isi file `.up.sql` yang sudah pernah di-run dan masuk ke database server (staging/prod).
> - Jika ingin mengubah struktur tabel, selalu buat file migrasi **baru**.

---

