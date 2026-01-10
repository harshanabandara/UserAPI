CREATE TYPE user_status AS ENUM ('ACTIVE', 'INACTIVE');
CREATE TABLE IF NOT EXISTS  users (
                                      user_id   UUID PRIMARY KEY DEFAULT gen_random_uuid(),

                                      first_name VARCHAR(50) NOT NULL,
                                      last_name  VARCHAR(50) NOT NULL,
                                      email      VARCHAR(256) NOT NULL,
                                      phone      VARCHAR(20),
                                      age        INTEGER,
                                      status     user_status DEFAULT 'ACTIVE',

                                      CONSTRAINT first_name_len CHECK (char_length(first_name) BETWEEN 2 AND 50),
                                      CONSTRAINT last_name_len  CHECK (char_length(last_name)  BETWEEN 2 AND 50),
                                      CONSTRAINT email_format  CHECK (email ~* '^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$'),
                                      CONSTRAINT phone_format  CHECK (phone IS NULL OR phone ~ '^(?:\+94|0)[0-9]{9}$'),
                                      CONSTRAINT age_positive  CHECK (age IS NULL OR age > 0)
);