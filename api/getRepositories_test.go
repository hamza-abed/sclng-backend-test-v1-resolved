package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Scalingo/sclng-backend-test-v1/api/response"
	"github.com/stretchr/testify/require"
)

func TestGetRepositoriesWithoutFilters(t *testing.T) {
	server := newTestServer(t)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "", nil)
	require.NoError(t, err)
	server.getRepositories(rr, req, make(map[string]string))
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)
	t.Log(rr.Result())

	var rResponse response.RepositoriesResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&rResponse)

	require.NoError(t, err)

	require.Equal(t, 1, len(rResponse.Repositories))
	require.Equal(t, "repotest", rResponse.Repositories[0].Fullname)
}

func TestGetRepositoriesWithFilterLanguage(t *testing.T) {
	server := newTestServer(t)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "?langauge=Ruby", nil)
	require.NoError(t, err)
	server.getRepositories(rr, req, make(map[string]string))
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)
	t.Log(rr.Result())

	var rResponse response.RepositoriesResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&rResponse)

	require.NoError(t, err)

	require.Equal(t, 1, len(rResponse.Repositories))
	require.Equal(t, "repotest", rResponse.Repositories[0].Fullname)
}

func TestGetRepositoriesWithFilterLicence(t *testing.T) {
	server := newTestServer(t)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "?licence=mit", nil)
	require.NoError(t, err)
	server.getRepositories(rr, req, make(map[string]string))
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)
	t.Log(rr.Result())

	var rResponse response.RepositoriesResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&rResponse)

	require.NoError(t, err)

	require.Equal(t, 1, len(rResponse.Repositories))
	require.Equal(t, "repotest", rResponse.Repositories[0].Fullname)
}

func TestGetRepositoriesWithFilterLicenceAndLanguage(t *testing.T) {
	server := newTestServer(t)
	rr := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, "?licence=mit&language=Ruby", nil)
	require.NoError(t, err)
	server.getRepositories(rr, req, make(map[string]string))
	require.Equal(t, http.StatusOK, rr.Result().StatusCode)
	t.Log(rr.Result())

	var rResponse response.RepositoriesResponse
	err = json.NewDecoder(rr.Result().Body).Decode(&rResponse)

	require.NoError(t, err)

	require.Equal(t, 1, len(rResponse.Repositories))
	require.Equal(t, "repotest", rResponse.Repositories[0].Fullname)
}
