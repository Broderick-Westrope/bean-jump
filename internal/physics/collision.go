package physics

import "math"

type Rect struct {
	X, Y, Width, Height float64
}

type Circle struct {
	Center Vector2
	Radius float64
}

func CircleRectCollision(circle Circle, rect Rect) bool {
	// Find the closest point on the rectangle to the circle center
	closestX := math.Max(rect.X, math.Min(circle.Center.X, rect.X+rect.Width))
	closestY := math.Max(rect.Y, math.Min(circle.Center.Y, rect.Y+rect.Height))
	
	// Calculate distance from circle center to closest point
	dx := circle.Center.X - closestX
	dy := circle.Center.Y - closestY
	
	return (dx*dx + dy*dy) <= (circle.Radius * circle.Radius)
}

// Check if player is landing on top of platform (not from sides/bottom)
func IsLandingOnPlatform(playerPos Vector2, playerVelY float64, platform Rect) bool {
	// Player must be falling (positive Y velocity in our coordinate system)
	if playerVelY <= 0 {
		return false
	}
	
	// Player must be horizontally within platform bounds
	if playerPos.X < platform.X || playerPos.X > platform.X+platform.Width {
		return false
	}
	
	// Player must be above the platform
	return playerPos.Y <= platform.Y && playerPos.Y >= platform.Y-5
}