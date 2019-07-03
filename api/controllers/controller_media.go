package controllers

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
	"path"
	"strconv"
)

func Download(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if len(vars["file"]) == 0 {
		http.Error(w, "File not found.", 404)
		return
	}
	filename := vars["file"]
	OpenFile, err := os.Open(path.Join("files", filename))
	defer OpenFile.Close()
	if err != nil {
		http.Error(w, "File not found.", 404)
		return
	}
	FileHeader := make([]byte, 512)
	_, _ = OpenFile.Read(FileHeader)
	FileContentType := http.DetectContentType(FileHeader)

	FileStat, _ := OpenFile.Stat()
	FileSize := strconv.FormatInt(FileStat.Size(), 10)

	w.Header().Set("Content-Disposition", "attachment; filename="+filename)
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	_, _ = OpenFile.Seek(0, 0)
	_, _ = io.Copy(w, OpenFile)
	return
}
