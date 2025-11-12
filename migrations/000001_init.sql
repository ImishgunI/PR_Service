CREATE TABLE IF NOT EXISTS teams (
  team_name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
  user_id TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  team_name TEXT NOT NULL REFERENCES teams(team_name) ON DELETE CASCADE,
  is_active BOOLEAN NOT NULL DEFAULT TRUE
);

DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'pr_status') THEN
    CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
  END IF;
END$$;

CREATE TABLE IF NOT EXISTS pull_requests (
  pull_request_id TEXT PRIMARY KEY,
  pull_request_name TEXT NOT NULL,
  author_id TEXT NOT NULL REFERENCES users(user_id),
  status pr_status NOT NULL DEFAULT 'OPEN',
  need_more_reviewers BOOLEAN NOT NULL DEFAULT true,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT now(),
  merged_at TIMESTAMP WITH TIME ZONE NULL
);

CREATE TABLE IF NOT EXISTS pr_reviewers (
  pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
  reviewer_id TEXT NOT NULL REFERENCES users(user_id),
  PRIMARY KEY (pull_request_id, reviewer_id)
);

CREATE INDEX IF NOT EXISTS idx_users_team_active ON users(team_name, is_active);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_reviewer ON pr_reviewers(reviewer_id);
