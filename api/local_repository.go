package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/sclng-backend-test-v1/api/response"
	"github.com/Scalingo/sclng-backend-test-v1/repository"
)

func (server *Server) getRepositories(w http.ResponseWriter, r *http.Request, _ map[string]string) error {
	log := logger.Get(r.Context())
	w.Header().Add("Content-Type", "application/json")

	licenceFilter := r.URL.Query().Get("licence")
	languageFilter := r.URL.Query().Get("language")
	limit_str := r.URL.Query().Get("limit")
	offset_str := r.URL.Query().Get("offset")

	limit := 100
	offset := 0
	var err error
	if limit_str != "" {
		limit, err = strconv.Atoi(limit_str)
		if err != nil {
			log.WithError(err).Error("Fail to convert limit ", err)
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Mauvaise requête", http.StatusBadRequest)
			return err
		}
	}
	if offset_str != "" {
		offset, err = strconv.Atoi(offset_str)
		if err != nil {
			log.WithError(err).Error("Fail to convert offset ", err)
			w.WriteHeader(http.StatusBadRequest)
			http.Error(w, "Mauvaise requête", http.StatusBadRequest)
			return err
		}
	}

	// if limit and offset not spec put a default value
	var dbRepoResults []repository.DBRepoResult

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
		log.WithError(err).Error("Fail to while requesting db ", err)
		http.Error(w, "Erreur lors de la cnx au serveur", http.StatusInternalServerError)
	}

	repositories := response.GenerateRepositoriesResponseFromDBResult(dbRepoResults)
	err = json.NewEncoder(w).Encode(repositories)
	if err != nil {
		log.WithError(err).Error("Fail to encode JSON")
		http.Error(w, "Erreur lors de l'encodage JSON", http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	return nil
}
