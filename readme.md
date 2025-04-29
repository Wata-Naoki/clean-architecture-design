# Echo + GORM によるクリーンアーキテクチャの実装

## 環境構築

### プロジェクトの初期化

```bash
mkdir go-echo-gorm-clean-arch
cd go-echo-gorm-clean-arch
go mod init github.com/yourusername/go-echo-gorm-clean-arch
```

### 必要なパッケージのインストール

```bash
# HTTPフレームワーク: Echo
go get github.com/labstack/echo/v4

# データベースORM: GORM
go get gorm.io/gorm
go get gorm.io/driver/mysql   # MySQLドライバー
go get gorm.io/driver/sqlite  # SQLiteドライバー

# 環境変数管理
go get github.com/joho/godotenv

# バリデーションライブラリ
go get github.com/go-playground/validator/v10

# テスト用ライブラリ
go get github.com/stretchr/testify
```

## プロジェクト構造

```
go-echo-gorm-clean-arch/
├── cmd/
│   └── api/
│       └── main.go        # エントリーポイント
├── internal/
│   ├── domain/            # ドメイン層
│   │   └── model/
│   │       ├── user.go    # ユーザーエンティティ
│   │       └── error.go   # ドメインエラー
│   ├── usecase/           # ユースケース層
│   │   ├── user_usecase.go
│   │   └── user_usecase_test.go
│   ├── repository/        # リポジトリ層
│   │   ├── user_repository.go     # インターフェース
│   │   └── gorm_repository.go     # GORM実装
│   └── delivery/          # デリバリー層
│       └── http/
│           ├── handler/
│           │   └── user_handler.go
│           ├── middleware/
│           │   └── middleware.go
│           └── route/
│               └── route.go
├── pkg/
│   ├── config/            # 設定
│   │   └── config.go
│   └── database/          # DB接続
│       └── database.go
├── .env                   # 環境変数
├── .env.example
├── go.mod
└── go.sum
```

## 各レイヤーの概要

- **ドメイン層**: ビジネスエンティティとルールを定義
- **リポジトリ層**: データアクセス方法を抽象化
- **ユースケース層**: アプリケーション固有のビジネスロジックを実装
- **デリバリー層**: HTTPリクエストの受信と処理
- **設定層**: アプリケーション設定の管理
- **データベース層**: GORMを使用したデータベース接続と操作

## 各レイヤーの実装

### 1. ドメイン層 (Domain Layer)

**internal/domain/model/user.go**

```go
package model

import (
 "time"
)


type User struct {
 ID        uint      `json:"id" gorm:"primaryKey"`
 Name      string    `json:"name" gorm:"size:100;not null"`
 Email     string    `json:"email" gorm:"size:100;not null;uniqueIndex"`
 Password  string    `json:"-" gorm:"size:100;not null"`
 CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
 UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
```

**internal/domain/model/error.go**

```go
package model

import "errors"


var (
 ErrInternalServerError = errors.New("internal server error")
 ErrNotFound            = errors.New("your requested item is not found")
 ErrConflict            = errors.New("your item already exists")
 ErrBadRequest          = errors.New("bad request")
 ErrInvalidCredentials  = errors.New("invalid credentials")
)
```

### 2. リポジトリ層 (Repository Layer)

**internal/repository/user_repository.go**

```go
package repository

import (
 "context"
 
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/domain/model"
)

// UserRepository はユーザー関連のデータアクセスを定義するインターフェース
type UserRepository interface {
 GetByID(ctx context.Context, id uint) (*model.User, error)
 GetByEmail(ctx context.Context, email string) (*model.User, error)
 Create(ctx context.Context, user *model.User) error
 Update(ctx context.Context, user *model.User) error
 Delete(ctx context.Context, id uint) error
 List(ctx context.Context, limit, offset int) ([]*model.User, error)
}
```

**internal/repository/gorm_repository.go**

```go
package repository

import (
 "context"
 "errors"
 
 "gorm.io/gorm"
 
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/domain/model"
)

type gormUserRepository struct {
 db *gorm.DB
}

// NewGormUserRepository はGORMを使用したUserRepositoryの新しいインスタンスを作成
func NewGormUserRepository(db *gorm.DB) UserRepository {
 return &gormUserRepository{
  db: db,
 }
}

func (r *gormUserRepository) GetByID(ctx context.Context, id uint) (*model.User, error) {
 var user model.User
 result := r.db.WithContext(ctx).First(&user, id)
 
 if result.Error != nil {
  if errors.Is(result.Error, gorm.ErrRecordNotFound) {
   return nil, model.ErrNotFound
  }
  return nil, model.ErrInternalServerError
 }
 
 return &user, nil
}

func (r *gormUserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
 var user model.User
 result := r.db.WithContext(ctx).Where("email = ?", email).First(&user)
 
 if result.Error != nil {
  if errors.Is(result.Error, gorm.ErrRecordNotFound) {
   return nil, model.ErrNotFound
  }
  return nil, model.ErrInternalServerError
 }
 
 return &user, nil
}

func (r *gormUserRepository) Create(ctx context.Context, user *model.User) error {
 result := r.db.WithContext(ctx).Create(user)
 if result.Error != nil {
  return model.ErrInternalServerError
 }
 
 return nil
}

func (r *gormUserRepository) Update(ctx context.Context, user *model.User) error {
 result := r.db.WithContext(ctx).Save(user)
 if result.Error != nil {
  return model.ErrInternalServerError
 }
 
 if result.RowsAffected == 0 {
  return model.ErrNotFound
 }
 
 return nil
}

func (r *gormUserRepository) Delete(ctx context.Context, id uint) error {
 result := r.db.WithContext(ctx).Delete(&model.User{}, id)
 if result.Error != nil {
  return model.ErrInternalServerError
 }
 
 if result.RowsAffected == 0 {
  return model.ErrNotFound
 }
 
 return nil
}

func (r *gormUserRepository) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
 var users []*model.User
 result := r.db.WithContext(ctx).Limit(limit).Offset(offset).Find(&users)
 
 if result.Error != nil {
  return nil, model.ErrInternalServerError
 }
 
 return users, nil
}
```

### 3. ユースケース層 (Usecase Layer)

**internal/usecase/user_usecase.go**

```go
package usecase

import (
 "context"
 "time"
 
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/domain/model"
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/repository"
)

// UserUsecase はユーザー関連のビジネスロジックを定義するインターフェース
type UserUsecase interface {
 GetByID(ctx context.Context, id uint) (*model.User, error)
 Create(ctx context.Context, user *model.User) error
 Update(ctx context.Context, user *model.User) error
 Delete(ctx context.Context, id uint) error
 List(ctx context.Context, limit, offset int) ([]*model.User, error)
}

type userUsecase struct {
 userRepo repository.UserRepository
}

// NewUserUsecase はUserUsecaseの新しいインスタンスを作成
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
 return &userUsecase{
  userRepo: userRepo,
 }
}

func (u *userUsecase) GetByID(ctx context.Context, id uint) (*model.User, error) {
 return u.userRepo.GetByID(ctx, id)
}

func (u *userUsecase) Create(ctx context.Context, user *model.User) error {
 existingUser, err := u.userRepo.GetByEmail(ctx, user.Email)
 // エラーがなく、ユーザーが見つかった場合は既に存在している
 if err == nil && existingUser != nil {
  return model.ErrConflict
 }
 
 // NotFoundエラー以外のエラーの場合
 if err != nil && err != model.ErrNotFound {
  return err
 }
 
 // 現在時刻を設定
 now := time.Now()
 user.CreatedAt = now
 user.UpdatedAt = now
 
 return u.userRepo.Create(ctx, user)
}

func (u *userUsecase) Update(ctx context.Context, user *model.User) error {
 // IDに対応するユーザーが存在するか確認
 _, err := u.userRepo.GetByID(ctx, user.ID)
 if err != nil {
  return err
 }
 
 // 更新時刻を設定
 user.UpdatedAt = time.Now()
 
 return u.userRepo.Update(ctx, user)
}

func (u *userUsecase) Delete(ctx context.Context, id uint) error {
 return u.userRepo.Delete(ctx, id)
}

func (u *userUsecase) List(ctx context.Context, limit, offset int) ([]*model.User, error) {
 // デフォルト値の設定
 if limit <= 0 {
  limit = 10
 }
 if offset < 0 {
  offset = 0
 }
 
 return u.userRepo.List(ctx, limit, offset)
}
```

### 4. デリバリー層 (Delivery Layer)

**internal/delivery/http/handler/user_handler.go**

```go
package handler

import (
 "net/http"
 "strconv"
 
 "github.com/labstack/echo/v4"
 
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/domain/model"
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/usecase"
)

// UserHandler はユーザー関連のHTTPリクエストを処理するハンドラー
type UserHandler struct {
 UserUsecase usecase.UserUsecase
}

// NewUserHandler はUserHandlerの新しいインスタンスを作成
func NewUserHandler(g *echo.Group, userUsecase usecase.UserUsecase) {
 handler := &UserHandler{
  UserUsecase: userUsecase,
 }
 
 // ルートの登録
 g.GET("", handler.GetUsers)
 g.GET("/:id", handler.GetUser)
 g.POST("", handler.CreateUser)
 g.PUT("/:id", handler.UpdateUser)
 g.DELETE("/:id", handler.DeleteUser)
}

// GetUsers はユーザー一覧を取得するハンドラー
func (h *UserHandler) GetUsers(c echo.Context) error {
 limit, _ := strconv.Atoi(c.QueryParam("limit"))
 offset, _ := strconv.Atoi(c.QueryParam("offset"))
 
 ctx := c.Request().Context()
 users, err := h.UserUsecase.List(ctx, limit, offset)
 if err != nil {
  return c.JSON(getStatusCode(err), map[string]string{
   "message": err.Error(),
  })
 }
 
 return c.JSON(http.StatusOK, map[string]interface{}{
  "data": users,
 })
}

// GetUser は指定されたIDのユーザーを取得するハンドラー
func (h *UserHandler) GetUser(c echo.Context) error {
 idParam := c.Param("id")
 id, err := strconv.ParseUint(idParam, 10, 32)
 if err != nil {
  return c.JSON(http.StatusBadRequest, map[string]string{
   "message": "invalid id parameter",
  })
 }
 
 ctx := c.Request().Context()
 user, err := h.UserUsecase.GetByID(ctx, uint(id))
 if err != nil {
  return c.JSON(getStatusCode(err), map[string]string{
   "message": err.Error(),
  })
 }
 
 return c.JSON(http.StatusOK, map[string]interface{}{
  "data": user,
 })
}

// CreateUser は新しいユーザーを作成するハンドラー
func (h *UserHandler) CreateUser(c echo.Context) error {
 var user model.User
 if err := c.Bind(&user); err != nil {
  return c.JSON(http.StatusBadRequest, map[string]string{
   "message": err.Error(),
  })
 }
 
 if user.Name == "" || user.Email == "" || user.Password == "" {
  return c.JSON(http.StatusBadRequest, map[string]string{
   "message": "name, email and password are required",
  })
 }
 
 ctx := c.Request().Context()
 err := h.UserUsecase.Create(ctx, &user)
 if err != nil {
  return c.JSON(getStatusCode(err), map[string]string{
   "message": err.Error(),
  })
 }
 
 return c.JSON(http.StatusCreated, map[string]interface{}{
  "data": user,
 })
}

// UpdateUser は既存のユーザーを更新するハンドラー
func (h *UserHandler) UpdateUser(c echo.Context) error {
 idParam := c.Param("id")
 id, err := strconv.ParseUint(idParam, 10, 32)
 if err != nil {
  return c.JSON(http.StatusBadRequest, map[string]string{
   "message": "invalid id parameter",
  })
 }
 
 var user model.User
 if err := c.Bind(&user); err != nil {
  return c.JSON(http.StatusBadRequest, map[string]string{
   "message": err.Error(),
  })
 }
 
 user.ID = uint(id)
 
 ctx := c.Request().Context()
 err = h.UserUsecase.Update(ctx, &user)
 if err != nil {
  return c.JSON(getStatusCode(err), map[string]string{
   "message": err.Error(),
  })
 }
 
 return c.JSON(http.StatusOK, map[string]interface{}{
  "data": user,
 })
}

// DeleteUser はユーザーを削除するハンドラー
func (h *UserHandler) DeleteUser(c echo.Context) error {
 idParam := c.Param("id")
 id, err := strconv.ParseUint(idParam, 10, 32)
 if err != nil {
  return c.JSON(http.StatusBadRequest, map[string]string{
   "message": "invalid id parameter",
  })
 }
 
 ctx := c.Request().Context()
 err = h.UserUsecase.Delete(ctx, uint(id))
 if err != nil {
  return c.JSON(getStatusCode(err), map[string]string{
   "message": err.Error(),
  })
 }
 
 return c.NoContent(http.StatusNoContent)
}

// getStatusCode はエラータイプに基づいてHTTPステータスコードを返す
func getStatusCode(err error) int {
 if err == nil {
  return http.StatusOK
 }
 
 switch err {
 case model.ErrInternalServerError:
  return http.StatusInternalServerError
 case model.ErrNotFound:
  return http.StatusNotFound
 case model.ErrConflict:
  return http.StatusConflict
 case model.ErrInvalidCredentials:
  return http.StatusUnauthorized
 default:
  return http.StatusInternalServerError
 }
}
```

**internal/delivery/http/middleware/middleware.go**

```go
package middleware

import (
 "github.com/labstack/echo/v4"
 "github.com/labstack/echo/v4/middleware"
)

// SetupMiddleware はHTTPミドルウェアを設定
func SetupMiddleware(e *echo.Echo) {
 e.Use(middleware.Logger())
 e.Use(middleware.Recover())
 e.Use(middleware.CORS())
}
```

**internal/delivery/http/route/route.go**

```go
package route

import (
 "github.com/labstack/echo/v4"
 
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/delivery/http/handler"
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/delivery/http/middleware"
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/usecase"
)

// SetupRoutes はHTTPルートをセットアップ
func SetupRoutes(e *echo.Echo, userUsecase usecase.UserUsecase) {
 // ミドルウェアの設定
 middleware.SetupMiddleware(e)
 
 // APIグループ
 apiV1 := e.Group("/api/v1")
 
 // ユーザーエンドポイント
 usersGroup := apiV1.Group("/users")
 handler.NewUserHandler(usersGroup, userUsecase)
}
```

### 5. データベース接続 (Database Connection)

**pkg/database/database.go**

```go
package database

import (
 "fmt"
 "log"
 
 "gorm.io/driver/mysql"
 "gorm.io/driver/sqlite"
 "gorm.io/gorm"
 "gorm.io/gorm/logger"
 
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/domain/model"
)

// NewMySQLDatabase はMySQLデータベース接続を作成
func NewMySQLDatabase(username, password, host, port, dbname string) (*gorm.DB, error) {
 dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
  username, password, host, port, dbname)
 
 db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Info),
 })
 if err != nil {
  return nil, err
 }
 
 // モデルのマイグレーション
 err = db.AutoMigrate(&model.User{})
 if err != nil {
  return nil, err
 }
 
 return db, nil
}

// NewSQLiteDatabase はSQLiteデータベース接続を作成
func NewSQLiteDatabase(dbPath string) (*gorm.DB, error) {
 db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
  Logger: logger.Default.LogMode(logger.Info),
 })
 if err != nil {
  return nil, err
 }
 
 // モデルのマイグレーション
 err = db.AutoMigrate(&model.User{})
 if err != nil {
  return nil, err
 }
 
 return db, nil
}
```

### 6. 設定 (Config)

**pkg/config/config.go**

```go
package config

import (
 "log"
 "os"
 
 "github.com/joho/godotenv"
)

// Config はアプリケーション設定を保持する構造体
type Config struct {
 AppPort       string
 DBDriver      string
 DBHost        string
 DBPort        string
 DBUser        string
 DBPassword    string
 DBName        string
 SQLitePath    string
}

// LoadConfig は.envファイルから設定を読み込み
func LoadConfig() *Config {
 err := godotenv.Load()
 if err != nil {
  log.Println("Warning: .env file not found. Using environment variables.")
 }
 
 return &Config{
  AppPort:       getEnv("APP_PORT", "8080"),
  DBDriver:      getEnv("DB_DRIVER", "sqlite"),  // "mysql" or "sqlite"
  DBHost:        getEnv("DB_HOST", "localhost"),
  DBPort:        getEnv("DB_PORT", "3306"),
  DBUser:        getEnv("DB_USER", "root"),
  DBPassword:    getEnv("DB_PASSWORD", "password"),
  DBName:        getEnv("DB_NAME", "cleanarch"),
  SQLitePath:    getEnv("SQLITE_PATH", "database.db"),
 }
}

// getEnv は環境変数を取得し、設定されていない場合はデフォルト値を返す
func getEnv(key, defaultValue string) string {
 value := os.Getenv(key)
 if value == "" {
  return defaultValue
 }
 return value
}
```

### 7. メイン関数

**cmd/api/main.go**

```go
package main

import (
 "fmt"
 "log"
 
 "github.com/labstack/echo/v4"
 
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/delivery/http/route"
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/repository"
 "github.com/yourusername/go-echo-gorm-clean-arch/internal/usecase"
 "github.com/yourusername/go-echo-gorm-clean-arch/pkg/config"
 "github.com/yourusername/go-echo-gorm-clean-arch/pkg/database"
)

func main() {
 // 設定ロード
 cfg := config.LoadConfig()
 
 // データベース接続
 var db, err = setupDatabase(cfg)
 if err != nil {
  log.Fatalf("Failed to connect to database: %v", err)
 }
 
 // リポジトリ層
 userRepo := repository.NewGormUserRepository(db)
 
 // ユースケース層
 userUsecase := usecase.NewUserUsecase(userRepo)
 
 // Echoインスタンス作成
 e := echo.New()
 
 // ルート設定
 route.SetupRoutes(e, userUsecase)
 
 // サーバー起動
 serverAddr := fmt.Sprintf(":%s", cfg.AppPort)
 log.Printf("Server running on %s", serverAddr)
 if err := e.Start(serverAddr); err != nil {
  log.Fatalf("Failed to start server: %v", err)
 }
}

// setupDatabase はDBドライバーに基づいてデータベース接続を設定
func setupDatabase(cfg *config.Config) (*gorm.DB, error) {
 if cfg.DBDriver == "mysql" {
  return database.NewMySQLDatabase(
   cfg.DBUser,
   cfg.DBPassword,
   cfg.DBHost,
   cfg.DBPort,
   cfg.DBName,
  )
 }
 
 // デフォルトはSQLite
 return database.NewSQLiteDatabase(cfg.SQLitePath)
}
```

### 8. 環境変数ファイル

**.env.example**

```
APP_PORT=8080
DB_DRIVER=sqlite       # mysql or sqlite
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=password
DB_NAME=cleanarch
SQLITE_PATH=database.db
```

## 実行方法

1. **.env**ファイルを作成:

```bash
cp .env.example .env
```

2. アプリケーションを起動:

```bash
go run cmd/api/main.go
```

3. APIエンドポイントにアクセス:

```bash
# ユーザーを作成
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name":"John Doe","email":"john@example.com","password":"secret"}'

# ユーザー一覧を取得
curl http://localhost:8080/api/v1/users

# IDでユーザーを取得
curl http://localhost:8080/api/v1/users/1

# ユーザーを更新
curl -X PUT http://localhost:8080/api/v1/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name":"John Updated","email":"john@example.com"}'

# ユーザーを削除
curl -X DELETE http://localhost:8080/api/v1/users/1
```
