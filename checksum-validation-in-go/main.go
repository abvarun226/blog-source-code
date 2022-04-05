package main

import (
	"crypto/md5"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("secret.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	fmt.Printf("md5 checksum: %s\n", checksumMD5(f))
	fmt.Printf("sha256 checksum: %s\n", checksumSHA256(f))
}

func checksumMD5(f io.Reader) string {
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func checksumSHA256(f io.Reader) string {
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}
