CREATE TABLE "pastes" (
    "id" TEXT NOT NULL,
    "content" TEXT NOT NULL,
    "modification_token" TEXT NOT NULL,
    "created" BIGINT NOT NULL,
    "metadata" TEXT NOT NULL,
    PRIMARY KEY ("id")
);