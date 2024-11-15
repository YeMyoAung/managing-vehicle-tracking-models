package models

type Model interface {
    // Validate validates the model before it is saved
    Validate() error
    // Build builds the model before it is saved
    Build() error
    // Check checks the model if it is come from the database
    Check() error
}
