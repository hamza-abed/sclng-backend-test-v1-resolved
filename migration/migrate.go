package migration

import (
	"database/sql"

	"github.com/Scalingo/sclng-backend-test-v1/queries"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	"github.com/sirupsen/logrus"
)

func Migrate(db *sql.DB, logger logrus.FieldLogger, eraseExistingSchema bool) error {
	var err error
	if eraseExistingSchema {
		err = dropAllSchemas(db, logger)
	}
	err = createSchema(db, logger)
	err = createAllTables(db, logger)
	return err
}

func dropAllSchemas(db *sql.DB, logger logrus.FieldLogger) error {
	logger.Info("removing all tables ...")
	return util.RunCreateSQLStatement(queries.DropSchemaPublic, db, logger)
}

func createSchema(db *sql.DB, logger logrus.FieldLogger) error {
	logger.Info("creating schema public ...")
	return util.RunCreateSQLStatement(queries.CreateSchemaPublic, db, logger)
}

func createAllTables(db *sql.DB, logger logrus.FieldLogger) error {
	var err error
	logger.Info("creating DB tables ...")
	err = util.RunCreateSQLStatement(queries.CreateTableRepositoryOwner, db, logger)
	err = util.RunCreateSQLStatement(queries.CreateTableRepositoryLicence, db, logger)
	err = util.RunCreateSQLStatement(queries.CreateTableRepository, db, logger)
	err = util.RunCreateSQLStatement(queries.CreateTableLanguage, db, logger)
	err = util.RunCreateSQLStatement(queries.CreateTableRepositoryLanguage, db, logger)
	return err
}

//@todo: create indexes
