# Requirements Document: Clean Architecture Refactoring

## Introduction

Dokumen ini menjelaskan requirements untuk refactoring project Authorizer agar mengikuti clean architecture pattern yang lebih baik, dengan menggunakan project Stackforge sebagai acuan. Refactoring ini bertujuan untuk meningkatkan separation of concerns, testability, dan maintainability dari codebase tanpa mengubah fungsionalitas yang ada.

## Glossary

- **Authorizer**: Sistem authorization dan authentication yang akan direfactor
- **Stackforge**: Project referensi yang mengimplementasikan clean architecture dengan baik
- **Resource_Layer**: Layer di Authorizer yang mencampur dependency injection dengan business logic (akan dihapus)
- **Domain_Layer**: Layer yang berisi business entities, repository interfaces, dan domain services (pure business logic)
- **Infrastructure_Layer**: Layer yang berisi implementasi konkret untuk database, cache, auth, config, dan logger
- **Usecase_Layer**: Layer yang berisi application business logic dan orchestration
- **Delivery_Layer**: Layer yang berisi HTTP handlers dan routing
- **Dependency_Injection**: Proses menyediakan dependencies ke komponen yang membutuhkan
- **Clean_Architecture**: Arsitektur software yang memisahkan concerns berdasarkan layers dengan dependency rule yang jelas
- **Echo**: Web framework yang digunakan oleh Authorizer
- **Gin**: Web framework yang digunakan oleh Stackforge
- **JWKS**: JSON Web Key Set untuk validasi JWT tokens
- **Redis**: In-memory cache yang digunakan untuk session management
- **PostgreSQL**: Relational database yang digunakan untuk persistence

## Requirements

### Requirement 1: Menghapus Resource Layer

**User Story:** Sebagai developer, saya ingin menghapus layer "resource" yang mencampur dependency injection dengan business logic, sehingga kode lebih mudah dipahami dan di-maintain.

#### Acceptance Criteria

1. THE System SHALL menghapus semua file di direktori `internal/app/resource/`
2. THE System SHALL memindahkan dependency injection dari resource layer ke `cmd/api/main.go`
3. WHEN resource layer dihapus, THEN semua handler initialization SHALL dilakukan langsung di main.go
4. THE System SHALL memastikan tidak ada import statement yang mereferensi package `internal/app/resource`

### Requirement 2: Memindahkan Dependency Injection ke Main.go

**User Story:** Sebagai developer, saya ingin dependency injection dilakukan di main.go seperti di Stackforge, sehingga dependency graph lebih eksplisit dan mudah di-trace.

#### Acceptance Criteria

1. THE System SHALL menginisialisasi semua infrastructure components (database, redis, config, logger) di main.go
2. THE System SHALL menginisialisasi semua repositories di main.go dengan dependencies yang diperlukan
3. THE System SHALL menginisialisasi semua use cases di main.go dengan repository dependencies
4. THE System SHALL menginisialisasi semua handlers di main.go dengan use case dependencies
5. WHEN dependencies diinisialisasi, THEN urutan inisialisasi SHALL mengikuti dependency order (infrastructure → repository → usecase → handler)
6. THE System SHALL memastikan tidak ada circular dependencies dalam dependency graph

### Requirement 3: Memisahkan Domain Service dari Infrastructure Service

**User Story:** Sebagai developer, saya ingin memisahkan domain service (pure business logic) dari infrastructure service (implementasi teknis), sehingga domain logic tidak tergantung pada infrastructure details.

#### Acceptance Criteria

1. THE System SHALL membuat direktori `internal/domain/service/` untuk domain services
2. THE System SHALL memindahkan pure business logic dari `internal/service/` ke `internal/domain/service/`
3. THE System SHALL membuat direktori `internal/infrastructure/auth/` untuk JWT-related infrastructure services
4. THE System SHALL memindahkan JWT implementation details ke infrastructure layer
5. WHEN JWT service direfactor, THEN domain service SHALL hanya bergantung pada domain entities dan repository interfaces
6. WHEN JWT service direfactor, THEN infrastructure service SHALL mengimplementasikan interface yang didefinisikan di domain layer

### Requirement 4: Merapikan Struktur Folder

**User Story:** Sebagai developer, saya ingin struktur folder yang konsisten dengan clean architecture principles, sehingga mudah menemukan dan mengorganisir kode.

#### Acceptance Criteria

1. THE System SHALL memindahkan config dari root level ke `internal/infrastructure/config/`
2. THE System SHALL membuat direktori `internal/infrastructure/logger/` untuk logging infrastructure
3. THE System SHALL memastikan struktur folder mengikuti pattern: `internal/{domain,usecase,infrastructure,delivery}/`
4. THE System SHALL memastikan domain layer tidak memiliki dependencies ke infrastructure atau delivery layers
5. THE System SHALL memastikan usecase layer hanya bergantung pada domain layer
6. THE System SHALL memastikan infrastructure dan delivery layers dapat bergantung pada domain dan usecase layers

### Requirement 5: Meningkatkan Testability

**User Story:** Sebagai developer, saya ingin dependency injection yang lebih baik sehingga komponen lebih mudah di-test dengan mocking.

#### Acceptance Criteria

1. THE System SHALL memastikan semua dependencies di-inject melalui constructor functions
2. THE System SHALL memastikan semua external dependencies (database, redis, HTTP client) di-abstract dengan interfaces
3. WHEN komponen memiliki dependencies, THEN dependencies tersebut SHALL berupa interfaces bukan concrete types
4. THE System SHALL memastikan setiap layer dapat di-test secara independen dengan mock dependencies
5. THE System SHALL menyediakan contoh test untuk minimal satu komponen di setiap layer

### Requirement 6: Mempertahankan Fungsionalitas Existing

**User Story:** Sebagai user, saya ingin semua fitur yang ada tetap berfungsi setelah refactoring, sehingga tidak ada breaking changes.

#### Acceptance Criteria

1. THE System SHALL mempertahankan semua HTTP endpoints yang ada dengan path dan method yang sama
2. THE System SHALL mempertahankan request/response format yang sama untuk semua endpoints
3. THE System SHALL mempertahankan authentication dan authorization logic yang sama
4. THE System SHALL mempertahankan database schema dan queries yang sama
5. THE System SHALL mempertahankan Redis caching behavior yang sama
6. WHEN refactoring selesai, THEN semua existing integration tests SHALL tetap pass tanpa modifikasi

### Requirement 7: Refactor Handler Initialization

**User Story:** Sebagai developer, saya ingin handler diinisialisasi langsung dengan usecase dependencies tanpa melalui resource layer, sehingga dependency flow lebih jelas.

#### Acceptance Criteria

1. THE System SHALL memodifikasi semua handler constructors untuk menerima usecase dependencies
2. THE System SHALL menghapus dependency ke resource layer dari semua handlers
3. WHEN handler diinisialisasi, THEN handler SHALL menerima usecase interface sebagai dependency
4. THE System SHALL memastikan handler tidak memiliki direct access ke repositories atau infrastructure services
5. THE System SHALL memastikan handler hanya memanggil usecase methods untuk business logic

### Requirement 8: Refactor Middleware

**User Story:** Sebagai developer, saya ingin middleware yang lebih clean dengan dependencies yang jelas, sehingga middleware mudah di-test dan di-maintain.

#### Acceptance Criteria

1. THE System SHALL memastikan JWT middleware menerima JWT service sebagai dependency
2. THE System SHALL memastikan permission middleware menerima authorization service sebagai dependency
3. WHEN middleware diinisialisasi, THEN middleware SHALL menerima dependencies melalui closure atau constructor
4. THE System SHALL menghapus global state atau singleton patterns dari middleware
5. THE System SHALL memastikan middleware dapat di-test dengan mock dependencies

### Requirement 9: Backward Compatibility

**User Story:** Sebagai API consumer, saya ingin API tetap kompatibel dengan client yang ada, sehingga tidak perlu mengubah client code.

#### Acceptance Criteria

1. THE System SHALL mempertahankan semua HTTP response status codes yang ada
2. THE System SHALL mempertahankan semua error message formats yang ada
3. THE System SHALL mempertahankan semua HTTP headers yang digunakan (Authorization, Content-Type, dll)
4. THE System SHALL mempertahankan CORS configuration yang ada
5. WHEN API dipanggil dengan request yang valid, THEN response SHALL identik dengan response sebelum refactoring

### Requirement 10: Documentation Update

**User Story:** Sebagai developer baru, saya ingin dokumentasi yang up-to-date dengan struktur baru, sehingga mudah memahami codebase.

#### Acceptance Criteria

1. THE System SHALL mengupdate README.md dengan struktur folder yang baru
2. THE System SHALL menambahkan dokumentasi tentang dependency injection pattern yang digunakan
3. THE System SHALL menambahkan diagram arsitektur yang menunjukkan layer dependencies
4. THE System SHALL menambahkan contoh cara menambahkan fitur baru dengan struktur yang baru
5. THE System SHALL menambahkan dokumentasi tentang cara menjalankan tests

