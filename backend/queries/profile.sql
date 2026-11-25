-- name: GetTalentProfile :one
SELECT * FROM talent_profiles WHERE user_id = $1;

-- name: UpsertTalentProfile :one
INSERT INTO talent_profiles (user_id, talent_type, display_name, city, age, height_cm, bio, rate_min, rate_max, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
ON CONFLICT (user_id) DO UPDATE SET
    talent_type = EXCLUDED.talent_type,
    display_name = EXCLUDED.display_name,
    city = EXCLUDED.city,
    age = EXCLUDED.age,
    height_cm = EXCLUDED.height_cm,
    bio = EXCLUDED.bio,
    rate_min = EXCLUDED.rate_min,
    rate_max = EXCLUDED.rate_max,
    updated_at = NOW()
RETURNING *;

-- name: GetClientProfile :one
SELECT * FROM client_profiles WHERE user_id = $1;

-- name: UpsertClientProfile :one
INSERT INTO client_profiles (user_id, display_name, org_name)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE SET
    display_name = EXCLUDED.display_name,
    org_name = EXCLUDED.org_name
RETURNING *;

-- name: CreateProfilePhoto :one
INSERT INTO profile_photos (talent_id, url, position)
VALUES ($1, $2, $3)
RETURNING *;

-- name: CountProfilePhotos :one
SELECT COUNT(*) FROM profile_photos WHERE talent_id = $1;

-- name: ListProfilePhotos :many
SELECT * FROM profile_photos WHERE talent_id = $1 ORDER BY position;
