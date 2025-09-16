package game

import "doodle-jump/internal/physics"

const (
	// Game constants
	GameWidth  = 80.0
	GameHeight = 120.0

	// Physics constants
	Gravity     = 150.0
	JumpSpeed   = -100.0
	MoveSpeed   = 30.0
	AirFriction = 20.0

	// Player constants
	PlayerRadius = 1.5
	PlayerStartX = GameWidth / 2
	PlayerStartY = GameHeight - 10

	// Platform constants
	PlatformWidth   = 12.0
	PlatformHeight  = 2.0
	PlatformSpacing = 15.0
	MaxPlatforms    = 20
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
	Player    Player
	Platforms []Platform
	Camera    physics.Vector2
	Score     int
	GameOver  bool
	HighestY  float64
}

func NewGame() *Game {
	return &Game{
		Player: Player{
			Position: physics.Vector2{X: PlayerStartX, Y: PlayerStartY},
			Velocity: physics.Vector2{X: 0, Y: 0},
			Radius:   PlayerRadius,
		},
		Platforms: generateInitialPlatforms(),
		Camera:    physics.Vector2{X: 0, Y: 0},
		Score:     0,
		GameOver:  false,
		HighestY:  PlayerStartY,
	}
}

func generateInitialPlatforms() []Platform {
	platforms := []Platform{
		// Starting platform
		{
			Position: physics.Vector2{X: GameWidth/2 - PlatformWidth/2, Y: PlayerStartY + 5},
			Width:    PlatformWidth,
			Height:   PlatformHeight,
		},
	}

	// Generate platforms going upward
	for i := 1; i < MaxPlatforms; i++ {
		y := PlayerStartY - float64(i)*PlatformSpacing
		x := float64(i * 7 % int(GameWidth-PlatformWidth)) // Simple pattern for demo

		platforms = append(platforms, Platform{
			Position: physics.Vector2{X: x, Y: y},
			Width:    PlatformWidth,
			Height:   PlatformHeight,
		})
	}

	return platforms
}
