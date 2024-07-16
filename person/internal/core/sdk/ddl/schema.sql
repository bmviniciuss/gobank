CREATE SCHEMA IF NOT EXISTS "person";

CREATE TABLE IF NOT EXISTS "person"."person" (
  "id" SERIAL PRIMARY KEY,
  "uuid" uuid UNIQUE NOT NULL,
  "name" VARCHAR(255) NOT NULL,
  "document" VARCHAR(255) UNIQUE NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT now(),
  "updated_at" timestamptz NOT NULL DEFAULT now(),
  "active" bool NOT NULL DEFAULT true
);

CREATE INDEX IF NOT EXISTS idx_person_person_uuid  ON "person"."person" ("uuid");
CREATE INDEX IF NOT EXISTS idx_person_person_document ON "person"."person" ("document");
