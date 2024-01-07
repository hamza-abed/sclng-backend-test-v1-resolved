package api

import (
	"encoding/json"
	"net/http"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/api/response"
	"github.com/Scalingo/sclng-backend-test-v1/repository"
)

func (server *Server) getRepositories(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	licenceFilter := r.URL.Query().Get("licence")
	languageFilter := r.URL.Query().Get("language")
	limit_str := r.URL.Query().Get("limit")
	offset_str := r.URL.Query().Get("offset")

	limit := 100
	offset := 0
	if limit_str != "" {
		limit = 5
	}
	if offset_str != "" {
		offset = 5
	}
	log.Info("licence : ", licenceFilter)
	log.Info("language : ", languageFilter)
	log.Info("limit : ", limit)
	log.Info("offset : ", offset)
	// get all repositories

	// if limit and offset not spec put a default value
	var dbRepoResults []repository.DBRepoResult
	var err error
	if licenceFilter != "" && languageFilter != "" { // filter by licence and language
		dbRepoResults, err = server.repoRepository.SearchInRepositoryByLanguageAndLicence(limit, offset, licenceFilter, languageFilter)
	} else if licenceFilter != "" { // filter only by licence
		dbRepoResults, err = server.repoRepository.SearchInRepositoryByLicence(limit, offset, licenceFilter)
	} else if languageFilter != "" { // filter only by language
		dbRepoResults, err = server.repoRepository.SearchInRepositoryByLanguage(limit, offset, languageFilter)
	} else { // no filter
		dbRepoResults, err = server.repoRepository.SearchInRepository(limit, offset)
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		server.log.Error("error : ", err)
	}

	repositories := response.GenerateRepositoriesResponseFromDBResult(dbRepoResults)
	err = json.NewEncoder(w).Encode(repositories)
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
	}
	return nil
}
