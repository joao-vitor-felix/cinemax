package domain

import (
	"time"

	"github.com/google/uuid"
)

type Gender string

const (
	Male           Gender = "male"
	Female         Gender = "female"
	Other          Gender = "other"
	PreferNotToSay Gender = "prefer_not_to_say"
)

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female, Other, PreferNotToSay:
		return true
	}
	return false
}

type User struct {
	ID              uuid.UUID
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	PasswordHash    string
	DateOfBirth     string
	Gender          Gender
	ProfilePhotoURL *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (u *User) IsAgeValid(targetAge int) bool {
	dob, _ := time.Parse("2006-01-02", u.DateOfBirth)
	now := time.Now()
	age := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		age--
	}
	return age >= targetAge
}

func NewUser(firstName, lastName, email, phone, dateOfBirth string, gender Gender) (*User, error) {
	user := &User{
		ID:          uuid.New(),
		FirstName:   firstName,
		LastName:    lastName,
		Email:       email,
		Phone:       phone,
		DateOfBirth: dateOfBirth,
		Gender:      gender,
	}

	if !user.Gender.IsValid() {
		return nil, InvalidGenderError
	}
	return user, nil
}
