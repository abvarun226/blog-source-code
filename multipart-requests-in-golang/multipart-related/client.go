package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"time"
)

func main() {
	flag.Parse()

	positionalArgs := flag.Args()
	if len(positionalArgs) < 1 {
		log.Fatalf("This program requires at least 1 positional argument.")
	}
	
	// Metadata content.
	metadata := `{"title": "hello world", "description": "Multipart related upload test"}`
	
	// New multipart writer.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Metadata part.
	metadataHeader := textproto.MIMEHeader{}
	metadataHeader.Set("Content-Type", "application/json")
	metadataHeader.Set("Content-ID", "metadata")
	part, err := writer.CreatePart(metadataHeader)
	if err != nil {
		log.Fatalf("Error writing metadata headers: %v", err)
	}
	part.Write([]byte(metadata))

	// Media Files.
	for _, mediaFilename := range positionalArgs {
		mediaData, errRead := ioutil.ReadFile(mediaFilename)
		if errRead != nil {
			log.Fatalf("Error reading media file: %v", errRead)
		}
		mediaHeader := textproto.MIMEHeader{}
		mediaHeader.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%v\".", mediaFilename))
		mediaHeader.Set("Content-ID", "media")
		mediaHeader.Set("Content-Filename", mediaFilename)

		mediaPart, err := writer.CreatePart(mediaHeader)
		if err != nil {
			log.Fatalf("Error writing media headers: %v", errRead)
		}

		if _, err := io.Copy(mediaPart, bytes.NewReader(mediaData)); err != nil {
			log.Fatalf("Error writing media: %v", errRead)
		}
	}

	// Close multipart writer.
	if err := writer.Close(); err != nil {
		log.Fatalf("Error closing multipart writer: %v", err)
	}

	// Request Content-Type with boundary parameter.
	contentType := fmt.Sprintf("multipart/related; boundary=%s", writer.Boundary())

	// Initialize HTTP Request and headers.
	uploadURL := "http://localhost:8080/upload"
	r, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(body.Bytes()))
	if err != nil {
		log.Fatalf("Error initializing a request: %v", err)
	}
	r.Header.Set("Content-Type", contentType)
	r.Header.Set("Accept", "*/*")
	
	// HTTP Client.
	client := &http.Client{Timeout: 180 * time.Second}
	rsp, err := client.Do(r)
	if err != nil {
		log.Fatalf("Error making a request: %v", err)
	}

	// Check response status code.
	if rsp.StatusCode != http.StatusOK {
		log.Printf("Request failed with response code: %d", rsp.StatusCode)
	} else {
		log.Print("Request was a success")
	}
}