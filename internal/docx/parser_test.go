package docx_test

import (
	"archive/zip"
	"bytes"
	"encoding/xml"
	"os"
	"strings"
	"testing"

	"github.com/insigmo/course_builder/internal/docx"
)

func xmlEscape(s string) string {
	var b strings.Builder
	xml.EscapeText(&b, []byte(s))
	return b.String()
}

func makeDocx(t *testing.T, bodyXML string) string {
	t.Helper()
	docXML := `<?xml version="1.0" encoding="UTF-8"?>` +
		`<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">` +
		`<w:body>` + bodyXML + `</w:body></w:document>`

	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	f, err := zw.Create("word/document.xml")
	if err != nil {
		t.Fatal(err)
	}
	f.Write([]byte(docXML))
	zw.Close()

	// Write to temp file
	tmp, err := os.CreateTemp("", "test-*.docx")
	if err != nil {
		t.Fatal(err)
	}
	tmp.Write(buf.Bytes())
	tmp.Close()
	return tmp.Name()
}

func para(text string) string {
	return `<w:p><w:r><w:t>` + xmlEscape(text) + `</w:t></w:r></w:p>`
}

func boldPara(text string) string {
	return `<w:p><w:r><w:rPr><w:b/></w:rPr><w:t>` + xmlEscape(text) + `</w:t></w:r></w:p>`
}

func heading1Para(text string) string {
	return `<w:p><w:pPr><w:pStyle w:val="heading1"/></w:pPr><w:r><w:t>` + xmlEscape(text) + `</w:t></w:r></w:p>`
}

func quizParas() string {
	return para("??Что такое горутина?") +
		para("+Лёгковесный поток") +
		para("-Процесс ОС") +
		para("-Канал")
}

func TestParse_SimpleParagraph(t *testing.T) {
	f := makeDocx(t, para("Привет мир"))
	defer os.Remove(f)

	res, err := docx.Parse(f)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !strings.Contains(res.HTML, "Привет мир") {
		t.Errorf("expected paragraph text in HTML, got: %s", res.HTML)
	}
	if !strings.Contains(res.HTML, "<p>") {
		t.Errorf("expected <p> tag in HTML, got: %s", res.HTML)
	}
}

func TestParse_BoldText(t *testing.T) {
	f := makeDocx(t, boldPara("Жирный текст"))
	defer os.Remove(f)

	res, err := docx.Parse(f)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !strings.Contains(res.HTML, "<strong>") {
		t.Errorf("expected <strong> tag, got: %s", res.HTML)
	}
}

func TestParse_Heading(t *testing.T) {
	f := makeDocx(t, heading1Para("Заголовок первого уровня"))
	defer os.Remove(f)

	res, err := docx.Parse(f)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if !strings.Contains(res.HTML, "<h2>") {
		t.Errorf("expected <h2> tag, got: %s", res.HTML)
	}
}

func TestParse_Quiz(t *testing.T) {
	f := makeDocx(t, quizParas())
	defer os.Remove(f)

	res, err := docx.Parse(f)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if len(res.Quizzes) != 1 {
		t.Fatalf("expected 1 quiz, got %d", len(res.Quizzes))
	}
	q := res.Quizzes[0]
	if q.Question != "Что такое горутина?" {
		t.Errorf("wrong question: %q", q.Question)
	}
	if len(q.Options) != 3 {
		t.Errorf("expected 3 options, got %d", len(q.Options))
	}
	if q.Answer != 0 {
		t.Errorf("expected answer index 0, got %d", q.Answer)
	}
}

func TestParse_InvalidFile(t *testing.T) {
	_, err := docx.Parse("/nonexistent/file.docx")
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestParse_EmptyBody(t *testing.T) {
	f := makeDocx(t, "")
	defer os.Remove(f)

	res, err := docx.Parse(f)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}
	if res.HTML != "" {
		t.Errorf("expected empty HTML for empty body, got: %q", res.HTML)
	}
}
