package models

import (
    "errors"
    "regexp"
    "time"

    "github.com/dgrijalva/jwt-go"
    "github.com/yemyoaung/managing-vehicle-tracking-common"
    "go.mongodb.org/mongo-driver/bson/primitive"
)

var (
    ErrIDMissing        = errors.New("id is missing")
    ErrCreatedAtMissing = errors.New("created_at is missing")
    ErrUpdatedAtMissing = errors.New("updated_at is missing")
    ErrEmailEmpty       = errors.New("email is required")
    ErrInvalidEmail     = errors.New("email is invalid")
    ErrPasswordEmpty    = errors.New("password is required")
    ErrRoleEmpty        = errors.New("role is required")
    ErrInvalidRole      = errors.New("role is invalid")
)

const (
    emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

var (
    emailRe = regexp.MustCompile(emailRegex)
)

type Email string

func (e Email) Validate() error {
    if e == "" {
        return ErrEmailEmpty
    }
    if !emailRe.MatchString(string(e)) {
        return ErrInvalidEmail
    }
    return nil
}

type Role string

// Validate checks if the role is valid
// Since golang doesn't have enums, we need to validate the role
func (r Role) Validate() error {
    if r == "" {
        return ErrRoleEmpty
    }
    if r != AdminRole && r != UserRole {
        return ErrInvalidRole
    }
    return nil
}

const (
    AdminRole Role = "admin"
    UserRole  Role = "user"
)

type User struct {
    ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
    Email     Email              `json:"email" bson:"email"`
    Password  string             `json:"-" bson:"password"`
    Role      Role               `json:"role" bson:"role"`
    CreatedAt time.Time          `json:"created_at" bson:"created_at"`
    UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
    DeletedAt *time.Time         `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
}

func NewUser() *User {
    return &User{}
}

func (u *User) SetEmail(email string) (*User, error) {
    u.Email = Email(email)
    if err := u.Email.Validate(); err != nil {
        return nil, err
    }
    return u, nil
}

func (u *User) SetPassword(password string) (*User, error) {
    if password == "" {
        return nil, ErrPasswordEmpty
    }
    hashPassword, err := common.HashPassword(password)
    if err != nil {
        return nil, err
    }
    u.Password = hashPassword
    return u, nil
}

func (u *User) SetRole(role Role) (*User, error) {
    if err := role.Validate(); err != nil {
        return nil, err
    }
    u.Role = role
    return u, nil
}

func (u *User) Validate() error {
    if err := u.Email.Validate(); err != nil {
        return err
    }
    if u.Password == "" {
        return ErrPasswordEmpty
    }
    if err := u.Role.Validate(); err != nil {
        return err
    }
    return nil
}

func (u *User) Build() error {
    if u.CreatedAt.IsZero() {
        u.CreatedAt = time.Now()
    }
    u.UpdatedAt = time.Now()
    return u.Validate()
}

func (u *User) Check() error {
    if u.ID.IsZero() {
        return ErrIDMissing
    }
    if u.CreatedAt.IsZero() {
        return ErrCreatedAtMissing
    }
    if u.UpdatedAt.IsZero() {
        return ErrUpdatedAtMissing
    }
    return u.Validate()
}

func (u *User) Claim() *jwt.StandardClaims {
    return &jwt.StandardClaims{
        Id:        u.ID.Hex(),
        Subject:   string(u.Email),
        Audience:  string(u.Role),
        IssuedAt:  time.Now().Unix(),
        ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
        Issuer:    "auth-service",
        NotBefore: time.Now().Unix(),
    }
}

// AuthUser represents the authenticated user
type AuthUser struct {
    Data struct {
        Id        string    `json:"id"`
        Email     string    `json:"email"`
        Role      string    `json:"role"`
        CreatedAt time.Time `json:"created_at"`
        UpdatedAt time.Time `json:"updated_at"`
    } `json:"data"`
}
