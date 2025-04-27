package main

import (
	"strings"
)

type checkedBody struct {
	CleanedBody string `json:"cleaned_body"`
}

func badWordReplacer(body string) checkedBody {
	sliceBody := strings.Split(body, " ")

	for i, word := range sliceBody {
		lowerWord := strings.ToLower(word)
		if lowerWord == "kerfuffle" || lowerWord == "sharbert" || lowerWord == "fornax" {
			word = "****"
			sliceBody[i] = word
		}
	}

	cleanedBody := checkedBody{
		CleanedBody: strings.Join(sliceBody, " "),
	}

	return cleanedBody
}
