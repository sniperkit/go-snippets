--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

SET default_with_oids = false;

--
-- Name: commit_diff_deltas; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE commit_diff_deltas (
    id bigint NOT NULL,
    commit_id bigint NOT NULL,
    file_status character varying NOT NULL,
    is_file_binary boolean,
    similarity integer,
    old_file_path character varying NOT NULL,
    new_file_path character varying NOT NULL
);


--
-- Name: commit_diff_deltas_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE commit_diff_deltas_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: commit_diff_deltas_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE commit_diff_deltas_id_seq OWNED BY commit_diff_deltas.id;


--
-- Name: commits; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE commits (
    id bigint NOT NULL,
    repository_id bigint NOT NULL,
    author_id bigint,
    committer_id bigint,
    vcs_id character varying NOT NULL,
    message text,
    author_date timestamp with time zone,
    commit_date timestamp with time zone,
    file_changed_count integer,
    insertions_count integer,
    deletions_count integer
);


--
-- Name: commits_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE commits_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: commits_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE commits_id_seq OWNED BY commits.id;


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY commit_diff_deltas ALTER COLUMN id SET DEFAULT nextval('commit_diff_deltas_id_seq'::regclass);


--
-- Name: id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY commits ALTER COLUMN id SET DEFAULT nextval('commits_id_seq'::regclass);


--
-- Name: commit_diff_deltas_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY commit_diff_deltas
    ADD CONSTRAINT commit_diff_deltas_pk PRIMARY KEY (id);


--
-- Name: commits_pk; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY commits
    ADD CONSTRAINT commits_pk PRIMARY KEY (id);


--
-- Name: fki_commit_diff_deltas_fk_commits; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_commit_diff_deltas_fk_commits ON commit_diff_deltas USING btree (commit_id);


--
-- Name: fki_commits_fk_repositories; Type: INDEX; Schema: public; Owner: -
--

CREATE INDEX fki_commits_fk_repositories ON commits USING btree (repository_id);


--
-- Name: commit_diff_deltas_fk_commits; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY commit_diff_deltas
    ADD CONSTRAINT commit_diff_deltas_fk_commits FOREIGN KEY (commit_id) REFERENCES commits(id);


--
-- Name: commits_fk_repositories; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY commits
    ADD CONSTRAINT commits_fk_repositories FOREIGN KEY (repository_id) REFERENCES repositories(id);


--
-- PostgreSQL database dump complete
--

