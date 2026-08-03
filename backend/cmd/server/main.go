package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"
	"github.com/rootage-ses-quiz/backend/internal/config"
	"github.com/rootage-ses-quiz/backend/internal/handler"
	"github.com/rootage-ses-quiz/backend/internal/middleware"
	"github.com/rootage-ses-quiz/backend/internal/pkg"
	"github.com/rootage-ses-quiz/backend/internal/repository"
	"github.com/rootage-ses-quiz/backend/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	if err := db.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}
	log.Println("Connected to database")

	// マイグレーション・シードデータ自動実行
	runMigrations(db)
	runSeed(db)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	questionRepo := repository.NewQuestionRepository(db)
	quizRepo := repository.NewQuizRepository(db)
	documentRepo := repository.NewDocumentRepository(db)
	chatRepo := repository.NewChatRepository(db)
	guideRepo := repository.NewGuideRepository(db)
	resetRepo := repository.NewPasswordResetRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo, resetRepo, cfg.JWTSecret, cfg.AllowedEmailDomain, cfg.ResendAPIKey, cfg.BaseURL)
	quizSvc := service.NewQuizService(quizRepo, questionRepo, categoryRepo)
	openaiClient := service.NewOpenAIClient(cfg.OpenAIAPIKey)
	if cfg.OpenAIAPIKey == "" {
		log.Println("WARNING: OPENAI_API_KEY is not set, RAG and document indexing will not work")
	}
	ragSvc := service.NewRAGService(documentRepo, chatRepo, openaiClient)
	adminSvc := service.NewAdminService(categoryRepo, questionRepo, documentRepo, guideRepo, userRepo, openaiClient)

	// 管理者アカウントは環境変数が設定されているときだけ作成する（既存アカウントは上書きしない）
	if err := authSvc.EnsureAdminUser(context.Background(), cfg.AdminEmail, cfg.AdminPassword); err != nil {
		log.Fatalf("failed to create admin user: %v", err)
	}
	if cfg.AdminEmail == "" {
		log.Println("ADMIN_EMAIL is not set, skipping admin user creation")
	}

	// シード文書のインデックス作成（シード完了後にバックグラウンドで実行）
	if cfg.OpenAIAPIKey != "" {
		go indexUnprocessedDocuments(adminSvc)
	}

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	quizHandler := handler.NewQuizHandler(quizSvc)
	ragHandler := handler.NewRAGHandler(ragSvc)
	adminHandler := handler.NewAdminHandler(adminSvc)

	// Router
	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		pkg.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	// Auth routes (public)
	mux.HandleFunc("POST /api/auth/register", authHandler.Register)
	mux.HandleFunc("POST /api/auth/login", authHandler.Login)
	mux.HandleFunc("POST /api/auth/forgot-password", authHandler.ForgotPassword)
	mux.HandleFunc("POST /api/auth/reset-password", authHandler.ResetPassword)

	// Category routes (authenticated)
	mux.HandleFunc("GET /api/categories", middleware.Auth(cfg.JWTSecret, quizHandler.ListCategories))

	// Quiz routes (authenticated)
	mux.HandleFunc("POST /api/quiz/start", middleware.Auth(cfg.JWTSecret, quizHandler.Start))
	mux.HandleFunc("POST /api/quiz/answer", middleware.Auth(cfg.JWTSecret, quizHandler.Answer))
	mux.HandleFunc("POST /api/quiz/finish", middleware.Auth(cfg.JWTSecret, quizHandler.Finish))
	mux.HandleFunc("GET /api/quiz/history", middleware.Auth(cfg.JWTSecret, quizHandler.History))
	mux.HandleFunc("GET /api/quiz/review", middleware.Auth(cfg.JWTSecret, quizHandler.Review))
	mux.HandleFunc("GET /api/quiz/stats", middleware.Auth(cfg.JWTSecret, quizHandler.Stats))

	// User dashboard (authenticated)
	mux.HandleFunc("GET /api/users/me/stats", middleware.Auth(cfg.JWTSecret, quizHandler.Dashboard))

	// RAG routes (authenticated)
	mux.HandleFunc("POST /api/rag/ask", middleware.Auth(cfg.JWTSecret, ragHandler.Ask))
	mux.HandleFunc("GET /api/rag/history", middleware.Auth(cfg.JWTSecret, ragHandler.History))

	// Admin routes (authenticated + admin)
	mux.HandleFunc("GET /api/admin/questions", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.ListQuestions)))
	mux.HandleFunc("POST /api/admin/questions", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.CreateQuestion)))
	mux.HandleFunc("PUT /api/admin/questions/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.UpdateQuestion)))
	mux.HandleFunc("DELETE /api/admin/questions/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.DeleteQuestion)))
	mux.HandleFunc("POST /api/admin/questions/import", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.ImportQuestions)))

	mux.HandleFunc("GET /api/admin/categories", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.ListCategories)))
	mux.HandleFunc("POST /api/admin/categories", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.CreateCategory)))
	mux.HandleFunc("PUT /api/admin/categories/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.UpdateCategory)))
	mux.HandleFunc("DELETE /api/admin/categories/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.DeleteCategory)))

	mux.HandleFunc("GET /api/admin/documents", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.ListDocuments)))
	mux.HandleFunc("POST /api/admin/documents", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.UploadDocument)))
	mux.HandleFunc("POST /api/admin/documents/text", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.CreateDocumentFromText)))
	mux.HandleFunc("DELETE /api/admin/documents/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.DeleteDocument)))
	mux.HandleFunc("POST /api/admin/documents/{id}/reindex", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.ReindexDocument)))

	// Admin summary
	mux.HandleFunc("GET /api/admin/summary", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.GetSummary)))

	// User progress (admin)
	mux.HandleFunc("GET /api/admin/user-progress", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.GetUserProgress)))
	mux.HandleFunc("POST /api/admin/users/{id}/reset-password", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.ResetPassword)))

	// Guide routes (user)
	mux.HandleFunc("GET /api/guides", middleware.Auth(cfg.JWTSecret, adminHandler.ListPublishedGuides))
	mux.HandleFunc("GET /api/guides/{id}", middleware.Auth(cfg.JWTSecret, adminHandler.GetGuide))

	// Guide admin routes
	mux.HandleFunc("GET /api/admin/guide-categories", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.ListGuideCategories)))
	mux.HandleFunc("POST /api/admin/guide-categories", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.CreateGuideCategory)))
	mux.HandleFunc("PUT /api/admin/guide-categories/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.UpdateGuideCategory)))
	mux.HandleFunc("DELETE /api/admin/guide-categories/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.DeleteGuideCategory)))
	mux.HandleFunc("GET /api/admin/guides", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.ListGuides)))
	mux.HandleFunc("POST /api/admin/guides", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.CreateGuide)))
	mux.HandleFunc("PUT /api/admin/guides/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.UpdateGuide)))
	mux.HandleFunc("DELETE /api/admin/guides/{id}", middleware.Auth(cfg.JWTSecret, middleware.Admin(adminHandler.DeleteGuide)))

	// フロントエンドの静的ファイル配信（本番用）
	if _, err := os.Stat("public"); err == nil {
		fs := http.FileServer(http.Dir("public"))
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// /api で始まるパスはAPIハンドラーに任せる
			path := r.URL.Path
			// 静的ファイルが存在するか確認
			if _, err := os.Stat("public" + path); err == nil && path != "/" {
				fs.ServeHTTP(w, r)
				return
			}
			// SPA: 存在しないパスは全てindex.htmlを返す
			http.ServeFile(w, r, "public/index.html")
		})
		log.Println("Serving frontend from ./public")
	}

	// Apply global middleware
	var h http.Handler = mux
	h = middleware.Recovery(h)
	h = middleware.Logging(h)
	h = middleware.CORS(cfg.CorsOrigin)(h)

	srv := &http.Server{
		Addr:         ":" + cfg.ServerPort,
		Handler:      h,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.ServerPort)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("Server stopped")
}

func runSeed(db *sql.DB) {
	log.Println("Checking seed data...")
	files := []struct {
		path  string
		check string // このクエリが0件ならシードを実行
	}{
		{"seed/seed.sql", "SELECT COUNT(*) FROM categories"},
		{"seed/seed_content.sql", "SELECT COUNT(*) FROM questions WHERE difficulty = 'intermediate'"},
		{"seed/seed_java.sql", "SELECT COUNT(*) FROM questions WHERE explanation LIKE '%Java%Silver%'"},
		{"seed/seed_guides.sql", "SELECT COUNT(*) FROM guide_categories"},
		{"seed/seed_documents.sql", "SELECT COUNT(*) FROM documents"},
	}

	for _, f := range files {
		var count int
		if err := db.QueryRow(f.check).Scan(&count); err == nil && count > 0 {
			continue
		}
		data, err := os.ReadFile(f.path)
		if err != nil {
			log.Printf("Seed file %s not found, skipping", f.path)
			continue
		}
		if _, err := db.Exec(string(data)); err != nil {
			log.Printf("Seed %s error: %v", f.path, err)
		} else {
			log.Printf("Seed applied: %s", f.path)
		}
	}
}

func indexUnprocessedDocuments(adminSvc *service.AdminService) {
	ctx := context.Background()
	docs, err := adminSvc.ListDocuments(ctx)
	if err != nil {
		return
	}
	for _, doc := range docs {
		if err := adminSvc.ReindexDocument(ctx, doc.ID); err != nil {
			log.Printf("Index error for %s: %v", doc.Title, err)
		} else {
			log.Printf("Indexed document: %s", doc.Title)
		}
	}
}

func runMigrations(db *sql.DB) {
	files, err := os.ReadDir("migrations")
	if err != nil {
		log.Println("No migrations directory found, skipping")
		return
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}
		data, err := os.ReadFile("migrations/" + f.Name())
		if err != nil {
			log.Printf("Failed to read migration %s: %v", f.Name(), err)
			continue
		}
		if _, err := db.Exec(string(data)); err != nil {
			log.Printf("Migration %s: %v (may already be applied)", f.Name(), err)
		} else {
			log.Printf("Migration applied: %s", f.Name())
		}
	}
}
