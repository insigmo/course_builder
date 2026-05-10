package video_test

import (
	"strings"
	"testing"

	"github.com/insigmo/course_builder/internal/video"
)

func TestTag_ContainsSource(t *testing.T) {
	tag := video.Tag("lessons/1.mp4", ".mp4")
	if !strings.Contains(tag, `src="lessons/1.mp4"`) {
		t.Errorf("Tag() missing src, got: %s", tag)
	}
	if !strings.Contains(tag, `type="video/mp4"`) {
		t.Errorf("Tag() missing type, got: %s", tag)
	}
	if !strings.Contains(tag, `class="video-wrapper"`) {
		t.Errorf("Tag() missing video-wrapper class, got: %s", tag)
	}
}

func TestTag_MimeTypes(t *testing.T) {
	cases := []struct {
		ext      string
		wantMime string
	}{
		{".mp4", "video/mp4"},
		{".webm", "video/webm"},
		{".mkv", "video/x-matroska"},
		{".mov", "video/quicktime"},
		{".avi", "video/x-msvideo"},
		{".ts", "video/mp2t"},
		{".unknown", "video/mp4"}, // fallback
	}
	for _, tc := range cases {
		t.Run(tc.ext, func(t *testing.T) {
			tag := video.Tag("file"+tc.ext, tc.ext)
			if !strings.Contains(tag, `type="`+tc.wantMime+`"`) {
				t.Errorf("Tag(%q) mime = want %q, got tag: %s", tc.ext, tc.wantMime, tag)
			}
		})
	}
}

func TestAvailable_NoPanic(t *testing.T) {
	// Just ensure it doesn't panic — result depends on environment
	_ = video.Available()
}
