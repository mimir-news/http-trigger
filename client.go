package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/mimir-news/pkg/httputil"
	"github.com/mimir-news/pkg/httputil/auth"
	"github.com/mimir-news/pkg/id"
	"github.com/mimir-news/pkg/schema/user"
)

const defaultTimeout time.Duration = 5 * time.Second

var emptyCredentials = user.Credentials{}

type client struct {
	httpClient *http.Client
	opts       options
	token      user.Token
}

func newClient(opts options) *client {
	c := &client{
		httpClient: &http.Client{Timeout: defaultTimeout},
		opts:       opts,
	}

	err := c.login()
	if err != nil {
		log.Fatal("Login failed:", err)
	}
	log.Println("Login OK")

	return c
}

func (c *client) login() error {
	req := c.makeRequest(c.opts.loginURL, http.MethodPost, c.opts.credentials)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return wrapErrorResponse(resp)
	}

	var token user.Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		return err
	}

	c.token = token
	return nil
}

func (c *client) trigger() error {
	req := c.makeRequest(c.opts.triggerURL, c.opts.triggerMethod, nil)
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return wrapErrorResponse(resp)
	}

	return nil
}

func (c *client) makeRequest(url, method string, body interface{}) *http.Request {
	req, err := http.NewRequest(method, url, createBody(body))
	if err != nil {
		log.Fatal("Creating request failed:", err)
	}

	req.Header.Set(auth.ClientIDKey, c.opts.clientID)
	req.Header.Set(httputil.RequestIDHeader, id.New())
	req.Header.Set("User-Agent", c.opts.userAgent)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	if c.token.Token != "" {
		bearerToken := auth.AuthTokenPrefix + c.token.Token
		req.Header.Set(auth.AuthHeaderKey, bearerToken)
	}

	return req
}

func createBody(body interface{}) io.Reader {
	var bodyReader io.Reader
	if body != nil {
		bytesBody, err := json.Marshal(body)
		if err != nil {
			log.Fatal(err)
		}
		bodyReader = bytes.NewBuffer(bytesBody)
	}

	return bodyReader
}

func wrapErrorResponse(resp *http.Response) error {
	var errResp httputil.ErrorResponse
	err := json.NewDecoder(resp.Body).Decode(&errResp)
	if err != nil {
		log.Fatalln("Failed to wrap error response", err)
	}

	bytesErr, _ := json.MarshalIndent(errResp, "", "    ")
	return errors.New(string(bytesErr))
}
