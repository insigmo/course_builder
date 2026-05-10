package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/user/course-builder/internal/builder"
	"github.com/user/course-builder/internal/config"
)

var Version = "dev"

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: course-builder <course-directory>")
		fmt.Fprintf(os.Stderr, "Version: %s\n", Version)
		os.Exit(1)
	}

	rootDir := filepath.Clean(args[0])
	info, err := os.Stat(rootDir)
	if err != nil || !info.IsDir() {
		fmt.Fprintf(os.Stderr, "\u274c Папка не найдена: %s\n", rootDir)
		os.Exit(1)
	}

	cfg := config.DefaultConfig()
	if err := builder.Run(rootDir, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "\u274c %v\n", err)
		os.Exit(1)
	}
}
