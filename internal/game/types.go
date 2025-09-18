package game

import (
	"bean-jump/internal/physics"
	"math"
	"math/rand"
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
	PlatformWidth            = 10.0
	PlatformHeight           = 1.0
	MinPlatformVerticalGap   = 5.0  // Minimum vertical distance between platforms
	MaxPlatformVerticalGap   = 28.0 // Maximum vertical distance between platforms
	MaxPlatformHorizontalGap = 25.0 // Maximum horizontal distance player can jump
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
	Boost    int
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
		Camera:        physics.Vector2{X: 0, Y: 0},
		Score:         0,
		GameOver:      false,
		HighestY:      PlayerStartY,
		LeftKeyTime:   999, // Start with high values (no recent input)
		RightKeyTime:  999,
		InputMomentum: 0,
		Platforms: []Platform{
			// Starting platform
			{
				Position: physics.Vector2{X: GameWidth/2 - PlatformWidth/2, Y: PlayerStartY + 5},
				Width:    PlatformWidth,
				Height:   PlatformHeight,
			},
		},
	}
}

func generateNewPlatforms(count int, last Platform) []Platform {
	output := make([]Platform, 0, count)

	// Add more platforms above using smart placement
	currentY := last.Position.Y

	for range count {
		// Random vertical spacing within safe range
		spacing := MinPlatformVerticalGap + rand.Float64()*(MaxPlatformVerticalGap-MinPlatformVerticalGap)
		currentY -= spacing

		// Generate new X position that's reachable from last platform
		var newX float64
		maxAttempts := 10
		for attempt := range maxAttempts {
			candidateX := rand.Float64() * (GameWidth - PlatformWidth)

			// Check horizontal distance (accounting for screen wrapping)
			horizontalDist := math.Abs(candidateX + PlatformWidth/2 - (last.Position.X + PlatformWidth/2))
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
				maxOffset := MaxPlatformHorizontalGap * 0.7  // be conservative
				offset := (rand.Float64()*2 - 1) * maxOffset // -maxOffset to +maxOffset
				newX = last.Position.X + offset

				// Keep within bounds
				newX = max(newX, 0)
				newX = min(newX, GameWidth-PlatformWidth)
			}
		}

		boost := 0
		odds := 5 + int(math.Abs(math.Floor(currentY/100))) // 1/5 chance of getting zero (ie. no boost) and the odds get worse as the player goes up
		if rand.Intn(odds) == 0 {
			const lambda = 0.5
			// Generate random variable on an exponential curve.
			u := rand.Float64()
			x := -math.Log(1-u) / lambda

			// Convert to int
			boost = int(math.Floor(x)) % 9 // mod 9 so that we never get a value higher than 8 whilst preserving the exponential bias
			boost++                        // increment 1 so that the max value is 9 and the min is 1 (the odds of getting zero are previously determined)
		}

		newPlatform := Platform{
			Position: physics.Vector2{X: newX, Y: currentY},
			Width:    PlatformWidth,
			Height:   PlatformHeight,
			Boost:    boost,
		}

		output = append(output, newPlatform)
		last = newPlatform
	}

	return output
}
