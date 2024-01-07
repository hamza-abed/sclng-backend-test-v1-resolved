package queries

const DropSchemaPublic = `Drop schema if exists public CASCADE`

const CreateSchemaPublic = `create schema IF NOT EXISTS public`

const CreateTableRepositoryOwner = `CREATE TABLE IF NOT EXISTS rowner(
    id bigint NOT NULL,
    name VARCHAR(255),
    CONSTRAINT rowner_pkey PRIMARY KEY (id)
)`

const CreateTableRepositoryLicence = `CREATE TABLE IF NOT EXISTS licence(
	key VARCHAR(255) NOT NULL,
	name VARCHAR(255),
	CONSTRAINT licence_pkey PRIMARY KEY (key)
)`

const CreateTableRepository = `CREATE TABLE IF NOT EXISTS repo (
    id bigint NOT NULL,
    full_name VARCHAR(255),
    name VARCHAR(255),
    created_at TIMESTAMP,
    rowner_id bigint NOT NULL,
    licence_key VARCHAR(255),
    CONSTRAINT repo_pkey PRIMARY KEY (id),
    CONSTRAINT fk_rowner_id_in_gr_from_o FOREIGN KEY (rowner_id)
			REFERENCES rowner (id) MATCH SIMPLE,
    CONSTRAINT fk_licence_id_in_gr_from_licence FOREIGN KEY (licence_key)
			REFERENCES licence (key) MATCH SIMPLE
)`

const CreateTableLanguage = `CREATE TABLE IF NOT EXISTS lang(
	name VARCHAR(255) NOT NULL,
	CONSTRAINT lang_pkey PRIMARY KEY (name)
)`

const CreateTableRepositoryLanguage = `CREATE TABLE IF NOT EXISTS repo_lang(
    repo_id bigint NOT NULL,
    langage_name VARCHAR(255) NOT NULL,
    bytes bigint,
    CONSTRAINT repo_lang_pkey PRIMARY KEY (repo_id,langage_name),
    CONSTRAINT fk_repo_id_in_rl_from_r FOREIGN KEY (repo_id)
			REFERENCES repo (id) MATCH SIMPLE,
    CONSTRAINT fk_lang_id_in_rl_from_l FOREIGN KEY (langage_name)
			REFERENCES lang (name) MATCH SIMPLE
)`
