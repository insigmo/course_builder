package config_test

import (
	"testing"

	"github.com/insigmo/course_builder/internal/config"
)

func TestDefaultConfig_IsVideo(t *testing.T) {
	cfg := config.DefaultConfig()
	videoCases := []struct {
		ext  string
		want bool
	}{
		{".mp4", true},
		{".MP4", true}, // case-insensitive
		{".webm", true},
		{".mkv", true},
		{".avi", true},
		{".mov", true},
		{".ts", true},
		{".docx", false},
		{".pdf", false},
		{".txt", false},
		{"", false},
	}
	for _, tc := range videoCases {
		t.Run(tc.ext, func(t *testing.T) {
			if got := cfg.IsVideo(tc.ext); got != tc.want {
				t.Errorf("IsVideo(%q) = %v, want %v", tc.ext, got, tc.want)
			}
		})
	}
}

func TestDefaultConfig_ShouldIgnore(t *testing.T) {
	cfg := config.DefaultConfig()
	cases := []struct {
		name string
		ext  string
		want bool
	}{
		{"file.url", ".url", true},
		{"file.lnk", ".lnk", true},
		{"thumbs.db", ".db", true},
		{".ds_store", "", true},
		{"video.mp4", ".mp4", false},
		{"lesson.docx", ".docx", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := cfg.ShouldIgnore(tc.name, tc.ext); got != tc.want {
				t.Errorf("ShouldIgnore(%q, %q) = %v, want %v", tc.name, tc.ext, got, tc.want)
			}
		})
	}
}

func TestDefaultConfig_IsContent(t *testing.T) {
	cfg := config.DefaultConfig()
	cases := []struct {
		ext  string
		want bool
	}{
		{".mp4", true},
		{".webm", true},
		{".docx", true},
		{".doc", true},
		{".url", false},
		{".lnk", false},
		{".pdf", false},
		{".txt", false},
		{"", false},
	}
	for _, tc := range cases {
		t.Run(tc.ext, func(t *testing.T) {
			if got := cfg.IsContent(tc.ext); got != tc.want {
				t.Errorf("IsContent(%q) = %v, want %v", tc.ext, got, tc.want)
			}
		})
	}
}
