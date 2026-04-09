package service

import (
	"context"
	"errors"
	"time"

	"github.com/rootage-ses-quiz/backend/internal/dto"
	"github.com/rootage-ses-quiz/backend/internal/model"
	"github.com/rootage-ses-quiz/backend/internal/repository"
)

type QuizService struct {
	quizRepo     *repository.QuizRepository
	questionRepo *repository.QuestionRepository
	categoryRepo *repository.CategoryRepository
}

func NewQuizService(quizRepo *repository.QuizRepository, questionRepo *repository.QuestionRepository, categoryRepo *repository.CategoryRepository) *QuizService {
	return &QuizService{quizRepo: quizRepo, questionRepo: questionRepo, categoryRepo: categoryRepo}
}

func (s *QuizService) ListCategories(ctx context.Context) ([]model.Category, error) {
	return s.categoryRepo.List(ctx)
}

func (s *QuizService) StartSession(ctx context.Context, userID, categoryID, difficulty string) (*dto.QuizStartResponse, error) {
	_, err := s.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return nil, errors.New("カテゴリが見つかりません")
	}

	questions, err := s.questionRepo.GetRandomByCategoryAndDifficulty(ctx, categoryID, difficulty, 10)
	if err != nil {
		return nil, errors.New("問題の取得に失敗しました")
	}
	if len(questions) == 0 {
		return nil, errors.New("このカテゴリには問題がありません")
	}

	session := &model.QuizSession{
		UserID:     userID,
		CategoryID: categoryID,
		Difficulty: difficulty,
		Total:      len(questions),
	}
	if err := s.quizRepo.CreateSession(ctx, session); err != nil {
		return nil, errors.New("セッションの作成に失敗しました")
	}

	quizQuestions := make([]dto.QuizQuestion, len(questions))
	for i, q := range questions {
		quizQuestions[i] = dto.QuizQuestion{
			ID:      q.ID,
			Text:    q.Text,
			Choices: q.Choices,
		}
	}

	return &dto.QuizStartResponse{
		SessionID: session.ID,
		Questions: quizQuestions,
	}, nil
}

func (s *QuizService) SubmitAnswer(ctx context.Context, userID string, req dto.AnswerRequest) (*dto.AnswerResponse, error) {
	session, err := s.quizRepo.GetSession(ctx, req.SessionID)
	if err != nil {
		return nil, errors.New("セッションが見つかりません")
	}
	if session.UserID != userID {
		return nil, errors.New("不正なセッションです")
	}
	if session.FinishedAt != nil {
		return nil, errors.New("このセッションは既に終了しています")
	}

	question, err := s.questionRepo.GetByID(ctx, req.QuestionID)
	if err != nil {
		return nil, errors.New("問題が見つかりません")
	}

	isCorrect := req.SelectedIndex == question.CorrectIndex

	answer := &model.QuizAnswer{
		SessionID:     req.SessionID,
		QuestionID:    req.QuestionID,
		SelectedIndex: req.SelectedIndex,
		IsCorrect:     isCorrect,
	}
	if err := s.quizRepo.AddAnswer(ctx, answer); err != nil {
		return nil, errors.New("回答の保存に失敗しました")
	}

	return &dto.AnswerResponse{
		IsCorrect:    isCorrect,
		CorrectIndex: question.CorrectIndex,
		Explanation:  question.Explanation,
	}, nil
}

func (s *QuizService) FinishSession(ctx context.Context, userID, sessionID string) (*dto.QuizResultResponse, error) {
	session, err := s.quizRepo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, errors.New("セッションが見つかりません")
	}
	if session.UserID != userID {
		return nil, errors.New("不正なセッションです")
	}

	answers, err := s.quizRepo.GetSessionAnswers(ctx, sessionID)
	if err != nil {
		return nil, errors.New("回答の取得に失敗しました")
	}

	score := 0
	for _, a := range answers {
		if a.IsCorrect {
			score++
		}
	}

	if err := s.quizRepo.FinishSession(ctx, sessionID, score); err != nil {
		return nil, errors.New("セッションの終了に失敗しました")
	}

	// 回答に紐づく問題を一括取得（N+1回避）
	qIDs := make([]string, len(answers))
	for i, a := range answers {
		qIDs[i] = a.QuestionID
	}
	questionMap, err := s.questionRepo.GetByIDs(ctx, qIDs)
	if err != nil {
		return nil, errors.New("問題データの取得に失敗しました")
	}

	answerDetails := make([]dto.AnswerDetail, len(answers))
	for i, a := range answers {
		detail := dto.AnswerDetail{
			QuestionID:    a.QuestionID,
			SelectedIndex: a.SelectedIndex,
			IsCorrect:     a.IsCorrect,
		}
		if q, ok := questionMap[a.QuestionID]; ok {
			detail.QuestionText = q.Text
			detail.Choices = q.Choices
			detail.CorrectIndex = q.CorrectIndex
			detail.Explanation = q.Explanation
		}
		answerDetails[i] = detail
	}

	total := session.Total
	var percentage float64
	if total > 0 {
		percentage = float64(score) / float64(total) * 100
	}

	return &dto.QuizResultResponse{
		SessionID:  sessionID,
		Score:      score,
		Total:      total,
		Percentage: percentage,
		Answers:    answerDetails,
	}, nil
}

func (s *QuizService) GetHistory(ctx context.Context, userID string) ([]dto.QuizHistoryItem, error) {
	sessions, err := s.quizRepo.GetSessionsByUser(ctx, userID)
	if err != nil {
		return nil, errors.New("履歴の取得に失敗しました")
	}

	items := make([]dto.QuizHistoryItem, len(sessions))
	for i, sess := range sessions {
		var categoryName string
		cat, err := s.categoryRepo.GetByID(ctx, sess.CategoryID)
		if err == nil {
			categoryName = cat.Name
		}

		var percentage float64
		if sess.Total > 0 {
			percentage = float64(sess.Score) / float64(sess.Total) * 100
		}

		items[i] = dto.QuizHistoryItem{
			SessionID:    sess.ID,
			CategoryID:   sess.CategoryID,
			CategoryName: categoryName,
			Score:        sess.Score,
			Total:        sess.Total,
			Percentage:   percentage,
			StartedAt:    sess.StartedAt,
			FinishedAt:   sess.FinishedAt,
		}
	}

	return items, nil
}

func (s *QuizService) GetStats(ctx context.Context, userID string) (*dto.QuizStatsResponse, error) {
	totalSessions, err := s.quizRepo.GetTotalSessions(ctx, userID)
	if err != nil {
		return nil, errors.New("統計情報の取得に失敗しました")
	}

	categoryStats, err := s.quizRepo.GetCategoryStats(ctx, userID)
	if err != nil {
		return nil, errors.New("カテゴリ別統計の取得に失敗しました")
	}

	items := make([]dto.CategoryStatsItem, len(categoryStats))
	weakCategory := ""
	lowestPct := 101.0

	for i, s := range categoryStats {
		var pct float64
		if s.TotalAnswers > 0 {
			pct = float64(s.Correct) / float64(s.TotalAnswers) * 100
		}
		items[i] = dto.CategoryStatsItem{
			CategoryID:   s.CategoryID,
			CategoryName: s.CategoryName,
			TotalAnswers: s.TotalAnswers,
			Correct:      s.Correct,
			Percentage:   pct,
		}
		if pct < lowestPct {
			lowestPct = pct
			weakCategory = s.CategoryName
		}
	}

	return &dto.QuizStatsResponse{
		TotalSessions: totalSessions,
		CategoryStats: items,
		WeakCategory:  weakCategory,
	}, nil
}

func (s *QuizService) GetDashboard(ctx context.Context, userID string) (*dto.UserDashboardResponse, error) {
	totalSessions, err := s.quizRepo.GetTotalSessions(ctx, userID)
	if err != nil {
		return nil, errors.New("統計情報の取得に失敗しました")
	}

	totalQuestions, totalCorrect, err := s.quizRepo.GetTotalQuestionsAndCorrect(ctx, userID)
	if err != nil {
		return nil, errors.New("回答集計の取得に失敗しました")
	}

	var overallPct float64
	if totalQuestions > 0 {
		overallPct = float64(totalCorrect) / float64(totalQuestions) * 100
	}

	categoryStats, err := s.quizRepo.GetCategoryStats(ctx, userID)
	if err != nil {
		return nil, errors.New("カテゴリ別統計の取得に失敗しました")
	}

	items := make([]dto.CategoryStatsItem, len(categoryStats))
	weakCategory := ""
	lowestPct := 101.0
	for i, cs := range categoryStats {
		var pct float64
		if cs.TotalAnswers > 0 {
			pct = float64(cs.Correct) / float64(cs.TotalAnswers) * 100
		}
		items[i] = dto.CategoryStatsItem{
			CategoryID:   cs.CategoryID,
			CategoryName: cs.CategoryName,
			TotalAnswers: cs.TotalAnswers,
			Correct:      cs.Correct,
			Percentage:   pct,
		}
		if pct < lowestPct {
			lowestPct = pct
			weakCategory = cs.CategoryName
		}
	}

	dailyRaw, err := s.quizRepo.GetDailyStats(ctx, userID)
	if err != nil {
		return nil, errors.New("日別統計の取得に失敗しました")
	}

	dailyStats := make([]dto.DailyStatsItem, len(dailyRaw))
	for i, d := range dailyRaw {
		var pct float64
		if d.Total > 0 {
			pct = float64(d.Correct) / float64(d.Total) * 100
		}
		dailyStats[i] = dto.DailyStatsItem{
			Date:       d.Date,
			Total:      d.Total,
			Correct:    d.Correct,
			Percentage: pct,
		}
	}

	// 連続学習日数の計算
	streak := calcStreak(dailyRaw)

	// アチーブメント（カテゴリ×難易度のベストスコア）
	bestScores, err := s.quizRepo.GetBestScores(ctx, userID)
	if err != nil {
		return nil, errors.New("アチーブメントの取得に失敗しました")
	}

	achievements := make([]dto.AchievementItem, len(bestScores))
	for i, b := range bestScores {
		achievements[i] = dto.AchievementItem{
			CategoryID:   b.CategoryID,
			CategoryName: b.CategoryName,
			Difficulty:   b.Difficulty,
			BestScore:    b.BestScore,
			Total:        b.Total,
			Percentage:   b.Percentage,
			Cleared:      b.Percentage >= 80,
		}
	}

	return &dto.UserDashboardResponse{
		TotalSessions:     totalSessions,
		TotalQuestions:    totalQuestions,
		TotalCorrect:     totalCorrect,
		OverallPercentage: overallPct,
		CategoryStats:    items,
		WeakCategory:     weakCategory,
		DailyStats:       dailyStats,
		CurrentStreak:    streak,
		Achievements:     achievements,
	}, nil
}

func calcStreak(dailyStats []repository.DailyStats) int {
	if len(dailyStats) == 0 {
		return 0
	}

	today := time.Now().Format("2006-01-02")
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	last := dailyStats[len(dailyStats)-1].Date
	if last != today && last != yesterday {
		return 0
	}

	streak := 1
	for i := len(dailyStats) - 2; i >= 0; i-- {
		cur, _ := time.Parse("2006-01-02", dailyStats[i+1].Date)
		prev, _ := time.Parse("2006-01-02", dailyStats[i].Date)
		if cur.Sub(prev).Hours() == 24 {
			streak++
		} else {
			break
		}
	}
	return streak
}

func (s *QuizService) GetReview(ctx context.Context, userID string) ([]dto.ReviewItem, error) {
	details, err := s.quizRepo.GetWrongAnswers(ctx, userID)
	if err != nil {
		return nil, errors.New("復習データの取得に失敗しました")
	}

	seen := make(map[string]bool)
	var items []dto.ReviewItem
	for _, d := range details {
		if seen[d.QuestionID] {
			continue
		}
		seen[d.QuestionID] = true
		items = append(items, dto.ReviewItem{
			QuestionID:   d.QuestionID,
			QuestionText: d.QuestionText,
			Choices:      d.Choices,
			CorrectIndex: d.CorrectIndex,
			Explanation:  d.Explanation,
			CategoryName: d.CategoryName,
		})
	}

	return items, nil
}
