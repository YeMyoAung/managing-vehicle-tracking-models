package models

import (
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
    ErrVehicleNameEmpty     = errors.New("vehicle name is required")
    ErrVehicleModelEmpty    = errors.New("vehicle model is required")
    ErrVehicleStatusEmpty   = errors.New("vehicle status is required")
    ErrLicenseNumberEmpty   = errors.New("license number is required")
    ErrInvalidVehicleStatus = errors.New("invalid vehicle status")
)

type VehicleStatus string

// Valid checks if the status is valid
// Since golang doesn't have enums, we need to validate the status 
func (v VehicleStatus) Valid() error {
    if v == "" {
        return ErrVehicleStatusEmpty
    }
    if v != VehicleStatusActive && v != VehicleStatusInactive && v != VehicleStatusRepair && v != VehicleStatusSold && v != VehicleStatusRented {
        return ErrInvalidVehicleStatus
    }
    return nil
}

var (
    VehicleStatusActive   VehicleStatus = "active"
    VehicleStatusInactive VehicleStatus = "inactive"
    VehicleStatusRepair   VehicleStatus = "repair"
    VehicleStatusSold     VehicleStatus = "sold"
    VehicleStatusRented   VehicleStatus = "rented"
)

type Vehicle struct {
    ID            primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    VehicleName   string             `json:"vehicle_name" bson:"vehicle_name"`
    VehicleModel  string             `json:"vehicle_model" bson:"vehicle_model"`
    VehicleStatus VehicleStatus      `json:"vehicle_status" bson:"vehicle_status"`
    // since mileage can be a float value and it can be a large number, we will use float64
    Mileage       float64    `json:"mileage" bson:"mileage"`
    LicenseNumber string     `json:"license_number" bson:"license_number"`
    CreatedAt     time.Time  `json:"created_at" bson:"created_at"`
    UpdatedAt     time.Time  `json:"updated_at" bson:"updated_at"`
    DeletedAt     *time.Time `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
    // we can add created_by, updated_by, deleted_by fields here
    // to track who created, updated, deleted the vehicle
    // but for now, we will skip it.
}

func NewVehicle() *Vehicle {
    return &Vehicle{}
}

func (v *Vehicle) SetVehicleName(name string) *Vehicle {
    v.VehicleName = name
    return v
}

func (v *Vehicle) SetVehicleModel(model string) *Vehicle {
    v.VehicleModel = model
    return v
}

func (v *Vehicle) SetVehicleStatus(status VehicleStatus) *Vehicle {
    v.VehicleStatus = status
    return v
}

func (v *Vehicle) SetMileage(mileage float64) *Vehicle {
    v.Mileage = mileage
    return v
}

func (v *Vehicle) SetLicenseNumber(license string) *Vehicle {
    v.LicenseNumber = license
    return v
}

func (v *Vehicle) Validate() error {
    if v.VehicleName == "" {
        return ErrVehicleNameEmpty
    }
    if v.VehicleModel == "" {
        return ErrVehicleModelEmpty
    }
    if err := v.VehicleStatus.Valid(); err != nil {
        return err
    }
    // we will skip mileage validation for now, because if we have a new vehicle, mileage will be 0
    if v.LicenseNumber == "" {
        return ErrLicenseNumberEmpty
    }
    return nil
}

func (v *Vehicle) Build() error {
    if v.CreatedAt.IsZero() {
        v.CreatedAt = time.Now()
    }
    v.UpdatedAt = time.Now()
    return v.Validate()
}

func (v *Vehicle) Check() error {
    if v.ID.IsZero() {
        return ErrIDMissing
    }
    if v.CreatedAt.IsZero() {
        return ErrCreatedAtMissing
    }
    if v.UpdatedAt.IsZero() {
        return ErrUpdatedAtMissing
    }
    return v.Validate()
}
