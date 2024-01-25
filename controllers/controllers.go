package controllers

import (
	"api/interceptors"
	"api/models"
	"api/setup"
	"api/utils"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Controller struct {
	*setup.ServiceDependencies
}

func NewController(opts *setup.ServiceDependencies) *Controller {
	return &Controller{opts}
}

// Home godoc
// @Summary  Welcome to the API
// @Description Welcome to the API
// @Produce			application/json
// @Tags   home
// @Accept   json
// @Success  200 {object} string
// @Failure  400 {object} controllers.ErrorResponse{}
// @Router   / [GET]
func (c *Controller) Home(w http.ResponseWriter, r *http.Request) {
	HttpResponse(w, nil, "Welcome to the API", 0)
}

// RegisterUser godoc
// @Summary  Create a random user
// @Description Generate a random user and return the user details
// @Produce			application/json
// @Tags   user
// @Accept   json
// @Success 200 {object} models.RegistrationResponse{} "Successful response"
// @Failure  400 {object} controllers.ErrorResponse{}
// @Router   /user/create [POST]
func (c *Controller) RegisterUser(w http.ResponseWriter, r *http.Request) {
	registrationResponse, err := c.UserService.Register(r.Context())
	HttpResponse(w, err, registrationResponse, 0)
	return
}

// LoginUser godoc
// @Summary  Login a user with email and password
// @Description Login a user with email and password
// @Produce			application/json
// @Tags   user
// @Accept   json
// @Param			user body models.LoginPayload{} true "Login Payload"
// @Success  200 {object} models.LoginResponse{}
// @Failure  400 {object} controllers.ErrorResponse{}
// @Router   /login [POST]
func (c *Controller) LoginUser(w http.ResponseWriter, r *http.Request) {
	var payload models.LoginPayload // Declare payload as a non-pointer
	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		fmt.Println(err)
		HttpResponse(w, errors.New("invalid email and password"), nil, 0)
		return
	}
	err = payload.Validate()
	if err != nil {
		fmt.Println(err)
		HttpResponse(w, err, nil, 0)
		return
	}
	loginResponse, err := c.UserService.Login(r.Context(), payload.Email, payload.Password)
	HttpResponse(w, err, loginResponse, 0)
}

// GetUser godoc
// @Summary  Get a user
// @Description Get a user
// @Produce			application/json
// @Tags   user
// @Accept   json
// @Security BearerToken
// @Param Authorization header string true "Bearer Token" default(bearer)
// @Success  200 {object} models.User{}
// @Failure  400 {object} controllers.ErrorResponse{}
// @Router   /user [GET]
func (c *Controller) GetUser(w http.ResponseWriter, r *http.Request) {
	account, err := interceptors.GetAuthenticatedAccount(r.Context())
	if err != nil {
		HttpResponse(w, errors.New("unauthorized account"), nil, 401)
		return
	}
	account.Age, _ = utils.CalculateAge(account.DateOfBirth)
	HttpResponse(w, err, account, 0)
	return
}

// DiscoverUsers godoc
// @Summary  Discover users
// @Description Discover users
// @Produce			application/json
// @Tags   discover
// @Accept   json
// @Security BearerToken
// @Param Authorization header string true "Bearer Token" default(bearer)
// @Param min_age query int false "Minimum Age"
// @Param max_age query int false "Maximum Age"
// @Param max_distance query int false "Maximum Distance"
// @Success  200 {object} []models.User{}
// @Failure  400 {object} controllers.ErrorResponse{}
// @Router   /discover [GET]
func (c *Controller) DiscoverUsers(w http.ResponseWriter, r *http.Request) {
	account, err := interceptors.GetAuthenticatedAccount(r.Context())
	if err != nil {
		HttpResponse(w, errors.New("unauthorized account"), nil, 401)
		return
	}
	var filter models.UserFilter
	q := r.URL.Query()
	filter.MinHeight, _ = strconv.Atoi(q.Get("min_height"))
	filter.MaxHeight, _ = strconv.Atoi(q.Get("max_height"))
	filter.MinAge, _ = strconv.Atoi(q.Get("min_age"))
	filter.MaxAge, _ = strconv.Atoi(q.Get("max_age"))
	//filter.Latitude, _ = strconv.ParseFloat(q.Get("latitude"), 64)
	//filter.Longitude, _ = strconv.ParseFloat(q.Get("longitude"), 64)
	filter.MaxDistance, _ = strconv.Atoi(q.Get("max_distance"))
	filter.DesiredEthnicity = strings.ToLower(q.Get("desired_ethnicity"))
	filter.DesiredPets = strings.ToLower(q.Get("desired_pets"))
	filter.DesiredSexuality = strings.ToLower(q.Get("desired_sexuality"))
	filter.DesiredDrinking = strings.ToLower(q.Get("desired_drinking"))
	filter.DesiredSmoking = strings.ToLower(q.Get("desired_smoking"))
	filter.DesiredDrugs = strings.ToLower(q.Get("desired_drugs"))
	filter.DesiredIntentions = strings.ToLower(q.Get("desired_intentions"))
	filter.DesiredReligion = strings.ToLower(q.Get("desired_religion"))
	if err := filter.Validate(); err != nil {
		HttpResponse(w, err, nil, http.StatusBadRequest)
		return
	}
	profiles, err := c.UserService.Discover(r.Context(), *account, filter)
	HttpResponse(w, err, profiles, 0)
	return
}

// SwipeUser godoc
// @Summary  Swipe a user profile
// @Description Swipe a user profile
// @Produce			application/json
// @Tags   swipe
// @Accept   json
// @Security BearerToken
// @Param Authorization header string true "Bearer Token" default(bearer)
// @Param			user body models.SwipePayload{} true "Login Payload"
// @Success  200 {object} []models.SwipeResponse{}
// @Failure  400 {object} controllers.ErrorResponse{}
func (c *Controller) SwipeUser(w http.ResponseWriter, r *http.Request) {
	account, err := interceptors.GetAuthenticatedAccount(r.Context())
	if err != nil {
		HttpResponse(w, errors.New("unauthorized account"), nil, 401)
		return
	}
	var payload models.SwipePayload
	err = json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		HttpResponse(w, errors.New("invalid payload"), nil, 400)
		return
	}
	err = payload.Validate()
	if err != nil {
		HttpResponse(w, err, nil, 400)
		return
	}
	swipeResponse, err := c.SwipeService.Swipe(r.Context(), account.ID, payload)
	HttpResponse(w, err, swipeResponse, 0)
	return
}
