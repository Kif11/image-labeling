package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type ResponseData struct {
	Message map[string]string `json:"message"`
}

func main() {
	err := filepath.Walk("./captionTest", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := strings.ToLower(filepath.Ext(path))
			imgExtensions := []string{".jpg", ".png", ".bmp", ".gif", ".webp", ".tif", ".tiff"}
			for _, v := range imgExtensions {
				if ext != v {
					continue
				}

				fmt.Printf("[+] Processing image %s\n", path)
				processImage(path)
			}
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Encountered error while processing images: %v", err)
	}
}

func processImage(path string) {
	tmpImg := "/tmp/" + filepath.Base(path)
	cmdStr := []string{"ffmpeg", "-y", "-i", path, "-q:v", "1", "-vf", `scale="1920:trunc(ow/a/2)*2"`, tmpImg}
	fmt.Printf("%s\n", strings.Join(cmdStr, " "))
	cmd := exec.Command(cmdStr[0], cmdStr[:1]...)
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to process image with ffmpeg: %v", err)
		return
	}

	d, err := uploadImage(tmpImg)
	if err != nil {
		log.Printf("Failed to upload image: %v", err)
		return
	}

	cmd = exec.Command("exiftool.exe", "-description="+d, path)
	err = cmd.Run()
	if err != nil {
		log.Printf("Failed to set image description with exiftool: %v", err)
		return
	}

	err = os.Remove(tmpImg)
	if err != nil {
		log.Printf("Failed to remove temporary image file: %v", err)
	}
}

func uploadImage(path string) (string, error) {
	var b bytes.Buffer
	w := multipart.NewWriter(&b)
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}

	fw, err := w.CreateFormFile("image", path)
	if err != nil {
		return "", err
	}

	_, err = io.Copy(fw, f)
	if err != nil {
		return "", err
	}

	err = w.Close()
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "http://localhost:5000/upload", &b)
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	r := ResponseData{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", err
	}

	err = res.Body.Close()
	if err != nil {
		return "", err
	}

	return r.Message["<MORE_DETAILED_CAPTION>"], nil
}
