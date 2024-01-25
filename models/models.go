package models

import (
	"errors"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"time"
)

type Gender string

const (
	Male      Gender = "male"
	Female    Gender = "female"
	NonBinary Gender = "non-binary"
)

func (g Gender) IsValid() bool {
	switch g {
	case Male, Female, NonBinary:
		return true
	default:
		return false
	}
}

var (
	validIntentions     = []interface{}{"life partner", "shorter time", "none", "figuring out", "other"}
	validGenders        = []interface{}{"male", "female", "non-binary"}
	validEthnicities    = []interface{}{"white", "black", "asian", "latino", "other"}
	validPets           = []interface{}{"dog", "cat", "bird", "reptile", "prefer not to say", "none"}
	validSexualities    = []interface{}{"straight", "gay", "bisexual", "pansexual", "asexual", "other"}
	validDrinkingHabits = []interface{}{"yes", "no", "none"}
	validSmokingHabits  = []interface{}{"yes", "no", "none"}
	validDrugHabits     = []interface{}{"yes", "no", "none"}
	validReligions      = []interface{}{"christian", "muslim", "hindu", "buddhist", "other"}
)

func (g Gender) String() string {
	return string(g)
}

type Status string

const (
	Activated   Status = "activated"
	Deactivated Status = "deactivated"
)

func (s Status) IsValid() bool {
	switch s {
	case Activated, Deactivated:
		return true
	default:
		return false
	}
}

func (s Status) String() string {
	return string(s)
}

type User struct {
	ID               string        `bson:"id,omitempty" json:"id,omitempty"`
	Name             string        `bson:"name" json:"name,omitempty"`
	Email            string        `bson:"email" json:"email,omitempty"`
	Password         string        `bson:"password" json:"-"`
	Age              int           `json:"age,omitempty"`
	DateOfBirth      string        `bson:"date_of_birth" json:"date_of_birth,omitempty"`
	Location         []float64     `bson:"location" json:"location,omitempty"`
	Height           float64       `bson:"height" json:"height,omitempty"`
	Ethnicity        string        `bson:"ethnicity" json:"ethnicity,omitempty"`
	Gender           Gender        `bson:"gender" json:"gender,omitempty"`
	Distance         float64       `json:"distance,omitempty"`
	Pets             string        `bson:"pets" json:"pets,omitempty"`
	Religion         Religion      `bson:"religion" json:"religion,omitempty"`
	Drinking         DrinkingHabit `bson:"drinking" json:"drinking,omitempty"`
	Smoking          SmokingHabit  `bson:"smoking" json:"smoking,omitempty"`
	Drugs            DrugHabit     `bson:"drugs" json:"drugs,omitempty"`
	DatingIntentions string        `bson:"dating_intentions" json:"dating_intentions,omitempty"`
	Kids             int           `bson:"kids" json:"kids,omitempty"`
	Occupation       string        `bson:"occupation" json:"occupation,omitempty"`
	SwipeCount       int           `bson:"swipe_count" json:"swipe_count,omitempty"`
	Attractiveness   int           `bson:"attractiveness" json:"attractiveness,omitempty"`
	Bio              string        `bson:"bio" json:"bio,omitempty"`
	SwipingRate      float64       `json:"swiping_rate,omitempty"`
	DailySwipeBudget int           `json:"daily_swipe_budget,omitempty"`
}

func (a User) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Name, validation.Required),
		validation.Field(&a.Email, validation.Required, is.Email),
		validation.Field(&a.Password, validation.Required),
	)
}

type RegistrationPayload struct {
	Name        string `json:"name,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	Gender      Gender `json:"gender,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}

func (r RegistrationPayload) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Password, validation.Required),
		validation.Field(&r.Gender, validation.Required),
		validation.Field(&r.DateOfBirth, validation.Required),
	)
}

type RegistrationResponse struct {
	ID       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Age      int    `json:"age,omitempty"`
	Gender   Gender `json:"gender,omitempty"`
}

type LoginPayload struct {
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

func (l LoginPayload) Validate() error {
	return validation.ValidateStruct(&l,
		validation.Field(&l.Email, validation.Required, is.Email),
		validation.Field(&l.Password, validation.Required),
	)
}

type LoginResponse struct {
	Token string `json:"token,omitempty"`
}

type UpdateUserPayload struct {
	Location         []float64     `json:"location"`
	Pets             string        `json:"pets"`
	Sexuality        Sexuality     `json:"sexuality"`
	Religion         Religion      `json:"religion"`
	Drinking         DrinkingHabit `json:"drinking"`
	Smoking          SmokingHabit  `json:"smoking"`
	Drugs            DrugHabit     `json:"drugs"`
	DatingIntentions Intentions    `json:"dating_intentions"`
	Kids             string        `json:"kids"`
	Occupation       string        `json:"occupation"`
	Bio              string        `json:"bio"`
}

func (u UpdateUserPayload) Validate() error {
	return validation.ValidateStruct(
		&u,
	)
}

type UpdatePasswordPayload struct {
	OldPassword string `json:"old_password,omitempty"`
	NewPassword string `json:"new_password,omitempty"`
}

func (u UpdatePasswordPayload) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.OldPassword, validation.Required),
		validation.Field(&u.NewPassword, validation.Required),
	)
}

type UserFilter struct {
	MinAge            int    `json:"min_age"`
	MaxAge            int    `json:"max_age"`
	MinHeight         int    `json:"min_height"`
	MaxHeight         int    `json:"max_height"`
	MaxDistance       int    `json:"max_distance,omitempty"`
	DesiredEthnicity  string `json:"desired_ethnicity,omitempty"`
	DesiredPets       string `json:"desired_pets,omitempty"`
	DesiredSexuality  string `json:"desired_sexuality,omitempty"`
	DesiredDrinking   string `json:"desired_drinking,omitempty"`
	DesiredSmoking    string `json:"desired_smoking,omitempty"`
	DesiredDrugs      string `json:"desired_drugs,omitempty"`
	DesiredIntentions string `json:"desired_intentions,omitempty"`
	DesiredReligion   string `json:"desired_religion,omitempty"`
}

func (uf UserFilter) Validate() error {
	return validation.ValidateStruct(&uf,
		validation.Field(&uf.MinHeight, validation.Min(0)),
		validation.Field(&uf.MaxHeight, validation.Min(0), validation.By(func(value interface{}) error {
			if uf.MinHeight > value.(int) {
				return errors.New("MaxHeight must be greater than or equal to MinHeight")
			}
			return nil
		})),
		validation.Field(&uf.MinAge, validation.Min(0)),
		validation.Field(&uf.MaxAge, validation.Min(0), validation.By(func(value interface{}) error {
			if uf.MinAge > value.(int) {
				return errors.New("MaxAge must be greater than or equal to MinAge")
			}
			return nil
		})),
		//validation.Field(&uf.Latitude, validation.Min(-90), validation.Max(90)),
		//validation.Field(&uf.Longitude, validation.Min(-180), validation.Max(180)),
		validation.Field(&uf.MaxDistance, validation.Min(0)),
		validation.Field(&uf.DesiredEthnicity, validation.Length(0, 255), validation.In(validEthnicities...)),
		validation.Field(&uf.DesiredPets, validation.Length(0, 255), validation.In(validPets...)),
		validation.Field(&uf.DesiredSexuality, validation.Length(0, 255), validation.In(validSexualities...)),
		validation.Field(&uf.DesiredDrinking, validation.Length(0, 255), validation.In(validDrinkingHabits...)),
		validation.Field(&uf.DesiredSmoking, validation.Length(0, 255), validation.In(validSmokingHabits...)),
		validation.Field(&uf.DesiredDrugs, validation.Length(0, 255), validation.In(validDrugHabits...)),
		validation.Field(&uf.DesiredIntentions, validation.Length(0, 255), validation.In(validIntentions...)),
		validation.Field(&uf.DesiredReligion, validation.Length(0, 255), validation.In(validReligions...)),
	)
}

type Swipe struct {
	ID         string    `bson:"id,omitempty" json:"id,omitempty"`
	UserID     string    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	ProspectID string    `bson:"prospect_id,omitempty" json:"prospect_id,omitempty"`
	Interested bool      `bson:"interested" json:"interested"`
	SwipeTime  time.Time `bson:"swipe_time" json:"swipe_time"`
}

type SwipePayload struct {
	ProspectID string `json:"user_id,omitempty"`
	Interested bool   `json:"interested"`
}

func (sp SwipePayload) Validate() error {
	return validation.ValidateStruct(&sp,
		validation.Field(&sp.ProspectID, validation.Required, validation.Length(1, 255)),
		validation.Field(&sp.Interested, validation.Required),
	)
}

type SwipeResponse struct {
	Matched bool   `json:"matched"`
	MatchID string `json:"match_id"`
}

type Match struct {
	ID       string   `bson:"_id,omitempty" json:"id,omitempty"`
	Profiles []string `bson:"profiles,omitempty" json:"profiles,omitempty"`
	Matched  bool     `bson:"matched" json:"matched"`
}

type SwipeFilter struct {
	UserID       string    `bson:"user_id,omitempty" json:"user_id,omitempty"`
	ProspectID   *string   `bson:"prospect_id,omitempty" json:"prospect_id,omitempty"`
	Interested   *bool     `bson:"interested,omitempty" json:"interested,omitempty"`
	MinSwipeDate time.Time `bson:"min_swipe_date,omitempty" json:"min_swipe_date,omitempty"`
	MaxSwipeDate time.Time `bson:"max_swipe_date,omitempty" json:"max_swipe_date,omitempty"`
}

type MatchFilter struct {
	UserID string `bson:"user_id,omitempty" json:"user_id,omitempty"`
}

type MatchedUserInfo struct {
	CurrentUser  *User   `bson:"current_user,omitempty" json:"user,omitempty"`
	MatchedUsers []*User `bson:"matched_users,omitempty" json:"matched,omitempty"`
}

type Ethnicity string

const (
	White  Ethnicity = "white"
	Black  Ethnicity = "black"
	Asian  Ethnicity = "asian"
	Latino Ethnicity = "latino"
	Other  Ethnicity = "other"
)

type Pet string

const (
	Dog            Pet = "dog"
	Cat            Pet = "cat"
	Bird           Pet = "bird"
	Reptile        Pet = "reptile"
	PreferNotToSay Pet = "prefer not to say"
	None           Pet = "none"
)

type Sexuality string

const (
	Straight       Sexuality = "straight"
	Gay            Sexuality = "gay"
	Bisexual       Sexuality = "bisexual"
	Pansexual      Sexuality = "pansexual"
	Asexual        Sexuality = "asexual"
	OtherSexuality Sexuality = "other"
)

type DrinkingHabit string

const (
	YesDrinkingHabit DrinkingHabit = "yes"
	NoDrinkingHabit  DrinkingHabit = "no"
)

type SmokingHabit string

const (
	YesSmokingHabit SmokingHabit = "yes"
	NoSmokingHabit  SmokingHabit = "no"
)

type DrugHabit string

const (
	YesDrugHabit DrugHabit = "yes"
	NoDrugHabit  DrugHabit = "no"
)

type Religion string

const (
	Christian     Religion = "christian"
	Muslim        Religion = "muslim"
	Hindu         Religion = "hindu"
	Buddhist      Religion = "buddhist"
	OtherReligion Religion = "other"
)

type Intentions string

const (
	LifePartner    Intentions = "life partner"
	ShorterTime    Intentions = "shorter time"
	NoneIntention  Intentions = "none"
	FiguringOut    Intentions = "figuring out"
	OtherIntention Intentions = "other"
)

type EducationLevel string

const (
	HighSchool       EducationLevel = "High School"
	College          EducationLevel = "College"
	AssociatesDegree EducationLevel = "Associates Degree"
	BachelorsDegree  EducationLevel = "Bachelors Degree"
	MastersDegree    EducationLevel = "Masters Degree"
	PhDPostDoctoral  EducationLevel = "PhD/Post Doctoral"
	OpenToAll        EducationLevel = "Open to All"
)

type AgeRange struct {
	Min         int  `json:"min"`
	Max         int  `json:"max"`
	DealBreaker bool `json:"deal_breaker"`
}

type HeightRange struct {
	Min         int  `json:"min"`
	Max         int  `json:"max"`
	DealBreaker bool `json:"deal_breaker"`
}

type SmokingPreference struct {
	DealBreaker bool         `json:"deal_breaker"`
	Status      SmokingHabit `json:"status"`
}

type DrinkingPreference struct {
	DealBreaker bool          `json:"deal_breaker"`
	Status      DrinkingHabit `json:"status"`
}

type DrugPreference struct {
	DealBreaker bool      `json:"deal_breaker"`
	Status      DrugHabit `json:"status"`
}

type Preferences struct {
	InterestedIn Gender             `json:"interested_in"` // men, women, non-binary, everyone
	MaxDistance  int                `json:"max_distance"`
	AgeRange     AgeRange           `json:"age_range"`
	Ethnicity    Ethnicity          `json:"ethnicity"`
	Religion     Religion           `json:"religion"`
	Height       HeightRange        `json:"height"`
	Children     string             `json:"children"`     // Doesn't have children, has children, open to all
	Drinking     DrinkingPreference `json:"drinking"`     // Never, sometimes, often, open to all
	FamilyPlans  string             `json:"family_plans"` // Doesn't want children, wants children, open to all, not sure yet, might want children
	Drugs        DrugPreference     `json:"drugs"`        // Never, sometimes, often, open to all
	Smoking      SmokingPreference  `json:"smoking"`      // Never, sometimes, often, open to all
	Education    EducationLevel     `json:"education"`    // High School, Some College, Associates Degree, Bachelors Degree, Masters Degree, PhD/Post Doctoral, open to all
}
