package response

import "github.com/Scalingo/sclng-backend-test-v1/repository"

type Repository struct {
	ID             int64                       `json:"-"`
	Fullname       string                      `json:"full_name"`
	OwnerName      string                      `json:"owner"`
	RepositoryName string                      `json:"repository"`
	Licence        string                      `json:"-"`
	Languages      map[string]map[string]int64 `json:"languages"`
}
type languages struct {
}
type RepositoriesResponse struct {
	Repositories []Repository `json:"repositories"`
}

func GenerateRepositoriesResponseFromDBResult(dbRepoResults []repository.DBRepoResult) *RepositoriesResponse {
	return &RepositoriesResponse{Repositories: createRepositoriesFromDBResult(dbRepoResults)}
}
func createRepositoriesFromDBResult(dbRepoResults []repository.DBRepoResult) []Repository {
	var repos []Repository
	for _, dbrepo := range dbRepoResults {
		if !existsInRepositories(dbrepo.RepoID, repos) {
			repo := Repository{ID: dbrepo.RepoID, Fullname: dbrepo.RepoFullName, OwnerName: dbrepo.OwnerName,
				RepositoryName: dbrepo.RepoName, Licence: dbrepo.LicenceName,
				Languages: extractLanguagesForRepository(dbrepo.RepoID, dbRepoResults)}
			repos = append(repos, repo)
		}
	}
	return repos
}

func extractLanguagesForRepository(repoId int64, dbRepoResults []repository.DBRepoResult) map[string]map[string]int64 {
	var languages map[string]map[string]int64 = make(map[string]map[string]int64)
	for _, repo := range dbRepoResults {
		if repo.RepoID == repoId {
			languages[repo.LanguageName] = make(map[string]int64)
			languages[repo.LanguageName]["bytes"] = repo.LanguageBytes
		}
	}
	return languages
}

func existsInRepositories(repoId int64, repos []Repository) bool {

	for _, repo := range repos {
		if repo.ID == repoId {
			return true
		}
	}
	return false
}
