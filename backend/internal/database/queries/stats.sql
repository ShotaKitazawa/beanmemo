-- name: StatsByOrigin :many
SELECT
  origin AS label,
  COUNT(*) AS count,
  AVG(rating) AS avg_rating
FROM records
WHERE user_id = ? AND origin IS NOT NULL AND origin != ''
GROUP BY origin
ORDER BY avg_rating DESC;

-- name: StatsByRoastLevel :many
SELECT
  roast_level AS label,
  COUNT(*) AS count,
  AVG(rating) AS avg_rating
FROM records
WHERE user_id = ? AND roast_level IS NOT NULL AND roast_level != ''
GROUP BY roast_level
ORDER BY avg_rating DESC;

-- name: StatsByBrewMethod :many
SELECT
  brew_method AS label,
  COUNT(*) AS count,
  AVG(rating) AS avg_rating
FROM records
WHERE user_id = ? AND brew_method IS NOT NULL AND brew_method != ''
GROUP BY brew_method
ORDER BY avg_rating DESC;

-- name: AvgRatingByOrigin :one
SELECT AVG(rating) FROM records WHERE user_id = ? AND origin = ?;

-- name: AvgRatingByName :one
SELECT AVG(rating) FROM records WHERE user_id = ? AND name = ?;
