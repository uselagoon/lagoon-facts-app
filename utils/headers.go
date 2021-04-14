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

func GetURLHeaders(url string) map[string]interface {} {
	response, err := http.Head(url)

	if err != nil {
		response, err = http.Head("http://nginx:8080")
		if err != nil {
			log.Fatal("Error: Unable to fetch URL (", url, ") with error: ", err)
		}
	}

	if response.StatusCode != http.StatusOK {
		log.Fatal("Error: HTTP Status = ", response.Status)
	}

	headers := make(map[string]interface{})

	for k, v := range response.Header {
		headers[strings.ToLower(k)] = v[0]
	}

	return headers
}

func GetURLHeaderByKey(url string, key string) HTTPHeaderOutput {
	headers := GetURLHeaders(url)
	key = strings.ToLower(key)

	if value, ok := headers[key]; ok {
		return HTTPHeaderOutput{
			Name:  key,
			Value: value.(string),
		}
	}

	return HTTPHeaderOutput{}
}