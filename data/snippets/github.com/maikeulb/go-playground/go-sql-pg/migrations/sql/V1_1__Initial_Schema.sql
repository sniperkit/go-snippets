CREATE TABLE states (
    id serial PRIMARY KEY,
    name varchar(50) NULL
);

CREATE TABLE parks (
    id serial PRIMARY KEY,
    name varchar(50) NOT NULL,
    description TEXT NOT NULL,
    nearest_city varchar(50) NOT NULL,
    visitors integer NOT NULL,
    established timestamp NOT NULL,
    state_id integer NULL
);

ALTER TABLE parks
ADD CONSTRAINT fk_states_parks
FOREIGN KEY (state_id)
REFERENCES states(id)
ON DELETE CASCADE;
