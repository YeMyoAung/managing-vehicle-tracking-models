package models

import (
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
    ErrInvalidVehicleID     = errors.New("invalid vehicle id")
    ErrVehicleIDEmpty       = errors.New("vehicle id is empty")
    ErrLocationEmpty        = errors.New("location is empty")
    ErrMileageEmpty         = errors.New("mileage is empty")
    ErrFuelConditionEmpty   = errors.New("fuel condition is empty")
    ErrInvalidFuelCondition = errors.New("invalid fuel condition")
)

type FuelCondition string

func (f FuelCondition) Valid() error {
    if f == "" {
        return ErrFuelConditionEmpty
    }

    if f != FuelConditionEmpty && f != FuelConditionLow && f != FuelConditionHalf && f != FuelConditionFull {
        return ErrInvalidFuelCondition
    }

    return nil
}

const (
    FuelConditionEmpty FuelCondition = "empty"
    FuelConditionLow   FuelCondition = "low"
    FuelConditionHalf  FuelCondition = "half"
    FuelConditionFull  FuelCondition = "full"
)

type TrackingData struct {
    ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
    VehicleID primitive.ObjectID `json:"vehicle_id"`
    Location  string             `json:"location"`
    // since mileage can be a float value and it can be a large number, we will use float64
    Mileage       float64       `json:"mileage"`
    Status        VehicleStatus `json:"status"`
    FuelCondition FuelCondition `json:"fuel_condition"`
    CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
    UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
    DeletedAt     *time.Time    `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func NewTrackingData() *TrackingData {
    return &TrackingData{}
}

func (t *TrackingData) SetVehicleID(hex string) (*TrackingData, error) {
    var err error
    t.VehicleID, err = primitive.ObjectIDFromHex(hex)
    if err != nil {
        return nil, ErrInvalidVehicleID
    }
    return t, nil
}

func (t *TrackingData) SetLocation(location string) *TrackingData {
    t.Location = location
    return t
}

func (t *TrackingData) SetMileage(mileage float64) *TrackingData {
    t.Mileage = mileage
    return t
}

func (t *TrackingData) SetStatus(status VehicleStatus) *TrackingData {
    t.Status = status
    return t
}

func (t *TrackingData) SetFuelCondition(fuelCondition FuelCondition) *TrackingData {
    t.FuelCondition = fuelCondition
    return t
}

func (t *TrackingData) Validate() error {
    if t.VehicleID.IsZero() {
        return ErrVehicleIDEmpty
    }
    if t.Location == "" {
        return ErrLocationEmpty
    }
    // since mileage only can be a positive number and it can't be zero
    if t.Mileage == 0 {
        return ErrMileageEmpty
    }

    if err := t.Status.Valid(); err != nil {
        return err
    }

    if err := t.FuelCondition.Valid(); err != nil {
        return err
    }

    return nil
}

func (t *TrackingData) Build() error {
    if t.CreatedAt.IsZero() {
        t.CreatedAt = time.Now()
    }
    t.UpdatedAt = time.Now()
    return t.Validate()
}

func (t *TrackingData) Check() error {
    if t.ID.IsZero() {
        return ErrIDMissing
    }
    if t.CreatedAt.IsZero() {
        return ErrCreatedAtMissing
    }
    if t.UpdatedAt.IsZero() {
        return ErrUpdatedAtMissing
    }
    return t.Validate()
}

// TrackingDataRequest is not a model, but we will put it here for consistency
type TrackingDataRequest struct {
    VehicleID     string        `json:"vehicle_id" validate:"required"`
    Location      string        `json:"location" validate:"required"`
    Mileage       float64       `json:"mileage" validate:"required"`
    Status        VehicleStatus `json:"status" validate:"required"`
    FuelCondition FuelCondition `json:"fuel_condition" validate:"required"`
}

func (t *TrackingDataRequest) Validate() error {
    if t.VehicleID == "" {
        return ErrVehicleIDEmpty
    }
    _, err := primitive.ObjectIDFromHex(t.VehicleID)
    if err != nil {
        return ErrInvalidVehicleID
    }
    if t.Location == "" {
        return ErrLocationEmpty
    }
    if t.Mileage == 0 {
        return ErrMileageEmpty
    }
    if err := t.Status.Valid(); err != nil {
        return err
    }
    if err := t.FuelCondition.Valid(); err != nil {
        return err
    }
    return nil
}

func (t *TrackingDataRequest) ToTrackingData() (*TrackingData, error) {
    vehicleID, err := primitive.ObjectIDFromHex(t.VehicleID)
    if err != nil {
        return nil, ErrInvalidVehicleID
    }
    return &TrackingData{
        VehicleID:     vehicleID,
        Location:      t.Location,
        Mileage:       t.Mileage,
        Status:        t.Status,
        FuelCondition: t.FuelCondition,
    }, nil
}
