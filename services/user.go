package services

import (
	"api/constants"
	"api/models"
	"api/repository"
	"api/store"
	"api/utils"
	"context"
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/gbrlsnchs/jwt/v3"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

var (
	ErrDuplicateProfile       = errors.New("sorry, account already exists")
	ErrCreateProfileFailed    = errors.New("sorry, failed to create profile")
	ErrProfileNotFoundById    = errors.New("sorry, account not found by id")
	ErrProfileNotFoundByEmail = errors.New("sorry, account not found by email")
	ErrProfileNotFoundByPhone = errors.New("sorry, account not found by phone")
	ErrGenerateTokenFailed    = errors.New("sorry, failed to generate token")
	ErrCalculateAgeFailed     = errors.New("sorry, failed to calculate age")
)

type UserService struct {
	eventStore     store.EventStore
	userRepository repository.UserRepository
	logger         *logrus.Logger
	jwtSecret      string
}

type TokenPayload struct {
	Id string `json:"id"`
	jwt.Payload
}

func NewUserService(eventStore store.EventStore, userRepository repository.UserRepository, logger *logrus.Logger, jwtSecret string) *UserService {
	return &UserService{
		eventStore:     eventStore,
		userRepository: userRepository,
		logger:         logger,
		jwtSecret:      jwtSecret,
	}
}

// Register registers a new user.
func (u UserService) Register(
	ctx context.Context,
) (*models.RegistrationResponse, error) {
	newUser := generateRandomUserPayload()
	createdUser, err := u.userRepository.CreateUser(ctx, newUser)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to create user")
		if errors.Is(err, repository.ErrDuplicateFound) {
			return nil, ErrDuplicateProfile
		}
		return nil, ErrCreateProfileFailed
	}

	// Generate and set token for the created user
	_, err = u.GenerateToken(createdUser)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to generate token for user")
		return nil, ErrGenerateTokenFailed
	}

	age, err := utils.CalculateAge(newUser.DateOfBirth)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to calculate age")
		return nil, ErrCalculateAgeFailed
	}

	return &models.RegistrationResponse{
		ID:       createdUser.ID,
		Name:     newUser.Name,
		Password: constants.DefaultPassword,
		Age:      age,
		Gender:   newUser.Gender,
		Email:    newUser.Email,
	}, nil
}

// GenerateToken generates a new token for the given user.
func (u UserService) GenerateToken(profile *models.User) (*string, error) {
	exp := time.Now().Add(constants.DefaultTokenExpires)
	payload := &TokenPayload{
		Payload: jwt.Payload{
			Issuer:         constants.TokenIssuer,
			Subject:        constants.TokenSubject,
			Audience:       jwt.Audience{constants.TokenAudience},
			IssuedAt:       jwt.NumericDate(time.Now()),
			ExpirationTime: jwt.NumericDate(exp),
			JWTID:          constants.TokenJWTID,
		},
		Id: profile.ID,
	}
	token, err := jwt.Sign(payload, jwt.NewHS256([]byte(u.jwtSecret)))
	if err != nil {
		u.logger.WithError(err).Error("failed to generate token")
		return nil, err
	}
	tokenString := string(token)
	return &tokenString, nil
}

// Login logs in a user.
func (u UserService) Login(ctx context.Context, email, password string) (*models.LoginResponse, error) {
	lowercaseEmail := strings.ToLower(email)
	profile, err := u.userRepository.GetUserByEmail(ctx, lowercaseEmail)
	if err != nil {
		return nil, errors.New("invalid email or password. Please try again")
	}
	if !utils.VerifyPasscode(profile.Password, password) {
		return nil, errors.New("invalid email or password. Please try again")
	}
	token, err := u.GenerateToken(profile)
	if err != nil {
		return nil, errors.New("failed to generate authentication token")
	}
	response := &models.LoginResponse{
		Token: *token,
	}
	return response, nil
}

// GetProfile returns the profile of the given user.
func (u UserService) GetProfile(ctx context.Context, id string) (*models.User, error) {
	profile, err := u.userRepository.GetUserById(ctx, id)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to get user profile by ID")
		return nil, ErrProfileNotFoundById
	}
	return profile, nil
}

// Discover returns a list of profiles that match the given filter.
func (u UserService) Discover(ctx context.Context, user models.User, filter models.UserFilter) ([]*models.User, error) {
	profiles, err := u.userRepository.Discover(ctx, filter, user)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to discover profiles")
		return nil, errors.New("failed to discover profiles")
	}
	if len(profiles) == 0 {
		return []*models.User{}, nil
	}
	return profiles, nil
}

// VerifyAuthToken verifies the given token.
func (u UserService) VerifyAuthToken(ctx context.Context, token string) (*models.User, error) {
	secret := jwt.NewHS256([]byte(u.jwtSecret))
	var payloadBody TokenPayload
	_, err := jwt.Verify([]byte(token), secret, &payloadBody)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to verify token")
		return nil, errors.New("invalid token: verification failed")
	}

	profile, err := u.userRepository.GetUserById(ctx, payloadBody.Id)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to retrieve user from token")
		return nil, errors.New("invalid token: user retrieval failed")
	}
	return profile, nil
}

// SeedDefaultUsers seeds the database with default users.
func (u UserService) SeedDefaultUsers(ctx context.Context) error {
	currentUserCount, err := u.userRepository.GetUserCount(ctx)
	if err != nil {
		u.logger.WithContext(ctx).WithError(err).Error("failed to get user count")
		return err
	}

	if currentUserCount < constants.DefaultUserCount {
		usersToInsert := make([]*models.User, constants.DefaultUserCount-currentUserCount)
		for i := range usersToInsert {
			usersToInsert[i] = generateRandomUserPayload()
		}
		err := u.userRepository.InsertUsers(ctx, usersToInsert)
		if err != nil {
			u.logger.WithContext(ctx).WithError(err).Error("failed to insert default users")
			return err
		}
		u.logger.WithContext(ctx).Info("successfully seeded default users")
	}
	return nil
}

// generateRandomUserPayload generates a random user payload.
func generateRandomUserPayload() *models.User {
	encryptedPassword := utils.EncryptPassword(constants.DefaultPassword)
	name := faker.FirstName() + " " + faker.LastName()
	lowercaseFirstName := strings.ToLower(faker.FirstName())
	return &models.User{
		Name:             name,
		Password:         encryptedPassword,
		DateOfBirth:      utils.GetRandomDOB().Format("2006-01-02"),
		Gender:           models.Gender(utils.GetRandomGender()),
		Email:            lowercaseFirstName + "@gmail.com",
		Location:         utils.GetRandomLocationInNorthLondon(),
		Ethnicity:        utils.GetRandomEthnicity(),
		Pets:             utils.GetRandomPet(),
		DatingIntentions: utils.GetRandomDatingIntentions(),
		Height:           utils.GetRandomHeight(),
	}
}
