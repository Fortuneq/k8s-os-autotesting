package utils

import (
	"encoding/base64"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func LoadFileFromURLToPath(from string, to string, username string, password string) (fileAbsolutePath string, err error) {
	outputFile, err := os.Create(to)
	if err != nil {
		return "", err
	}
	defer outputFile.Close()
	client := &http.Client{
		CheckRedirect: redirectPolicyFunc,
	}

	req, err := http.NewRequest("GET", from, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(username, password))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if err == nil {
		_, err = io.Copy(outputFile, resp.Body)
		if err != nil {
			return "", err
		}
		fileAbsolutePath, err = filepath.Abs(to)
		return
	}
	return
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func redirectPolicyFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", "Basic "+basicAuth("username1", "password123"))
	return nil
}
