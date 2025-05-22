package utils

import (
	"path/filepath"
	"strings"
)

func IgnoreFile(path string) bool {
	base := filepath.Base(path)

	// macos sandbox tmp files
	if strings.Contains(base, ".sb-") {
		return true
	}

	// ignore classic prefixes
	if strings.HasPrefix(base, "~$") || // word, excel
		strings.HasPrefix(base, ".#") || // emacs, LibreOffice
		strings.HasPrefix(base, "~") || // backup unix file
		strings.HasPrefix(base, "._") || // HFS+ metadata
		strings.HasSuffix(base, "~") || // save file
		strings.HasSuffix(base, ".swp") || // swap Vim
		strings.HasSuffix(base, ".tmp") || // tmp files
		strings.HasSuffix(base, ".lock") { // locks
		return true
	}

	return false
}
