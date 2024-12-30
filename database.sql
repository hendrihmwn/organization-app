-- PostgreSQL as the database.

CREATE TABLE IF NOT EXISTS organizations (
                                             id  bigserial,
                                             parent_id  bigint,
                                             level integer,
                                             name       VARCHAR(100),
                                             created_at timestamp default NOW() not null,
                                             updated_at timestamp default NOW() not null,
                                             deleted_at timestamp,
                                       PRIMARY KEY (id)
);

CREATE INDEX IF NOT EXISTS parent_id ON organizations(parent_id);
