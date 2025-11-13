INSERT INTO teams (team_name)
VALUES
  ('backend'),
  ('frontend'),
  ('ml');

INSERT INTO users (user_id, username, team_name, is_active)
VALUES
  ('u1', 'Alice', 'backend', TRUE),
  ('u2', 'Bob', 'backend', TRUE),
  ('u3', 'Charlie', 'backend', FALSE),
  ('u4', 'Diana', 'frontend', TRUE),
  ('u5', 'Eve', 'ml', TRUE),
  ('u6', 'Frank', 'ml', TRUE);

INSERT INTO pull_requests (pull_request_id, pull_request_name, author_id, status, need_more_reviewers)
VALUES
  ('pr-1001', 'Add login endpoint', 'u1', 'OPEN', TRUE),
  ('pr-1002', 'Refactor auth middleware', 'u4', 'OPEN', TRUE);

INSERT INTO pr_reviewers (pull_request_id, reviewer_id)
VALUES
  ('pr-1001', 'u2'),
  ('pr-1001', 'u3'),
  ('pr-1002', 'u5');
