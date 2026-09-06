package service

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/rootage-ses-quiz/backend/internal/dto"
	"github.com/rootage-ses-quiz/backend/internal/model"
	"github.com/rootage-ses-quiz/backend/internal/pkg"
	"github.com/rootage-ses-quiz/backend/internal/repository"
)

type AdminService struct {
	categoryRepo *repository.CategoryRepository
	questionRepo *repository.QuestionRepository
	documentRepo *repository.DocumentRepository
	guideRepo    *repository.GuideRepository
	userRepo     *repository.UserRepository
	openaiClient *OpenAIClient
}

func NewAdminService(categoryRepo *repository.CategoryRepository, questionRepo *repository.QuestionRepository, documentRepo *repository.DocumentRepository, guideRepo *repository.GuideRepository, userRepo *repository.UserRepository, openaiClient *OpenAIClient) *AdminService {
	return &AdminService{
		categoryRepo: categoryRepo,
		questionRepo: questionRepo,
		documentRepo: documentRepo,
		guideRepo:    guideRepo,
		userRepo:     userRepo,
		openaiClient: openaiClient,
	}
}

// Category CRUD

func (s *AdminService) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.categoryRepo.List(ctx)
}

func (s *AdminService) CreateCategory(ctx context.Context, req dto.CreateCategoryRequest) (*model.Category, error) {
	cat := &model.Category{Name: req.Name, Description: req.Description}
	if err := s.categoryRepo.Create(ctx, cat); err != nil {
		return nil, fmt.Errorf("カテゴリの作成に失敗しました: %w", err)
	}
	return cat, nil
}

func (s *AdminService) UpdateCategory(ctx context.Context, id string, req dto.UpdateCategoryRequest) error {
	cat := &model.Category{ID: id, Name: req.Name, Description: req.Description}
	return s.categoryRepo.Update(ctx, cat)
}

// DeleteCategory は論理削除（is_active=false）にする。
// 物理削除すると quiz_sessions / quiz_answers が CASCADE で巻き込まれ、
// 全ユーザーの学習履歴が消えるため。
func (s *AdminService) DeleteCategory(ctx context.Context, id string) error {
	return s.categoryRepo.Deactivate(ctx, id)
}

// Question CRUD

func (s *AdminService) ListQuestions(ctx context.Context, categoryID string) ([]model.Question, error) {
	return s.questionRepo.List(ctx, categoryID)
}

func (s *AdminService) CreateQuestion(ctx context.Context, req dto.CreateQuestionRequest) (*model.Question, error) {
	q := &model.Question{
		CategoryID:   req.CategoryID,
		Text:         req.Text,
		Choices:      req.Choices,
		CorrectIndex: req.CorrectIndex,
		Explanation:  req.Explanation,
		Difficulty:   req.Difficulty,
	}
	if err := s.questionRepo.Create(ctx, q); err != nil {
		return nil, fmt.Errorf("問題の作成に失敗しました: %w", err)
	}
	return q, nil
}

func (s *AdminService) UpdateQuestion(ctx context.Context, id string, req dto.UpdateQuestionRequest) error {
	q := &model.Question{
		ID:           id,
		Text:         req.Text,
		Choices:      req.Choices,
		CorrectIndex: req.CorrectIndex,
		Explanation:  req.Explanation,
		Difficulty:   req.Difficulty,
	}
	return s.questionRepo.Update(ctx, q)
}

// DeleteQuestion は論理削除（is_active=false）にする。
// 物理削除すると quiz_answers が CASCADE で消え、回答履歴と正答率が壊れるため。
func (s *AdminService) DeleteQuestion(ctx context.Context, id string) error {
	return s.questionRepo.Deactivate(ctx, id)
}

func (s *AdminService) ImportQuestions(ctx context.Context, file io.Reader) (*dto.ImportQuestionsResponse, error) {
	// カテゴリ一覧を取得して名前→IDのマッピングを作成
	categories, err := s.categoryRepo.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("カテゴリの取得に失敗しました: %w", err)
	}
	categoryMap := make(map[string]string)
	for _, c := range categories {
		categoryMap[c.Name] = c.ID
	}

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("CSVの読み込みに失敗しました: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("CSVにデータ行がありません")
	}

	correctMap := map[string]int{"A": 0, "B": 1, "C": 2, "D": 3}

	var imported int
	var importErrors []dto.ImportError

	// 1行目はヘッダーなのでスキップ
	for i, row := range records[1:] {
		rowNum := i + 2 // ヘッダー行 + 0-indexed

		if len(row) < 8 {
			importErrors = append(importErrors, dto.ImportError{Row: rowNum, Message: "列数が不足しています（8列必要）"})
			continue
		}

		categoryName := strings.TrimSpace(row[0])
		text := strings.TrimSpace(row[1])
		choiceA := strings.TrimSpace(row[2])
		choiceB := strings.TrimSpace(row[3])
		choiceC := strings.TrimSpace(row[4])
		choiceD := strings.TrimSpace(row[5])
		correct := strings.ToUpper(strings.TrimSpace(row[6]))
		explanation := strings.TrimSpace(row[7])
		difficulty := "beginner"
		if len(row) >= 9 && strings.TrimSpace(row[8]) != "" {
			difficulty = strings.TrimSpace(row[8])
		}

		// バリデーション
		categoryID, ok := categoryMap[categoryName]
		if !ok {
			importErrors = append(importErrors, dto.ImportError{Row: rowNum, Message: fmt.Sprintf("カテゴリ「%s」が見つかりません", categoryName)})
			continue
		}

		if !isValidDifficulty(difficulty) {
			importErrors = append(importErrors, dto.ImportError{Row: rowNum, Message: fmt.Sprintf("難易度「%s」が不正です（beginner / intermediate / advanced のいずれか）", difficulty)})
			continue
		}

		if text == "" {
			importErrors = append(importErrors, dto.ImportError{Row: rowNum, Message: "問題文が空です"})
			continue
		}

		correctIndex, ok := correctMap[correct]
		if !ok {
			importErrors = append(importErrors, dto.ImportError{Row: rowNum, Message: fmt.Sprintf("正解「%s」が不正です（A/B/C/Dで指定してください）", correct)})
			continue
		}

		choices, err := json.Marshal([]string{choiceA, choiceB, choiceC, choiceD})
		if err != nil {
			importErrors = append(importErrors, dto.ImportError{Row: rowNum, Message: "選択肢のJSON変換に失敗しました"})
			continue
		}

		q := &model.Question{
			CategoryID:   categoryID,
			Text:         text,
			Choices:      choices,
			CorrectIndex: correctIndex,
			Explanation:  explanation,
			Difficulty:   difficulty,
		}
		if err := s.questionRepo.Create(ctx, q); err != nil {
			importErrors = append(importErrors, dto.ImportError{Row: rowNum, Message: "保存に失敗しました"})
			continue
		}
		imported++
	}

	return &dto.ImportQuestionsResponse{
		Imported: imported,
		Errors:   importErrors,
	}, nil
}

func isValidDifficulty(d string) bool {
	switch d {
	case "beginner", "intermediate", "advanced":
		return true
	default:
		return false
	}
}

// Document management

func (s *AdminService) ListDocuments(ctx context.Context) ([]model.Document, error) {
	return s.documentRepo.List(ctx)
}

// ListUnindexedDocuments はチャンク未作成の文書だけを返す（起動時の一括処理用）。
func (s *AdminService) ListUnindexedDocuments(ctx context.Context) ([]model.Document, error) {
	return s.documentRepo.ListWithoutChunks(ctx)
}

func (s *AdminService) UploadDocument(ctx context.Context, title, filename string, file io.Reader) (*model.Document, error) {
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("ファイルの読み込みに失敗しました: %w", err)
	}
	return s.createDocumentWithIndex(ctx, title, filename, string(content))
}

func (s *AdminService) CreateDocumentFromText(ctx context.Context, title, content string) (*model.Document, error) {
	if title == "" || content == "" {
		return nil, fmt.Errorf("タイトルと内容は必須です")
	}
	return s.createDocumentWithIndex(ctx, title, "text-input", content)
}

// createDocumentWithIndex は文書とそのチャンクを1トランザクションで作成する。
// 埋め込み生成を先に済ませ、全て成功したときだけ documents 行を確定するため、
// ベクトル化失敗で「本文だけ残ってチャンク0件」の孤児文書が生じない。
func (s *AdminService) createDocumentWithIndex(ctx context.Context, title, filename, content string) (*model.Document, error) {
	chunks, err := s.buildChunks(ctx, "", content)
	if err != nil {
		return nil, fmt.Errorf("インデックスの作成に失敗しました: %w", err)
	}

	tx, err := s.documentRepo.BeginTx(ctx)
	if err != nil {
		return nil, fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
	}
	defer tx.Rollback()

	doc := &model.Document{Title: title, Filename: filename, Content: content}
	if err := s.documentRepo.CreateWith(ctx, tx, doc); err != nil {
		return nil, fmt.Errorf("文書の保存に失敗しました: %w", err)
	}
	for i := range chunks {
		chunks[i].DocumentID = doc.ID
	}
	if err := s.documentRepo.CreateChunks(ctx, tx, chunks); err != nil {
		return nil, fmt.Errorf("インデックスの作成に失敗しました: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("文書の保存に失敗しました: %w", err)
	}
	return doc, nil
}

func (s *AdminService) DeleteDocument(ctx context.Context, id string) error {
	return s.documentRepo.Delete(ctx, id)
}

// ReindexDocument は既存チャンクを新しいチャンクで置き換える。
// 先に全チャンクの埋め込みを生成し（DB未変更）、成功後にトランザクションで
// 削除と挿入をまとめて行う。OpenAI が途中でエラーを返しても既存チャンクは無傷。
func (s *AdminService) ReindexDocument(ctx context.Context, id string) error {
	doc, err := s.documentRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("文書が見つかりません: %w", err)
	}

	chunks, err := s.buildChunks(ctx, doc.ID, doc.Content)
	if err != nil {
		return err
	}

	tx, err := s.documentRepo.BeginTx(ctx)
	if err != nil {
		return fmt.Errorf("トランザクションの開始に失敗しました: %w", err)
	}
	defer tx.Rollback()

	if err := s.documentRepo.DeleteChunksByDocument(ctx, tx, id); err != nil {
		return fmt.Errorf("既存チャンクの削除に失敗しました: %w", err)
	}
	if err := s.documentRepo.CreateChunks(ctx, tx, chunks); err != nil {
		return fmt.Errorf("チャンクの保存に失敗しました: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("チャンクの保存に失敗しました: %w", err)
	}
	return nil
}

// buildChunks は本文をチャンクに分割し、各チャンクの埋め込みを生成する。
// DB には一切触れないため、途中で失敗しても既存データに影響しない。
func (s *AdminService) buildChunks(ctx context.Context, docID, content string) ([]model.DocumentChunk, error) {
	textChunks := ChunkText(content, 1000)
	chunks := make([]model.DocumentChunk, 0, len(textChunks))
	for i, text := range textChunks {
		embedding, err := s.openaiClient.GenerateEmbedding(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("チャンク%dのベクトル化に失敗しました: %w", i+1, err)
		}
		chunks = append(chunks, model.DocumentChunk{
			DocumentID: docID,
			ChunkIndex: i,
			Content:    text,
			Embedding:  embedding,
		})
	}
	return chunks, nil
}

// Guide Category management

func (s *AdminService) ListGuideCategories(ctx context.Context) ([]model.GuideCategory, error) {
	return s.guideRepo.ListCategories(ctx)
}

func (s *AdminService) CreateGuideCategory(ctx context.Context, req dto.CreateGuideCategoryRequest) (*model.GuideCategory, error) {
	c := &model.GuideCategory{Name: req.Name, SortOrder: req.SortOrder}
	if err := s.guideRepo.CreateCategory(ctx, c); err != nil {
		return nil, fmt.Errorf("ガイドカテゴリの作成に失敗しました: %w", err)
	}
	return c, nil
}

func (s *AdminService) UpdateGuideCategory(ctx context.Context, id string, req dto.UpdateGuideCategoryRequest) error {
	c := &model.GuideCategory{ID: id, Name: req.Name, SortOrder: req.SortOrder}
	return s.guideRepo.UpdateCategory(ctx, c)
}

func (s *AdminService) DeleteGuideCategory(ctx context.Context, id string) error {
	return s.guideRepo.DeleteCategory(ctx, id)
}

// Guide management

func (s *AdminService) ListGuides(ctx context.Context) ([]model.Guide, error) {
	return s.guideRepo.ListAll(ctx)
}

func (s *AdminService) ListPublishedGuides(ctx context.Context) ([]model.Guide, error) {
	return s.guideRepo.ListPublished(ctx)
}

func (s *AdminService) GetGuide(ctx context.Context, id string) (*model.Guide, error) {
	return s.guideRepo.GetByID(ctx, id)
}

func (s *AdminService) CreateGuide(ctx context.Context, req dto.CreateGuideRequest) (*model.Guide, error) {
	g := &model.Guide{
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Content:     req.Content,
		IsPublished: req.IsPublished,
		SortOrder:   req.SortOrder,
	}
	if err := s.guideRepo.Create(ctx, g); err != nil {
		return nil, fmt.Errorf("ガイドの作成に失敗しました: %w", err)
	}
	return g, nil
}

func (s *AdminService) UpdateGuide(ctx context.Context, id string, req dto.UpdateGuideRequest) error {
	g := &model.Guide{
		ID:          id,
		CategoryID:  req.CategoryID,
		Title:       req.Title,
		Content:     req.Content,
		IsPublished: req.IsPublished,
		SortOrder:   req.SortOrder,
	}
	return s.guideRepo.Update(ctx, g)
}

func (s *AdminService) DeleteGuide(ctx context.Context, id string) error {
	return s.guideRepo.Delete(ctx, id)
}

func (s *AdminService) ResetPassword(ctx context.Context, userID, newPassword string) error {
	if err := validatePassword(newPassword); err != nil {
		return err
	}
	hash, err := pkg.HashPassword(newPassword)
	if err != nil {
		return fmt.Errorf("パスワードのハッシュ化に失敗しました")
	}
	return s.userRepo.UpdatePassword(ctx, userID, hash)
}

// Admin Summary

func (s *AdminService) GetSummary(ctx context.Context) (*dto.AdminSummary, error) {
	var summary dto.AdminSummary
	db := s.categoryRepo.DB()

	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM users WHERE role = 'user'`).Scan(&summary.TotalUsers); err != nil {
		return nil, fmt.Errorf("ユーザー数の取得に失敗しました: %w", err)
	}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM questions`).Scan(&summary.TotalQuestions); err != nil {
		return nil, fmt.Errorf("問題数の取得に失敗しました: %w", err)
	}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM documents`).Scan(&summary.TotalDocuments); err != nil {
		return nil, fmt.Errorf("文書数の取得に失敗しました: %w", err)
	}
	if err := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM guides WHERE is_published = true`).Scan(&summary.TotalGuides); err != nil {
		return nil, fmt.Errorf("ガイド数の取得に失敗しました: %w", err)
	}

	return &summary, nil
}

// User Progress

func (s *AdminService) GetUserProgress(ctx context.Context) ([]dto.UserProgressItem, error) {
	summaries, err := s.userRepo.ListUserSummaries(ctx)
	if err != nil {
		return nil, fmt.Errorf("ユーザー一覧の取得に失敗しました: %w", err)
	}

	allStats, err := s.userRepo.ListAllCategoryStats(ctx)
	if err != nil {
		return nil, fmt.Errorf("カテゴリ別統計の取得に失敗しました: %w", err)
	}

	// user_id → []stats のマップを構築
	statsMap := make(map[string][]dto.UserCategoryStatsItem)
	for _, s := range allStats {
		var pct float64
		if s.TotalAnswers > 0 {
			pct = float64(s.Correct) / float64(s.TotalAnswers) * 100
		}
		statsMap[s.UserID] = append(statsMap[s.UserID], dto.UserCategoryStatsItem{
			CategoryID:   s.CategoryID,
			CategoryName: s.CategoryName,
			TotalAnswers: s.TotalAnswers,
			Correct:      s.Correct,
			Percentage:   pct,
		})
	}

	items := make([]dto.UserProgressItem, len(summaries))
	for i, u := range summaries {
		var pct float64
		if u.TotalQuestions > 0 {
			pct = float64(u.TotalCorrect) / float64(u.TotalQuestions) * 100
		}
		items[i] = dto.UserProgressItem{
			ID:             u.ID,
			Name:           u.Name,
			Email:          u.Email,
			CreatedAt:      u.CreatedAt,
			LastQuizAt:     u.LastQuizAt,
			TotalCorrect:   u.TotalCorrect,
			TotalQuestions: u.TotalQuestions,
			Percentage:     pct,
			CategoryStats:  statsMap[u.ID],
		}
	}

	return items, nil
}
