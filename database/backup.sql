--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

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

--
-- Name: public; Type: SCHEMA; Schema: -; Owner: user_sa
--

-- *not* creating schema, since initdb creates it


ALTER SCHEMA public OWNER TO user_sa;

--
-- Name: uuid-ossp; Type: EXTENSION; Schema: -; Owner: -
--

CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;


--
-- Name: EXTENSION "uuid-ossp"; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION "uuid-ossp" IS 'generate universally unique identifiers (UUIDs)';


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: api_clients; Type: TABLE; Schema: public; Owner: user_sa
--

CREATE TABLE public.api_clients (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    name character varying(100) NOT NULL,
    api_key_hash character varying(255) NOT NULL,
    is_active boolean DEFAULT true,
    rate_limit integer DEFAULT 60,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(64) NOT NULL,
    updated_at timestamp with time zone DEFAULT now(),
    updated_by character varying(64) NOT NULL,
    client_id character varying(15)
);


ALTER TABLE public.api_clients OWNER TO user_sa;

--
-- Name: auth_logs; Type: TABLE; Schema: public; Owner: user_sa
--

CREATE TABLE public.auth_logs (
    id bigint NOT NULL,
    api_client_id uuid,
    request_ip inet,
    endpoint character varying(255),
    status character varying(16),
    message text,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(64)
);


ALTER TABLE public.auth_logs OWNER TO user_sa;

--
-- Name: auth_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: user_sa
--

CREATE SEQUENCE public.auth_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.auth_logs_id_seq OWNER TO user_sa;

--
-- Name: auth_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user_sa
--

ALTER SEQUENCE public.auth_logs_id_seq OWNED BY public.auth_logs.id;


--
-- Name: transaction_hist; Type: TABLE; Schema: public; Owner: user_sa
--

CREATE TABLE public.transaction_hist (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    transaction_id uuid NOT NULL,
    event_type character varying(64) NOT NULL,
    event_data jsonb,
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(64)
);


ALTER TABLE public.transaction_hist OWNER TO user_sa;

--
-- Name: transaction_logs; Type: TABLE; Schema: public; Owner: user_sa
--

CREATE TABLE public.transaction_logs (
    id bigint NOT NULL,
    transaction_id uuid NOT NULL,
    message text,
    source character varying(64),
    created_at timestamp with time zone DEFAULT now(),
    created_by character varying(64)
);


ALTER TABLE public.transaction_logs OWNER TO user_sa;

--
-- Name: transaction_logs_id_seq; Type: SEQUENCE; Schema: public; Owner: user_sa
--

CREATE SEQUENCE public.transaction_logs_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER SEQUENCE public.transaction_logs_id_seq OWNER TO user_sa;

--
-- Name: transaction_logs_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: user_sa
--

ALTER SEQUENCE public.transaction_logs_id_seq OWNED BY public.transaction_logs.id;


--
-- Name: transactions; Type: TABLE; Schema: public; Owner: user_sa
--

CREATE TABLE public.transactions (
    id uuid DEFAULT public.uuid_generate_v4() NOT NULL,
    number_billing character varying(20) NOT NULL,
    request_id character varying(32),
    customer_pan character varying(20),
    amount numeric(12,2),
    transaction_datetime timestamp(6) with time zone,
    retrieval_reference_number character varying(12),
    customer_name character varying(100),
    merchant_id character varying(15),
    merchant_name character varying(100),
    merchant_city character varying(100),
    currency_code character varying(5),
    payment_status character varying(10),
    payment_description character varying(100),
    created_at timestamp(6) with time zone DEFAULT now() NOT NULL,
    created_by character varying(64) NOT NULL
);


ALTER TABLE public.transactions OWNER TO user_sa;

--
-- Name: auth_logs id; Type: DEFAULT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.auth_logs ALTER COLUMN id SET DEFAULT nextval('public.auth_logs_id_seq'::regclass);


--
-- Name: transaction_logs id; Type: DEFAULT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.transaction_logs ALTER COLUMN id SET DEFAULT nextval('public.transaction_logs_id_seq'::regclass);


--
-- Data for Name: api_clients; Type: TABLE DATA; Schema: public; Owner: user_sa
--

COPY public.api_clients (id, name, api_key_hash, is_active, rate_limit, created_at, created_by, updated_at, updated_by, client_id) FROM stdin;
3f797680-664f-4a0f-aa75-347aa2866d51	PT Tokopersia	lDD1rfzuIH8wXDCUV9lFqChrrx4DPRlR3xiUo6EXG2XDOtxSJzWVo3jLsFSLT0Mls+mC77gDizdz2sYRFXfidA==	t	60	2025-08-10 12:31:53.276+00	admin	2025-08-10 12:31:57.107+00	admin	api-tokopersia
\.


--
-- Data for Name: auth_logs; Type: TABLE DATA; Schema: public; Owner: user_sa
--

COPY public.auth_logs (id, api_client_id, request_ip, endpoint, status, message, created_at, created_by) FROM stdin;
\.


--
-- Data for Name: transaction_hist; Type: TABLE DATA; Schema: public; Owner: user_sa
--

COPY public.transaction_hist (id, transaction_id, event_type, event_data, created_at, created_by) FROM stdin;
92d7bc99-e449-4497-b44a-a2b7d95f9a24	0cf73a13-584e-4369-9d4b-e073f621e8c2	INSERT	{"id": "0cf73a13-584e-4369-9d4b-e073f621e8c2", "amount": 10000, "created_by": "", "request_id": "XwVjF5zfuHhrDZuw", "merchant_id": "008800223497", "customer_pan": "9360001110000000019", "currency_code": "360", "customer_name": "John Doe", "merchant_city": "Jakarta Pusat", "merchant_name": "Sukses Makmur Bendungan Hilir", "number_billing": "20250809171042571927", "payment_status": "000", "transaction_date": "2021-02-25T13:36:13Z", "retrieval_ref_num": "123456789012", "payment_description": "Payment Success"}	2025-08-09 17:10:39.780537+00	
1685d890-36de-4e46-9313-e719e8474a61	6cee0dea-d1a8-43ac-81aa-9d8ead797958	INSERT	{"id": "6cee0dea-d1a8-43ac-81aa-9d8ead797958", "amount": 120000, "created_by": "", "request_id": "XwVjF5zfuHhrDZTT", "merchant_id": "888800223497", "customer_pan": "9360001110000000019", "currency_code": "360", "customer_name": "John Doe", "merchant_city": "Jakarta Selatan", "merchant_name": "Kaos Keren Anak Muda", "number_billing": "20250810062129758175", "payment_status": "000", "transaction_date": "2025-08-10T13:36:13Z", "retrieval_ref_num": "123456789000", "payment_description": "Payment Success"}	2025-08-10 06:21:28.159967+00	
415b06af-7f7c-4429-874a-59dbce519f39	591d9698-5110-4e96-b847-4c09847e539a	INSERT	{"id": "591d9698-5110-4e96-b847-4c09847e539a", "amount": 130000, "created_by": "3f797680-664f-4a0f-aa75-347aa2866d51", "request_id": "XwVjF5zfuHhrDZTT", "merchant_id": "888800223497", "customer_pan": "9360001110000123019", "currency_code": "360", "customer_name": "Rafi Ahmad", "merchant_city": "Jakarta Selatan", "merchant_name": "Kaos Keren Anak Muda", "number_billing": "20250810062354721797", "payment_status": "000", "transaction_date": "2025-08-10T13:36:13Z", "retrieval_ref_num": "123456789000", "payment_description": "Payment Success"}	2025-08-10 06:23:54.836712+00	3f797680-664f-4a0f-aa75-347aa2866d51
\.


--
-- Data for Name: transaction_logs; Type: TABLE DATA; Schema: public; Owner: user_sa
--

COPY public.transaction_logs (id, transaction_id, message, source, created_at, created_by) FROM stdin;
\.


--
-- Data for Name: transactions; Type: TABLE DATA; Schema: public; Owner: user_sa
--

COPY public.transactions (id, number_billing, request_id, customer_pan, amount, transaction_datetime, retrieval_reference_number, customer_name, merchant_id, merchant_name, merchant_city, currency_code, payment_status, payment_description, created_at, created_by) FROM stdin;
0cf73a13-584e-4369-9d4b-e073f621e8c2	20250809171042571927	XwVjF5zfuHhrDZuw	9360001110000000019	10000.00	2021-02-25 13:36:13+00	123456789012	John Doe	008800223497	Sukses Makmur Bendungan Hilir	Jakarta Pusat	360	000	Payment Success	2025-08-09 17:10:39.780537+00	
6cee0dea-d1a8-43ac-81aa-9d8ead797958	20250810062129758175	XwVjF5zfuHhrDZTT	9360001110000000019	120000.00	2025-08-10 13:36:13+00	123456789000	John Doe	888800223497	Kaos Keren Anak Muda	Jakarta Selatan	360	000	Payment Success	2025-08-10 06:21:28.159967+00	
591d9698-5110-4e96-b847-4c09847e539a	20250810062354721797	XwVjF5zfuHhrDZTT	9360001110000123019	130000.00	2025-08-10 13:36:13+00	123456789000	Rafi Ahmad	888800223497	Kaos Keren Anak Muda	Jakarta Selatan	360	000	Payment Success	2025-08-10 06:23:54.836712+00	3f797680-664f-4a0f-aa75-347aa2866d51
\.


--
-- Name: auth_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user_sa
--

SELECT pg_catalog.setval('public.auth_logs_id_seq', 1, false);


--
-- Name: transaction_logs_id_seq; Type: SEQUENCE SET; Schema: public; Owner: user_sa
--

SELECT pg_catalog.setval('public.transaction_logs_id_seq', 1, false);


--
-- Name: api_clients api_clients_pkey; Type: CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.api_clients
    ADD CONSTRAINT api_clients_pkey PRIMARY KEY (id);


--
-- Name: auth_logs auth_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.auth_logs
    ADD CONSTRAINT auth_logs_pkey PRIMARY KEY (id);


--
-- Name: transaction_hist transaction_events_pkey; Type: CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.transaction_hist
    ADD CONSTRAINT transaction_events_pkey PRIMARY KEY (id);


--
-- Name: transaction_logs transaction_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.transaction_logs
    ADD CONSTRAINT transaction_logs_pkey PRIMARY KEY (id);


--
-- Name: transactions transactions_number_billing_key; Type: CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_number_billing_key UNIQUE (number_billing);


--
-- Name: transactions transactions_pkey; Type: CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.transactions
    ADD CONSTRAINT transactions_pkey PRIMARY KEY (id);


--
-- Name: idx_auth_logs_client_id; Type: INDEX; Schema: public; Owner: user_sa
--

CREATE INDEX idx_auth_logs_client_id ON public.auth_logs USING btree (api_client_id);


--
-- Name: idx_auth_logs_ip; Type: INDEX; Schema: public; Owner: user_sa
--

CREATE INDEX idx_auth_logs_ip ON public.auth_logs USING btree (request_ip);


--
-- Name: idx_event_transaction_id; Type: INDEX; Schema: public; Owner: user_sa
--

CREATE INDEX idx_event_transaction_id ON public.transaction_hist USING btree (transaction_id);


--
-- Name: idx_log_transaction_id; Type: INDEX; Schema: public; Owner: user_sa
--

CREATE INDEX idx_log_transaction_id ON public.transaction_logs USING btree (transaction_id);


--
-- Name: auth_logs auth_logs_api_client_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.auth_logs
    ADD CONSTRAINT auth_logs_api_client_id_fkey FOREIGN KEY (api_client_id) REFERENCES public.api_clients(id);


--
-- Name: transaction_hist transaction_events_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.transaction_hist
    ADD CONSTRAINT transaction_events_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON DELETE CASCADE;


--
-- Name: transaction_logs transaction_logs_transaction_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: user_sa
--

ALTER TABLE ONLY public.transaction_logs
    ADD CONSTRAINT transaction_logs_transaction_id_fkey FOREIGN KEY (transaction_id) REFERENCES public.transactions(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

