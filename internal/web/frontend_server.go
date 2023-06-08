package web

import (
	"embed"
	"errors"
	"io"
	"io/fs"
	"mime"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

//go:embed frontend/*
var frontend embed.FS

func frontendHandler(notFoundHandler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := strings.TrimSpace(strings.TrimSuffix(request.URL.Path, "/"))

		isFirstLevel := strings.Count(path, "/") <= 1

		file, err := frontend.Open(filepath.Join("frontend", path))
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				if isFirstLevel {
					serveIndexFile(writer, request)
				} else {
					notFoundHandler(writer, request)
				}
				return
			}
			writeErr(writer, err)
			return
		}
		defer func() {
			_ = file.Close()
		}()

		fileInfo, err := file.Stat()
		if err != nil {
			writeErr(writer, err)
			return
		}

		if fileInfo.IsDir() {
			if isFirstLevel {
				serveIndexFile(writer, request)
			} else {
				notFoundHandler(writer, request)
			}
			return
		}

		content, err := io.ReadAll(file)
		if err != nil {
			writeErr(writer, err)
			return
		}

		writer.Header().Set("Content-Type", mime.TypeByExtension(fileInfo.Name()))
		writer.Header().Set("Content-Length", strconv.Itoa(len(content)))
		_, _ = writer.Write(content)
	}
}

func serveIndexFile(writer http.ResponseWriter, _ *http.Request) {
	indexFile, err := frontend.ReadFile("frontend/index.html")
	if err != nil {
		writeErr(writer, err)
		return
	}
	writer.Header().Set("Content-Type", "text/html")
	writer.Header().Set("Content-Length", strconv.Itoa(len(indexFile)))
	_, _ = writer.Write(indexFile)
}
