package htmlrender

import (
	_ "embed"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"html"
	"strings"
)

//go:embed template.html
var tmpl string

// Render injects course data into the embedded HTML template.
func Render(title string, lessons interface{}, totalSteps int) (string, error) {
	data, err := json.Marshal(lessons)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}
	b64 := base64.StdEncoding.EncodeToString(data)

	out := tmpl
	out = strings.ReplaceAll(out, "__BASE64_DATA__", b64)
	out = strings.ReplaceAll(out, "__TOTAL_STEPS__", fmt.Sprintf("%d", totalSteps))
	out = strings.ReplaceAll(out, "__COURSE_TITLE__", html.EscapeString(title))
	out = strings.ReplaceAll(out, "__COURSE_TITLE_JS__", jsEscape(title))
	return out, nil
}

func jsEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	return s
}
