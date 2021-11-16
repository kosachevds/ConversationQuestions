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
	RawUrl   string `json:"raw_url"`
	Content  string `json:"content"`
}

type MessageResponse struct {
	Message          string `json:"message"`
	DocumentationUrl string `json:"documentation_url"`
}

func downloadAnswersFile(gistId, accessToken string) {
	downloadUrl := "https://api.github.com/gists/" + gistId

	req, err := http.NewRequest("GET", downloadUrl, nil)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	req.Header.Set("Authorization", "token "+accessToken)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	defer res.Body.Close()
	if res.StatusCode != 200 {
		printErrorMessage(res)
		os.Exit(1)
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	var jsonRes Response
	err = json.Unmarshal(body, &jsonRes)
	if err != nil {
		fmt.Printf("%s", err)
		os.Exit(1)
	}

	for _, file := range jsonRes.Files {
		fmt.Printf("Downloading %s\n", file.FileName)
		fmt.Printf("Downloading %s\n", file.RawUrl)
		ioutil.WriteFile(file.FileName, []byte(file.Content), 0660)
	}
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
