CREATE TABLE IF NOT EXISTS entries(
   id SERIAL PRIMARY KEY,
   text VARCHAR,
   creator_id VARCHAR (50) NOT NULL,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP,
   is_deleted BOOLEAN DEFAULT false
);

CREATE INDEX entries_creator_id_id_idx ON entries (creator_id, id);
