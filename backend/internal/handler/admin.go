package handler

import (
	"net/http"

	"github.com/rootage-ses-quiz/backend/internal/dto"
	"github.com/rootage-ses-quiz/backend/internal/pkg"
	"github.com/rootage-ses-quiz/backend/internal/service"
)

type AdminHandler struct {
	adminSvc *service.AdminService
}

func NewAdminHandler(adminSvc *service.AdminService) *AdminHandler {
	return &AdminHandler{adminSvc: adminSvc}
}

// Category endpoints

func (h *AdminHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.adminSvc.ListCategories(r.Context())
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "カテゴリの取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, categories)
}

func (h *AdminHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateCategoryRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	cat, err := h.adminSvc.CreateCategory(r.Context(), req)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, cat)
}

func (h *AdminHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateCategoryRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	if err := h.adminSvc.UpdateCategory(r.Context(), id, req); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "更新しました"})
}

func (h *AdminHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.adminSvc.DeleteCategory(r.Context(), id); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "削除しました"})
}

// Question endpoints

func (h *AdminHandler) ListQuestions(w http.ResponseWriter, r *http.Request) {
	categoryID := r.URL.Query().Get("category_id")
	questions, err := h.adminSvc.ListQuestions(r.Context(), categoryID)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "問題の取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, questions)
}

func (h *AdminHandler) CreateQuestion(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateQuestionRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	q, err := h.adminSvc.CreateQuestion(r.Context(), req)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, q)
}

func (h *AdminHandler) UpdateQuestion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateQuestionRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	if err := h.adminSvc.UpdateQuestion(r.Context(), id, req); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "更新しました"})
}

func (h *AdminHandler) DeleteQuestion(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.adminSvc.DeleteQuestion(r.Context(), id); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "削除しました"})
}

// Question import

func (h *AdminHandler) ImportQuestions(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "ファイルサイズが大きすぎます")
		return
	}

	file, _, err := r.FormFile("file")
	if err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "CSVファイルが必要です")
		return
	}
	defer file.Close()

	result, err := h.adminSvc.ImportQuestions(r.Context(), file)
	if err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusOK, result)
}

// Document endpoints

func (h *AdminHandler) CreateDocumentFromText(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateDocumentTextRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}

	doc, err := h.adminSvc.CreateDocumentFromText(r.Context(), req.Title, req.Content)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, doc)
}

func (h *AdminHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	docs, err := h.adminSvc.ListDocuments(r.Context())
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "文書の取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, docs)
}

func (h *AdminHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil { // 10MB limit
		pkg.WriteError(w, http.StatusBadRequest, "ファイルサイズが大きすぎます")
		return
	}

	title := r.FormValue("title")
	if title == "" {
		pkg.WriteError(w, http.StatusBadRequest, "タイトルは必須です")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "ファイルが必要です")
		return
	}
	defer file.Close()

	doc, err := h.adminSvc.UploadDocument(r.Context(), title, header.Filename, file)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	pkg.WriteJSON(w, http.StatusCreated, doc)
}

func (h *AdminHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.adminSvc.DeleteDocument(r.Context(), id); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "削除しました"})
}

func (h *AdminHandler) ReindexDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.adminSvc.ReindexDocument(r.Context(), id); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "再インデックスが完了しました"})
}

// Admin Summary

func (h *AdminHandler) GetSummary(w http.ResponseWriter, r *http.Request) {
	summary, err := h.adminSvc.GetSummary(r.Context())
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "サマリーの取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, summary)
}

// User Progress

func (h *AdminHandler) GetUserProgress(w http.ResponseWriter, r *http.Request) {
	progress, err := h.adminSvc.GetUserProgress(r.Context())
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "ユーザー進捗の取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, progress)
}

// Password Reset

func (h *AdminHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	userID := r.PathValue("id")
	var req struct {
		Password string `json:"password"`
	}
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}
	if err := h.adminSvc.ResetPassword(r.Context(), userID, req.Password); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "パスワードをリセットしました"})
}

// Guide Category endpoints

func (h *AdminHandler) ListGuideCategories(w http.ResponseWriter, r *http.Request) {
	cats, err := h.adminSvc.ListGuideCategories(r.Context())
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "ガイドカテゴリの取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, cats)
}

func (h *AdminHandler) CreateGuideCategory(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGuideCategoryRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}
	cat, err := h.adminSvc.CreateGuideCategory(r.Context(), req)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusCreated, cat)
}

func (h *AdminHandler) UpdateGuideCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateGuideCategoryRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}
	if err := h.adminSvc.UpdateGuideCategory(r.Context(), id, req); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "更新しました"})
}

func (h *AdminHandler) DeleteGuideCategory(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.adminSvc.DeleteGuideCategory(r.Context(), id); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "削除しました"})
}

// Guide endpoints

func (h *AdminHandler) ListGuides(w http.ResponseWriter, r *http.Request) {
	guides, err := h.adminSvc.ListGuides(r.Context())
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "ガイドの取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, guides)
}

func (h *AdminHandler) ListPublishedGuides(w http.ResponseWriter, r *http.Request) {
	guides, err := h.adminSvc.ListPublishedGuides(r.Context())
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, "ガイドの取得に失敗しました")
		return
	}
	pkg.WriteJSON(w, http.StatusOK, guides)
}

func (h *AdminHandler) GetGuide(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	guide, err := h.adminSvc.GetGuide(r.Context(), id)
	if err != nil {
		pkg.WriteError(w, http.StatusNotFound, "ガイドが見つかりません")
		return
	}
	if !guide.IsPublished {
		// 管理者以外は下書きを見れない
		role := r.Context().Value("role")
		if role != "admin" {
			pkg.WriteError(w, http.StatusNotFound, "ガイドが見つかりません")
			return
		}
	}
	pkg.WriteJSON(w, http.StatusOK, guide)
}

func (h *AdminHandler) CreateGuide(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateGuideRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}
	guide, err := h.adminSvc.CreateGuide(r.Context(), req)
	if err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusCreated, guide)
}

func (h *AdminHandler) UpdateGuide(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	var req dto.UpdateGuideRequest
	if err := pkg.ReadJSON(r, &req); err != nil {
		pkg.WriteError(w, http.StatusBadRequest, "リクエストの形式が正しくありません")
		return
	}
	if err := h.adminSvc.UpdateGuide(r.Context(), id, req); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "更新しました"})
}

func (h *AdminHandler) DeleteGuide(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.adminSvc.DeleteGuide(r.Context(), id); err != nil {
		pkg.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	pkg.WriteJSON(w, http.StatusOK, map[string]string{"message": "削除しました"})
}
