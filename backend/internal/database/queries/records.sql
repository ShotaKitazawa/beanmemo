-- name: ListRecords :many
SELECT * FROM records WHERE user_id = ? ORDER BY created_at DESC;

-- name: ListRecordsByOrigin :many
SELECT * FROM records WHERE user_id = ? AND origin = ? ORDER BY created_at DESC;

-- name: ListRecordsByRoastLevel :many
SELECT * FROM records WHERE user_id = ? AND roast_level = ? ORDER BY created_at DESC;

-- name: ListRecordsByRatingMin :many
SELECT * FROM records WHERE user_id = ? AND rating >= ? ORDER BY created_at DESC;

-- name: ListRecordsByBrewMethod :many
SELECT * FROM records WHERE user_id = ? AND brew_method = ? ORDER BY created_at DESC;

-- name: GetRecord :one
SELECT * FROM records WHERE id = ? AND user_id = ? LIMIT 1;

-- name: GetRelatedRecords :many
SELECT * FROM records WHERE user_id = ? AND id != ? AND name = ? ORDER BY created_at DESC LIMIT 10;

-- name: CreateRecord :execlastid
INSERT INTO records (
  user_id, name, rating, origin, roast_level, shop, price,
  purchased_at, tasting_note, brew_method, recipe, is_note_filled
) VALUES (
  ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
);

-- name: UpdateRecord :exec
UPDATE records SET
  name          = ?,
  rating        = ?,
  origin        = ?,
  roast_level   = ?,
  shop          = ?,
  price         = ?,
  purchased_at  = ?,
  tasting_note  = ?,
  brew_method   = ?,
  recipe        = ?,
  is_note_filled = ?,
  updated_at    = CURRENT_TIMESTAMP
WHERE id = ? AND user_id = ?;

-- name: DeleteRecord :exec
DELETE FROM records WHERE id = ? AND user_id = ?;

-- name: CountRecords :one
SELECT COUNT(*) FROM records WHERE user_id = ?;

-- name: ListAllTastingNotes :many
SELECT tasting_note FROM records WHERE user_id = ? AND tasting_note IS NOT NULL AND tasting_note != '';
