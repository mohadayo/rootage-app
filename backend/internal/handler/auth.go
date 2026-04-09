package handler

import (
	"net/http"

	"github.com/rootage-ses-quiz/backend/internal/dto"
	"github.com/rootage-ses-quiz/backend/internal/pkg"
	"github.com/rootage-ses-quiz/backend/internal/service"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	resp, err := h.authSvc.Register(r.Context(), req)
	if err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, resp)
}

func (h *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ForgotPasswordRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	h.authSvc.RequestPasswordReset(r.Context(), req.Email)
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "登録されているメールアドレスの場合、リセット用のメールを送信しました"})
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	if err := h.authSvc.ResetPasswordWithToken(r.Context(), req.Token, req.Password); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "パスワードを再設定しました"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	resp, err := h.authSvc.Login(r.Context(), req)
	if err != nil {
		pkg.WriteError(w, http.StatusUnauthorized, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, resp)
}
