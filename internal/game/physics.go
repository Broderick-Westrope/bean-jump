package game

import (
	"bean-jump/internal/physics"
	"math"
)

const dt = 1.0 / 60.0 // 60 FPS simulation

func (g *Game) Update(leftPressed, rightPressed bool) {
	if g.GameOver {
		return
	}

	// Update input timing - track frames since last key press
	if leftPressed {
		g.LeftKeyTime = 0     // Reset timer - key was just pressed/held
		g.InputMomentum = 1.0 // Full momentum when key is pressed
	} else {
		g.LeftKeyTime++ // Increment frames since last left press
	}

	if rightPressed {
		g.RightKeyTime = 0    // Reset timer - key was just pressed/held
		g.InputMomentum = 1.0 // Full momentum when key is pressed
	} else {
		g.RightKeyTime++ // Increment frames since last right press
	}

	// Determine effective input state using momentum and grace period
	effectiveLeft := leftPressed || (g.LeftKeyTime < InputGraceFrames && g.InputMomentum > 0.1)
	effectiveRight := rightPressed || (g.RightKeyTime < InputGraceFrames && g.InputMomentum > 0.1)

	// Handle horizontal movement with smoothed input
	if effectiveLeft {
		// Apply initial boost only on actual first press (not grace period)
		if leftPressed {
			g.Player.Velocity.X = min(g.Player.Velocity.X, -InitialMoveSpeed)
		}

		// Apply acceleration
		g.Player.Velocity.X -= MoveAcceleration * dt
		g.Player.Velocity.X = max(g.Player.Velocity.X, -MaxMoveSpeed)

	}
	if effectiveRight {
		// Apply initial boost only on actual first press (not grace period)
		if rightPressed {
			g.Player.Velocity.X = max(g.Player.Velocity.X, InitialMoveSpeed)
		}

		// Apply acceleration
		g.Player.Velocity.X += MoveAcceleration * dt
		g.Player.Velocity.X = min(g.Player.Velocity.X, MaxMoveSpeed)

	}
	if !effectiveLeft && !effectiveRight {
		// Apply friction only when no effective input
		g.Player.Velocity = physics.ApplyFriction(g.Player.Velocity, AirFriction, dt)
	}

	// Decay momentum when no keys are pressed
	if !leftPressed && !rightPressed {
		g.InputMomentum *= MomentumDecay
	}

	// Apply gravity
	g.Player.Velocity = physics.ApplyGravity(g.Player.Velocity, Gravity, dt)

	// Update position
	g.Player.Position.X += g.Player.Velocity.X * dt
	g.Player.Position.Y += g.Player.Velocity.Y * dt

	// Wrap around screen horizontally
	if g.Player.Position.X < -g.Player.Radius {
		g.Player.Position.X = GameWidth + g.Player.Radius
	} else if g.Player.Position.X > GameWidth+g.Player.Radius {
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
	if g.Player.Position.Y > g.Camera.Y+GameHeight+20 {
		g.GameOver = true
	}

	// Generate new platforms if needed
	g.platformMaintenance()
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
				if platform.Boost != 0 {
					g.Player.Velocity.Y -= float64(35 * platform.Boost)
				}
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

// platformMaintenance removes platforms that are far below the camera and creates new platforms when the camera is nearing the highest existing plaform.
func (g *Game) platformMaintenance() {
	// Remove platforms that are far below camera
	var activePlatforms []Platform
	for _, platform := range g.Platforms {
		if platform.Position.Y < g.Camera.Y+GameHeight+50 {
			activePlatforms = append(activePlatforms, platform)
		}
	}

	// Find the highest platform
	highestPlatformHeight := GameHeight
	var highestPlatform Platform
	for _, platform := range activePlatforms {
		if platform.Position.Y < highestPlatformHeight {
			highestPlatformHeight = platform.Position.Y
			highestPlatform = platform
		}
	}

	// If the highest platform is getting close to camera view, generate 5 more
	if highestPlatformHeight > g.Camera.Y-MaxPlatformVerticalGap*5 {
		activePlatforms = append(activePlatforms, generateNewPlatforms(5, highestPlatform)...)
	}

	g.Platforms = activePlatforms
}
