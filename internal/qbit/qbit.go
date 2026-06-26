package qbit

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

const (
	_QBIT_HOST = "http://localhost:8080"
	_QBIT_USER = "QBIT_USER"
	_QBIT_PASS = "QBIT_PASS"
)

var (
	user string
	pass string
	once sync.Once
)

func SourceEnv(path string) error {
	var e error
	once.Do(func() {
		err := godotenv.Load(path)
		if err != nil {
			e = err
			return
		}

		var ok bool
		user, ok = os.LookupEnv(_QBIT_USER)
		if !ok {
			e = fmt.Errorf("env variable not found: '%s'", _QBIT_USER)
			return
		}

		pass, ok = os.LookupEnv(_QBIT_PASS)
		if !ok {
			e = fmt.Errorf("env variable not found: '%s'", _QBIT_PASS)
			return
		}
	})

	return e
}

func Login() (string, error) {
	body := url.Values{}
	body.Set("username", user)
	body.Set("password", pass)

	resp, err := http.Post(
		_QBIT_HOST+"/api/v2/auth/login",
		"application/x-www-form-urlencoded",
		strings.NewReader(body.Encode()),
	)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	for _, c := range resp.Cookies() {
		// sid usually comes in form QBT_SID_PORT
		if strings.Contains(c.Name, "QBT_SID_") {
			return c.Value, nil
		}
	}

	respBody, _ := io.ReadAll(resp.Body)
	bodyStr := strings.TrimSpace(string(respBody))
	return "", fmt.Errorf("SID cookie not found: body='%s' cookies='%v'", bodyStr, resp.Cookies())
}

func Recheck(sid, hash string) error {
	if sid == "" || hash == "" {
		return fmt.Errorf("invalid arguments: sid='%s', hash='%s'", sid, hash)
	}

	body := url.Values{}
	body.Set("hashes", hash)

	req, err := http.NewRequest(http.MethodPost, _QBIT_HOST+"/api/v2/torrents/recheck", strings.NewReader(body.Encode()))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: "SID", Value: sid})

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

func Version(sid string) (string, error) {
	req, err := http.NewRequest(http.MethodGet, _QBIT_HOST+"/api/v2/app/version", nil)
	if err != nil {
		return "", err
	}
	req.AddCookie(&http.Cookie{Name: "SID", Value: sid})

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
