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

--
-- Data for Name: features; Type: TABLE DATA; Schema: public; Owner: devmine
--

COPY features (id, name, category, default_weight) FROM stdin;
1	average_forks	other	1
2	average_stars	other	1
3	commits_count	other	1
4	contributions_count	other	1
5	followers_count	other	1
6	hireable	other	1
\.


--
-- Name: features_id_seq; Type: SEQUENCE SET; Schema: public; Owner: devmine
--

SELECT pg_catalog.setval('features_id_seq', 6, true);


--
-- PostgreSQL database dump complete
--

