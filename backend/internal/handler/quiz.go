package handler

import (
	"net/http"

	"github.com/rootage-ses-quiz/backend/internal/dto"
	"github.com/rootage-ses-quiz/backend/internal/middleware"
	"github.com/rootage-ses-quiz/backend/internal/pkg"
	"github.com/rootage-ses-quiz/backend/internal/service"
)

type QuizHandler struct {
	quizSvc *service.QuizService
}

func NewQuizHandler(quizSvc *service.QuizService) *QuizHandler {
	return &QuizHandler{quizSvc: quizSvc}
}

func (h *QuizHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.quizSvc.ListCategories(r.Context())
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "カテゴリの取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, categories)
}

func (h *QuizHandler) Start(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.QuizStartRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	resp, err := h.quizSvc.StartSession(r.Context(), userID, req.CategoryID, req.Difficulty)
	if err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, resp)
}

func (h *QuizHandler) Answer(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.AnswerRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	resp, err := h.quizSvc.SubmitAnswer(r.Context(), userID, req)
	if err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, resp)
}

func (h *QuizHandler) Finish(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.QuizFinishRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	resp, err := h.quizSvc.FinishSession(r.Context(), userID, req.SessionID)
	if err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, resp)
}

func (h *QuizHandler) History(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	history, err := h.quizSvc.GetHistory(r.Context(), userID)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, history)
}

func (h *QuizHandler) Stats(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	stats, err := h.quizSvc.GetStats(r.Context(), userID)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, stats)
}

func (h *QuizHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	dashboard, err := h.quizSvc.GetDashboard(r.Context(), userID)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, dashboard)
}

func (h *QuizHandler) Review(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	items, err := h.quizSvc.GetReview(r.Context(), userID)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, items)
}
