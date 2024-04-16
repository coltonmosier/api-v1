package models

type DeviceType struct {
	ID     int32  `json:"id" example:"1"`
	Name   string `json:"name" example:"computer"`
	Status string `json:"status" example:"active|inactive"`
}

type Manufacturer struct {
	ID     int32  `json:"id" example:"1"`
	Name   string `json:"name" example:"Apple"`
	Status string `json:"status" example:"active|inactive"`
}

type Equipment struct {
	AutoID         int32  `json:"auto_id" example:"1"`
	DeviceTypeID   int32  `json:"device_type_id" example:"1"`
	ManufacturerID int32  `json:"manufacturer_id" example:"1"`
	SerialNumber   string `json:"serial_number" example:"SN-123456"`
}

// Message is an interface for response message can be string, models.DeviceType, models.Manufacturer, models.Equipment
type Message interface{}

// JsonResponse is a struct for response
type JsonResponse struct {
    // Status is a string for response status
	Status  string  `json:"Status" example:"SUCCESS|ERROR"`
    // Message is an interface for response message can be string, models.DeviceType, models.Manufacturer, models.Equipment
	Message Message `json:"MSG"`
    // Action is a string for response action
	Action  string  `json:"Action"`
}
