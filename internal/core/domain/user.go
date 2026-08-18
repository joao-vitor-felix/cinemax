package domain

import (
	"log"
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
	Phone           *string
	PasswordHash    *string
	DateOfBirth     *string
	Gender          Gender
	ProfilePhotoURL *string
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type OAuthUser struct {
	FirstName string
	LastName  string
	Email     string
}

func (u *User) IsAgeValid(targetAge int) bool {
	if u.DateOfBirth == nil {
		return false
	}
	dob, err := time.Parse("2006-01-02", *u.DateOfBirth)
	if err != nil {
		log.Default().Printf("Error parsing date of birth: %v", err)
		return false
	}
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
		Phone:       &phone,
		DateOfBirth: &dateOfBirth,
		Gender:      gender,
	}

	if !user.Gender.IsValid() {
		return nil, InvalidGenderError
	}
	return user, nil
}

func NewOAuthUser(firstName, lastName, email string) *OAuthUser {
	return &OAuthUser{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
}

func (u *User) HasPassword() bool {
	return u.PasswordHash != nil
}

func (u *User) IsAllowedToBook() bool {
	if u.DateOfBirth == nil {
		return false
	}
	if _, err := time.Parse("2006-01-02", *u.DateOfBirth); err != nil {
		return false
	}
	return true
}
