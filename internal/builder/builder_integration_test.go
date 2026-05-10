package builder_test

import (
	"archive/zip"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/insigmo/course_builder/internal/builder"
	"github.com/insigmo/course_builder/internal/config"
)

func xmlEscape(s string) string {
	var b strings.Builder
	xml.EscapeText(&b, []byte(s))
	return b.String()
}

// makeDocxSimple creates a minimal valid .docx with given text.
func makeDocxSimple(t *testing.T, text string) []byte {
	t.Helper()
	escaped := xmlEscape(text)
	docXML := `<?xml version="1.0" encoding="UTF-8"?>` +
		`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:body><w:p><w:r><w:t>` + escaped + `</w:t></w:r></w:p></w:body></w:document>`

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte(docXML))
	zw.Close()
	return buf.Bytes()
}

// extractCourseData decodes the base64 JSON embedded in the output HTML.
func extractCourseData(t *testing.T, htmlContent string) string {
	t.Helper()
	re := regexp.MustCompile(`const b64Data = "([^"]+)"`)
	m := re.FindStringSubmatch(htmlContent)
	if m == nil {
		t.Fatal("could not find b64Data in output HTML")
	}
	decoded, err := base64.StdEncoding.DecodeString(m[1])
	if err != nil {
		t.Fatalf("base64 decode failed: %v", err)
	}
	return string(decoded)
}

func TestRun_TextOnlyCourse(t *testing.T) {
	dir := t.TempDir()

	if err := os.WriteFile(filepath.Join(dir, "урок.docx"), makeDocxSimple(t, "Привет, мир!"), 0644); err != nil {
		t.Fatal(err)
	}

	cfg := config.DefaultConfig()
	if err := builder.Run(dir, cfg); err != nil {
		t.Fatalf("builder.Run failed: %v", err)
	}

	outPath := filepath.Join(dir, filepath.Base(dir)+".html")
	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("output HTML not found: %v", err)
	}

	html := string(data)
	if !strings.Contains(html, "<html") {
		t.Error("output does not look like HTML")
	}

	courseJSON := extractCourseData(t, html)
	if !strings.Contains(courseJSON, "Привет, мир!") {
		t.Errorf("course data does not contain docx text, got: %s", courseJSON[:min(200, len(courseJSON))])
	}
}

func TestRun_EmptyDir(t *testing.T) {
	dir := t.TempDir()
	cfg := config.DefaultConfig()
	err := builder.Run(dir, cfg)
	if err == nil {
		t.Error("expected error for empty directory, got nil")
	}
}

func TestRun_IgnoresJunkFiles(t *testing.T) {
	dir := t.TempDir()

	os.WriteFile(filepath.Join(dir, "прочти перед изучением!.docx"), []byte("junk"), 0644)
	os.WriteFile(filepath.Join(dir, "урок.docx"), makeDocxSimple(t, "Реальный урок"), 0644)

	cfg := config.DefaultConfig()
	if err := builder.Run(dir, cfg); err != nil {
		t.Fatalf("builder.Run failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(dir, "прочти перед изучением!.docx")); !os.IsNotExist(err) {
		t.Error("junk file should have been deleted")
	}
}

func TestRun_SubfolderLessons(t *testing.T) {
	dir := t.TempDir()

	lesson1 := filepath.Join(dir, "1. Введение")
	lesson2 := filepath.Join(dir, "2. Основы")
	os.MkdirAll(lesson1, 0755)
	os.MkdirAll(lesson2, 0755)

	os.WriteFile(filepath.Join(lesson1, "intro.docx"), makeDocxSimple(t, "Введение в курс"), 0644)
	os.WriteFile(filepath.Join(lesson2, "basics.docx"), makeDocxSimple(t, "Основы Go"), 0644)

	cfg := config.DefaultConfig()
	if err := builder.Run(dir, cfg); err != nil {
		t.Fatalf("builder.Run failed: %v", err)
	}

	outPath := filepath.Join(dir, filepath.Base(dir)+".html")
	data, _ := os.ReadFile(outPath)
	courseJSON := extractCourseData(t, string(data))

	if !strings.Contains(courseJSON, "Введение в курс") {
		t.Error("missing content from lesson 1")
	}
	if !strings.Contains(courseJSON, "Основы Go") {
		t.Error("missing content from lesson 2")
	}
}

func TestRun_NumberedFilenameTitle(t *testing.T) {
	dir := t.TempDir()
	subdir := filepath.Join(dir, "Урок 1")
	os.MkdirAll(subdir, 0755)
	os.WriteFile(filepath.Join(subdir, "1.docx"), makeDocxSimple(t, "Контент урока"), 0644)

	cfg := config.DefaultConfig()
	if err := builder.Run(dir, cfg); err != nil {
		t.Fatalf("builder.Run failed: %v", err)
	}

	outPath := filepath.Join(dir, filepath.Base(dir)+".html")
	data, _ := os.ReadFile(outPath)
	courseJSON := extractCourseData(t, string(data))

	var raw interface{}
	if err := json.Unmarshal([]byte(courseJSON), &raw); err != nil {
		t.Fatalf("failed to parse course JSON: %v", err)
	}
	if !strings.Contains(courseJSON, "1 \\u0443\\u0440\\u043e\\u043a") && !strings.Contains(courseJSON, "1 урок") {
		t.Errorf("expected title '1 урок', course JSON: %s", courseJSON[:min(300, len(courseJSON))])
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
