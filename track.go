package models

import (
    "errors"
    "time"

    "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
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
    VehicleID string             `json:"vehicle_id"`
    Location  string             `json:"location"`
    // since mileage can be a float value and it can be a large number, we will use float64
    Mileage       float64       `json:"mileage"`
    Status        VehicleStatus `json:"status"`
    FuelCondition FuelCondition `json:"fuel_condition"`
    CreatedAt     time.Time     `json:"created_at" bson:"created_at"`
    UpdatedAt     time.Time     `json:"updated_at" bson:"updated_at"`
    DeletedAt     *time.Time    `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func (t *TrackingData) Validate() error {
    if t.VehicleID == "" {
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
