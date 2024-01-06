package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/Scalingo/sclng-backend-test-v1/model"
	"github.com/Scalingo/sclng-backend-test-v1/queries"
	"github.com/sirupsen/logrus"
)

type IRepoRepository interface {
	IsRepositorySaved(repoId int64) (bool, error)
	existsInDBRequestByID(sql string, idOrKey interface{}) (bool, error)
	isOwnerSaved(ownerId int64) (bool, error)
	isLicenceSaved(licenceKey string) (bool, error)
	isLanguageSaved(languageName string) (bool, error)
	SaveRepository(repo *model.Repository) error
}

type LocalRepoRepository struct {
	db     *sql.DB
	logger logrus.FieldLogger
}

func NewRepoRepository(db *sql.DB, logger logrus.FieldLogger) IRepoRepository {
	return &LocalRepoRepository{
		db: db, logger: logger,
	}
}

func (repoRepository *LocalRepoRepository) IsRepositorySaved(repoId int64) (bool, error) {
	return repoRepository.existsInDBRequestByID(queries.IsGithubRepositoryExists, repoId)
}

func (repoRepository *LocalRepoRepository) SaveRepository(repo *model.Repository) error {

	// check if owner exists
	ownerSaved, err := repoRepository.isOwnerSaved(repo.Owner.ID)
	if err != nil {
		repoRepository.logger.Error("error while checking if owner exists", err)
	}
	// check if licence exits get its id
	licenceSaved, err := repoRepository.isLicenceSaved(repo.Licence.Key)
	if err != nil {
		repoRepository.logger.Error("error while checking if licence exists", err)
	}
	tx, err := repoRepository.db.Begin()
	if err != nil {
		repoRepository.logger.Error("cant begin transaction", err)
	}
	defer func() {
		// Rollback if one error detected in transaction
		if r := recover(); r != nil {
			repoRepository.logger.Error("Transaction rollback cuased by error", r)
			fmt.Println("Rollback de la transaction en raison d'une erreur:", r)
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

	if !licenceSaved {
		// Create Licence
		_, err = tx.Exec(queries.CreateLicence, repo.Licence.Key, repo.Licence.Name)
		if err != nil {
			panic(err)
		}
	}
	// Create Repository
	_, err = tx.Exec(queries.CreateRepository, repo.ID, repo.FullName, repo.Name, repo.CreatedAt, repo.Owner.ID, repo.Licence.Key)
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

func (repoRepository *LocalRepoRepository) existsInDBRequestByID(sql string, idOrKey interface{}) (bool, error) {
	var exist bool
	if err := repoRepository.db.QueryRow(sql, idOrKey).Scan(&exist); err != nil {
		if err != nil {
			return false, err
		}
		return false, err
	}
	return exist, nil
}

func (repoRepository *LocalRepoRepository) isOwnerSaved(ownerId int64) (bool, error) {
	return repoRepository.existsInDBRequestByID(queries.IsOwnerExists, ownerId)
}

func (repoRepository *LocalRepoRepository) isLicenceSaved(licenceKey string) (bool, error) {
	return repoRepository.existsInDBRequestByID(queries.IsLicenceExists, licenceKey)
}

func (repoRepository *LocalRepoRepository) isLanguageSaved(languageName string) (bool, error) {
	return repoRepository.existsInDBRequestByID(queries.IsLanguageExists, languageName)
}
