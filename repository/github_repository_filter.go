package repository

import (
	"database/sql"

	"github.com/Scalingo/sclng-backend-test-v1/queries"
)

type DBRepoResult struct {
	RepoID        int64
	RepoFullName  string
	RepoName      string
	OwnerId       int64
	OwnerName     string
	LicenceName   string
	LanguageName  string
	LanguageBytes int64
}

func (repoRepository *LocalRepoRepository) SearchInRepository(limit int, offset int) ([]DBRepoResult, error) {
	rows, err := repoRepository.db.Query(queries.GetAllRepositoriesWithoutFilter, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return returnDbResultsFromRows(rows)
}

func (repoRepository *LocalRepoRepository) SearchInRepositoryByLicence(limit int, offset int, licence string) ([]DBRepoResult, error) {
	rows, err := repoRepository.db.Query(queries.GetRepositoryByLicence, limit, offset, licence)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return returnDbResultsFromRows(rows)
}

func (repoRepository *LocalRepoRepository) SearchInRepositoryByLanguage(limit int, offset int, language string) ([]DBRepoResult, error) {
	rows, err := repoRepository.db.Query(queries.GetRepositoryByLanguage, limit, offset, language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return returnDbResultsFromRows(rows)
}

func (repoRepository *LocalRepoRepository) SearchInRepositoryByLanguageAndLicence(limit int, offset int, licence string, language string) ([]DBRepoResult, error) {
	rows, err := repoRepository.db.Query(queries.GetRepositoryByLicenceAndLanguage, limit, offset, licence, language)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return returnDbResultsFromRows(rows)
}

func returnDbResultsFromRows(rows *sql.Rows) ([]DBRepoResult, error) {
	var dbRepoResults []DBRepoResult

	for rows.Next() {
		var dbRepoResult DBRepoResult
		err := rows.Scan(&dbRepoResult.RepoID, &dbRepoResult.RepoFullName, &dbRepoResult.RepoName, &dbRepoResult.OwnerId, &dbRepoResult.OwnerName, &dbRepoResult.LicenceName, &dbRepoResult.LanguageName, &dbRepoResult.LanguageBytes)
		if err != nil {
			return nil, err
		}
		dbRepoResults = append(dbRepoResults, dbRepoResult)
	}

	return dbRepoResults, nil
}
