--
-- PostgreSQL database dump
--

-- Dumped from database version 17.4
-- Dumped by pg_dump version 17.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: payments; Type: TABLE; Schema: public; Owner: payments_user
--

CREATE TABLE public.payments (
    id integer NOT NULL,
    user_id integer NOT NULL,
    amount numeric(10,2) NOT NULL,
    payment_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    status character varying(50) DEFAULT 'pending'::character varying NOT NULL
);


ALTER TABLE public.payments OWNER TO payments_user;

--
-- Name: payments_id_seq; Type: SEQUENCE; Schema: public; Owner: payments_user
--

CREATE SEQUENCE public.payments_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.payments_id_seq OWNER TO payments_user;

--
-- Name: payments_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: payments_user
--

ALTER SEQUENCE public.payments_id_seq OWNED BY public.payments.id;


--
-- Name: users; Type: TABLE; Schema: public; Owner: payments_user
--

CREATE TABLE public.users (
    id integer NOT NULL,
    username character varying(100) NOT NULL,
    email character varying(150),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP
);


ALTER TABLE public.users OWNER TO payments_user;

--
-- Name: users_id_seq; Type: SEQUENCE; Schema: public; Owner: payments_user
--

CREATE SEQUENCE public.users_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.users_id_seq OWNER TO payments_user;

--
-- Name: users_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: payments_user
--

ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;


--
-- Name: payments id; Type: DEFAULT; Schema: public; Owner: payments_user
--

ALTER TABLE ONLY public.payments ALTER COLUMN id SET DEFAULT nextval('public.payments_id_seq'::regclass);


--
-- Name: users id; Type: DEFAULT; Schema: public; Owner: payments_user
--

ALTER TABLE ONLY public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq'::regclass);


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: payments_user
--

COPY public.payments (id, user_id, amount, payment_date, status) FROM stdin;
1	1	100.50	2025-11-30 15:38:01.670739	completed
2	1	75.25	2025-11-30 15:38:01.670739	pending
3	2	200.00	2025-11-30 15:38:01.670739	completed
4	2	50.00	2025-11-30 15:38:01.670739	failed
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: payments_user
--

COPY public.users (id, username, email, created_at) FROM stdin;
1	alice	alice@example.com	2025-11-30 15:37:53.277048
2	bob	bob@example.com	2025-11-30 15:37:53.277048
\.


--
-- Name: payments_id_seq; Type: SEQUENCE SET; Schema: public; Owner: payments_user
--

SELECT pg_catalog.setval('public.payments_id_seq', 4, true);


--
-- Name: users_id_seq; Type: SEQUENCE SET; Schema: public; Owner: payments_user
--

SELECT pg_catalog.setval('public.users_id_seq', 2, true);


--
-- Name: payments payments_pkey; Type: CONSTRAINT; Schema: public; Owner: payments_user
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_pkey PRIMARY KEY (id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: payments_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: users users_username_key; Type: CONSTRAINT; Schema: public; Owner: payments_user
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_username_key UNIQUE (username);


--
-- Name: payments payments_user_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: payments_user
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payments_user_id_fkey FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: SCHEMA public; Type: ACL; Schema: -; Owner: pg_database_owner
--

GRANT USAGE ON SCHEMA public TO payments_user;


--
-- PostgreSQL database dump complete
--

