package client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/appkins/terraform-provider-synology/synology/client/api"
	"github.com/appkins/terraform-provider-synology/synology/client/util/form"
	"github.com/google/go-querystring/query"
)

// Login runs a login flow to retrieve session token from Synology.
func (c *synologyClient) Login(user, password, sessionName string) error {
	u := c.baseURL

	q := u.Query()
	q.Add("api", "SYNO.API.Auth")
	q.Add("version", "7")
	q.Add("method", "login")
	q.Add("account", user)
	q.Add("passwd", password)
	q.Add("session", sessionName)
	q.Add("format", "cookie")
	u.RawQuery = q.Encode()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}()

	return nil
}

func (c synologyClient) Post(r api.Request, response api.Response) error {
	u := c.baseURL

	// Prepare a form that you will submit to that URL.
	if b, err := form.Marshal(r); err != nil {
		return err
	} else {

		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(b))
		if err != nil {
			return err
		}

		return c.Do(req, response)
	}
}

// Get performs an HTTP request to remote Synology instance.
//
// Returns error in case of any transport errors.
// For API-level errors, check response object.
func (c synologyClient) Get(r api.Request, response api.Response) error {
	u := c.baseURL

	if q, err := query.Values(r); err != nil {
		return err
	} else {
		u.RawQuery = q.Encode()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return err
	}

	return c.Do(req, response)
}

func GetResponse[T any | api.ApiError](c SynologyClient, r api.Request) (response api.Response, err error) {

	if err = c.Get(r, response); err != nil {
		return
	}
	return
}

func Do[T any | api.ApiError](client *http.Client, req *http.Request) (*T, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_, _ = io.ReadAll(resp.Body)
		_ = resp.Body.Close()
	}()

	synoResponse := api.ApiResponse[T]{}
	if err := json.NewDecoder(resp.Body).Decode(&synoResponse); err != nil {
		return nil, err
	}
	return &synoResponse.Data, nil

	// response.SetError(handleErrors(synoResponse, response, api.GlobalErrors))
	// return nil
}

func handleErrors[T any | api.ApiError](response api.ApiResponse[T], errorDescriber api.ErrorDescriber, knownErrors api.ErrorSummary) api.ApiError {
	err := api.ApiError{
		Code: response.Error.Code,
	}
	if response.Error.Code == 0 {
		return err
	}

	combinedKnownErrors := append(errorDescriber.ErrorSummaries(), knownErrors)
	err.Summary = api.DescribeError(err.Code, combinedKnownErrors...)
	for _, e := range response.Error.Errors {
		item := api.ErrorItem{
			Code:    e.Code,
			Summary: api.DescribeError(e.Code, combinedKnownErrors...),
		}
		if len(e.Details) > 0 {
			item.Details = make(api.ErrorFields)
			for k, v := range e.Details {
				item.Details[k] = v
			}
		}
		err.Errors = append(err.Errors, item)
	}

	return err
}
