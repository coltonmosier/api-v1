-- DEVICETYPE QUERIES
-- name: GetDeviceTypesActive :many
SELECT id, name FROM device_type
ORDER BY id;

-- name: GetDeviceTypeByName :one
SELECT id, name, status FROM device_type
WHERE name = ?
ORDER BY id;

-- name: GetDeviceTypeById :one
SELECT id, name, status FROM device_type
WHERE id = ?
ORDER BY id;

-- name: CreateDeviceType :exec
INSERT INTO device_type (name) VALUES (?);

-- name: UpdateDeviceType :exec
UPDATE device_type SET name = ?
WHERE id = ?;

-- name: UpdateDeviceTypeStatus :exec
UPDATE device_type SET status = ?
WHERE id = ?;

-- name: DeleteDeviceType :exec
DELETE FROM device_type
WHERE id = ?;




-- MANUFACTURER QUERIES
-- name: GetManufacturersActive :many
SELECT id, name FROM manufacturer
ORDER BY id;

-- name: GetManufacturerByName :one
SELECT id, name FROM manufacturer
WHERE name = ?
ORDER BY id;

-- name: GetManufacturerById :one
SELECT id, name FROM manufacturer
WHERE id = ?
ORDER BY id;

-- name: CreateManufacturer :exec
INSERT INTO manufacturer (name) VALUES (?);

-- name: UpdateManufacturer :exec
UPDATE manufacturer SET name = ?
WHERE id = ?;

-- name: UpdateManufacturerStatus :exec
UPDATE manufacturer SET status = ?
WHERE id = ?;

-- name: DeleteManufacturer :exec
DELETE FROM manufacturer
WHERE id = ?;




-- SERIALNUMBER QUERIES
-- name: GetSerialNumbers :many
SELECT serial_number FROM serial_numbers;

-- name: GetSerialNumberBySerialNumber :one
SELECT serial_number FROM serial_numbers
WHERE serial_number = ?;

-- name: GetSerialNumberLikeSerialNumber :many
SELECT serial_number FROM serial_numbers
WHERE serial_number LIKE ?;

-- name: UpdateSerialNumber :exec
UPDATE serial_numbers SET serial_number = ?
WHERE serial_number = ?;




-- EQUIPMENT QUERIES
-- name: GetAllEquipment :many
SELECT * FROM serial_numbers
LIMIT ? OFFSET ?;

-- name: GetEquipmentByDeviceType :many
SELECT * FROM serial_numbers
WHERE device_type_id = ?
LIMIT ? OFFSET ?;

-- name: GetEquipmentByManufacturer :many
SELECT * FROM serial_numbers
WHERE manufacturer_id = ?
LIMIT ? OFFSET ?;

-- name: GetEquipmentBySerialNumber :one
SELECT * FROM serial_numbers
WHERE serial_number = ?
LIMIT ? OFFSET ?;

-- name: GetEquipmentLikeSerialNumber :many
SELECT * FROM serial_numbers
WHERE serial_number LIKE ?
LIMIT ? OFFSET ?;

-- name: GetEquipmentByDeviceTypeAndManufacturer :many
SELECT * FROM serial_numbers
WHERE device_type_id = ? AND manufacturer_id = ?
LIMIT ? OFFSET ?;

-- name: GetEquipmentByDeviceTypeAndSerialNumber :many
SELECT * FROM serial_numbers
WHERE device_type_id = ? AND serial_number = ?
LIMIT ? OFFSET ?;

-- name: GetEquipmentByManufacturerAndSerialNumber :many
SELECT * FROM serial_numbers
WHERE manufacturer_id = ? AND serial_number = ?
LIMIT ? OFFSET ?;

-- name: GetEquipmentByDeviceTypeManufacturerAndSerialNumber :many
SELECT * FROM serial_numbers
WHERE device_type_id = ? AND manufacturer_id = ? AND serial_number = ?
LIMIT ? OFFSET ?;

