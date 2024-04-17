package models

// @description DeviceType is a struct for device type
type DeviceType struct {
	// ID is an int32 for device type id
	ID int32 `json:"id" example:"1"`
	// Name is a string for device type name
	Name string `json:"name" example:"computer"`
	// Status is a string for device type status either active or inactive
	Status string `json:"status" example:"active|inactive"`
}

// @description Manufacturer is a struct for manufacturer
type Manufacturer struct {
	// ID is an int32 for manufacturer id
	ID int32 `json:"id"      example:"1"`
	// Name is a string for manufacturer name
	Name string `json:"name"    example:"Apple"`
	// Status is a string for manufacturer status either active or inactive
	Status string `json:"status"  example:"active|inactive"`
}

// @description Equipment is a struct for equipment
type Equipment struct {
	AutoID         int32  `json:"auto_id" example:"1"`               // AutoID is an int32 for equipment auto id
	DeviceTypeID   int32  `json:"device_type_id" example:"1"`        // DeviceTypeID is an int32 for device id
	ManufacturerID int32  `json:"manufacturer_id" example:"1"`       // ManufacturerID is an int32 for manufacturer id
	SerialNumber   string `json:"serial_number" example:"SN-123456"` // SerialNumber is a string for equipment serial number
	Status         string `json:"status" example:"active|inactive"`
}

// Message is an interface for response message can be string, models.DeviceType, models.Manufacturer, models.Equipment
type Message interface{}

// @description JsonResponse is a struct for response JSON message
type JsonResponse struct {
	// Status is a string for response status
	Status string `json:"Status" example:"SUCCESS|ERROR"`
	// Message is an interface for response message can be string, models.DeviceType, models.Manufacturer, models.Equipment
	Message Message `json:"MSG"`
	// Action is a string for response action
	Action string `json:"Action"`
}
