package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Response struct {
	Files map[string]FileDetails `json:"files"`
}

type FileDetails struct {
	FileName string `json:"filename"`
	Content  string `json:"content"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

func downloadAnswersFile(gistId, accessToken string) error {
	downloadUrl := "https://api.github.com/gists/" + gistId

	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "token "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		printErrorMessage(res)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	var jsonRes Response
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		return err
	}

	for _, file := range jsonRes.Files {
		ioutil.WriteFile(file.FileName, []byte(file.Content), 0660)
	}
	return nil
}

func printErrorMessage(resp *http.Response) {
	var jsonRes MessageResponse
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	err = json.Unmarshal(body, &jsonRes)
	fmt.Printf("Sorry, %s\n", jsonRes.Message)
}
