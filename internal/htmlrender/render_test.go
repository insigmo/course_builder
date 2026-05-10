package htmlrender_test

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"

	"github.com/insigmo/course_builder/internal/htmlrender"
)

type fakeStep struct {
	Title string `json:"title"`
	HTML  string `json:"html"`
}

type fakeLesson struct {
	Title string     `json:"title"`
	Steps []fakeStep `json:"steps"`
}

func TestRender_ReturnsHTML(t *testing.T) {
	lessons := []fakeLesson{{Title: "Введение", Steps: []fakeStep{{Title: "Шаг 1", HTML: "<p>Текст</p>"}}}}
	out, err := htmlrender.Render("Мой курс", lessons, 1)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(out, "<html") {
		t.Error("output does not contain <html>")
	}
	if !strings.Contains(out, "</html>") {
		t.Error("output does not contain </html>")
	}
}

func TestRender_CourseTitle(t *testing.T) {
	lessons := []fakeLesson{}
	out, err := htmlrender.Render("Мой Тестовый Курс", lessons, 0)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(out, "Мой Тестовый Курс") {
		t.Error("course title not found in output HTML")
	}
}

func TestRender_CourseTitleEscapesHTML(t *testing.T) {
	out, err := htmlrender.Render(`<script>alert("xss")</script>`, []fakeLesson{}, 0)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if strings.Contains(out, "<script>alert") {
		t.Error("HTML in course title should be escaped")
	}
}

func TestRender_Base64DataEmbedded(t *testing.T) {
	lessons := []fakeLesson{{Title: "Урок 1", Steps: []fakeStep{{Title: "Шаг", HTML: "<p>hello</p>"}}}}
	out, err := htmlrender.Render("Курс", lessons, 1)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	// Extract b64 blob and decode
	const marker = `const b64Data = "`
	idx := strings.Index(out, marker)
	if idx == -1 {
		t.Fatal("b64Data marker not found in output")
	}
	start := idx + len(marker)
	end := strings.Index(out[start:], `"`)
	if end == -1 {
		t.Fatal("closing quote for b64Data not found")
	}
	blob := out[start : start+end]
	decoded, err := base64.StdEncoding.DecodeString(blob)
	if err != nil {
		t.Fatalf("base64 decode failed: %v", err)
	}
	if !json.Valid(decoded) {
		t.Error("decoded b64Data is not valid JSON")
	}
	if !strings.Contains(string(decoded), "Урок 1") {
		t.Errorf("decoded JSON does not contain lesson title, got: %s", string(decoded))
	}
}

func TestRender_TotalStepsEmbedded(t *testing.T) {
	out, err := htmlrender.Render("Курс", []fakeLesson{}, 42)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(out, "42") {
		t.Error("total steps count (42) not found in output")
	}
}

func TestRender_StyleAndScriptsInlined(t *testing.T) {
	out, err := htmlrender.Render("Курс", []fakeLesson{}, 0)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if !strings.Contains(out, "<style>") {
		t.Error("inlined <style> not found")
	}
	if !strings.Contains(out, "<script>") {
		t.Error("inlined <script> not found")
	}
}

func TestRender_NoPlaceholdersLeft(t *testing.T) {
	out, err := htmlrender.Render("Курс", []fakeLesson{}, 0)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	for _, placeholder := range []string{
		"STYLE_PLACEHOLDER",
		"APP_JS_PLACEHOLDER",
		"PLAYER_JS_PLACEHOLDER",
		"SETTINGS_JS_PLACEHOLDER",
		"__BASE64_DATA__",
		"__TOTAL_STEPS__",
		"__COURSE_TITLE__",
	} {
		if strings.Contains(out, placeholder) {
			t.Errorf("placeholder %q was not replaced in output", placeholder)
		}
	}
}

func TestRender_EmptyLessons(t *testing.T) {
	out, err := htmlrender.Render("Пустой курс", []fakeLesson{}, 0)
	if err != nil {
		t.Fatalf("Render() error: %v", err)
	}
	if out == "" {
		t.Error("Render() returned empty string for empty lessons")
	}
}
