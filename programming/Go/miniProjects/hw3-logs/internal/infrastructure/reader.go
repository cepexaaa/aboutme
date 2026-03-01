package infrastructure

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FileReader struct{}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (r *FileReader) ReadFiles(paths []string) ([]io.Reader, []string, error) {
	var readers []io.Reader
	var filenames []string

	for _, path := range paths {
		if strings.HasPrefix(path, "http://") || strings.HasPrefix(path, "https://") {
			slog.Info("Reading a remote file", "url", path)
			reader, err := r.readRemoteFile(path)
			if err != nil {
				return nil, nil, err
			}
			readers = append(readers, reader)
			filenames = append(filenames, path)
		} else {
			slog.Info("Reading a local file", "path", path)
			reader, err := r.readLocalFile(path)
			if err != nil {
				return nil, nil, err
			}
			readers = append(readers, reader...)
			filenames = append(filenames, path)
		}
	}

	return readers, filenames, nil
}

func (r *FileReader) readLocalFile(path string) ([]io.Reader, error) {
	matches, err := filepath.Glob(path)
	if err != nil {
		return nil, fmt.Errorf("files didn't find in this path: %v", err)
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("files didn't find in path: %s", path)
	}

	res := make([]io.Reader, len(matches))
	for i, m := range matches {
		file, err := os.Open(m)
		if err != nil {
			return nil, fmt.Errorf("error opening file: %v", err)
		}
		res[i] = file
	}

	return res, nil
}

func (r *FileReader) readRemoteFile(url string) (io.Reader, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error loading file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("file not found (404): %s", url)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error HTTP %d: %s", resp.StatusCode, url)
	}

	return resp.Body, nil
}
