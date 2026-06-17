# QA Test Report: Product & Service Management API (Fase 3 - Comprehensive)

> [!NOTE]
> Laporan pengujian ini dibuat secara menyeluruh oleh persona **QA Automation Engineer** (`qa_test`) untuk memverifikasi keselarasan sistem dengan 14 Aturan Bisnis (**BR-001** sampai **BR-014**) dan 12 Skenario Uji (**TC-001** sampai **TC-012**) yang tertuang dalam dokumen `@[Technical Test - Backend Developer.pdf]`.

---

## 1. Agent Instruction & Prompt (Fase 3 Requirement)

Sesuai kriteria kelulusan Fase 3, instruksi pengujian ini tersimpan secara terstruktur dan dapat digunakan kembali (*reusable*) di dalam `.antigravity/rules.json` pada bagian persona `qa_test`.

---

## 2. Test Case Summary (12 Skenario Minimal)

Berikut adalah status pengujian dari ke-12 skenario wajib:

| Test Case ID | Skenario | Endpoint | Method | Expected Status | Status Hasil |
|---|---|---|---|---|---|
| **TC-001** | Create item valid | `/companies/{cid}/items` | `POST` | `201 ITEM_CREATED` | 🟢 **PASS** |
| **TC-002** | Create item dengan duplicate code | `/companies/{cid}/items` | `POST` | `409 ITEM_CODE_ALREADY_EXISTS` | 🟢 **PASS** |
| **TC-003** | Create item dengan type invalid | `/companies/{cid}/items` | `POST` | `400 VALIDATION_ERROR` | 🟢 **PASS** |
| **TC-004** | Create item dengan price negatif | `/companies/{cid}/items` | `POST` | `400 VALIDATION_ERROR` | 🟢 **PASS** |
| **TC-005** | Get list item | `/companies/{cid}/items` | `GET` | `200 ITEM_LIST_RETRIEVED` | 🟢 **PASS** |
| **TC-006** | Get detail item valid | `/companies/{cid}/items/{id}` | `GET` | `200 ITEM_DETAIL_RETRIEVED` | 🟢 **PASS** |
| **TC-007** | Get detail item tidak ditemukan | `/companies/{cid}/items/{id}` | `GET` | `404 ITEM_NOT_FOUND` | 🟢 **PASS** |
| **TC-008** | Update item valid | `/companies/{cid}/items/{id}` | `PATCH` | `200 ITEM_UPDATED` | 🟢 **PASS** |
| **TC-009** | Update item menjadi duplicate code | `/companies/{cid}/items/{id}` | `PATCH` | `409 ITEM_CODE_ALREADY_EXISTS` | 🟢 **PASS** |
| **TC-010** | Archive item valid | `/companies/{cid}/items/{id}/archive` | `PATCH` | `200 ITEM_ARCHIVED` | 🟢 **PASS** |
| **TC-011** | Archive item yang sudah archived | `/companies/{cid}/items/{id}/archive` | `PATCH` | `409 ITEM_ALREADY_ARCHIVED` | 🟢 **PASS** |
| **TC-012** | Update archived item | `/companies/{cid}/items/{id}` | `PATCH` | `409 ITEM_ALREADY_ARCHIVED` | 🟢 **PASS** |

---

## 3. API Request Details (Generated Curl & Expected Responses)

Pengujian dieksekusi menggunakan Base URL: `http://localhost:8080` dan `company_id: c0a80121-7ac0-11d1-898c-00c04fd8d5c1`.

### TC-001 & TC-002: Create & Duplicate Code Validation
```bash
# Create Item (TC-001)
curl -s -X POST http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items \
-H "Content-Type: application/json" \
-d '{"code":"ITEM-999", "name":"Kopi Susu Gula Aren", "type":"PRODUCT", "price":15000, "category_name":"Beverages"}'

# Duplicate (TC-002) - Kirim ulang payload yang sama
```
* **Expected (TC-001)**: `201` dengan JSON `{"success":true,"code":"ITEM_CREATED"}`
* **Expected (TC-002)**: `409` dengan JSON `{"success":false,"code":"ITEM_CODE_ALREADY_EXISTS"}`

### TC-003 & TC-004: Invalid Type & Negative Price
```bash
# Invalid Type (TC-003)
curl -s -X POST http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items \
-H "Content-Type: application/json" \
-d '{"code":"ITEM-003", "name":"Incorrect Type", "type":"FOOD", "price":15000}'

# Negative Price (TC-004)
curl -s -X POST http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items \
-H "Content-Type: application/json" \
-d '{"code":"ITEM-004", "name":"Negative Price", "type":"PRODUCT", "price":-5000}'
```
* **Expected (TC-003 / TC-004)**: `400` dengan JSON `{"success":false,"code":"VALIDATION_ERROR"}`

### TC-005, TC-006 & TC-007: List & Detail Checks
```bash
# Get List (TC-005)
curl -s -X GET http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items

# Get Detail Valid (TC-006)
curl -s -X GET http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/{created_item_id}

# Get Detail Not Found (TC-007)
curl -s -X GET http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/00000000-0000-0000-0000-000000000000
```
* **Expected (TC-005)**: `200 ITEM_LIST_RETRIEVED`
* **Expected (TC-006)**: `200 ITEM_DETAIL_RETRIEVED`
* **Expected (TC-007)**: `404 ITEM_NOT_FOUND`

### TC-008 & TC-009: Update Valid & Duplicate Check
```bash
# Update Valid (TC-008)
curl -s -X PATCH http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/{created_item_id} \
-H "Content-Type: application/json" \
-d '{"code":"ITEM-999-REV", "name":"Kopi Susu Gula Aren Update", "type":"PRODUCT", "price":17000, "status":"ACTIVE"}'
```
* **Expected (TC-008)**: `200 ITEM_UPDATED`
* **Expected (TC-009)**: `409 ITEM_CODE_ALREADY_EXISTS` (jika diubah ke kode item milik produk lain)

### TC-010, TC-011 & TC-012: Lifecycle of Archived Items
```bash
# Archive Item (TC-010)
curl -s -X PATCH http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/{created_item_id}/archive

# Re-archive (TC-011) - Jalankan curl archive di atas sekali lagi

# Update Archived Item (TC-012)
curl -s -X PATCH http://localhost:8080/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/{created_item_id} \
-H "Content-Type: application/json" \
-d '{"code":"ITEM-999-REV", "name":"Try to Update Archived", "type":"PRODUCT", "price":17000, "status":"ACTIVE"}'
```
* **Expected (TC-010)**: `200 ITEM_ARCHIVED`
* **Expected (TC-011)**: `409 ITEM_ALREADY_ARCHIVED`
* **Expected (TC-012)**: `409 ITEM_ALREADY_ARCHIVED`

---

## 4. Expected vs Actual Response Comparison

Semua status aktual yang diperoleh dari server lokal adalah **PASS** (100% Cocok). Perilaku sistem sudah benar-benar sesuai dengan response format standard:
- Success: `{"success": true, "code": "...", "message": "...", "data": ...}`
- Error: `{"success": false, "code": "...", "message": "...", "errors": ...}`

---

## 5. Penjelasan False Positive Prevention (Fase 3 Requirement)

Untuk menjamin kualitas dan mencegah terjadinya *False Positive* (status ditandai lulus padahal implementasi bermasalah), QA Automation Engineer menerapkan strategi berikut:

1. **Assert pada Response Payload, Bukan Hanya HTTP Status**:
   Banyak server melempar `200 OK` namun mengembalikan body berisi error or data kosong. Kami selalu melakukan *assertion* terhadap field `success: true` dan kecocokan kode error terstruktur (seperti `ITEM_CODE_ALREADY_EXISTS`).
2. **UUID Format Validation**:
   Validasi ketat pada *path parameter* menggunakan middleware regex di level controller memastikan database terlindungi dari malformed queries yang berujung pada error 500 yang ambigu.
3. **Data Isolation/Multitenancy Check**:
   Kami menguji kueri dengan `company_id` acak untuk memastikan data tidak bocor antar *tenant* (mengembalikan `404 ITEM_NOT_FOUND` / `403 FORBIDDEN` sesuai BR-007).
4. **State Transition Strict Rules Assertion**:
   Kami menguji bahwa transisi dari `ARCHIVED` kembali ke `ACTIVE` tidak bisa dilakukan melalui rute Update biasa.

---

## 6. Gap Analysis

1. **Temuan Awal (Resolved)**: Route prefix `/api/v1` bertolak belakang dengan PRD. Masalah ini telah di-patch oleh `golang_coder` sehingga rute sekarang sinkron ke root `/companies/...`.
2. **Temuan Validasi (Resolved)**: Nilai `price: 0` sempat gagal divalidasi karena penggunaan `binding:"required"` tipe primitif. Kami sudah mengubah tipe data `Price` menjadi `*float64` di request struct sehingga harga 0 kini sah dan divalidasi dengan benar sesuai aturan bisnis.
3. **Temuan UUID Exception (Resolved)**: Parameter invalid memicu error `500` dari postgres. Hal ini berhasil dimitigasi dengan pengenalan `UUIDParamValidator` middleware yang memberikan respons `400` secara elegan.
