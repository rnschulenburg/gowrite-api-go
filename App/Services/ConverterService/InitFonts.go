package ConverterService

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
)
import "embed"

//go:embed fonts/*
var fontFS embed.FS
var FontDir = filepath.Join(os.TempDir(), "gowrite-fonts")

func InitFonts() error {

	log.Println("InitFonts start")

	// immer neu erstellen (fonts sind klein)
	err := os.RemoveAll(FontDir)
	if err != nil {
		return err
	}

	err = os.MkdirAll(FontDir, 0755)
	if err != nil {
		return err
	}

	return fs.WalkDir(fontFS, "fonts", func(path string, d fs.DirEntry, err error) error {

		if d.IsDir() {
			return nil
		}

		data, err := fontFS.ReadFile(path)
		if err != nil {
			return err
		}

		target := filepath.Join(FontDir, filepath.Base(path))

		//log.Println("extract font:", target)

		return os.WriteFile(target, data, 0644)
	})
}
