-- +goose Up
ALTER TABLE memberships
  ADD CONSTRAINT memberships_profile_plan_key UNIQUE (profile_id, plan_code);

ALTER TABLE course_access
  ADD CONSTRAINT course_access_profile_course_key UNIQUE (profile_id, course_id);

-- +goose Down
ALTER TABLE course_access
  DROP CONSTRAINT IF EXISTS course_access_profile_course_key;

ALTER TABLE memberships
  DROP CONSTRAINT IF EXISTS memberships_profile_plan_key;
