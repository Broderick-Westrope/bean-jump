package game

import (
	"doodle-jump/internal/physics"
	"math"
)

const dt = 1.0 / 60.0 // 60 FPS simulation

func (g *Game) Update(leftPressed, rightPressed bool) {
	if g.GameOver {
		return
	}
	
	// Handle horizontal movement
	if leftPressed {
		g.Player.Velocity.X = -MoveSpeed
	} else if rightPressed {
		g.Player.Velocity.X = MoveSpeed
	} else {
		// Apply air friction when no input
		g.Player.Velocity = physics.ApplyFriction(g.Player.Velocity, AirFriction, dt)
	}
	
	// Apply gravity
	g.Player.Velocity = physics.ApplyGravity(g.Player.Velocity, Gravity, dt)
	
	// Update position
	g.Player.Position.X += g.Player.Velocity.X * dt
	g.Player.Position.Y += g.Player.Velocity.Y * dt
	
	// Wrap around screen horizontally
	if g.Player.Position.X < -g.Player.Radius {
		g.Player.Position.X = GameWidth + g.Player.Radius
	} else if g.Player.Position.X > GameWidth + g.Player.Radius {
		g.Player.Position.X = -g.Player.Radius
	}
	
	// Check platform collisions
	g.checkPlatformCollisions()
	
	// Update camera to follow player
	g.updateCamera()
	
	// Update score based on height
	if g.Player.Position.Y < g.HighestY {
		g.HighestY = g.Player.Position.Y
		g.Score = int(math.Max(0, (PlayerStartY-g.HighestY)/10))
	}
	
	// Check game over (fell too far below camera)
	if g.Player.Position.Y > g.Camera.Y + GameHeight + 20 {
		g.GameOver = true
	}
	
	// Generate new platforms if needed
	g.generatePlatformsIfNeeded()
}

func (g *Game) checkPlatformCollisions() {
	// Only check collisions when falling (positive velocity)
	if g.Player.Velocity.Y <= 0 {
		return
	}
	
	playerCircle := physics.Circle{
		Center: g.Player.Position,
		Radius: g.Player.Radius,
	}
	
	for _, platform := range g.Platforms {
		// Skip platforms that are too far away
		if math.Abs(platform.Position.Y-g.Player.Position.Y) > 10 {
			continue
		}
		
		if physics.IsLandingOnPlatform(g.Player.Position, g.Player.Velocity.Y, platform.ToRect()) {
			// Check if player is actually colliding with platform
			if physics.CircleRectCollision(playerCircle, platform.ToRect()) {
				// Land on platform - jump!
				g.Player.Velocity.Y = JumpSpeed
				break
			}
		}
	}
}

func (g *Game) updateCamera() {
	// Camera follows player, but only moves up, never down
	targetCameraY := g.Player.Position.Y - GameHeight/2
	if targetCameraY < g.Camera.Y {
		g.Camera.Y = targetCameraY
	}
}

func (g *Game) generatePlatformsIfNeeded() {
	// Find the highest platform
	highestPlatform := GameHeight
	for _, platform := range g.Platforms {
		if platform.Position.Y < highestPlatform {
			highestPlatform = platform.Position.Y
		}
	}
	
	// If the highest platform is getting close to camera view, generate more
	if highestPlatform > g.Camera.Y - PlatformSpacing*5 {
		// Add more platforms above
		for i := 0; i < 5; i++ {
			newY := highestPlatform - PlatformSpacing*float64(i+1)
			newX := float64((len(g.Platforms)+i)*11 % int(GameWidth-PlatformWidth))
			
			g.Platforms = append(g.Platforms, Platform{
				Position: physics.Vector2{X: newX, Y: newY},
				Width:    PlatformWidth,
				Height:   PlatformHeight,
			})
		}
	}
	
	// Remove platforms that are far below camera
	var activePlatforms []Platform
	for _, platform := range g.Platforms {
		if platform.Position.Y < g.Camera.Y + GameHeight + 50 {
			activePlatforms = append(activePlatforms, platform)
		}
	}
	g.Platforms = activePlatforms
}