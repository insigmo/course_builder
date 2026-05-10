package prefix_test

import (
	"testing"

	"github.com/insigmo/course_builder/internal/prefix"
)

func TestExtractLeading(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"square bracket", "[01] Введение.mp4", []string{"[01]"}},
		{"round bracket", "(02) Урок.mp4", []string{"(02)"}},
		{"curly bracket", "{03} Тема.mp4", []string{"{03}"}},
		{"multiple tokens", "[sw] [01] Введение.mp4", []string{"[sw]", "[01]"}},
		{"no token", "Введение.mp4", nil},
		{"empty", "", nil},
		{"whitespace only", "  ", nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := prefix.ExtractLeading(tt.input)
			if len(got) != len(tt.want) {
				t.Fatalf("ExtractLeading(%q) = %v, want %v", tt.input, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("token[%d]: got %q, want %q", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestDetectRepeated(t *testing.T) {
	names := []string{
		"[sw] 01 Введение.mp4",
		"[sw] 02 Урок.mp4",
		"[sw] 03 Практика.mp4",
		"[other] Что-то.mp4",
	}
	got := prefix.DetectRepeated(names)
	if _, ok := got["[sw]"]; !ok {
		t.Error("expected [sw] to be detected as repeated")
	}
	if _, ok := got["[other]"]; ok {
		t.Error("[other] appears only once, should not be in repeated")
	}
}

func TestStrip(t *testing.T) {
	removable := map[string]struct{}{"[sw]": {}, "[01]": {}}
	tests := []struct {
		input string
		want  string
	}{
		{"[sw] Введение.mp4", "Введение.mp4"},
		{"[sw] [01] Урок.mp4", "Урок.mp4"},
		{"[keep] Урок.mp4", "[keep] Урок.mp4"},
		{"Без префикса.mp4", "Без префикса.mp4"},
		{"  [sw]  Пробелы.mp4", "Пробелы.mp4"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := prefix.Strip(tt.input, removable)
			if got != tt.want {
				t.Errorf("Strip(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestCleanTitle(t *testing.T) {
	removable := map[string]struct{}{"[sw]": {}}
	knownExts := map[string]struct{}{".mp4": {}, ".docx": {}, ".mkv": {}}
	tests := []struct {
		input string
		want  string
	}{
		{"[sw] Введение.mp4", "Введение"},
		{"Урок по Go.mp4", "Урок по Go"},
		{"1.mp4", "1 урок"},
		{"2.mp4", "2 урок"},
		{"42.mp4", "42 урок"},
		{"10 Основы.mp4", "10 Основы"},
		{"Теория.docx", "Теория"},
		{"video.mkv", "video"},
		{"readme.txt", "readme.txt"}, // .txt not in knownExts → kept
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := prefix.CleanTitle(tt.input, removable, knownExts)
			if got != tt.want {
				t.Errorf("CleanTitle(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
