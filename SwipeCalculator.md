
---

# Swipe Matching Calculator

## Overview
This Swipe Matching Calculator is designed to evaluate potential matches in a muzz application.
It calculates compatibility scores between users based on their preferences and characteristics. The core functionality of the system is encapsulated in the `PerformSwipe` function.

## Key Components
The system comprises several functions, each responsible for a specific aspect of the matching process:

- `calculateCompatibilityScore`: Evaluates how well the preferences of one user align with the characteristics of another user.
- `calculateProximityScore`: Computes a score based on the geographical distance between two users.
- `calculateSwipeCostScore`: Determines the cost of swiping, factoring in the user's swipe budget.
- `calculateSwipeScore`: Integrates various scores (compatibility, proximity, attractiveness, and swipe cost) to produce a comprehensive swipe score.
- `updateSwipeRating`: Adjusts the user's swipe rating based on their match success rate.
- `performSwipe`: The main function that orchestrates the swipe process, utilizing the aforementioned functions.

## Process Description
### performSwipe Function
The `PerformSwipe` function is the heart of the swipe calculator. It performs the following steps:

1. **Compatibility Calculation**: It first calls `calculateCompatibilityScore` to assess how well the viewing user's preferences align with the attributes of the viewed user.

2. **Swipe Score Calculation**: Next, it invokes `calculateSwipeScore` to compute a swipe score. This score is a weighted sum of compatibility, proximity, attractiveness, and swipe cost scores.

3. **Match Decision**: The function then compares the swipe score against a predefined threshold. If the score is above this threshold, it counts as a successful match; otherwise, it's an unsuccessful match.

4. **Swipe Rating Update**: Finally, `updateSwipeRating` is called to adjust the viewing user's swipe rating based on their recent match success or failure.

### Return Values
`performSwipe` returns two values:
- The swipe score, indicating the strength of the match.
- The new swipe rating of the viewing user, reflecting their updated match success rate.

## Usage
To use this system, instantiate `User` and `Preferences` objects with appropriate data, and 
then call the `PerformSwipe` function with these objects. 
The function will process the data and return the match results.

---