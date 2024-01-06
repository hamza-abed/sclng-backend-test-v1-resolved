package migration

import (
	"database/sql"

	"github.com/Scalingo/sclng-backend-test-v1/queries"
	"github.com/Scalingo/sclng-backend-test-v1/util"
	"github.com/sirupsen/logrus"
)

func Migrate(db *sql.DB, logger logrus.FieldLogger, eraseExistingSchema bool) {
	if eraseExistingSchema {
		dropAllSchemas(db, logger)
	}
	createSchema(db, logger)
	createAllTables(db, logger)
}

func dropAllSchemas(db *sql.DB, logger logrus.FieldLogger) {
	logger.Info("removing all tables ...")
	util.RunCreateSQLStatement(queries.DropSchemaPublic, db, logger)
}

func createSchema(db *sql.DB, logger logrus.FieldLogger) {
	logger.Info("creating schema public ...")
	util.RunCreateSQLStatement(queries.CreateSchemaPublic, db, logger)
}

func createAllTables(db *sql.DB, logger logrus.FieldLogger) {
	logger.Info("creating DB tables ...")
	util.RunCreateSQLStatement(queries.CreateTableRepositoryOwner, db, logger)
	util.RunCreateSQLStatement(queries.CreateTableRepositoryLicence, db, logger)
	util.RunCreateSQLStatement(queries.CreateTableRepository, db, logger)
	util.RunCreateSQLStatement(queries.CreateTableLanguage, db, logger)
	util.RunCreateSQLStatement(queries.CreateTableRepositoryLanguage, db, logger)
}

//@todo: create indexes
