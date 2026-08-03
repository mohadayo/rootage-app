package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/rootage-ses-quiz/backend/internal/dto"
	"github.com/rootage-ses-quiz/backend/internal/model"
	"github.com/rootage-ses-quiz/backend/internal/pkg"
	"github.com/rootage-ses-quiz/backend/internal/repository"
)

type AuthService struct {
	userRepo           *repository.UserRepository
	resetRepo          *repository.PasswordResetRepository
	jwtSecret          string
	allowedEmailDomain string
	resendAPIKey       string
	baseURL            string
}

func NewAuthService(userRepo *repository.UserRepository, resetRepo *repository.PasswordResetRepository, jwtSecret, allowedEmailDomain, resendAPIKey, baseURL string) *AuthService {
	return &AuthService{
		userRepo:           userRepo,
		resetRepo:          resetRepo,
		jwtSecret:          jwtSecret,
		allowedEmailDomain: allowedEmailDomain,
		resendAPIKey:       resendAPIKey,
		baseURL:            baseURL,
	}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) (*dto.AuthResponse, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, errors.New("メール、パスワード、名前は必須です")
	}
	if len(req.Password) < 6 {
		return nil, errors.New("パスワードは6文字以上必要です")
	}
	if s.allowedEmailDomain != "" {
		suffix := "@" + s.allowedEmailDomain
		if !strings.HasSuffix(strings.ToLower(req.Email), strings.ToLower(suffix)) {
			return nil, fmt.Errorf("%s のメールアドレスのみ登録できます", suffix)
		}
	}

	existing, _ := s.userRepo.FindByEmail(ctx, req.Email)
	if existing != nil {
		return nil, errors.New("このメールアドレスは既に登録されています")
	}

	hash, err := pkg.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("パスワードのハッシュ化に失敗しました")
	}

	user := &model.User{
		Email:        req.Email,
		PasswordHash: hash,
		Name:         req.Name,
		Role:         "user",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, errors.New("ユーザーの作成に失敗しました")
	}

	token, err := pkg.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, errors.New("トークンの生成に失敗しました")
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}, nil
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, errors.New("メールとパスワードは必須です")
	}

	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
		}
		return nil, errors.New("ログインに失敗しました")
	}

	if !pkg.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("メールアドレスまたはパスワードが正しくありません")
	}

	token, err := pkg.GenerateToken(user.ID, user.Email, user.Role, s.jwtSecret)
	if err != nil {
		return nil, errors.New("トークンの生成に失敗しました")
	}

	return &dto.AuthResponse{
		Token: token,
		User: dto.UserInfo{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
			Role:  user.Role,
		},
	}, nil
}

// EnsureAdminUser は管理者アカウントが存在しなければ作成する。
// 既に同じメールアドレスのユーザーがいる場合は何もしない（運用中に変更されたパスワードを上書きしないため）。
func (s *AuthService) EnsureAdminUser(ctx context.Context, email, password string) error {
	if email == "" || password == "" {
		return nil
	}
	if len(password) < 8 {
		return errors.New("ADMIN_PASSWORD は8文字以上にしてください")
	}

	email = strings.ToLower(strings.TrimSpace(email))
	existing, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return err
	}
	if existing != nil {
		return nil
	}

	hash, err := pkg.HashPassword(password)
	if err != nil {
		return errors.New("パスワードのハッシュ化に失敗しました")
	}

	return s.userRepo.Create(ctx, &model.User{
		Email:        email,
		PasswordHash: hash,
		Name:         "管理者",
		Role:         "admin",
	})
}

func (s *AuthService) RequestPasswordReset(ctx context.Context, email string) error {
	user, err := s.userRepo.FindByEmail(ctx, email)
	if err != nil || user == nil {
		return nil
	}

	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return errors.New("トークンの生成に失敗しました")
	}
	token := hex.EncodeToString(b)

	if err := s.resetRepo.Create(ctx, user.ID, token, time.Now().Add(time.Hour)); err != nil {
		return errors.New("リセットトークンの保存に失敗しました")
	}

	resetURL := s.baseURL + "/reset-password?token=" + token
	if err := s.sendResetEmail(user.Email, resetURL); err != nil {
		log.Printf("Failed to send reset email to %s: %v", user.Email, err)
	}

	return nil
}

func (s *AuthService) ResetPasswordWithToken(ctx context.Context, token, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("パスワードは6文字以上必要です")
	}

	resetToken, err := s.resetRepo.FindValidToken(ctx, token)
	if err != nil {
		return errors.New("リセットリンクが無効または期限切れです")
	}

	hash, err := pkg.HashPassword(newPassword)
	if err != nil {
		return errors.New("パスワードのハッシュ化に失敗しました")
	}

	if err := s.userRepo.UpdatePassword(ctx, resetToken.UserID, hash); err != nil {
		return errors.New("パスワードの更新に失敗しました")
	}

	s.resetRepo.MarkUsed(ctx, resetToken.ID)
	return nil
}

func (s *AuthService) sendResetEmail(to, resetURL string) error {
	if s.resendAPIKey == "" {
		return errors.New("RESEND_API_KEY is not set")
	}

	body := map[string]any{
		"from":    "rootage <onboarding@resend.dev>",
		"to":      []string{to},
		"subject": "パスワードリセット - rootage",
		"html": fmt.Sprintf(
			`<p>以下のリンクからパスワードを再設定してください。</p><p><a href="%s">パスワードを再設定する</a></p><p>このリンクは1時間有効です。</p><p>心当たりがない場合はこのメールを無視してください。</p>`,
			resetURL,
		),
	}

	jsonBody, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+s.resendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Resend API error: %d", resp.StatusCode)
	}
	return nil
}
