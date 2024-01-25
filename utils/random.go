package utils

import (
	"api/constants"
	"math/rand"
	"time"
)

func GetRandomGender() string {
	genders := []string{"male", "female", "non-binary"}
	return genders[rand.Intn(len(genders))]
}

func GetRandomDOB() time.Time {
	yearsAgo := rand.Intn(30) + 18 // Random age between 18 and 48
	return time.Now().AddDate(-yearsAgo, 0, 0)
}

func GetRandomLocationInNorthLondon() []float64 {
	// Generate random latitude and longitude within North London bounding box
	lat := rand.Float64()*(constants.NorthLondonCoordinates.MaxLat-constants.NorthLondonCoordinates.MinLat) + constants.NorthLondonCoordinates.MinLat
	lon := rand.Float64()*(constants.NorthLondonCoordinates.MaxLon-constants.NorthLondonCoordinates.MinLon) + constants.NorthLondonCoordinates.MinLon
	return []float64{lat, lon}
}

func GetRandomEthnicity() string {
	ethnicities := []string{"white", "black", "asian", "latino", "other"}
	return ethnicities[rand.Intn(len(ethnicities))]
}

func GetRandomPet() string {
	pets := []string{"dog", "cat", "bird", "reptile", "prefer not to say", "none"}
	return pets[rand.Intn(len(pets))]
}

func GetRandomSexuality() string {
	sexualities := []string{"straight", "gay", "bisexual", "pansexual", "asexual", "other"}
	return sexualities[rand.Intn(len(sexualities))]
}

func GetRandomDrinkingHabit() string {
	habits := []string{"yes", "sometimes", "no", "none"}
	return habits[rand.Intn(len(habits))]
}

func GetRandomSmokingHabit() string {
	habits := []string{"yes", "sometimes", "no", "none"}
	return habits[rand.Intn(len(habits))]
}

func GetRandomDrugsHabit() string {
	habits := []string{"yes", "sometimes", "no", "none"}
	return habits[rand.Intn(len(habits))]
}

func GetRandomDatingIntentions() string {
	intentions := []string{"life partner", "shorter time", "none", "figuring out", "other"}
	return intentions[rand.Intn(len(intentions))]
}

func GetRandomReligion() string {
	religions := []string{"christian", "muslim", "hindu", "buddhist", "other"}
	return religions[rand.Intn(len(religions))]
}

func GetRandomHeight() float64 {
	return rand.Float64()*(190.0-150.0) + 150.0
}
