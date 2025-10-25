package domain

import (
	"time"
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
	ID              string
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

func (u *User) IsAgeValid() bool {
	dob, _ := time.Parse("2006-01-02", u.DateOfBirth)
	now := time.Now()
	years := now.Year() - dob.Year()
	if now.YearDay() < dob.YearDay() {
		years--
	}
	return years >= 13
}

func NewUser(user User) (*User, error) {
	if !user.Gender.IsValid() {
		return nil, InvalidGenderError
	}
	if !user.IsAgeValid() {
		return nil, UserTooYoungError
	}
	return &user, nil
}
