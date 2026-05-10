package video

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// Available returns true if ffmpeg is present in PATH.
func Available() bool {
	_, err := exec.LookPath("ffmpeg")
	return err == nil
}

// ConvertToMP4 converts any video to .mp4 via ffmpeg.
func ConvertToMP4(inputPath string) (string, error) {
	if !Available() {
		return "", fmt.Errorf("ffmpeg not found in PATH")
	}
	dir := filepath.Dir(inputPath)
	stem := strings.TrimSuffix(filepath.Base(inputPath), filepath.Ext(inputPath))
	out := filepath.Join(dir, stem+".mp4")
	cmd := exec.Command("ffmpeg",
		"-i", inputPath,
		"-c:v", "libx264",
		"-c:a", "aac",
		"-movflags", "+faststart",
		"-y", out,
	)
	if b, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("ffmpeg: %w\n%s", err, string(b))
	}
	return out, nil
}

// Tag builds an HTML <video> tag.
func Tag(relPath, ext string) string {
	return fmt.Sprintf(
		`<div class="video-wrapper"><video controls preload="metadata"><source src="%s" type="%s">Ваш браузер не поддерживает видео.</video></div>`,
		relPath, mime(ext),
	)
}

func mime(ext string) string {
	switch strings.ToLower(ext) {
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mkv":
		return "video/x-matroska"
	case ".mov":
		return "video/quicktime"
	case ".avi":
		return "video/x-msvideo"
	case ".ts":
		return "video/mp2t"
	default:
		return "video/mp4"
	}
}
