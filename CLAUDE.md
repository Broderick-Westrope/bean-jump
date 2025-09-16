# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

DoodleJump TUI is a terminal-based DoodleJump clone built with Go and Bubble Tea, featuring a sophisticated physics simulation inspired by the go-go-go project. The game maintains a 9:16 portrait aspect ratio and uses advanced input smoothing to handle Bubble Tea's key event limitations.

## Common Commands

### Building and Running
```bash
# Run the game directly
go run cmd/main.go

# Build executable
go build -o doodle-jump cmd/main.go

# Install dependencies
go mod tidy
```

### Development
```bash
# Check compilation without building
go build ./...

# Format code
go fmt ./...
```

## Architecture Overview

### Core Components

**Physics Engine (`/internal/physics/`)**
- `vector.go`: Vector2 math, gravity, and friction calculations
- `collision.go`: Circle-rectangle collision detection for player-platform interactions
- Fixed timestep simulation at 60 FPS (dt = 1.0/60.0)

**Game Logic (`/internal/game/`)**
- `types.go`: Game constants, data structures, and platform generation algorithms
- `physics.go`: Input handling with momentum-based smoothing, collision detection, camera following

**TUI Layer (`/internal/tui/`)**  
- `model.go`: Bubble Tea model handling input, rendering, and game state
- `gameview.go`: 9:16 aspect ratio constrained view with automatic centering
- Runs at 60 FPS using 16ms tick intervals

### Key Architectural Decisions

**Input Smoothing System**: Handles Bubble Tea's intermittent key events using momentum persistence and grace frames. Critical for smooth movement feel.

**Aspect Ratio Management**: GameView component enforces 9:16 portrait ratio regardless of terminal size, maintaining authentic mobile game proportions.

**Platform Generation**: Smart procedural generation with reachability validation, considering screen wrapping and horizontal jump distance constraints.

**Physics Integration**: Separates physics simulation from rendering, allowing smooth 60 FPS updates independent of terminal refresh rates.

## Key Constants and Configuration

Located in `internal/game/types.go`:

**Physics Tuning**:
- `Gravity = 150.0`: Downward acceleration
- `JumpSpeed = -100.0`: Upward velocity when hitting platforms  
- `MoveAcceleration = 250.0`: Horizontal acceleration rate
- `InitialMoveSpeed = 5.0`: Instant responsiveness boost

**Input Smoothing**:
- `InputGraceFrames = 5`: Frames to maintain movement after key release
- `MomentumDecay = 0.85`: How quickly momentum decays per frame

**Platform Generation**:
- `MinPlatformVerticalGap = 10.0` to `MaxPlatformVerticalGap = 25.0`
- `MaxPlatformHorizontalGap = 25.0`: Must be reachable considering player movement speed

## Bubble Tea Considerations

The game works around Bubble Tea's key event limitations where holding a key generates intermittent events rather than continuous input. The input system uses:

1. **Frame-based tracking**: `LeftKeyTime`/`RightKeyTime` count frames since last key press
2. **Momentum persistence**: `InputMomentum` maintains movement intent during input gaps  
3. **Grace periods**: Continue acceleration for several frames after key "release"

When modifying input handling, always test with held keys to ensure smooth movement without stuttering or unwanted friction application.

## Physics Simulation Details

The game uses a custom 2D physics engine adapted from go-go-go's collision system:

- **Gravity**: Applied every frame as downward acceleration
- **Collision**: Circle (player) vs Rectangle (platform) with top-landing detection
- **Camera**: Follows player upward only, never downward
- **Screen Wrapping**: Player wraps horizontally at screen edges
- **Platform Cleanup**: Removes platforms far below camera for performance

All physics calculations use consistent units where the game world is 80x120 units, mapped to the terminal's 9:16 aspect ratio view.