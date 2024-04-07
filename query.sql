-- DEVICETYPE QUERIES
-- name: GetDeviceTypesActive :many
SELECT id, name FROM device_type
WHERE status = 'active'
ORDER BY id;

-- name: GetDeviceTypeByName :one
SELECT id, name, status FROM device_type
WHERE name = $1
ORDER BY id;

-- name: GetDeviceTypeById :one
SELECT id, name, status FROM device_type
WHERE id = $1
ORDER BY id;

-- name: CreateDeviceType :exec
INSERT INTO device_type (name) VALUES ($1);

-- name: UpdateDeviceType :exec
UPDATE device_type SET name = $1
WHERE id = $2;

-- name: UpdateDeviceTypeStatus :exec
UPDATE device_type SET status = $1
WHERE id = $2;

-- name: DeleteDeviceType :exec
DELETE FROM device_type
WHERE id = $1;




-- MANUFACTURER QUERIES
-- name: GetManufacturersActive :many
SELECT id, name FROM manufacturer
WHERE status = 'active'
ORDER BY id;

-- name: GetManufacturerByName :one
SELECT id, name FROM manufacturer
WHERE name = $1
ORDER BY id;

-- name: GetManufacturerById :one
SELECT id, name FROM manufacturer
WHERE id = $1
ORDER BY id;

-- name: CreateManufacturer :exec
INSERT INTO manufacturer (name) VALUES ($1);

-- name: UpdateManufacturer :exec
UPDATE manufacturer SET name = $1
WHERE id = $2;

-- name: UpdateManufacturerStatus :exec
UPDATE manufacturer SET status = $1
WHERE id = $2;

-- name: DeleteManufacturer :exec
DELETE FROM manufacturer
WHERE id = $1;




-- SERIALNUMBER QUERIES
-- name: GetSerialNumbers :many
SELECT serial_number FROM serial_numbers;

-- name: GetSerialNumberBySerialNumber :one
SELECT serial_number FROM serial_numbers
WHERE serial_number = $1;

-- name: GetSerialNumberLikeSerialNumber :many
SELECT serial_number FROM serial_numbers
WHERE serial_number LIKE $1;

-- name: UpdateSerialNumber :exec
UPDATE serial_numbers SET serial_number = $1
WHERE serial_number = $2;




-- EQUIPMENT QUERIES
-- name: GetAllEquipment :many
SELECT * FROM serial_numbers;

-- name: GetEquipmentByDeviceType :many
SELECT * FROM serial_numbers
WHERE device_type_id = $1;

-- name: GetEquipmentByManufacturer :many
SELECT * FROM serial_numbers
WHERE manufacturer_id = $1;

-- name: GetEquipmentBySerialNumber :one
SELECT * FROM serial_numbers
WHERE serial_number = $1;

-- name: GetEquipmentLikeSerialNumber :many
SELECT * FROM serial_numbers
WHERE serial_number LIKE $1;

-- name: GetEquipmentByDeviceTypeAndManufacturer :many
SELECT * FROM serial_numbers
WHERE device_type_id = $1 AND manufacturer_id = $2;

-- name: GetEquipmentByDeviceTypeAndSerialNumber :many
SELECT * FROM serial_numbers
WHERE device_type_id = $1 AND serial_number = $2;

-- name: GetEquipmentByManufacturerAndSerialNumber :many
SELECT * FROM serial_numbers
WHERE manufacturer_id = $1 AND serial_number = $2;

-- name: GetEquipmentByDeviceTypeManufacturerAndSerialNumber :many
SELECT * FROM serial_numbers
WHERE device_type_id = $1 AND manufacturer_id = $2 AND serial_number = $3;

