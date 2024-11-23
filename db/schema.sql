CREATE TABLE IF NOT EXISTS "solves" (
	"userid" INTEGER NOT NULL,
	"chalid" INTEGER NOT NULL,
	timestamp DATETIME,
	PRIMARY KEY ("userid", "chalid"),
	FOREIGN KEY("userid") REFERENCES "users"("id"),
	FOREIGN KEY("chalid") REFERENCES "challenges"("id")
);

CREATE TABLE IF NOT EXISTS "notifications" (
	"id" INTEGER NOT NULL,
	"user" INTEGER,
	"text" TEXT,
	PRIMARY KEY ("id"),
	FOREIGN KEY("user") REFERENCES "users"("id")
);

CREATE TABLE IF NOT EXISTS "users" (
	"id" INTEGER NOT NULL,
	"username" VARCHAR(32) UNIQUE NOT NULL,
	"email" TEXT UNIQUE NOT NULL,
	"salt" VARCHAR(16) NOT NULL,
	"password" VARCHAR(64) NOT NULL,
	"apikey" VARCHAR(256) NOT NULL,
	"score" INTEGER NOT NULL,
	"is_admin" BOOLEAN NOT NULL CHECK("is_admin" IN (0, 1)),
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "submissions" (
	"id" INTEGER NOT NULL,
	"userid" INTEGER,
	"chalid" INTEGER,
	"status" VARCHAR(1),
	"flag" TEXT,
	"timestamp" DATETIME,
	FOREIGN KEY("chalid") REFERENCES "challenges"("id"),
	FOREIGN KEY("userid") REFERENCES "users"("id"),
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "challenges" (
	"id" INTEGER NOT NULL,
	"name" TEXT UNIQUE NOT NULL,
	"description" TEXT NOT NULL,
	"difficulty" TEXT NOT NULL,
	"points" INTEGER NOT NULL,
	"max_points" INTEGER NOT NULL,
	"solves" INTEGER NOT NULL,
	"host" TEXT,
	"port" TEXT,
	"category" TEXT NOT NULL,
	"files" TEXT,
	"flag" TEXT UNIQUE NOT NULL,
	"hint1" TEXT,
	"hint2" TEXT,
	"hidden" BOOLEAN NOT NULL DEFAULT 1 CHECK("hidden" IN (0, 1)),
	"is_extra" BOOLEAN NOT NULL DEFAULT 0,
	PRIMARY KEY("id")
);

CREATE TABLE IF NOT EXISTS "config" (
	"key" TEXT NOT NULL,
	"type" TEXT NOT NULL DEFAULT 'text',
	"value" INTEGER,
	"desc" TEXT,
	PRIMARY KEY("key")
);