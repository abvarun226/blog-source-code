package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		parseErr := r.ParseMultipartForm(32 << 20)
		if parseErr != nil {
			http.Error(w, "failed to parse multipart message", http.StatusBadRequest)
			return
		}

		if r.MultipartForm == nil || r.MultipartForm.File == nil {
			http.Error(w, "expecting multipart form file", http.StatusBadRequest)
			return
		}

		if err := verifyRequest(r); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		metadata, errMeta := getMetadata(r)
		if errMeta != nil {
			http.Error(w, "failed to get metadata", http.StatusBadRequest)
			return
		}
		log.Print(string(metadata))

		for _, h := range r.MultipartForm.File["media"] {
			file, err := h.Open()
			if err != nil {
				http.Error(w, "failed to get media form file", http.StatusBadRequest)
				return
			}
			uploadMedia(file, h.Filename)
		}
	})

	http.ListenAndServe(":8080", mux)
}

func uploadMedia(file multipart.File, filename string) {
	defer file.Close()
	tmpfile, _ := os.Create("./" + filename)
	defer tmpfile.Close()
	io.Copy(tmpfile, file)
}

func getMetadata(r *http.Request) ([]byte, error) {
	f, _, err := r.FormFile("metadata")
	if err != nil {
		return nil, fmt.Errorf("failed to get metadata form file: %v", err)
	}

	metadata, errRead := ioutil.ReadAll(f)
	if errRead != nil {
		return nil, fmt.Errorf("failed to read metadata: %v", errRead)
	}

	return metadata, nil
}

func verifyRequest(r *http.Request) error {
	if _, ok := r.MultipartForm.File["media"]; !ok {
		return fmt.Errorf("media is absent")
	}

	if _, ok := r.MultipartForm.File["metadata"]; !ok {
		return fmt.Errorf("metadata is absent")
	}

	return nil
}
