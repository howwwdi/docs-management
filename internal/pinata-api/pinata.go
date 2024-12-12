package pinata_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
)

const (
	baseUrl = "https://api.pinata.cloud/pinning/pinFileToIPFS"
)

func UploadFile(filename string, data []byte) (string, error) {

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	_, err = part.Write(data)
	if err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}
	writer.Close()

	req, err := http.NewRequest("POST", baseUrl, body)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("pinata_api_key", os.Getenv("PINATA_API_KEY"))
	req.Header.Set("pinata_secret_api_key", os.Getenv("PINATA_SECRET"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}
	defer resp.Body.Close()

	bodyResp, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Failed to read response: %v", err)
	}

	if resp.StatusCode == http.StatusOK {
		log.Println("File uploaded successfully")
		var result map[string]interface{}
		if err := json.Unmarshal(bodyResp, &result); err != nil {
			log.Fatalf("Failed to unmarshal response: %v", err)
		}
		log.Printf("IPFS hash: %s\n", result["IpfsHash"])
		return result["IpfsHash"].(string), nil
	} else {
		return "", fmt.Errorf("failed to upload file: %s", bodyResp)
	}
}

func DownloadFromPinata(cid, filename string) ([]byte, error) {
	reqUrl := fmt.Sprintf("https://gateway.pinata.cloud/ipfs/%s", cid)

	proxyUrl, _ := url.Parse(os.Getenv("PROXY_URL"))

	client := http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
		},
	}

	resp, err := client.Get(reqUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Pinata: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}

	data, _ := io.ReadAll(resp.Body)

	log.Printf("File downloaded successfully: %s\n", filename)
	return data, nil
}

func DeleteFromPinata(cid string) error {
	url := fmt.Sprintf("https://api.pinata.cloud/pinning/unpin/%s", cid)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("pinata_api_key", os.Getenv("PINATA_API_KEY"))
	req.Header.Set("pinata_secret_api_key", os.Getenv("PINATA_SECRET"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	} else {
		return fmt.Errorf("failed to delete file: status code %d", resp.StatusCode)
	}
}
