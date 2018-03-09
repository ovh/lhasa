-- +migrate Down

-- Remove tables
DROP TABLE IF EXISTS "applications";
DROP TABLE IF EXISTS "person_infos";
DROP TABLE IF EXISTS "dependency_infos";

-- Remove indexes
DROP INDEX IF EXISTS "applications_pkey";
DROP INDEX IF EXISTS "idx_unique_fullname";
DROP INDEX IF EXISTS "person_infos_pkey";
DROP INDEX IF EXISTS "dependency_infos_pkey";

-- +migrate Up

-- Tables
CREATE TABLE "applications" (
    "id" integer DEFAULT nextval('dependency_infos_id_seq'::regclass) NOT NULL primary key,
    "domain" text  NULL,                                                    
    "name" text  NULL,                                                      
    "type" text  NULL,                                                      
    "language" text  NULL,                                                  
    "repository_url" text  NULL,                                            
    "avatar_url" text  NULL,                                                
    "description" text  NULL,                                               
    "version" text  NULL,                                                   
    "artefact" text  NULL,                                                  
    "repository" text  NULL,                                                
    "package" text  NULL,                                                   
    "email" text  NULL,                                                     
    "issues" text  NULL,                                                    
    "docs" text  NULL,                                                      
    "cisco_conf" text  NULL,                                                
    "produces" text  NULL,                                                  
    "consumes" text  NULL,                                                  
    "vault_aliases" text  NULL,                                             
    "tags" text  NULL,                                                      
    "created_at" timestamp with time zone  NULL,                            
    "updated_at" timestamp with time zone  NULL,                            
    "deleted_at" timestamp with time zone  NULL
);

CREATE TABLE "person_infos" (
    "id" integer DEFAULT nextval('dependency_infos_id_seq'::regclass) NOT NULL primary key,
    "name" text  NULL,
    "country" text  NULL,
    "email" text  NULL,
    "role" text  NULL);

CREATE TABLE "dependency_infos" (
    "id" integer DEFAULT nextval('dependency_infos_id_seq'::regclass) NOT NULL primary key,
    "type" text  NULL,
    "name" text  NULL,
    "version" text  NULL,
    "critical" boolean  NULL);

-- Index
CREATE UNIQUE INDEX applications_pkey ON applications USING btree (id);
CREATE UNIQUE INDEX idx_unique_fullname ON applications USING btree (domain, name, version)
CREATE UNIQUE INDEX person_infos_pkey ON person_infos USING btree (id);
CREATE UNIQUE INDEX dependency_infos_pkey ON dependency_infos USING btree (id);
