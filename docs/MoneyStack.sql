CREATE TABLE "Categories" (
  "id" bigserial PRIMARY KEY,
  "category" varchar NOT NULL,
  "status" boolean DEFAULT true
);

CREATE TABLE "Expenses" (
  "id" bigserial PRIMARY KEY,
  "cat_id" int NOT NULL,
  "amount" bigint NOT NULL,
  "from_ac" int NOT NULL,
  "to_ac" int NOT NULL,
  "status" boolean DEFAULT true,
  "created_at" timestamptz DEFAULT (now())
);

CREATE TABLE "Account" (
  "id" bigserial PRIMARY KEY,
  "owner" varchar NOT NULL,
  "owe" bigint NOT NULL,
  "balance" bigint NOT NULL,
  "status" boolean DEFAULT true
);

CREATE TABLE "Transfers" (
  "id" bigserial PRIMARY KEY,
  "from_account_id" int NOT NULL,
  "to_account_id" int NOT NULL,
  "amount" bigint NOT NULL,
  "created_at" timestamptz DEFAULT (now())
);

CREATE INDEX ON "Categories" ("category");

CREATE INDEX ON "Expenses" ("from_ac");

CREATE INDEX ON "Expenses" ("to_ac");

CREATE INDEX ON "Account" ("owner");

CREATE INDEX ON "Transfers" ("from_account_id");

CREATE INDEX ON "Transfers" ("to_account_id");

COMMENT ON COLUMN "Account"."owe" IS 'It can be negative or postive';

COMMENT ON COLUMN "Transfers"."amount" IS 'It must be +ve';

ALTER TABLE "Expenses" ADD FOREIGN KEY ("cat_id") REFERENCES "Categories" ("id");

ALTER TABLE "Expenses" ADD FOREIGN KEY ("from_ac") REFERENCES "Account" ("id");

ALTER TABLE "Expenses" ADD FOREIGN KEY ("to_ac") REFERENCES "Account" ("id");

ALTER TABLE "Transfers" ADD FOREIGN KEY ("from_account_id") REFERENCES "Account" ("id");

ALTER TABLE "Transfers" ADD FOREIGN KEY ("to_account_id") REFERENCES "Account" ("id");
