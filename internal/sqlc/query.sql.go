// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: query.sql

package sqlc

import (
	"context"
)

const createDeviceType = `-- name: CreateDeviceType :exec
INSERT INTO device_type (name) VALUES (?)
`

func (q *Queries) CreateDeviceType(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, createDeviceType, name)
	return err
}

const createManufacturer = `-- name: CreateManufacturer :exec
INSERT INTO manufacturer (name) VALUES (?)
`

func (q *Queries) CreateManufacturer(ctx context.Context, name string) error {
	_, err := q.db.ExecContext(ctx, createManufacturer, name)
	return err
}

const deleteDeviceType = `-- name: DeleteDeviceType :exec
DELETE FROM device_type
WHERE id = ?
`

func (q *Queries) DeleteDeviceType(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteDeviceType, id)
	return err
}

const deleteManufacturer = `-- name: DeleteManufacturer :exec
DELETE FROM manufacturer
WHERE id = ?
`

func (q *Queries) DeleteManufacturer(ctx context.Context, id int32) error {
	_, err := q.db.ExecContext(ctx, deleteManufacturer, id)
	return err
}

const getAllEquipment = `-- name: GetAllEquipment :many
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
LIMIT ? OFFSET ?
`

type GetAllEquipmentParams struct {
	Limit  int32
	Offset int32
}

// EQUIPMENT QUERIES
func (q *Queries) GetAllEquipment(ctx context.Context, arg GetAllEquipmentParams) ([]SerialNumber, error) {
	rows, err := q.db.QueryContext(ctx, getAllEquipment, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SerialNumber
	for rows.Next() {
		var i SerialNumber
		if err := rows.Scan(
			&i.AutoID,
			&i.DeviceTypeID,
			&i.ManufacturerID,
			&i.SerialNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getDeviceTypeById = `-- name: GetDeviceTypeById :one
SELECT id, name, status FROM device_type
WHERE id = ?
ORDER BY id
`

func (q *Queries) GetDeviceTypeById(ctx context.Context, id int32) (DeviceType, error) {
	row := q.db.QueryRowContext(ctx, getDeviceTypeById, id)
	var i DeviceType
	err := row.Scan(&i.ID, &i.Name, &i.Status)
	return i, err
}

const getDeviceTypeByName = `-- name: GetDeviceTypeByName :one
SELECT id, name, status FROM device_type
WHERE name = ?
ORDER BY id
`

func (q *Queries) GetDeviceTypeByName(ctx context.Context, name string) (DeviceType, error) {
	row := q.db.QueryRowContext(ctx, getDeviceTypeByName, name)
	var i DeviceType
	err := row.Scan(&i.ID, &i.Name, &i.Status)
	return i, err
}

const getDeviceTypesActive = `-- name: GetDeviceTypesActive :many
SELECT id, name, status FROM device_type
ORDER BY id
`

// DEVICETYPE QUERIES
func (q *Queries) GetDeviceTypesActive(ctx context.Context) ([]DeviceType, error) {
	rows, err := q.db.QueryContext(ctx, getDeviceTypesActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []DeviceType
	for rows.Next() {
		var i DeviceType
		if err := rows.Scan(&i.ID, &i.Name, &i.Status); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquipmentByDeviceType = `-- name: GetEquipmentByDeviceType :many
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
WHERE device_type_id = ?
LIMIT ? OFFSET ?
`

type GetEquipmentByDeviceTypeParams struct {
	DeviceTypeID int32
	Limit        int32
	Offset       int32
}

func (q *Queries) GetEquipmentByDeviceType(ctx context.Context, arg GetEquipmentByDeviceTypeParams) ([]SerialNumber, error) {
	rows, err := q.db.QueryContext(ctx, getEquipmentByDeviceType, arg.DeviceTypeID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SerialNumber
	for rows.Next() {
		var i SerialNumber
		if err := rows.Scan(
			&i.AutoID,
			&i.DeviceTypeID,
			&i.ManufacturerID,
			&i.SerialNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquipmentByDeviceTypeAndManufacturer = `-- name: GetEquipmentByDeviceTypeAndManufacturer :many
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
WHERE device_type_id = ? AND manufacturer_id = ?
LIMIT ? OFFSET ?
`

type GetEquipmentByDeviceTypeAndManufacturerParams struct {
	DeviceTypeID   int32
	ManufacturerID int32
	Limit          int32
	Offset         int32
}

func (q *Queries) GetEquipmentByDeviceTypeAndManufacturer(ctx context.Context, arg GetEquipmentByDeviceTypeAndManufacturerParams) ([]SerialNumber, error) {
	rows, err := q.db.QueryContext(ctx, getEquipmentByDeviceTypeAndManufacturer,
		arg.DeviceTypeID,
		arg.ManufacturerID,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SerialNumber
	for rows.Next() {
		var i SerialNumber
		if err := rows.Scan(
			&i.AutoID,
			&i.DeviceTypeID,
			&i.ManufacturerID,
			&i.SerialNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquipmentByDeviceTypeAndSerialNumber = `-- name: GetEquipmentByDeviceTypeAndSerialNumber :many
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
WHERE device_type_id = ? AND serial_number = ?
`

type GetEquipmentByDeviceTypeAndSerialNumberParams struct {
	DeviceTypeID int32
	SerialNumber string
}

func (q *Queries) GetEquipmentByDeviceTypeAndSerialNumber(ctx context.Context, arg GetEquipmentByDeviceTypeAndSerialNumberParams) ([]SerialNumber, error) {
	rows, err := q.db.QueryContext(ctx, getEquipmentByDeviceTypeAndSerialNumber, arg.DeviceTypeID, arg.SerialNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SerialNumber
	for rows.Next() {
		var i SerialNumber
		if err := rows.Scan(
			&i.AutoID,
			&i.DeviceTypeID,
			&i.ManufacturerID,
			&i.SerialNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquipmentByDeviceTypeManufacturerAndSerialNumber = `-- name: GetEquipmentByDeviceTypeManufacturerAndSerialNumber :many
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
WHERE device_type_id = ? AND manufacturer_id = ? AND serial_number = ?
`

type GetEquipmentByDeviceTypeManufacturerAndSerialNumberParams struct {
	DeviceTypeID   int32
	ManufacturerID int32
	SerialNumber   string
}

func (q *Queries) GetEquipmentByDeviceTypeManufacturerAndSerialNumber(ctx context.Context, arg GetEquipmentByDeviceTypeManufacturerAndSerialNumberParams) ([]SerialNumber, error) {
	rows, err := q.db.QueryContext(ctx, getEquipmentByDeviceTypeManufacturerAndSerialNumber, arg.DeviceTypeID, arg.ManufacturerID, arg.SerialNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SerialNumber
	for rows.Next() {
		var i SerialNumber
		if err := rows.Scan(
			&i.AutoID,
			&i.DeviceTypeID,
			&i.ManufacturerID,
			&i.SerialNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquipmentByManufacturer = `-- name: GetEquipmentByManufacturer :many
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
WHERE manufacturer_id = ?
LIMIT ? OFFSET ?
`

type GetEquipmentByManufacturerParams struct {
	ManufacturerID int32
	Limit          int32
	Offset         int32
}

func (q *Queries) GetEquipmentByManufacturer(ctx context.Context, arg GetEquipmentByManufacturerParams) ([]SerialNumber, error) {
	rows, err := q.db.QueryContext(ctx, getEquipmentByManufacturer, arg.ManufacturerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SerialNumber
	for rows.Next() {
		var i SerialNumber
		if err := rows.Scan(
			&i.AutoID,
			&i.DeviceTypeID,
			&i.ManufacturerID,
			&i.SerialNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquipmentByManufacturerAndSerialNumber = `-- name: GetEquipmentByManufacturerAndSerialNumber :many
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
WHERE manufacturer_id = ? AND serial_number = ?
`

type GetEquipmentByManufacturerAndSerialNumberParams struct {
	ManufacturerID int32
	SerialNumber   string
}

func (q *Queries) GetEquipmentByManufacturerAndSerialNumber(ctx context.Context, arg GetEquipmentByManufacturerAndSerialNumberParams) ([]SerialNumber, error) {
	rows, err := q.db.QueryContext(ctx, getEquipmentByManufacturerAndSerialNumber, arg.ManufacturerID, arg.SerialNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SerialNumber
	for rows.Next() {
		var i SerialNumber
		if err := rows.Scan(
			&i.AutoID,
			&i.DeviceTypeID,
			&i.ManufacturerID,
			&i.SerialNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getEquipmentBySerialNumber = `-- name: GetEquipmentBySerialNumber :one
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
WHERE serial_number = ?
`

func (q *Queries) GetEquipmentBySerialNumber(ctx context.Context, serialNumber string) (SerialNumber, error) {
	row := q.db.QueryRowContext(ctx, getEquipmentBySerialNumber, serialNumber)
	var i SerialNumber
	err := row.Scan(
		&i.AutoID,
		&i.DeviceTypeID,
		&i.ManufacturerID,
		&i.SerialNumber,
	)
	return i, err
}

const getEquipmentLikeSerialNumber = `-- name: GetEquipmentLikeSerialNumber :many
SELECT auto_id, device_type_id, manufacturer_id, serial_number FROM serial_numbers
WHERE serial_number LIKE ?
LIMIT ? OFFSET ?
`

type GetEquipmentLikeSerialNumberParams struct {
	SerialNumber string
	Limit        int32
	Offset       int32
}

func (q *Queries) GetEquipmentLikeSerialNumber(ctx context.Context, arg GetEquipmentLikeSerialNumberParams) ([]SerialNumber, error) {
	rows, err := q.db.QueryContext(ctx, getEquipmentLikeSerialNumber, arg.SerialNumber, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SerialNumber
	for rows.Next() {
		var i SerialNumber
		if err := rows.Scan(
			&i.AutoID,
			&i.DeviceTypeID,
			&i.ManufacturerID,
			&i.SerialNumber,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getManufacturerById = `-- name: GetManufacturerById :one
SELECT id, name, status FROM manufacturer
WHERE id = ?
ORDER BY id
`

func (q *Queries) GetManufacturerById(ctx context.Context, id int32) (Manufacturer, error) {
	row := q.db.QueryRowContext(ctx, getManufacturerById, id)
	var i Manufacturer
	err := row.Scan(&i.ID, &i.Name, &i.Status)
	return i, err
}

const getManufacturerByName = `-- name: GetManufacturerByName :one
SELECT id, name, status FROM manufacturer
WHERE name = ?
ORDER BY id
`

func (q *Queries) GetManufacturerByName(ctx context.Context, name string) (Manufacturer, error) {
	row := q.db.QueryRowContext(ctx, getManufacturerByName, name)
	var i Manufacturer
	err := row.Scan(&i.ID, &i.Name, &i.Status)
	return i, err
}

const getManufacturersActive = `-- name: GetManufacturersActive :many
SELECT id, name, status FROM manufacturer
ORDER BY id
`

// MANUFACTURER QUERIES
func (q *Queries) GetManufacturersActive(ctx context.Context) ([]Manufacturer, error) {
	rows, err := q.db.QueryContext(ctx, getManufacturersActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Manufacturer
	for rows.Next() {
		var i Manufacturer
		if err := rows.Scan(&i.ID, &i.Name, &i.Status); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSerialNumberBySerialNumber = `-- name: GetSerialNumberBySerialNumber :one
SELECT serial_number FROM serial_numbers
WHERE serial_number = ?
`

func (q *Queries) GetSerialNumberBySerialNumber(ctx context.Context, serialNumber string) (string, error) {
	row := q.db.QueryRowContext(ctx, getSerialNumberBySerialNumber, serialNumber)
	var serial_number string
	err := row.Scan(&serial_number)
	return serial_number, err
}

const getSerialNumberLikeSerialNumber = `-- name: GetSerialNumberLikeSerialNumber :many
SELECT serial_number FROM serial_numbers
WHERE serial_number LIKE ?
`

func (q *Queries) GetSerialNumberLikeSerialNumber(ctx context.Context, serialNumber string) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getSerialNumberLikeSerialNumber, serialNumber)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var serial_number string
		if err := rows.Scan(&serial_number); err != nil {
			return nil, err
		}
		items = append(items, serial_number)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getSerialNumbers = `-- name: GetSerialNumbers :many
SELECT serial_number FROM serial_numbers
`

// SERIALNUMBER QUERIES
func (q *Queries) GetSerialNumbers(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getSerialNumbers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var serial_number string
		if err := rows.Scan(&serial_number); err != nil {
			return nil, err
		}
		items = append(items, serial_number)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateDeviceType = `-- name: UpdateDeviceType :exec
UPDATE device_type SET name = ?
WHERE id = ?
`

type UpdateDeviceTypeParams struct {
	Name string
	ID   int32
}

func (q *Queries) UpdateDeviceType(ctx context.Context, arg UpdateDeviceTypeParams) error {
	_, err := q.db.ExecContext(ctx, updateDeviceType, arg.Name, arg.ID)
	return err
}

const updateDeviceTypeStatus = `-- name: UpdateDeviceTypeStatus :exec
UPDATE device_type SET status = ?
WHERE id = ?
`

type UpdateDeviceTypeStatusParams struct {
	Status DeviceTypeStatus
	ID     int32
}

func (q *Queries) UpdateDeviceTypeStatus(ctx context.Context, arg UpdateDeviceTypeStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateDeviceTypeStatus, arg.Status, arg.ID)
	return err
}

const updateManufacturer = `-- name: UpdateManufacturer :exec
UPDATE manufacturer SET name = ?
WHERE id = ?
`

type UpdateManufacturerParams struct {
	Name string
	ID   int32
}

func (q *Queries) UpdateManufacturer(ctx context.Context, arg UpdateManufacturerParams) error {
	_, err := q.db.ExecContext(ctx, updateManufacturer, arg.Name, arg.ID)
	return err
}

const updateManufacturerStatus = `-- name: UpdateManufacturerStatus :exec
UPDATE manufacturer SET status = ?
WHERE id = ?
`

type UpdateManufacturerStatusParams struct {
	Status ManufacturerStatus
	ID     int32
}

func (q *Queries) UpdateManufacturerStatus(ctx context.Context, arg UpdateManufacturerStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateManufacturerStatus, arg.Status, arg.ID)
	return err
}

const updateSerialNumber = `-- name: UpdateSerialNumber :exec
UPDATE serial_numbers SET serial_number = ?
WHERE serial_number = ?
`

type UpdateSerialNumberParams struct {
	SerialNumber   string
	SerialNumber_2 string
}

func (q *Queries) UpdateSerialNumber(ctx context.Context, arg UpdateSerialNumberParams) error {
	_, err := q.db.ExecContext(ctx, updateSerialNumber, arg.SerialNumber, arg.SerialNumber_2)
	return err
}
