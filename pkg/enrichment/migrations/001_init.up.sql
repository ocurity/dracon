-- issues table
CREATE TABLE IF NOT EXISTS issues (
  -- protobuf: enriched issue
  "hash" VARCHAR(32) NOT NULL,
  source VARCHAR(512) NOT NULL,
  first_seen TIMESTAMP WITH TIME ZONE,
  occurrences INTEGER NOT NULL,
  false_positive BOOLEAN NOT NULL,
  updated_at TIMESTAMP WITH TIME ZONE,
  -- protobuf: issue
  "target" VARCHAR(512) NOT NULL,
  "type" VARCHAR(128) NOT NULL,
  "title" VARCHAR(512) NOT NULL,
  severity INTEGER NOT NULL,
  cvss FLOAT NOT NULL,
  confidence INTEGER NOT NULL,
  "description" TEXT NOT NULL,
  PRIMARY KEY ("hash")
);
CREATE INDEX idx_issues_source ON issues(source);
