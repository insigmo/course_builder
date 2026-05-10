package config

import "strings"

// Config holds runtime settings for the builder.
type Config struct {
	FilesToRemove map[string]struct{}
	VideoExts     map[string]struct{}
	IgnoreExts    map[string]struct{}
	IgnoreNames   map[string]struct{}
	KnownExts     map[string]struct{}
}

func setOf(vals ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(vals))
	for _, v := range vals {
		m[v] = struct{}{}
	}
	return m
}

// DefaultConfig returns sensible defaults matching the original Python script.
func DefaultConfig() *Config {
	return &Config{
		FilesToRemove: setOf(
			"[www.sw.band] прочти перед изучением!.docx",
			"прочти перед изучением!.docx",
			"Качай БЕСПЛАТНО курсы в нашем Telegram!.pdf",
			"Доступ к 200000 курсов на Sliwbl.Com.url",
			"Курсы SkillBox, Otus, GeekBrains, Вердикт, Маркарян и др..url",
		),
		VideoExts:   setOf(".mp4", ".webm", ".mkv", ".avi", ".mov", ".ts"),
		IgnoreExts:  setOf(".url", ".lnk"),
		IgnoreNames: setOf("thumbs.db", ".ds_store", ".desktop"),
		KnownExts: setOf(
			".docx", ".doc",
			".mp4", ".webm", ".mkv", ".avi", ".mov", ".ts",
			".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp",
			".url", ".lnk", ".desktop",
			".html", ".htm", ".txt", ".pdf",
		),
	}
}

func (c *Config) IsVideo(ext string) bool {
	_, ok := c.VideoExts[strings.ToLower(ext)]
	return ok
}

// ShouldIgnore returns true if the file should be completely skipped.
func (c *Config) ShouldIgnore(nameLow, ext string) bool {
	if _, ok := c.IgnoreExts[ext]; ok {
		return true
	}
	_, ok := c.IgnoreNames[nameLow]
	return ok
}

// IsContent returns true if this file carries course content (.docx or video).
func (c *Config) IsContent(ext string) bool {
	if _, ok := c.IgnoreExts[ext]; ok {
		return false
	}
	if ext == ".docx" || ext == ".doc" {
		return true
	}
	_, ok := c.VideoExts[ext]
	return ok
}
