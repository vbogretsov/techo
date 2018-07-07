package techo_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/suite"

	"github.com/vbogretsov/techo"
)

const (
	authrorization = "Authorization"
	token          = "Bearer XXX"
)

type Body struct {
	Message string `json:"message"`
}

type Suite struct {
	suite.Suite
	client   *techo.Client
	header   http.Header
	exampleB []byte
	exampleV Body
}

func (s *Suite) SetupSuite() {
	e := echo.New()

	cpheader := func(c echo.Context) {
		reqH := c.Request().Header
		resH := c.Response().Header()

		for k, v := range reqH {
			resH[k] = v
		}
	}

	cpbody := func(c echo.Context) error {
		body := Body{}
		if err := c.Bind(&body); err != nil {
			return err
		}

		return c.JSON(http.StatusOK, &body)
	}

	e.GET("/test", func(c echo.Context) error {
		cpheader(c)
		return c.JSON(http.StatusOK, &Body{Message: "Test"})
	})
	e.POST("/test", func(c echo.Context) error {
		cpheader(c)
		return cpbody(c)
	})
	e.PUT("/test", func(c echo.Context) error {
		cpheader(c)
		return cpbody(c)
	})
	e.PATCH("/test", func(c echo.Context) error {
		cpheader(c)
		return cpbody(c)
	})
	e.DELETE("/test", func(c echo.Context) error {
		cpheader(c)
		return cpbody(c)
	})

	s.client = techo.New(e, json.Marshal)
	s.header = make(http.Header)
	s.header.Set(authrorization, token)
	s.exampleV = Body{Message: "Test"}
	s.exampleB, _ = json.Marshal(&s.exampleV)
}

func (s *Suite) TestGet() {
	resp := s.client.Get("/test", s.header)
	require.Equal(s.T(), http.StatusOK, resp.Code)
	require.Equal(s.T(), resp.Body, s.exampleB)
	require.Equal(s.T(), resp.Header.Get(authrorization), token)
}

func (s *Suite) TestPost() {
	resp := s.client.Post("/test", s.header, s.exampleV)
	require.Equal(s.T(), http.StatusOK, resp.Code)
	require.Equal(s.T(), resp.Body, s.exampleB)
	require.Equal(s.T(), resp.Header.Get(authrorization), token)
}

func (s *Suite) TestPut() {
	resp := s.client.Put("/test", s.header, s.exampleV)
	require.Equal(s.T(), http.StatusOK, resp.Code)
	require.Equal(s.T(), resp.Body, s.exampleB)
	require.Equal(s.T(), resp.Header.Get(authrorization), token)
}

func (s *Suite) TestPatch() {
	resp := s.client.Patch("/test", s.header, s.exampleV)
	require.Equal(s.T(), http.StatusOK, resp.Code)
	require.Equal(s.T(), resp.Body, s.exampleB)
	require.Equal(s.T(), resp.Header.Get(authrorization), token)
}

func (s *Suite) TestDelete() {
	resp := s.client.Delete("/test", s.header, s.exampleV)
	require.Equal(s.T(), http.StatusOK, resp.Code)
	require.Equal(s.T(), resp.Body, s.exampleB)
	require.Equal(s.T(), resp.Header.Get(authrorization), token)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
