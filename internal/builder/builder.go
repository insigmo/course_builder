package builder

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/insigmo/course_builder/internal/config"
	"github.com/insigmo/course_builder/internal/docx"
	"github.com/insigmo/course_builder/internal/htmlrender"
	"github.com/insigmo/course_builder/internal/prefix"
	"github.com/insigmo/course_builder/internal/video"
)

// Step is one course step — contains a title, optional HTML body, and optional quizzes.
type Step struct {
	Title   string      `json:"title"`
	HTML    string      `json:"html"`
	Quizzes []docx.Quiz `json:"quizzes"`
}

// LessonNode is a node in the recursive lesson tree.
type LessonNode struct {
	Lesson   string       `json:"lesson"`
	Steps    []Step       `json:"steps"`
	Children []LessonNode `json:"children"`
}

// Stats accumulates build counters printed at the end.
type Stats struct {
	Lessons, Steps, TextSteps, VideoSteps, Videos, Quizzes int
}

var safeRE = regexp.MustCompile(`[<>:"/\\|?*]+`)

// Run is the main entry point: scans rootDir and writes a self-contained HTML course.
func Run(rootDir string, cfg *config.Config) error {
	exePath, _ := os.Executable()
	outputPath := filepath.Join(rootDir, safeFilename(filepath.Base(rootDir))+".html")

	allNames := gatherAllNames(rootDir)
	removable := prefix.DetectRepeated(allNames)
	if len(removable) > 0 {
		keys := make([]string, 0, len(removable))
		for k := range removable {
			keys = append(keys, k)
		}
		fmt.Printf("🔍 Повторяющиеся префиксы (удалены из заголовков): %s\n", strings.Join(keys, ", "))
	}

	deleteJunk(rootDir, cfg)

	stats := &Stats{}
	lessons, err := buildCourse(rootDir, outputPath, exePath, cfg, removable, stats)
	if err != nil {
		return err
	}
	if len(lessons) == 0 {
		return fmt.Errorf("не найдено валидных .docx или видеофайлов")
	}

	courseTitle := prefix.CleanTitle(filepath.Base(rootDir), removable, cfg.KnownExts)
	if courseTitle == "" {
		courseTitle = "Курс"
	}

	htmlOut, err := htmlrender.Render(courseTitle, lessons, stats.Steps)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(outputPath), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(outputPath, []byte(htmlOut), 0o644); err != nil {
		return err
	}

	fi, _ := os.Stat(outputPath)
	fmt.Println("✅ Курс успешно собран!")
	fmt.Printf("  📁 Глав:   %d\n", stats.Lessons)
	fmt.Printf("  📄 Шагов:  %d (текстовых: %d, видео: %d)\n", stats.Steps, stats.TextSteps, stats.VideoSteps)
	fmt.Printf("  🧩 Тестов: %d\n", stats.Quizzes)
	if fi != nil {
		fmt.Printf("  📦 Размер: %s\n", prettySize(fi.Size()))
	}
	fmt.Printf("  💾 Файл:   %s\n", outputPath)
	return nil
}

func buildCourse(rootDir, outputHTML, exePath string, cfg *config.Config, removable map[string]struct{}, stats *Stats) ([]LessonNode, error) {
	var lessons []LessonNode

	if rootFiles := validFiles(rootDir, exePath, outputHTML, cfg, removable); len(rootFiles) > 0 {
		intro, err := stepsFromFiles(rootFiles, outputHTML, cfg, removable, stats)
		if err != nil {
			return nil, err
		}
		if len(intro) > 0 {
			lessons = append(lessons, LessonNode{Lesson: "Введение", Steps: intro})
		}
	}

	for _, dir := range sortedDirs(rootDir, removable) {
		node, err := lessonNode(dir, outputHTML, exePath, cfg, removable, stats)
		if err != nil {
			return nil, err
		}
		if node != nil {
			lessons = append(lessons, *node)
		}
	}

	stats.Lessons = len(lessons)
	return lessons, nil
}

func lessonNode(folder, outputHTML, exePath string, cfg *config.Config, removable map[string]struct{}, stats *Stats) (*LessonNode, error) {
	name := prefix.CleanTitle(filepath.Base(folder), removable, cfg.KnownExts)

	var steps []Step
	if files := validFiles(folder, exePath, outputHTML, cfg, removable); len(files) > 0 {
		var err error
		steps, err = stepsFromFiles(files, outputHTML, cfg, removable, stats)
		if err != nil {
			return nil, err
		}
	}

	var children []LessonNode
	for _, sub := range sortedDirs(folder, removable) {
		child, err := lessonNode(sub, outputHTML, exePath, cfg, removable, stats)
		if err != nil {
			return nil, err
		}
		if child != nil {
			children = append(children, *child)
		}
	}

	if len(steps) == 0 && len(children) == 0 {
		return nil, nil
	}
	return &LessonNode{Lesson: name, Steps: steps, Children: children}, nil
}

type fileGroup struct {
	docxFiles  []string
	videoFiles []string
}

func stepsFromFiles(files []string, outputHTML string, cfg *config.Config, removable map[string]struct{}, stats *Stats) ([]Step, error) {
	var groupOrder []string
	groups := map[string]*fileGroup{}

	for _, f := range files {
		name := filepath.Base(f)
		ext := strings.ToLower(filepath.Ext(name))
		key := groupKey(name, removable)
		if _, ok := groups[key]; !ok {
			groups[key] = &fileGroup{}
			groupOrder = append(groupOrder, key)
		}
		if ext == ".docx" || ext == ".doc" {
			groups[key].docxFiles = append(groups[key].docxFiles, f)
		} else {
			groups[key].videoFiles = append(groups[key].videoFiles, f)
		}
	}

	dir := ""
	if len(files) > 0 {
		dir = filepath.Dir(files[0])
	}

	var steps []Step
	for _, key := range groupOrder {
		g := groups[key]

		dn := basenames(g.docxFiles)
		SortNames(dn, removable)
		g.docxFiles = rebuildPaths(dir, dn)

		vn := basenames(g.videoFiles)
		SortNames(vn, removable)
		g.videoFiles = rebuildPaths(dir, vn)

		var htmlParts []string
		var quizzes []docx.Quiz
		for _, d := range g.docxFiles {
			res, err := docx.Parse(d)
			if err != nil {
				fmt.Printf("⚠️  Пропущен DOCX %s: %v\n", filepath.Base(d), err)
				continue
			}
			if res.HTML != "" {
				htmlParts = append(htmlParts, res.HTML)
			}
			quizzes = append(quizzes, res.Quizzes...)
		}

		if len(g.videoFiles) == 0 {
			if len(htmlParts) == 0 && len(quizzes) == 0 {
				continue
			}
			title := ""
			if len(g.docxFiles) > 0 {
				title = prefix.CleanTitle(filepath.Base(g.docxFiles[0]), removable, cfg.KnownExts)
			}
			steps = append(steps, Step{
				Title:   title,
				HTML:    strings.Join(htmlParts, "\n"),
				Quizzes: quizzes,
			})
			stats.Steps++
			stats.TextSteps++
			stats.Quizzes += len(quizzes)
			continue
		}

		for i, vPath := range g.videoFiles {
			actualPath := vPath
			if ext := strings.ToLower(filepath.Ext(vPath)); ext != ".mp4" {
				converted, err := video.ConvertToMP4(vPath)
				if err != nil {
					fmt.Printf("⚠️  Конвертация %s: %v\n", filepath.Base(vPath), err)
				} else {
					fmt.Printf("  🔄 %s → %s\n", filepath.Base(vPath), filepath.Base(converted))
					actualPath = converted
				}
			}
			relPath, _ := filepath.Rel(filepath.Dir(outputHTML), actualPath)
			vHTML := video.Tag(relPath, strings.ToLower(filepath.Ext(actualPath)))

			var stepHTML []string
			stepHTML = append(stepHTML, vHTML)
			var stepQuizzes []docx.Quiz
			if i == 0 {
				stepHTML = append(stepHTML, htmlParts...)
				stepQuizzes = quizzes
			}

			steps = append(steps, Step{
				Title:   prefix.CleanTitle(filepath.Base(vPath), removable, cfg.KnownExts),
				HTML:    strings.Join(stepHTML, "\n"),
				Quizzes: stepQuizzes,
			})
			stats.Steps++
			stats.VideoSteps++
			stats.Videos++
		}
		if len(g.docxFiles) > 0 {
			stats.TextSteps++
		}
		stats.Quizzes += len(quizzes)
	}
	return steps, nil
}

func validFiles(folder, exePath, outputPath string, cfg *config.Config, removable map[string]struct{}) []string {
	entries, _ := os.ReadDir(folder)
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		ext := strings.ToLower(filepath.Ext(name))
		if cfg.ShouldIgnore(strings.ToLower(name), ext) || !cfg.IsContent(ext) {
			continue
		}
		full := filepath.Join(folder, name)
		if full == exePath || full == outputPath {
			continue
		}
		names = append(names, name)
	}
	SortNames(names, removable)
	result := make([]string, len(names))
	for i, n := range names {
		result[i] = filepath.Join(folder, n)
	}
	return result
}

func sortedDirs(dir string, removable map[string]struct{}) []string {
	entries, _ := os.ReadDir(dir)
	var names []string
	for _, e := range entries {
		if e.IsDir() && !strings.HasPrefix(e.Name(), ".") {
			names = append(names, e.Name())
		}
	}
	SortNames(names, removable)
	result := make([]string, len(names))
	for i, n := range names {
		result[i] = filepath.Join(dir, n)
	}
	return result
}

func gatherAllNames(rootDir string) []string {
	var names []string
	_ = filepath.WalkDir(rootDir, func(_ string, d fs.DirEntry, _ error) error {
		if d != nil && !strings.HasPrefix(d.Name(), ".") {
			names = append(names, d.Name())
		}
		return nil
	})
	return names
}

func groupKey(name string, removable map[string]struct{}) string {
	s := prefix.Strip(name, removable)
	if dot := strings.LastIndex(s, "."); dot >= 0 {
		s = s[:dot]
	}
	return strings.ToLower(strings.TrimSpace(s))
}

func deleteJunk(rootDir string, cfg *config.Config) {
	_ = filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, _ error) error {
		if d == nil || d.IsDir() {
			return nil
		}
		if _, ok := cfg.FilesToRemove[strings.ToLower(d.Name())]; ok {
			if err := os.Remove(path); err == nil {
				fmt.Printf("🗑️  Удалён: %s\n", d.Name())
			}
		}
		return nil
	})
}

func basenames(paths []string) []string {
	names := make([]string, len(paths))
	for i, p := range paths {
		names[i] = filepath.Base(p)
	}
	return names
}

func rebuildPaths(dir string, names []string) []string {
	result := make([]string, len(names))
	for i, n := range names {
		result[i] = filepath.Join(dir, n)
	}
	return result
}

func safeFilename(name string) string {
	name = strings.TrimSpace(name)
	if name == "" {
		return "course"
	}
	return safeRE.ReplaceAllString(name, "_")
}

func prettySize(b int64) string {
	units := []string{"B", "KB", "MB", "GB"}
	size := float64(b)
	idx := 0
	for size >= 1024 && idx < len(units)-1 {
		size /= 1024
		idx++
	}
	if idx == 0 {
		return fmt.Sprintf("%d B", b)
	}
	return fmt.Sprintf("%.1f %s", size, units[idx])
}
