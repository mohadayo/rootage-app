package service

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rootage-ses-quiz/backend/internal/dto"
	"github.com/rootage-ses-quiz/backend/internal/model"
	"github.com/rootage-ses-quiz/backend/internal/repository"
)

type RAGService struct {
	docRepo      *repository.DocumentRepository
	chatRepo     *repository.ChatRepository
	openaiClient *OpenAIClient
}

func NewRAGService(docRepo *repository.DocumentRepository, chatRepo *repository.ChatRepository, openaiClient *OpenAIClient) *RAGService {
	return &RAGService{docRepo: docRepo, chatRepo: chatRepo, openaiClient: openaiClient}
}

const (
	maxQuestionChars       = 4000
	maxHistoryContentChars = 4000
)

func (s *RAGService) Ask(ctx context.Context, userID, question string, history []dto.RAGHistoryMessage) (*dto.RAGResponse, error) {
	// 質問文の長さを制限する（青天井の OpenAI 課金を防ぐ）。
	question = truncateRunes(question, maxQuestionChars)

	// Generate embedding for the question
	embedding, err := s.openaiClient.GenerateEmbedding(ctx, question)
	if err != nil {
		return nil, fmt.Errorf("質問のベクトル化に失敗しました: %w", err)
	}

	// Search similar chunks (8件に増加して精度向上)
	chunks, err := s.docRepo.SearchSimilarWithTitle(ctx, embedding, 8)
	if err != nil {
		return nil, fmt.Errorf("類似検索に失敗しました: %w", err)
	}

	if len(chunks) == 0 {
		return &dto.RAGResponse{
			Answer:  "申し訳ありませんが、関連する情報が見つかりませんでした。ナレッジベースに文書が登録されているか確認してください。",
			Sources: []dto.RAGSource{},
		}, nil
	}

	// Build context from chunks
	var contextParts []string
	sources := make([]dto.RAGSource, len(chunks))
	for i, chunk := range chunks {
		contextParts = append(contextParts, fmt.Sprintf("[文書: %s, チャンク%d]\n%s", chunk.DocumentTitle, chunk.ChunkIndex+1, chunk.Content))
		sources[i] = dto.RAGSource{
			DocumentID:    chunk.DocumentID,
			DocumentTitle: chunk.DocumentTitle,
			ChunkIndex:    chunk.ChunkIndex,
			Content:       chunk.Content,
		}
	}

	systemPrompt := fmt.Sprintf(`あなたは社内ナレッジアシスタントです。社員が上司や総務に確認するような質問に、正確かつ丁寧に回答してください。

## 回答ルール
- 必ず以下のナレッジ情報に基づいて回答すること
- ナレッジにない情報は推測せず、「この情報はナレッジベースに登録されていません。総務または上司にご確認ください。」と回答すること
- 回答は具体的で実用的にすること（手続き方法、連絡先、期限など）
- 箇条書きや段落分けで読みやすく整理すること
- 回答の最後に参照した文書名を記載すること

--- 社内ナレッジ ---
%s`, strings.Join(contextParts, "\n\n"))

	// 会話履歴を含めたメッセージ構築
	var chatMessages []chatMessage
	// 過去の会話を最大5往復まで含める
	maxHistory := 10
	if len(history) > maxHistory {
		history = history[len(history)-maxHistory:]
	}
	for _, h := range history {
		// role はクライアント由来。system を混ぜられるとシステムプロンプトを
		// 乗っ取られる（プロンプトインジェクション）ため user / assistant のみ許可する。
		if h.Role != "user" && h.Role != "assistant" {
			continue
		}
		chatMessages = append(chatMessages, chatMessage{
			Role:    h.Role,
			Content: truncateRunes(h.Content, maxHistoryContentChars),
		})
	}
	chatMessages = append(chatMessages, chatMessage{Role: "user", Content: question})

	answer, err := s.openaiClient.ChatCompletionWithHistory(ctx, systemPrompt, chatMessages)
	if err != nil {
		return nil, fmt.Errorf("回答生成に失敗しました: %w", err)
	}

	// Save chat history
	chatSession := &model.ChatSession{UserID: userID}
	if err := s.chatRepo.CreateSession(ctx, chatSession); err == nil {
		userMsg := &model.ChatMessage{
			SessionID: chatSession.ID,
			Role:      "user",
			Content:   question,
		}
		s.chatRepo.AddMessage(ctx, userMsg)

		sourcesJSON, _ := json.Marshal(sources)
		assistantMsg := &model.ChatMessage{
			SessionID: chatSession.ID,
			Role:      "assistant",
			Content:   answer,
			Sources:   sourcesJSON,
		}
		s.chatRepo.AddMessage(ctx, assistantMsg)
	}

	return &dto.RAGResponse{
		Answer:  answer,
		Sources: sources,
	}, nil
}

func (s *RAGService) GetHistory(ctx context.Context, userID string) ([]dto.ChatHistoryItem, error) {
	sessions, err := s.chatRepo.GetSessionsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("チャット履歴の取得に失敗しました: %w", err)
	}

	var items []dto.ChatHistoryItem
	for _, sess := range sessions {
		msgs, err := s.chatRepo.GetMessages(ctx, sess.ID)
		if err != nil {
			continue
		}

		var messageDTOs []dto.ChatMessageDTO
		for _, m := range msgs {
			msgDTO := dto.ChatMessageDTO{
				Role:      m.Role,
				Content:   m.Content,
				CreatedAt: m.CreatedAt,
			}
			if m.Sources != nil {
				var sources []dto.RAGSource
				json.Unmarshal(m.Sources, &sources)
				msgDTO.Sources = sources
			}
			messageDTOs = append(messageDTOs, msgDTO)
		}

		items = append(items, dto.ChatHistoryItem{
			SessionID: sess.ID,
			Messages:  messageDTOs,
			CreatedAt: sess.CreatedAt,
		})
	}

	return items, nil
}

func truncateRunes(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max])
}

// ChunkText splits text into chunks of roughly maxChars characters
func ChunkText(text string, maxChars int) []string {
	if maxChars <= 0 {
		maxChars = 1000
	}

	// 改行コードを正規化する（CRLF / CR -> LF）。CRLF のままだと段落区切り "\n\n"
	// が見つからず全体が1チャンクになり、埋め込み API の入力上限を超えて失敗する。
	text = strings.ReplaceAll(text, "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	paragraphs := strings.Split(text, "\n\n")
	var chunks []string
	var current strings.Builder

	flush := func() {
		if current.Len() > 0 {
			chunks = append(chunks, current.String())
			current.Reset()
		}
	}

	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}

		// 段落単体が上限を超える場合は rune 単位でハードスプリットする。
		if len([]rune(para)) > maxChars {
			flush()
			chunks = append(chunks, hardSplit(para, maxChars)...)
			continue
		}

		if current.Len()+len(para) > maxChars && current.Len() > 0 {
			flush()
		}
		if current.Len() > 0 {
			current.WriteString("\n\n")
		}
		current.WriteString(para)
	}

	flush()
	return chunks
}

// hardSplit は文字列を rune 単位で maxChars ごとに分割する。
func hardSplit(s string, maxChars int) []string {
	runes := []rune(s)
	var out []string
	for i := 0; i < len(runes); i += maxChars {
		end := i + maxChars
		if end > len(runes) {
			end = len(runes)
		}
		out = append(out, string(runes[i:end]))
	}
	return out
}
