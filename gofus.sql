--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

DROP DATABASE gofus;
--
-- Name: gofus; Type: DATABASE; Schema: -; Owner: postgres
--

CREATE DATABASE gofus WITH TEMPLATE = template0 ENCODING = 'UTF8' LC_COLLATE = 'en_US.UTF-8' LC_CTYPE = 'en_US.UTF-8';


ALTER DATABASE gofus OWNER TO postgres;

\connect gofus

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: postgres
--

CREATE SCHEMA public;


ALTER SCHEMA public OWNER TO postgres;

--
-- Name: SCHEMA public; Type: COMMENT; Schema: -; Owner: postgres
--

COMMENT ON SCHEMA public IS 'standard public schema';


--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = public, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: players; Type: TABLE; Schema: public; Owner: gofus; Tablespace: 
--

CREATE TABLE players (
    id integer NOT NULL,
    owner_id integer NOT NULL,
    name character varying(255) NOT NULL,
    skin integer NOT NULL,
    first_color bigint NOT NULL,
    second_color bigint NOT NULL,
    third_color bigint NOT NULL,
    level integer NOT NULL,
    experience bigint NOT NULL,
    current_map integer NOT NULL,
    current_cell integer NOT NULL,
    gender boolean NOT NULL,
    breed integer NOT NULL
);


ALTER TABLE public.players OWNER TO gofus;

--
-- Name: players_id_seq; Type: SEQUENCE; Schema: public; Owner: gofus
--

CREATE SEQUENCE players_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.players_id_seq OWNER TO gofus;

--
-- Name: players_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: gofus
--

ALTER SEQUENCE players_id_seq OWNED BY players.id;


--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: gofus
--

CREATE SEQUENCE users_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.users_id_seq OWNER TO gofus;

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
-- Name: id; Type: DEFAULT; Schema: public; Owner: gofus
--

ALTER TABLE ONLY players ALTER COLUMN id SET DEFAULT nextval('players_id_seq'::regclass);


--
-- Name: players_pkey; Type: CONSTRAINT; Schema: public; Owner: gofus; Tablespace: 
--

ALTER TABLE ONLY players
    ADD CONSTRAINT players_pkey PRIMARY KEY (id);


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
-- Name: players_uniq_name; Type: INDEX; Schema: public; Owner: gofus; Tablespace: 
--

CREATE UNIQUE INDEX players_uniq_name ON players USING btree (name);


--
-- Name: public; Type: ACL; Schema: -; Owner: postgres
--

REVOKE ALL ON SCHEMA public FROM PUBLIC;
REVOKE ALL ON SCHEMA public FROM postgres;
GRANT ALL ON SCHEMA public TO postgres;
GRANT ALL ON SCHEMA public TO PUBLIC;


--
-- PostgreSQL database dump complete
--

