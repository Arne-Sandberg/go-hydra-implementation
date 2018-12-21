package main

import (
	"bytes"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"

	"golang.org/x/oauth2"
)

type hydra struct {
	client *http.Client
}

func (h *hydra) getOAuthClient() *oauth2.Config {
	return &oauth2.Config{
		ClientID:     os.Getenv("HYDRA_CLIENT_ID"),
		ClientSecret: os.Getenv("HYDRA_SECRET"),
		Endpoint: oauth2.Endpoint{
			TokenURL: "",
			AuthURL:  os.Getenv("HYDRA_PUBLIC_URL") + "/oauth2/auth",
		},
		RedirectURL: "http://127.0.0.1:5000/callback",
		Scopes:      []string{"openid", "offline"},
	}
}

func (h *hydra) generateAuthenticationEndpoint(conf *oauth2.Config) string {
	state := generateRandomKey(24)
	nonce := generateRandomKey(24)
	return conf.AuthCodeURL(state) + "&nonce=" + nonce + "&prompt=&max_age=0"
}

func (h *hydra) getLoginRequest(challenge string) (*http.Response, interface{}) {
	return h.getRequest("login", challenge)
}

func (h *hydra) acceptLoginRequest(challenge, data string) (*http.Response, interface{}) {
	return h.putRequest("login", "accept", challenge, data)
}

func (h *hydra) rejectLoginRequest(challenge, data string) (*http.Response, interface{}) {
	return h.putRequest("login", "reject", challenge, data)
}

func (h *hydra) getConsentRequest(challenge string) (*http.Response, interface{}) {
	return h.getRequest("consent", challenge)
}

func (h *hydra) acceptConsentRequest(challenge, data string) (*http.Response, interface{}) {
	return h.putRequest("consent", "accept", challenge, data)
}

func (h *hydra) rejectConsentRequest(challenge, data string) (*http.Response, interface{}) {
	return h.putRequest("consent", "reject", challenge, data)
}

func (h *hydra) getRequest(flow, challenge string) (*http.Response, interface{}) {
	return h.executeRequestNoData("GET",
		path.Join("/oauth2/auth/requests", flow, challenge),
	)
}

func (h *hydra) putRequest(flow, action, challenge string, data interface{}) (*http.Response, interface{}) {
	return h.executeRequest("PUT",
		path.Join("/oauth2/auth/requests", flow, challenge, action),
		data,
	)
}

func (h *hydra) executeRequestNoData(method, url string) (*http.Response, interface{}) {
	return h.executeRequest(method, url, nil)
}

func (h *hydra) executeRequest(method, url string, data interface{}) (*http.Response, interface{}) {
	return h.executeAnyURLRequest(method, os.Getenv("HYDRA_ADMIN_URL")+url, data)
}

func (h *hydra) executeAnyURLRequest(method, url string, data interface{}) (*http.Response, interface{}) {
	jsonData, _ := json.Marshal(data)

	request, _ := http.NewRequest(method, url, bytes.NewBuffer(jsonData))
	request.Header = map[string][]string{
		"Content-Type": {"application/json"},
		"Accept":       {"application/json"},
	}

	response, _ := h.client.Do(request)
	defer response.Body.Close()

	var body interface{}
	json.NewDecoder(response.Body).Decode(&body)

	return response, body
}

// Source: https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go/22892986
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var srcRand = rand.NewSource(time.Now().UnixNano())

func generateRandomKey(n int) string {
	b := make([]byte, n)
	// A srcRand.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, srcRand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = srcRand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
