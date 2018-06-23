SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

DROP TABLE IF EXISTS tmp_gh_repositories;

CREATE TABLE tmp_gh_repositories (
    name character varying NOT NULL,
    primary_language character varying NOT NULL,
    clone_url character varying NOT NULL,
    clone_path character varying NOT NULL,
    vcs character varying NOT NULL,
    github_id bigint NOT NULL,
    full_name character varying,
    description character varying,
    homepage character varying,
    fork boolean,
    default_branch character varying,
    master_branch character varying,
    html_url character varying,
    forks_count integer,
    open_issues_count integer,
    stargazers_count integer,
    subscribers_count integer,
    watchers_count integer,
    size_in_kb integer,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    pushed_at timestamp with time zone
);
