package qbit

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	_QBIT_USER = "QBIT_USER"
	_QBIT_PASS = "QBIT_PASS"
)

type Client struct {
	host string
	user string
	pass string
	sid  string
}

func NewClient(host string) *Client {
	return &Client{
		host: host,
	}
}

func (c *Client) authedRequest(method, url string, body io.Reader) *http.Request {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		// requests should never be an invalid format
		panic(err)
	}
	req.AddCookie(&http.Cookie{Name: "SID", Value: c.sid})

	return req
}

func (c *Client) SourceEnv(path string) error {
	err := godotenv.Load(path)
	if err != nil {
		return err
	}

	var ok bool
	c.user, ok = os.LookupEnv(_QBIT_USER)
	if !ok {
		return fmt.Errorf("env variable not found: '%s'", _QBIT_USER)
	}

	c.pass, ok = os.LookupEnv(_QBIT_PASS)
	if !ok {
		return fmt.Errorf("env variable not found: '%s'", _QBIT_PASS)
	}
	return nil
}

func (c *Client) Login() error {
	body := url.Values{}
	body.Set("username", c.user)
	body.Set("password", c.pass)

	resp, err := http.Post(
		c.host+"/api/v2/auth/login",
		"application/x-www-form-urlencoded",
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	for _, cookie := range resp.Cookies() {
		// sid usually comes in form QBT_SID_PORT
		if strings.Contains(cookie.Name, "QBT_SID_") {
			c.sid = cookie.Value
			return nil
		}
	}

	respBody, _ := io.ReadAll(resp.Body)
	bodyStr := strings.TrimSpace(string(respBody))
	return fmt.Errorf("SID cookie not found: body='%s' cookies='%v'", bodyStr, resp.Cookies())
}

func (c *Client) Recheck(hash string) error {
	if c.sid == "" {
		return fmt.Errorf("must call .Login() first")
	}

	if hash == "" {
		return fmt.Errorf("invalid arguments: hash='%s'", hash)
	}

	body := url.Values{}
	body.Set("hashes", hash)

	req := c.authedRequest(http.MethodPost, c.host+"/api/v2/torrents/recheck", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code '%d'", resp.StatusCode)
	}

	return nil
}

func (c *Client) AddTags(hash string, tags []string) error {
	if c.sid == "" {
		return fmt.Errorf("must call .Login() first")
	}

	if hash == "" {
		return fmt.Errorf("invalid arguments: hash='%s'", hash)
	}

	body := url.Values{}
	body.Set("hashes", hash)
	body.Set("tags", strings.Join(tags, ","))

	req := c.authedRequest(http.MethodPost, c.host+"/api/v2/torrents/addTags", strings.NewReader(body.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("got status code '%d'", resp.StatusCode)
	}

	return nil
}

func (c *Client) Version() (string, error) {
	if c.sid == "" {
		return "", fmt.Errorf("must call .Login() first")
	}

	req := c.authedRequest(http.MethodGet, c.host+"/api/v2/app/version", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("got status code '%d'", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)
	version := strings.TrimSpace(string(body))

	if version == "" {
		return "", fmt.Errorf("empty response")
	}
	return version, nil
}

func (c *Client) Content(hash string) ([]Content, error) {
	if c.sid == "" {
		return nil, fmt.Errorf("must call .Login() first")
	}

	req := c.authedRequest(http.MethodGet, c.host+"/api/v2/torrents/contents", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("got status code '%d'", resp.StatusCode)
	}

	body, _ := io.ReadAll(resp.Body)

	var content []Content
	err = json.Unmarshal(body, &content)
	if err != nil {
		return nil, err
	}

	return content, nil
}
