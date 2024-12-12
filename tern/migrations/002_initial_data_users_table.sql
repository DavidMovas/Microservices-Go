-- Write your migrate up statements here

INSERT INTO users("email","first_name","last_name","password_hash","user_active")
VALUES
    ('admin@gmail.com','Admin','User','$2a$12$1zGLuYDDNvATh4RA4avbKuheAMpb1svexSzrQm7up.bnpwQHs0jNe',1);

---- create above / drop below ----

-- Write your migrate down statements here. If this migration is irreversible
-- Then delete the separator line above.
