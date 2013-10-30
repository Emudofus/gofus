--
-- PostgreSQL database dump
--

SET statement_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;

SET search_path = public, pg_catalog;

DROP INDEX public.players_uniq_name;
ALTER TABLE ONLY public.players DROP CONSTRAINT players_pkey;
ALTER TABLE public.players ALTER COLUMN id DROP DEFAULT;
DROP SEQUENCE public.players_id_seq;
DROP TABLE public.players;
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
    breed integer NOT NULL,
    vitality smallint NOT NULL,
    wisdom smallint NOT NULL,
    strength smallint NOT NULL,
    intelligence smallint NOT NULL,
    chance smallint NOT NULL,
    agility smallint NOT NULL
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
-- Name: id; Type: DEFAULT; Schema: public; Owner: gofus
--

ALTER TABLE ONLY players ALTER COLUMN id SET DEFAULT nextval('players_id_seq'::regclass);


--
-- Name: players_pkey; Type: CONSTRAINT; Schema: public; Owner: gofus; Tablespace: 
--

ALTER TABLE ONLY players
    ADD CONSTRAINT players_pkey PRIMARY KEY (id);


--
-- Name: players_uniq_name; Type: INDEX; Schema: public; Owner: gofus; Tablespace: 
--

CREATE UNIQUE INDEX players_uniq_name ON players USING btree (name);


--
-- PostgreSQL database dump complete
--

