package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"main/tools"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"time"
)

func main() {
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			DisableKeepAlives: true,
			DialContext: (&net.Dialer{
				Timeout: 5 * time.Second,
			}).DialContext,
		},
	}

	// GET request
	getURL := "http://localhost:8080/hello"
	resp, err := client.Get(getURL)
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	tools.ErrorHandler(err)
	fmt.Println("GET response:", string(body))

	// POST JSON
	postJSONURL := "http://localhost:8080/post-json"
	jsonData := map[string]string{"name": "Febrio", "age": "20"}
	jsonValue, _ := json.Marshal(jsonData)
	req, err := http.NewRequest("POST", postJSONURL, bytes.NewBuffer(jsonValue))
	tools.ErrorHandler(err)
	req.Header.Set("Content-Type", "application/json")

	ctx, cancel := context.WithTimeout(req.Context(), 5*time.Second)
	defer cancel()

	req = req.WithContext(ctx)
	resp, err = client.Do(req)
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	tools.ErrorHandler(err)
	fmt.Println("POST JSON response:", string(body))

	// POST File
	postFileURL := "http://localhost:8080/post-file"
	file, err := os.Open("testfile.txt")
	tools.ErrorHandler(err)
	defer file.Close()

	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	fileWriter, err := bodyWriter.CreateFormFile("file", "testfile.txt")
	tools.ErrorHandler(err)
	_, err = io.Copy(fileWriter, file)
	tools.ErrorHandler(err)

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	req, err = http.NewRequest("POST", postFileURL, bodyBuf)
	tools.ErrorHandler(err)
	req.Header.Set("Content-Type", contentType)

	resp, err = client.Do(req)
	tools.ErrorHandler(err)
	defer resp.Body.Close()

	body, err = io.ReadAll(resp.Body)
	tools.ErrorHandler(err)
	fmt.Println("POST File response:", string(body))
}
