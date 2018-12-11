package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	"github.com/mimir-news/pkg/schema/user"
)

type options struct {
	credentials   user.Credentials
	clientID      string
	userAgent     string
	loginURL      string
	triggerURL    string
	triggerMethod string
}

func getOptions() options {
	credentialsFile := mustGetenv("CREDENTIALS_FILE")
	clientID := mustGetenv("CLIENT_ID")

	return options{
		credentials:   mustGetCredentials(credentialsFile),
		clientID:      clientID,
		userAgent:     createUserAgent(clientID),
		loginURL:      mustGetenv("LOGIN_URL"),
		triggerURL:    mustGetenv("TRIGGER_URL"),
		triggerMethod: mustGetenv("TRIGGER_METHOD"),
	}
}

func mustGetCredentials(filename string) user.Credentials {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	var creds user.Credentials
	err = json.Unmarshal(content, &creds)
	if err != nil {
		log.Fatal(err)
	}

	return creds
}

func mustGetenv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("No value for key: %s\n", key)
	}

	return val
}

func createUserAgent(clientID string) string {
	return fmt.Sprintf("%s %s %s/%s", clientID, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}
