package handler

import (
	"net/http"

	"github.com/rootage-ses-quiz/backend/internal/dto"
	"github.com/rootage-ses-quiz/backend/internal/middleware"
	"github.com/rootage-ses-quiz/backend/internal/pkg"
	"github.com/rootage-ses-quiz/backend/internal/service"
)

type RAGHandler struct {
	ragSvc *service.RAGService
}

func NewRAGHandler(ragSvc *service.RAGService) *RAGHandler {
	return &RAGHandler{ragSvc: ragSvc}
}

func (h *RAGHandler) Ask(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	var req dto.RAGAskRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	if req.Question == "" {
		pkg.WriteError(w, http.StatusBadRequest, "質問を入力してください")
		return
	}

	resp, err := h.ragSvc.Ask(r.Context(), userID, req.Question, req.History)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "回答の生成に失敗しました")
		return
	}

	pkg.WriteJSON(w, http.StatusOK, resp)
}

func (h *RAGHandler) History(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())

	history, err := h.ragSvc.GetHistory(r.Context(), userID)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "履歴の取得に失敗しました")
		return
	}

	pkg.WriteJSON(w, http.StatusOK, history)
}
