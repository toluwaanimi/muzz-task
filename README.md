# Muzz

## Documentation of Muzz

## Author 🚀

> ADEBAYO EMMANUEL TOLUWANIMI

---

## Overview

This submission is for the Muzz Backend Technical Exercise,
where I have built a mini API to power a simple dating app.
The API includes endpoints for user creation, login, discovering
potential matches, and responding to profiles with a swipe action.

## Technologies Used

- **Programming Language**: Go (Golang)
- **Database**: MongoDB
- **Containerization**: Docker

## Features and Endpoints

### 1. User Creation

- **Endpoint**: `/user/create`
- **Functionality**: Generates and stores a new user.
- **Response**:
  ```json
  {
    "results": {
        "id": "01hmz28c17hqhc11cck87rkkaa",
        "name": "Orland Wiegand",
        "email": "elliott@gmail.com",
        "password": "password",
        "age": 28,
        "gender": "male"
      }
  }
  ```

### 2. User Login

- **Endpoint**: `/login`
- **Functionality**: Authenticates a user and returns a token.
- **Request**:
  ```json
   {
      "email":"hallie@gmail.com",
      "password":"password"
   }
  ```
- **Response**:
  ```json
    {
    "results": {
           "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjAxaGt6N3BtamN6ZW00eHE0ajU3YnI5eXpwIiwiaXNzIjoiTXV6eiBEYXRpbmciLCJzdWIiOiJNdXp6IERhdGluZyBUb2tlbiIsImF1ZCI6Imh0dHBzOi8vbXV6ei5jb20iLCJleHAiOjE3MDYyMzEwNjgsImlhdCI6MTcwNjE0NDY2OCwianRpIjoiTXV6eiBEYXRpbmcifQ.FiCqrenBrLlCY_7ejYoJmFPXba_H0Ot2vCWJLehmow8"
         }
    }
  ```

### 3. Discover Potential Matches

- **Endpoint**: `/discover?max_distance=1000&min_age=24&max_age=25`
- **Authentication**: Bearer Token required.
- **Functionality**: Returns profiles of potential matches, excluding already swiped profiles.
- **Response Example**:
  ```json
  {
    "results": [
        {
            "id": "01hkz7pmjd18jxhzvhz0k9mczc",
            "name": "Dereck Lockman",
            "age": 24,
            "date_of_birth": "2000-01-12",
            "location": [
                51.58809654551533,
                -0.020549300462398043
            ],
            "height": 174.23906174067827,
            "ethnicity": "other",
            "gender": "non-binary",
            "distance": 20.561493527208636,
            "pets": "cat",
            "religion": "hindu",
            "drinking": "sometimes",
            "smoking": "yes",
            "drugs": "no",
            "dating_intentions": "life partner"
        },
        {
            "id": "01hkz7pmjc22z0ma37abc1pa71",
            "name": "Devon Goyette",
            "age": 25,
            "date_of_birth": "1999-01-12",
            "location": [
                51.616449857953754,
                0.10056658551719024
            ],
            "height": 186.55899067081998,
            "ethnicity": "black",
            "gender": "female",
            "distance": 33.7201271749244,
            "pets": "cat",
            "religion": "muslim",
            "drinking": "yes",
            "smoking": "yes",
            "drugs": "no",
            "dating_intentions": "none"
        },
    ]
  }
  ```

### 4. Swipe Response

- **Endpoint**: `/swipe`
- **Authentication**: Bearer Token required.
- **Functionality**: Allows users to swipe on other profiles, and returns if there's a match.
- **Request**:
  ```json
  {
      "user_id":"01hkz7pmjd698vqrcvfgsz88e8",
      "interested": true
  }
  ```
- **Response Example**:
  ```json
  {
       "results": {
        "matched": false,
        "match_id": ""
    }
  }
  ```

### 5. Get Current User

- **Endpoint**: `/user`
- **Functionality**: Returns the current user's profile.
- **Authentication**: Bearer Token required.
  - **Response Example**:
    ```json
    {
        "results": {
            "id": "01hkz7ppqwy5xt8j687dsx26vf",
            "name": "Mckayla Gleichner",
            "email": "daphney@gmail.com",
            "age": 42,
            "date_of_birth": "1982-01-12",
            "location": [
            51.62093416548472,
            -0.20231503371018789
            ],
            "height": 186.86204668281977,
            "ethnicity": "other",
            "gender": "female",
            "pets": "dog",
            "religion": "muslim",
            "drinking": "yes",
            "smoking": "none",
            "drugs": "yes",
            "dating_intentions": "other"
        }
    }
    ```
  

## How to Run the Application

### Without Docker Compose

1. **Setup MongoDB**: Ensure MongoDB is installed and running on your system.
2. **Environment Variables**: Set the required environment variables as specified in the `.env.example` file.
3. **Run the Application**: Execute `go run main.go` in the root directory.
4. **Default Port**: The application runs on port `4000` by default, but this can be changed in the `.env` file.
5. **API Documentation**: The API documentation is available at `http://localhost:4000/docs`.


### With Docker Compose

1. **Docker Setup**: Ensure Docker is installed on your system.
2. **Build and Run**: Use `docker-compose up` to build and start the application.
3. **Docker Port**: The application runs on port `8080` by default (mapped to local port 4000), but this can be changed in the `docker-compose.yml` file.

## Folder Structure

- **config**: Contains `.env` configuration.
- **constants**: Stores constant values.
- **controllers**: Houses controller functions for handling API requests.
- **docs**: Swagger setup for API documentation.
- **interceptors**: Implements interceptors for request processing.
- **middlewares**: Contains middleware for request handling.
- **models**: Data structures and models.
- **repository**:
    - **mongo**: MongoDB repository functions.
    - **repository.go**: Interface for database functions.
- **routes**: API route definitions.
- **services**: Business logic implementation.
- **store**: Event consumer implementation.
- **utils**: Utility functions and helpers.
- **setup**: Dependency setup and initialization.

## Example `.env` File

```env
PORT=8080
DATABASE_URL=mongodb://localhost:27017
DATABASE_NAME=muzzdb
JWT_SECRET=your_jwt_secret
CURRENT_DATABASE=mongodb
```

## Notes and Assumptions

## Swipe Score Calculation
In addition to the implemented features, the folder includes a conceptual document `SwipeCalculator.md`. 
This document proposes a methodology for calculating the swipe score, 
a crucial element in determining user compatibility.
It details the factors considered and the algorithm used for scoring user interactions,
providing a theoretical foundation for future enhancements.

- The API is designed for local testing and demonstration.
- Authentication is simplified for the scope of this exercise.
- User passwords are hashed for security.
- MongoDB is used for its flexibility and ease of scaling.

---

This README provides a comprehensive guide on the structure, functionality, and setup of the API for the Muzz Backend
Technical Exercise. It is designed to be clear and concise, providing all necessary information for understanding and
running the application and `SwipeCalculator.md` for future feature development.