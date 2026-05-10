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
var templateHTML string

//go:embed style.css
var styleCSS string

//go:embed app.js
var appJS string

//go:embed player.js
var playerJS string

//go:embed settings.js
var settingsJS string

// Render injects course data into the HTML template and returns the final HTML.
func Render(title string, lessons interface{}, totalSteps int) (string, error) {
	jsonBytes, err := json.Marshal(lessons)
	if err != nil {
		return "", fmt.Errorf("marshal lessons: %w", err)
	}
	b64 := base64.StdEncoding.EncodeToString(jsonBytes)

	out := templateHTML
	out = strings.ReplaceAll(out, "STYLE_PLACEHOLDER", styleCSS)
	out = strings.ReplaceAll(out, "APP_JS_PLACEHOLDER", appJS)
	out = strings.ReplaceAll(out, "PLAYER_JS_PLACEHOLDER", playerJS)
	out = strings.ReplaceAll(out, "SETTINGS_JS_PLACEHOLDER", settingsJS)
	out = strings.ReplaceAll(out, "__BASE64_DATA__", b64)
	out = strings.ReplaceAll(out, "__TOTAL_STEPS__", fmt.Sprintf("%d", totalSteps))
	out = strings.ReplaceAll(out, "__COURSE_TITLE__", html.EscapeString(title))
	out = strings.ReplaceAll(out, "__COURSE_TITLE_JS__", jsEscape(title))
	return out, nil
}

func jsEscape(s string) string {
	s = strings.ReplaceAll(s, `\`, `\\`)
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}
