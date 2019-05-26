-- Table: public.pictures

-- DROP TABLE public.pictures;

CREATE TABLE public.pictures
(
    id integer NOT NULL DEFAULT nextval('pictures_id_seq'::regclass),
    category character varying COLLATE pg_catalog."default" NOT NULL,
    url character varying COLLATE pg_catalog."default" NOT NULL,
    analysis json,
    last_seen date,
    CONSTRAINT pictures_pkey PRIMARY KEY (id)
)
WITH (
    OIDS = FALSE
)
TABLESPACE pg_default;

ALTER TABLE public.pictures
    OWNER to doadmin;

-- Index: uniq_url

-- DROP INDEX public.uniq_url;

CREATE UNIQUE INDEX uniq_url
    ON public.pictures USING btree
    (url COLLATE pg_catalog."default" varchar_pattern_ops)
    TABLESPACE pg_default;