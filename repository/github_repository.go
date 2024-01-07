package repository

import (
	"database/sql"
	"log"
	"sync"

	"github.com/Scalingo/sclng-backend-test-v1/model"
	"github.com/Scalingo/sclng-backend-test-v1/queries"
	"github.com/sirupsen/logrus"
)

type IRepoRepository interface {
	IsRepositorySaved(repoId int64) (bool, error)
	isOwnerSaved(ownerId int64) (bool, error)
	isLicenceSaved(licenceKey string) (bool, error)
	isLanguageSaved(languageName string) (bool, error)
	SaveRepository(repo *model.Repository) error
	GetMaxRepositoryId() (int64, error)
	SearchInRepository(limit int, offset int) ([]DBRepoResult, error)
	SearchInRepositoryByLicence(limit int, offset int, licence string) ([]DBRepoResult, error)
	SearchInRepositoryByLanguage(limit int, offset int, language string) ([]DBRepoResult, error)
	SearchInRepositoryByLanguageAndLicence(limit int, offset int, licence string, language string) ([]DBRepoResult, error)
}

type LocalRepoRepository struct {
	db     *sql.DB
	logger logrus.FieldLogger
	mutex  sync.Mutex
}

func NewRepoRepository(db *sql.DB, logger logrus.FieldLogger) IRepoRepository {
	return &LocalRepoRepository{
		db: db, logger: logger,
	}
}

func (repoRepository *LocalRepoRepository) IsRepositorySaved(repoId int64) (bool, error) {
	return repoRepository.existsInDBRequestByIDOrKey(queries.IsGithubRepositoryExists, repoId)
}

func (repoRepository *LocalRepoRepository) SaveRepository(repo *model.Repository) error {
	// protect transactions from deadlock
	repoRepository.mutex.Lock()
	defer repoRepository.mutex.Unlock()

	// check if owner exists
	ownerSaved, err := repoRepository.isOwnerSaved(repo.Owner.ID)
	if err != nil {
		repoRepository.logger.Error("error while checking if owner exists", err)
	}
	// check if licence exits get its id
	licenceSaved := false
	if repo.Licence != nil {
		licenceSaved, err = repoRepository.isLicenceSaved(repo.Licence.Key)
		if err != nil {
			repoRepository.logger.Error("error while checking if licence exists", err)
		}
	}

	tx, err := repoRepository.db.Begin()
	if err != nil {
		repoRepository.logger.Error("cant begin transaction", err)
	}
	defer func() {
		// Rollback if one error detected in transaction
		if r := recover(); r != nil {
			repoRepository.logger.Error("Transaction rollback cuased by error ", r, "  licenceSaved = ", licenceSaved, " ownerSaved=", ownerSaved)
			tx.Rollback()
		}
	}()
	if err != nil {
		log.Fatal(err)
	}

	if !ownerSaved {
		// Create Owner
		_, err = tx.Exec(queries.CreateOwner, repo.Owner.ID, repo.Owner.Name)
		if err != nil {
			panic(err)
		}
	}

	if repo.Licence != nil && !licenceSaved {
		// Create Licence
		_, err = tx.Exec(queries.CreateLicence, repo.Licence.Key, repo.Licence.Name)
		if err != nil {
			panic(err)
		}
	}

	// Create Repository
	if repo.Licence != nil {
		_, err = tx.Exec(queries.CreateRepository, repo.ID, repo.FullName, repo.Name, repo.CreatedAt, repo.Owner.ID, repo.Licence.Key)
	} else {
		_, err = tx.Exec(queries.CreateRepository, repo.ID, repo.FullName, repo.Name, repo.CreatedAt, repo.Owner.ID, nil)
	}

	if err != nil {
		panic(err)
	}
	if len(repo.Languages) > 0 {
		// Create languages
		for _, lang := range repo.Languages {
			langSaved, err := repoRepository.isLanguageSaved(lang.Name)
			if err != nil {
				repoRepository.logger.Error("error while checking if language exists", err)
			}
			if !langSaved {
				_, err = tx.Exec(queries.CreateLanguage, lang.Name)
				if err != nil {
					panic(err)
				}
			}
			// Create associations
			_, err = tx.Exec(queries.CreateRepositoryLanguage, repo.ID, lang.Name, lang.Bytes)
			if err != nil {
				panic(err)
			}
		}
	}

	return tx.Commit()
}

// This returns the most recent repository ID, and -1 if DB is empty
func (repoRepository *LocalRepoRepository) GetMaxRepositoryId() (int64, error) {
	var maxId int64 = -1
	if err := repoRepository.db.QueryRow(queries.GetMaxRepositoryId).Scan(&maxId); err != nil {
		if err != nil {
			repoRepository.logger.Error("error while request : ", queries.GetMaxRepositoryId, " ", err)
			return maxId, err
		}
		repoRepository.logger.Error("error while request : ", queries.GetMaxRepositoryId, " ", err)
		return maxId, err
	}
	return maxId, nil
}

func (repoRepository *LocalRepoRepository) existsInDBRequestByIDOrKey(sql string, key interface{}) (bool, error) {
	var exist bool
	if err := repoRepository.db.QueryRow(sql, key).Scan(&exist); err != nil {
		if err != nil {
			repoRepository.logger.Error("error while request : ", sql, " ", err)
			return false, err
		}
		repoRepository.logger.Error("error while request : ", sql, " ", err)
		return false, err
	}
	return exist, nil
}
func (repoRepository *LocalRepoRepository) isOwnerSaved(ownerId int64) (bool, error) {
	return repoRepository.existsInDBRequestByIDOrKey(queries.IsOwnerExists, ownerId)
}

func (repoRepository *LocalRepoRepository) isLicenceSaved(licenceKey string) (bool, error) {
	return repoRepository.existsInDBRequestByIDOrKey(queries.IsLicenceExists, licenceKey)
}

func (repoRepository *LocalRepoRepository) isLanguageSaved(languageName string) (bool, error) {
	return repoRepository.existsInDBRequestByIDOrKey(queries.IsLanguageExists, languageName)
}
