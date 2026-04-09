package service

import (
	"strings"
	"testing"
)

func TestChunkText_Basic(t *testing.T) {
	text := "paragraph1\n\nparagraph2\n\nparagraph3"
	chunks := ChunkText(text, 1000)

	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk for short text, got %d", len(chunks))
	}
	if !strings.Contains(chunks[0], "paragraph1") {
		t.Error("chunk should contain paragraph1")
	}
}

func TestChunkText_Split(t *testing.T) {
	text := "aaaa\n\nbbbb\n\ncccc"
	chunks := ChunkText(text, 8)

	if len(chunks) < 2 {
		t.Errorf("expected at least 2 chunks, got %d", len(chunks))
	}
}

func TestChunkText_EmptyText(t *testing.T) {
	chunks := ChunkText("", 1000)
	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks for empty text, got %d", len(chunks))
	}
}

func TestChunkText_OnlyWhitespace(t *testing.T) {
	chunks := ChunkText("\n\n\n\n", 1000)
	if len(chunks) != 0 {
		t.Errorf("expected 0 chunks for whitespace-only text, got %d", len(chunks))
	}
}

func TestChunkText_ZeroMaxChars(t *testing.T) {
	// maxChars <= 0 はデフォルトの1000にフォールバック
	chunks := ChunkText("short text", 0)
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk with default maxChars, got %d", len(chunks))
	}
}

func TestChunkText_SingleParagraph(t *testing.T) {
	text := strings.Repeat("a", 500)
	chunks := ChunkText(text, 1000)
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk, got %d", len(chunks))
	}
	if chunks[0] != text {
		t.Error("chunk should contain the full text")
	}
}

func TestChunkText_LargeParagraph(t *testing.T) {
	// maxCharsを超える単一段落はそのまま1チャンクになる
	text := strings.Repeat("a", 2000)
	chunks := ChunkText(text, 1000)
	if len(chunks) != 1 {
		t.Errorf("expected 1 chunk for unsplittable text, got %d", len(chunks))
	}
}

func TestChunkText_PreservesContent(t *testing.T) {
	text := "first paragraph\n\nsecond paragraph\n\nthird paragraph"
	chunks := ChunkText(text, 1000)

	joined := strings.Join(chunks, "\n\n")
	if !strings.Contains(joined, "first paragraph") {
		t.Error("chunks should contain 'first paragraph'")
	}
	if !strings.Contains(joined, "third paragraph") {
		t.Error("chunks should contain 'third paragraph'")
	}
}
