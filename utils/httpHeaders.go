package utils

import (
	"log"
	"net/http"
	"strings"
)

type HTTPHeaderOutput struct {
	Name     string
	Value 	 string
}

func GetURLHeaders(url string) (map[string]interface {}, error) {
	response, err := http.Head(url)

	if err != nil {
		response, err = http.Head("http://nginx:8080")
		if err != nil {
			log.Println("Error: Unable to fetch URL (", url, ") with error: ", err)
			return nil, err
		}
	}

	if response.StatusCode != http.StatusOK {
		log.Println("Error: HTTP Status = ", response.Status)
		return nil, err
	}

	headers := make(map[string]interface{})

	for k, v := range response.Header {
		headers[strings.ToLower(k)] = v[0]
	}

	return headers, nil
}

func GetURLHeaderByKey(url string, key string) (HTTPHeaderOutput, error) {
	headers, err := GetURLHeaders(url)
	key = strings.ToLower(key)

	if err != nil {
		log.Printf("Error: Cannot fetch header %v", err)
	}

	if value, ok := headers[key]; ok {
		return HTTPHeaderOutput{
			Name:  key,
			Value: value.(string),
		}, nil
	}

	return HTTPHeaderOutput{}, nil
}