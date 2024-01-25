package utils

import (
	"bytes"
	"encoding/json"
	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"strings"
	"time"
)

func GenerateId() string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	return strings.ToLower(ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String())
}

func CalculateAge(dateOfBirth string) (int, error) {
	dob, err := time.Parse("2006-01-02", dateOfBirth)
	if err != nil {
		return 0, err
	}
	age := time.Since(dob)

	years := age.Hours() / 24 / 365.25
	ageInYears := int(years)
	return ageInYears, nil
}

func UnPack(in interface{}, target interface{}) error {
	var e1 error
	var b []byte
	switch in := in.(type) {
	case []byte:
		b = in
	default:
		b, e1 = json.Marshal(in)
		if e1 != nil {
			return e1
		}
	}

	buf := bytes.NewBuffer(b)
	enc := json.NewDecoder(buf)
	enc.UseNumber()
	if err := enc.Decode(&target); err != nil {
		return err
	}
	return nil
}

func EncryptPassword(password string) string {
	encrypted, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encrypted)
}

func VerifyPasscode(hashedPasscode, providedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPasscode), []byte(providedPassword))
	return err == nil
}
