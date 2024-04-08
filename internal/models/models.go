package models

type DeviceType struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Manufacturer struct {
	ID     int32  `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Equipment struct {
	AutoID         int32  `json:"auto_id"`
	DeviceTypeID   int32  `json:"device_type_id"`
	ManufacturerID int32  `json:"manufacturer_id"`
	SerialNumber   string `json:"serial_number"`
}
