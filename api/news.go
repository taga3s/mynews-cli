package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mynews/config"
	"mynews/model"
	"net/http"
	"net/url"
	"strings"
	"time"
)

func Fetch() (NewsResponse, error) {
	uri := config.NEWS_API_URI
	key := config.NEWS_API_KEY

	newsResponse := NewsResponse{}

	res, err := http.Get(uri + fmt.Sprintf("/top-headlines?country=jp&apiKey=%s", key))
	if err != nil {
		return NewsResponse{}, err
	}
	defer res.Body.Close()

	byteArray, err := io.ReadAll(res.Body)
	if err != nil {
		return NewsResponse{}, err
	}

	if err := json.Unmarshal(byteArray, &newsResponse); err != nil {
		return NewsResponse{}, err
	}

	if newsResponse.Status != "ok" {
		return NewsResponse{}, fmt.Errorf("code: %s, message: %s", newsResponse.Code, newsResponse.Message)
	}

	return newsResponse, nil
}

func Notify(newsList []model.News) error {
	message := ""
	for i, v := range newsList {
		content := fmt.Sprintf("\n[ %d ]----------\n%s\n日時: %s\n%s\n--------------\n", i+1, v.Title, v.PublishedAt, v.URL)
		message += content
	}

	form := url.Values{}
	form.Add("message", message)
	body := strings.NewReader(form.Encode())

	uri := config.NOTIFY_API_URI
	req, err := http.NewRequest("POST", uri, body)
	if err != nil {
		return err
	}

	accessToken := config.NOTIFY_API_TOKEN
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := http.Client{Timeout: 30 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	byteArray, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	response := NotifyResponse{}

	err = json.Unmarshal(byteArray, &response)
	if err != nil {
		return err
	}

	err = checkStatus(response.Status)
	if err != nil {
		return err
	}
	return nil
}

func checkStatus(status int) error {
	switch status {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		return errors.New("Bad request")
	case http.StatusUnauthorized:
		return errors.New("Invalid access token")
	case http.StatusInternalServerError:
		return errors.New("Server-side error occurred")
	default:
		return errors.New("Unknown status code received")
	}
}
