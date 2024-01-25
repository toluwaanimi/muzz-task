package utils

import (
	"api/models"
	"math"
)

const earthRadiusKm = 6371 // Earth's radius in kilometers
const swipeCost = 0.5

func calculateCompatibilityScore(userPreference *models.Preferences, match *models.User) float64 {
	compatibilityScore, totalWeight := 0.0, 0.0
	dealBreakerFailed := false

	drinkWeight, smokeWeight, religionWeight := 0.3, 0.2, 0.5

	// Drinking preference
	if userPreference.Drinking.Status == match.Drinking {
		compatibilityScore += drinkWeight
	} else if userPreference.Drinking.DealBreaker {
		dealBreakerFailed = true
	}
	totalWeight += drinkWeight

	// Smoking preference
	if userPreference.Smoking.Status == match.Smoking {
		compatibilityScore += smokeWeight
	} else if userPreference.Smoking.DealBreaker {
		dealBreakerFailed = true
	}
	totalWeight += smokeWeight

	// Religion preference
	if userPreference.Religion == match.Religion {
		compatibilityScore += religionWeight
	} else if userPreference.Religion == models.OtherReligion && match.Religion != models.OtherReligion {
		dealBreakerFailed = true
	}
	totalWeight += religionWeight

	// Normalize the compatibility score
	normalizedScore := compatibilityScore / totalWeight

	// If any dealbreaker condition failed, return a reduced score
	if dealBreakerFailed {
		return normalizedScore * 0.5 // Example: Reduce score by half if dealbreaker is not met
	}

	return normalizedScore
}

func toRadians(deg float64) float64 {
	return deg * math.Pi / 180
}

func calculateProximityScore(userLocation, viewedUserLocation []float64) float64 {
	lat1, lon1 := toRadians(userLocation[0]), toRadians(userLocation[1])
	lat2, lon2 := toRadians(viewedUserLocation[0]), toRadians(viewedUserLocation[1])

	dlat := lat2 - lat1
	dlon := lon2 - lon1

	a := math.Pow(math.Sin(dlat/2), 2) + math.Cos(lat1)*math.Cos(lat2)*math.Pow(math.Sin(dlon/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := earthRadiusKm * c
	return 1 / (1 + distance)
}

func calculateSwipeCostScore(user models.User) float64 {
	if user.SwipeCount > user.DailySwipeBudget {
		return float64(user.SwipeCount-user.DailySwipeBudget) * swipeCost
	}
	return 0
}

func calculateSwipeScore(viewingUser, viewedUser models.User, compatibilityScore float64) float64 {
	// Adjust weights if necessary
	proximityWeight := 0.25
	attractivenessWeight := 0.25
	swipeCostWeight := 0.25
	compatibilityWeight := 0.25

	// normalization scales based on the expected data range
	proximityScore := calculateProximityScore(viewingUser.Location, viewedUser.Location) * 2.5
	attractivenessScore := float64(viewedUser.Attractiveness) / 10 * 2.5                //  attractiveness is scored out of 10
	swipeCostScore := math.Min(calculateSwipeCostScore(viewingUser)/swipeCost*2.5, 2.5) // Normalize based on your swipe cost logic
	normalizedCompatibilityScore := compatibilityScore * 2.5                            // Normalize compatibility score

	totalScore :=
		proximityScore*proximityWeight +
			attractivenessScore*attractivenessWeight +
			swipeCostScore*swipeCostWeight +
			normalizedCompatibilityScore*compatibilityWeight

	// Constrain total score to a maximum of 10
	return math.Min(math.Max(totalScore, 0), 10)
}

func updateSwipeRating(currentSwipeRating float64, successfulMatches, unsuccessfulMatches int) float64 {
	dailySuccessRate := float64(successfulMatches) / math.Max(1, float64(unsuccessfulMatches))

	successRateThresholdHigh := 0.8
	successRateThresholdLow := 0.2

	newSwipeRating := currentSwipeRating

	if dailySuccessRate > successRateThresholdHigh {
		newSwipeRating *= 0.8 // Reduce swipe rating for users with high success rates
	} else if dailySuccessRate < successRateThresholdLow {
		newSwipeRating *= 1.2 // Increase swipe rating for users with low success rates
	}

	return newSwipeRating
}

func PerformSwipe(viewingUser *models.User, viewingUserPreference *models.Preferences, viewedUser *models.User, successfulMatches, unsuccessfulMatches int, someThreshold float64) (float64, float64) {
	compatibilityScore := calculateCompatibilityScore(viewingUserPreference, viewedUser)
	swipeScore := calculateSwipeScore(*viewingUser, *viewedUser, compatibilityScore)

	if swipeScore >= someThreshold {
		successfulMatches++
	} else {
		unsuccessfulMatches++
	}

	newSwipeRating := updateSwipeRating(viewingUser.SwipingRate, successfulMatches, unsuccessfulMatches)
	return swipeScore, newSwipeRating
}
