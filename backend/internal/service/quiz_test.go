package service

import (
	"testing"
	"time"

	"github.com/rootage-ses-quiz/backend/internal/repository"
)

func TestCalcStreak_Empty(t *testing.T) {
	result := calcStreak(nil)
	if result != 0 {
		t.Errorf("calcStreak(nil) = %d, want 0", result)
	}
}

func TestCalcStreak_Today(t *testing.T) {
	today := time.Now().Format("2006-01-02")
	stats := []repository.DailyStats{{Date: today, Total: 10, Correct: 8}}

	result := calcStreak(stats)
	if result != 1 {
		t.Errorf("calcStreak with today = %d, want 1", result)
	}
}

func TestCalcStreak_Yesterday(t *testing.T) {
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	stats := []repository.DailyStats{{Date: yesterday, Total: 10, Correct: 8}}

	result := calcStreak(stats)
	if result != 1 {
		t.Errorf("calcStreak with yesterday = %d, want 1", result)
	}
}

func TestCalcStreak_Consecutive(t *testing.T) {
	d3 := time.Now().AddDate(0, 0, -2).Format("2006-01-02")
	d2 := time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	d1 := time.Now().Format("2006-01-02")

	stats := []repository.DailyStats{
		{Date: d3, Total: 5, Correct: 3},
		{Date: d2, Total: 10, Correct: 8},
		{Date: d1, Total: 10, Correct: 9},
	}

	result := calcStreak(stats)
	if result != 3 {
		t.Errorf("calcStreak consecutive 3 days = %d, want 3", result)
	}
}

func TestCalcStreak_BrokenStreak(t *testing.T) {
	// 3日前と今日だけ（2日前が欠けている）
	d3 := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
	d1 := time.Now().Format("2006-01-02")

	stats := []repository.DailyStats{
		{Date: d3, Total: 5, Correct: 3},
		{Date: d1, Total: 10, Correct: 9},
	}

	result := calcStreak(stats)
	if result != 1 {
		t.Errorf("calcStreak with gap = %d, want 1", result)
	}
}

func TestCalcStreak_OldData(t *testing.T) {
	// 1週間前のデータのみ（昨日でも今日でもない）
	old := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	stats := []repository.DailyStats{{Date: old, Total: 10, Correct: 5}}

	result := calcStreak(stats)
	if result != 0 {
		t.Errorf("calcStreak with old data = %d, want 0", result)
	}
}

func TestCalcStreak_LongStreak(t *testing.T) {
	var stats []repository.DailyStats
	for i := 9; i >= 0; i-- {
		d := time.Now().AddDate(0, 0, -i).Format("2006-01-02")
		stats = append(stats, repository.DailyStats{Date: d, Total: 10, Correct: 7})
	}

	result := calcStreak(stats)
	if result != 10 {
		t.Errorf("calcStreak 10 days = %d, want 10", result)
	}
}
