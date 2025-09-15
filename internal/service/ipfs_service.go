package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
)

type IpfsService struct {
	APIKey    string
	APISecret string
}

type UploadFileResponse struct {
	IpfsHash    string `json:"IpfsHash"`
	PinSize     int    `json:"PinSize"`
	Timestamp   string `json:"Timestamp"`
	IsDuplicate bool   `json:"isDuplicate"`
}

func NewIpfsService(APIKey string, APISecret string) *IpfsService {
	return &IpfsService{APIKey: APIKey, APISecret: APISecret}
}

func (is *IpfsService) UploadFile(ctx context.Context, fileName string, file io.Reader) (string, error) {
	buffer := bytes.Buffer{}
	writer := multipart.NewWriter(&buffer)
	formFile, formErr := writer.CreateFormFile("file", fileName)
	if formErr != nil {
		return "", formErr
	}
	_, copyErr := io.Copy(formFile, file)
	if copyErr != nil {
		return "", copyErr
	}
	closeErr := writer.Close()
	if closeErr != nil {
		return "", closeErr
	}

	req, reqErr := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.pinata.cloud/pinning/pinFileToIPFS", &buffer)
	if reqErr != nil {
		return "", reqErr
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("pinata_api_key", is.APIKey)
	req.Header.Set("pinata_secret_api_key", is.APISecret)

	client := http.Client{}
	response, clientErr := client.Do(req)
	if clientErr != nil {
		return "", clientErr
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(response.Body)
		return "", fmt.Errorf("pinata upload failed: status=%d body=%s", response.StatusCode, string(body))
	}

	var respStruct UploadFileResponse
	if decodeErr := json.NewDecoder(response.Body).Decode(&respStruct); decodeErr != nil {
		return "", fmt.Errorf("failed to decode API response: %w", decodeErr)
	}

	return respStruct.IpfsHash, nil
}
