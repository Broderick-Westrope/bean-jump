package game

import (
	"doodle-jump/internal/physics"
	"math/rand"
	"time"
)

const (
	// Game constants
	GameWidth  = 80.0
	GameHeight = 120.0

	// Physics constants
	Gravity          = 150.0
	JumpSpeed        = -100.0
	MaxMoveSpeed     = 25.0  // Maximum horizontal velocity
	MoveAcceleration = 250.0 // How quickly player accelerates horizontally
	InitialMoveSpeed = 5.0   // Instant velocity boost when first pressing key
	AirFriction      = 30.0  // Friction when no input is pressed

	// Input smoothing constants
	InputGraceFrames = 5    // Frames to maintain input momentum after key release
	MomentumDecay    = 0.85 // How quickly momentum decays per frame

	// Player constants
	PlayerRadius = 1.5
	PlayerStartX = GameWidth / 2
	PlayerStartY = GameHeight - 10

	// Platform constants
	PlatformWidth            = 12.0
	PlatformHeight           = 1.0
	MinPlatformVerticalGap   = 10.0 // Minimum vertical distance between platforms
	MaxPlatformVerticalGap   = 25.0 // Maximum vertical distance between platforms
	MaxPlatformHorizontalGap = 25.0 // Maximum horizontal distance player can jump
	MaxPlatforms             = 20
)

type Player struct {
	Position physics.Vector2
	Velocity physics.Vector2
	Radius   float64
}

type Platform struct {
	Position physics.Vector2
	Width    float64
	Height   float64
}

func (p Platform) ToRect() physics.Rect {
	return physics.Rect{
		X:      p.Position.X,
		Y:      p.Position.Y,
		Width:  p.Width,
		Height: p.Height,
	}
}

type Game struct {
	Player        Player
	Platforms     []Platform
	Camera        physics.Vector2
	Score         int
	GameOver      bool
	HighestY      float64
	LeftKeyTime   int     // Frames since left key was last pressed (0 = just pressed)
	RightKeyTime  int     // Frames since right key was last pressed (0 = just pressed)
	InputMomentum float64 // Momentum from recent input to smooth over gaps
}

func NewGame() *Game {
	return &Game{
		Player: Player{
			Position: physics.Vector2{X: PlayerStartX, Y: PlayerStartY},
			Velocity: physics.Vector2{X: 0, Y: 0},
			Radius:   PlayerRadius,
		},
		Platforms:     generateInitialPlatforms(),
		Camera:        physics.Vector2{X: 0, Y: 0},
		Score:         0,
		GameOver:      false,
		HighestY:      PlayerStartY,
		LeftKeyTime:   999, // Start with high values (no recent input)
		RightKeyTime:  999,
		InputMomentum: 0,
	}
}

func generateInitialPlatforms() []Platform {
	rand.Seed(time.Now().UnixNano())

	platforms := []Platform{
		// Starting platform
		{
			Position: physics.Vector2{X: GameWidth/2 - PlatformWidth/2, Y: PlayerStartY + 5},
			Width:    PlatformWidth,
			Height:   PlatformHeight,
		},
	}

	// Generate platforms going upward with smart placement
	currentY := PlayerStartY
	for i := 1; i < MaxPlatforms; i++ {
		// Random vertical spacing within safe range
		spacing := MinPlatformVerticalGap + rand.Float64()*(MaxPlatformVerticalGap-MinPlatformVerticalGap)
		currentY -= spacing

		// Generate candidate X positions and pick one that's reachable
		lastPlatform := platforms[len(platforms)-1]
		var newX float64

		// Try to place platform within reasonable horizontal distance
		maxAttempts := 10
		for attempt := 0; attempt < maxAttempts; attempt++ {
			// Random X position
			candidateX := rand.Float64() * (GameWidth - PlatformWidth)

			// Check if horizontal distance is reasonable
			horizontalDist := abs(candidateX + PlatformWidth/2 - (lastPlatform.Position.X + PlatformWidth/2))

			// Account for screen wrapping - player can wrap around
			wrapDist := GameWidth - horizontalDist
			if wrapDist < horizontalDist {
				horizontalDist = wrapDist
			}

			if horizontalDist <= MaxPlatformHorizontalGap {
				newX = candidateX
				break
			}

			// If we can't find a good spot, place it closer to the last platform
			if attempt == maxAttempts-1 {
				// Place within safe horizontal distance
				maxOffset := MaxPlatformHorizontalGap * 0.7  // Be conservative
				offset := (rand.Float64()*2 - 1) * maxOffset // -maxOffset to +maxOffset
				newX = lastPlatform.Position.X + offset

				// Keep within bounds
				if newX < 0 {
					newX = 0
				} else if newX > GameWidth-PlatformWidth {
					newX = GameWidth - PlatformWidth
				}
			}
		}

		platforms = append(platforms, Platform{
			Position: physics.Vector2{X: newX, Y: currentY},
			Width:    PlatformWidth,
			Height:   PlatformHeight,
		})
	}

	return platforms
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
