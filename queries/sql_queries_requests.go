package queries

// --------------- VERIFICATION ----------------//
const IsGithubRepositoryExists = `SELECT EXISTS(SELECT id from repo where id=$1)`

const IsOwnerExists = `SELECT EXISTS(SELECT id from rowner where id=$1)`

const IsLicenceExists = `SELECT EXISTS(SELECT key from licence where key=$1)`

const IsLanguageExists = `SELECT EXISTS(SELECT name from lang where name=$1)`

// --------------- SELECTION -------------------//
const GetRepositoryByFilter = `-- name: GetRepositoryByFilter :many
SELECT repo.*, lang.*, rowner.*,licence.*
FROM repo
JOIN repo_lang ON repo_lang.repo_id = repo.id
JOIN lang ON repo_lang.langage_name=lang.name
JOIN rowner ON rowner.id=repo.rowner_id
JOIN licence ON licence.key=repo.licence_key
WHERE licence.name = $1`

// --------------- CREATION --------------------//
const CreateOwner = `INSERT INTO rowner(id, name) VALUES ($1, $2)`

const CreateLicence = `INSERT INTO licence(key, name) VALUES ($1, $2)`

const CreateLanguage = `INSERT INTO lang(name) VALUES ($1)`

const CreateRepository = `INSERT INTO public.repo(
	id, full_name, name, created_at, rowner_id, licence_key)
	VALUES ($1, $2, $3, $4, $5,$6)`

const CreateRepositoryLanguage = `INSERT INTO repo_lang(
	repo_id, langage_name,bytes)
	VALUES ($1, $2,$3)`
