package main

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	// "encoding/json"

	"github.com/stretchr/testify/suite"
	"golang.org/x/oauth2"
)

type HydraTestSuite struct {
	suite.Suite
	hydra       *hydra
	oAuthClient *oauth2.Config
}

func oAuthClientToJSONForTest(o *oauth2.Config) interface{} {
	return struct {
		ClientID      string   `json:"client_id"`
		ClientSecret  string   `json:"client_secret"`
		Scope         string   `json:"scope"`
		RedirectURIS  []string `json:"redirect_uris"`
		GrantTypes    []string `json:"grant_types"`
		ResponseTypes []string `json:"response_types"`
	}{
		ClientID:      o.ClientID,
		ClientSecret:  o.ClientSecret,
		Scope:         strings.Join(o.Scopes, " "),
		RedirectURIS:  []string{o.RedirectURL},
		GrantTypes:    []string{"authorization_code", "refresh_token"},
		ResponseTypes: []string{"code", "id_token"},
	}
}

func (s *HydraTestSuite) SetupSuite() {
	s.hydra = &hydra{
		client: &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				return errors.New("Redirect")
			},
		},
	}

	s.oAuthClient = s.hydra.getOAuthClient()
	s.oAuthClient.ClientID = "test"
	s.oAuthClient.Endpoint.AuthURL = "http://hydra:4444/oauth2/auth"

	s.hydra.executeRequestNoData("DELETE", "/clients/"+s.oAuthClient.ClientID)
	s.hydra.executeRequest("POST", "/clients",
		oAuthClientToJSONForTest(s.oAuthClient))
}

func (s *HydraTestSuite) TearDownSuite() {
	s.hydra.executeRequestNoData("DELETE", "/clients/"+s.oAuthClient.ClientID)
}

func (s *HydraTestSuite) TestGetLoginRequest() {
	response, body := s.hydra.getLoginRequest(s.getTestLoginChallenge())
	s.Equal(http.StatusOK, response.StatusCode)
	s.NotNil(body)
}

func (s *HydraTestSuite) getTestLoginChallenge() string {
	url := s.hydra.generateAuthenticationEndpoint(s.oAuthClient)
	response, _ := s.hydra.executeAnyURLRequest("GET", url, nil)
	loginRedirectedURL, _ := response.Location()
	return loginRedirectedURL.Query().Get("login_challenge")
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(HydraTestSuite))
}
