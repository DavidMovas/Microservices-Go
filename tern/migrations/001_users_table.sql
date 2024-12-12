-- Write your migrate up statements here

CREATE SEQUENCE user_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

SET default_tablespace = '';
SET default_table_access_method = heap;

CREATE TABLE users (
    id integer DEFAULT nextval('user_id_seq'::regclass) NOT NULL,
    email character varying(255),
    first_name character varying(255),
    last_name character varying(255),
    password_hash character varying(60),
    user_active integer DEFAULT 0,
    created_at timestamp without time zone default now(),
    updated_at timestamp without time zone,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

SELECT setval('user_id_seq', 1, true);

---- create above / drop below ----

DROP TABLE users;
DROP SEQUENCE user_id_seq;

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
