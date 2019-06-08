package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/upload", func(w http.ResponseWriter, r *http.Request) {
		contentType, params, parseErr := mime.ParseMediaType(r.Header.Get("Content-Type"))
		if parseErr != nil || !strings.HasPrefix(contentType, "multipart/") {
			log.Printf("invalid parameters: %v", parseErr)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		multipartReader := multipart.NewReader(r.Body, params["boundary"])
		defer r.Body.Close()

		for {
			part, err := multipartReader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				http.Error(w, "unexpected error when retrieving a part of the message", http.StatusInternalServerError)
				return
			}
			defer part.Close()
		
			fileBytes, err := ioutil.ReadAll(part)
			if err != nil {
				http.Error(w, "failed to read content of the part", http.StatusInternalServerError)
				return
			}

			switch part.Header.Get("Content-ID") {
			case "metadata":
				log.Print(string(fileBytes))
		
			case "media":
				log.Printf("filesize = %d", len(fileBytes))
				f, _ := os.Create(part.Header.Get("Content-Filename"))
				f.Write(fileBytes)
				f.Close()
			}
		}
	})

	http.ListenAndServe(":8080", mux)
}