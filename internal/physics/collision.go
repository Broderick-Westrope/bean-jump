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

// Check if player trajectory intersects with platform during this frame
// Prevents tunneling through platforms when moving at high speed
func SweptCollisionCheck(prevPos, currentPos Vector2, playerRadius float64, platform Rect) bool {
	// Create player circle at both positions
	prevCircle := Circle{Center: prevPos, Radius: playerRadius}
	currentCircle := Circle{Center: currentPos, Radius: playerRadius}
	
	// If either position collides, we have intersection
	if CircleRectCollision(prevCircle, platform) || CircleRectCollision(currentCircle, platform) {
		return true
	}
	
	// Check if the movement path intersects the platform
	// This handles high-speed tunneling cases
	deltaY := currentPos.Y - prevPos.Y
	
	// Only check for downward movement intersecting top of platform
	if deltaY > 0 {
		// Player is moving down, check if path crosses platform top
		platformTop := platform.Y
		
		// Check if the movement crosses the platform's top edge
		if prevPos.Y <= platformTop && currentPos.Y >= platformTop {
			// Check horizontal bounds at the intersection point
			// Linear interpolation to find X position when crossing platform top
			t := (platformTop - prevPos.Y) / deltaY
			intersectionX := prevPos.X + t*(currentPos.X-prevPos.X)
			
			// Check if intersection point is within platform bounds (accounting for radius)
			return intersectionX >= platform.X-playerRadius && intersectionX <= platform.X+platform.Width+playerRadius
		}
	}
	
	return false
}