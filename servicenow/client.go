package servicenow

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

const OperationalStatusRetired string = "6"

// Client interacts with ServiceNow
type Client struct {
	Credentials
	Instance   string
	HTTPClient *http.Client
}

type ChangeRequestClient struct {
	*Client
	path string
}

// NewDefaultClient is how most users of service now should create a client.
// It connects to the production instance and sets a timeout on HTTP calls.
func NewDefaultClient(credentials Credentials) *Client {
	return &Client{
		Credentials: credentials,
		Instance:    "https://organization.service-now.com",
		HTTPClient: &http.Client{
			Timeout: 3 * time.Second,
		},
	}
}

// Retry executes the given request with an exponential backoff
func (c *Client) Retry(req *http.Request, tries int, wait time.Duration) (*http.Response, error) {
	response, err := c.HTTPClient.Do(req)

	if err != nil {
		return c.Retry(req, tries-1, wait*2)
	}

	if response == nil {
		time.Sleep(wait)
		if tries <= 1 {
			return nil, errors.New("error connecting to SNAPI")
		}

		return c.Retry(req, tries-1, wait*2)
	}

	return response, nil
}

// GetApplication calls ServiceNow using SNAPI to retrieve an application record
// The appID parameter takes a string in the format of APP00738
// It returns an error if the application record is nil for any reason.
func (c *Client) GetApplication(appID string) (*Application, *http.Response, string, error) {
	req, _ := http.NewRequest("GET", c.Instance+"/api/nords/v3/snapi?id="+appID, nil)

	req.SetBasicAuth(c.Username, c.Password)
	c.setRequestHeaders(req.Header)

	response, err := c.Retry(req, 5, 1*time.Second)

	if err != nil {
		return nil, nil, "", err
	}

	if response.StatusCode == http.StatusUnauthorized {
		return nil, nil, "", errors.New("unauthorized")
	}

	var envelope responseEnvelope

	body, _ := ioutil.ReadAll(response.Body)

	jsonString := string(body)

	err = json.Unmarshal(body, &envelope)

	if len(envelope.Result) == 0 {
		return nil, nil, "", errors.New("response body is empty")
	}

	app := &envelope.Result[0]
	return app, response, jsonString, err
}

func (c *Client) setRequestHeaders(h http.Header) {
	h.Set("Content-type", "application/json")
	h.Set("Accept", "application/json")
}
