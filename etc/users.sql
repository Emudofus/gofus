--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

ALTER TABLE ONLY public.users DROP CONSTRAINT users_pkey;
ALTER TABLE ONLY public.users DROP CONSTRAINT users_nickname_key;
ALTER TABLE ONLY public.users DROP CONSTRAINT users_name_key;
DROP TABLE public.users;
SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: users; Type: TABLE; Schema: public; Owner: gofus; Tablespace: 
--

CREATE TABLE users (
    id integer DEFAULT nextval('users_id_seq'::regclass) NOT NULL,
    name character varying(255),
    password character varying(255),
    nickname character varying(255),
    secret_question character varying(255),
    secret_answer character varying(255),
    rights numeric(20,0),
    community_id integer,
    subscription_end timestamp(6) without time zone,
    current_realm_server integer
);


ALTER TABLE public.users OWNER TO gofus;

--
-- Name: COLUMN users.rights; Type: COMMENT; Schema: public; Owner: gofus
--

COMMENT ON COLUMN users.rights IS 'UNSIGNED BIGINT';


--
-- Name: users_name_key; Type: CONSTRAINT; Schema: public; Owner: gofus; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_name_key UNIQUE (name);


--
-- Name: users_nickname_key; Type: CONSTRAINT; Schema: public; Owner: gofus; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_nickname_key UNIQUE (nickname);


--
-- Name: users_pkey; Type: CONSTRAINT; Schema: public; Owner: gofus; Tablespace: 
--

ALTER TABLE ONLY users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- PostgreSQL database dump complete
--

