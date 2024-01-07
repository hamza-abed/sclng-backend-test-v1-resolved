package factory

import (
	"github.com/Scalingo/sclng-backend-test-v1/model"
	"github.com/google/go-github/github"
	"github.com/sirupsen/logrus"
)

func CreateRepository(repo *github.Repository, languages []*model.Language, logger logrus.FieldLogger) model.Repository {
	return model.Repository{ID: repo.GetID(), FullName: repo.GetFullName(), Name: repo.GetName(), CreatedAt: repo.GetCreatedAt().Time, Owner: CreateOwner(repo.GetOwner()), Licence: CreateLicence(repo.GetLicense()), Languages: languages}
}

func CreateOwner(gowner *github.User) *model.Owner {
	return &model.Owner{ID: gowner.GetID(), Name: gowner.GetLogin()}
}

func CreateLicence(glicence *github.License) *model.Licence {
	if glicence == nil {
		return nil
	}
	return &model.Licence{Key: glicence.GetKey(), Name: glicence.GetName()}
}
