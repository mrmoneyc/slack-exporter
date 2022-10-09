package archiver

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mrmoneyc/slack-exporter/pkg/config"
	log "github.com/sirupsen/logrus"
)

func MakeArchive(cfg *config.Config, now string) error {
	currDir, err := os.Getwd()
	if err != nil {
		return err
	}

	exportPath := filepath.Join(currDir, cfg.ExportBasePath, now)
	zipPath := filepath.Join(currDir, cfg.ExportBasePath, fmt.Sprintf("%s.zip", now))

	if err = os.Chdir(exportPath); err != nil {
		return err
	}

	currDir, err = os.Getwd()
	if err != nil {
		return err
	}

	relPath, _ := filepath.Rel(currDir, exportPath)

	log.Infof("create archive file: %s", zipPath)
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	w := zip.NewWriter(zipFile)
	defer w.Close()

	walker := func(path string, fi os.FileInfo, err error) error {
		log.Debugf("adding: %s", path)

		if err != nil {
			return err
		}

		if fi.IsDir() {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		h := &zip.FileHeader{
			Name:   path,
			Method: zip.Deflate,
			Flags:  0x800,
		}
		archiveFile, err := w.CreateHeader(h)
		if err != nil {
			return err
		}

		_, err = io.Copy(archiveFile, f)
		if err != nil {
			return err
		}

		return nil
	}

	err = filepath.Walk(relPath, walker)
	if err != nil {
		return err
	}

	if err := os.RemoveAll(exportPath); err != nil {
		log.WithError(err).Warnf("cannot remove data export working directory")
	}

	return nil
}
