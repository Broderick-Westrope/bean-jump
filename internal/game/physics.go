package game

import (
	"doodle-jump/internal/physics"
	"math"
	"math/rand"
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
	var highestPlatformObj Platform
	for _, platform := range g.Platforms {
		if platform.Position.Y < highestPlatform {
			highestPlatform = platform.Position.Y
			highestPlatformObj = platform
		}
	}

	// If the highest platform is getting close to camera view, generate more
	if highestPlatform > g.Camera.Y-MaxPlatformVerticalGap*5 {
		// Add more platforms above using smart placement
		currentY := highestPlatform
		lastPlatform := highestPlatformObj

		for i := 0; i < 5; i++ {
			// Random vertical spacing within safe range
			spacing := MinPlatformVerticalGap + rand.Float64()*(MaxPlatformVerticalGap-MinPlatformVerticalGap)
			currentY -= spacing

			// Generate new X position that's reachable from last platform
			var newX float64
			maxAttempts := 10

			for attempt := 0; attempt < maxAttempts; attempt++ {
				candidateX := rand.Float64() * (GameWidth - PlatformWidth)

				// Check horizontal distance (accounting for screen wrapping)
				horizontalDist := math.Abs(candidateX + PlatformWidth/2 - (lastPlatform.Position.X + PlatformWidth/2))
				wrapDist := GameWidth - horizontalDist
				if wrapDist < horizontalDist {
					horizontalDist = wrapDist
				}

				if horizontalDist <= MaxPlatformHorizontalGap {
					newX = candidateX
					break
				}

				// Fallback: place within safe distance
				if attempt == maxAttempts-1 {
					maxOffset := MaxPlatformHorizontalGap * 0.7
					offset := (rand.Float64()*2 - 1) * maxOffset
					newX = lastPlatform.Position.X + offset

					// Keep within bounds
					if newX < 0 {
						newX = 0
					} else if newX > GameWidth-PlatformWidth {
						newX = GameWidth - PlatformWidth
					}
				}
			}

			newPlatform := Platform{
				Position: physics.Vector2{X: newX, Y: currentY},
				Width:    PlatformWidth,
				Height:   PlatformHeight,
			}

			g.Platforms = append(g.Platforms, newPlatform)
			lastPlatform = newPlatform
		}
	}

	// Remove platforms that are far below camera
	var activePlatforms []Platform
	for _, platform := range g.Platforms {
		if platform.Position.Y < g.Camera.Y+GameHeight+50 {
			activePlatforms = append(activePlatforms, platform)
		}
	}
	g.Platforms = activePlatforms
}
