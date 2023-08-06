-- +goose Up
--
-- PostgreSQL database dump
--

-- Dumped from database version 12.7
-- Dumped by pg_dump version 13.4

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', 'public', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

SET default_tablespace = '';

SET default_table_access_method = heap;

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";


CREATE TABLE public.user (
    userid uuid DEFAULT uuid_generate_v4 (),
    firstname varchar,
    lastname varchar,
    email varchar,
    password_hash varchar,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(userid)
);

ALTER TABLE public.user OWNER TO postgres;

ALTER TABLE ONLY public.user
    ADD CONSTRAINT email_unique UNIQUE (email);


CREATE TABLE public.portfolio (
    portfolioid uuid DEFAULT uuid_generate_v4 (),
    userid uuid,
    name varchar,
    extid varchar,
    tax_countryid varchar,
    financial_year timestamp,
    performence_calculation int,
    summary boolean,
    price_alert boolean,
    company_event_alert boolean,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(portfolioid),
    FOREIGN KEY(userid) REFERENCES public.user(userid)
);

CREATE TABLE public.broker (
    brokerid uuid DEFAULT uuid_generate_v4 (),
    extid varchar,
    name varchar,
    description varchar,
    slug varchar,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(brokerid)
);

CREATE TABLE public.holding (
    holdingid uuid DEFAULT uuid_generate_v4 (),
    portfolioid uuid,
    brokerid uuid,
    extid varchar,
    name varchar,
    description varchar,
    slug varchar,
    trade_date timestamp,
    trade_type varchar,
    quantity int,
    share_price float,
    exchange_rate float,
    exchange_currencyid uuid,
    brokerage_unit_price float,
    brokerage_currency uuid,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(holdingid),
    FOREIGN KEY(portfolioid) REFERENCES public.portfolio(portfolioid),
    FOREIGN KEY(brokerid) REFERENCES public.broker(brokerid)
);

CREATE TABLE public.stock (
    stockid uuid DEFAULT uuid_generate_v4 (),
    holdingid uuid,
    ticker_symbol varchar,
    company_name varchar,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(stockid),
    FOREIGN KEY(holdingid) REFERENCES public.holding(holdingid)
);

CREATE TABLE public.currency (
    currencyid uuid DEFAULT uuid_generate_v4 (),
    code varchar,
    name varchar,
    symbol varchar,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(currencyid)
);

CREATE TABLE public.cryptocurrency (
    cryptocurrencyid uuid DEFAULT uuid_generate_v4 (),
    holdingid uuid,
    symbol varchar,
    name varchar,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(cryptocurrencyid),
    FOREIGN KEY(holdingid) REFERENCES public.holding(holdingid)
);

CREATE TABLE public.transaction (
    transactionid uuid DEFAULT uuid_generate_v4 (),
    holdingid uuid,
    stockid uuid,
    cryptocurrencyid uuid,
    transaction_type varchar,
    quantity int,
    transaction_price float,
    transaction_date timestamp,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(transactionid),
    FOREIGN KEY(holdingid) REFERENCES public.holding(holdingid),
    FOREIGN KEY(stockid) REFERENCES public.stock(stockid),
    FOREIGN KEY(cryptocurrencyid) REFERENCES public.cryptocurrency(cryptocurrencyid)
);

CREATE TABLE public.dividend (
    dividendid uuid DEFAULT uuid_generate_v4 (),
    stockid uuid,
    dividend_per_share float,
    dividend_date timestamp,
    created_at timestamp,
    updated_at timestamp,
    PRIMARY KEY(dividendid),
    FOREIGN KEY(stockid) REFERENCES public.stock(stockid)
);

CREATE TABLE public.exchange_rate (
    exchange_rate_id uuid DEFAULT uuid_generate_v4 () PRIMARY KEY,    -- unique identifier for this record
    currency_id uuid REFERENCES public.currency(currencyid),          -- currency being converted from USD
    rate float NOT NULL,                                              -- exchange rate between USD and the currency
    rate_date date NOT NULL,                                          -- date this exchange rate is for
    created_at timestamp,                                             -- when the record was created
    updated_at timestamp                                              -- when the record was last updated
);

COMMENT ON COLUMN public.exchange_rate.currency_id IS 'Currency being converted from USD';
COMMENT ON COLUMN public.exchange_rate.rate IS 'Exchange rate between USD and the currency';
COMMENT ON COLUMN public.exchange_rate.rate_date IS 'Date this exchange rate is for';
COMMENT ON COLUMN public.exchange_rate.created_at IS 'When the record was created';
COMMENT ON COLUMN public.exchange_rate.updated_at IS 'When the record was last updated';