# DoodleJump TUI: Design Decisions & Challenges

A record of interesting technical decisions and challenges encountered while building a terminal-based DoodleJump clone with Go and Bubble Tea.

## Project Genesis & Inspiration

**The Physics Challenge**: Started from analyzing the sophisticated 2D physics engine in the go-go-go multiplayer stone game, which featured realistic collision detection, impulse-based responses, and smooth stone movement. The question was: "Could similar techniques work for different game genres in terminal environments?"

**Key Insight**: Physics simulation can be completely decoupled from rendering, allowing complex gameplay mechanics even in text-based interfaces.

## Architecture Decisions

### 1. Physics-First Design

**Decision**: Build a real physics engine rather than simple game logic.

**Rationale**: 
- Authentic feel requires gravity, acceleration, and realistic jumping
- Collision detection needed for platform interactions
- Camera following requires smooth mathematical interpolation

**Implementation**: Separated `/internal/physics/` from game logic, enabling 60 FPS simulation independent of terminal refresh rates.

**Result**: Game feels responsive and natural despite being text-based.

---

### 2. The Bubble Tea Input Challenge

**Problem**: Bubble Tea doesn't provide key press/hold/release events. Holding a key generates: press → silence (1 second) → rapid repeated presses every few milliseconds.

**Initial Naive Approach**: Treat each KeyMsg as momentary input
```go
// This creates stuttering movement
if leftPressed {
    velocity.X = -speed
} else {
    // Applies friction immediately when key events pause
    velocity = applyFriction(velocity)
}
```

**The Problem**: During the 1-second gap and between rapid-fire events, `leftPressed` becomes false, causing unwanted friction and movement interruption.

**Solution**: Momentum-Based Input Smoothing
```go
// Track frames since last key press
if leftPressed {
    g.LeftKeyTime = 0
    g.InputMomentum = 1.0
} else {
    g.LeftKeyTime++ // Count frames without input
}

// "Effective" input considers recent history
effectiveLeft := leftPressed || (g.LeftKeyTime < InputGraceFrames && g.InputMomentum > 0.1)
```

**Key Innovation**: Grace period of 5 frames (~83ms) bridges input gaps, while momentum decay prevents infinite movement.

**Result**: Smooth, continuous movement that feels responsive despite Bubble Tea's limitations.

---

### 3. Platform Generation Algorithm

**Challenge**: Generate infinite, random platforms that are always reachable.

**Constraints**:
- Player has limited horizontal movement speed
- Vertical jump height is fixed
- Screen wrapping allows teleporting across edges
- Must account for both direct and wrapped distances

**Algorithm Design**:
1. **Vertical Spacing**: Random between 10-25 units (jump height ~33 units provides safety margin)
2. **Horizontal Reachability**: Calculate both direct distance and wrap-around distance
3. **Smart Fallback**: If random placement fails 10 times, place conservatively near previous platform

```go
// Account for screen wrapping in distance calculation
horizontalDist := abs(candidateX - lastX)
wrapDist := GameWidth - horizontalDist
if wrapDist < horizontalDist {
    horizontalDist = wrapDist // Use shorter wrapped distance
}
```

**Result**: Infinite gameplay with guaranteed reachable platforms, creating challenging but fair jumps.

---

### 4. Aspect Ratio Constraint System

**Design Challenge**: Maintain authentic DoodleJump feel across different terminal sizes.

**Decision**: Enforce 9:16 portrait ratio (like mobile phone) rather than filling entire terminal.

**Implementation**:
```go
// Try both width-constrained and height-constrained sizing
targetRatio := 9.0 / 16.0
widthConstrained := maxWidth
heightFromWidth := int(float64(widthConstrained) / targetRatio)
heightConstrained := maxHeight  
widthFromHeight := int(float64(heightConstrained) * targetRatio)

// Choose the size that fits within terminal bounds
```

**Visual Result**: Game maintains proper proportions with centered layout and borders, creating consistent experience regardless of terminal size.

---

### 5. Movement Feel Tuning

**Challenge**: Balancing responsiveness with realistic physics.

**The Responsiveness Problem**: Pure acceleration from zero feels sluggish for the first few frames.

**Solution**: Hybrid approach
- **Initial Boost**: Instant velocity on first key press (5.0 units)
- **Continued Acceleration**: Gradual buildup to max speed (25.0 units)
- **Smart Detection**: Only apply boost on genuine first press, not during grace period

```go
if leftPressed && g.LeftKeyTime == 0 && abs(g.Player.Velocity.X) < InitialMoveSpeed {
    g.Player.Velocity.X = -InitialMoveSpeed // Instant responsiveness
}
// Then continue with normal acceleration
g.Player.Velocity.X -= MoveAcceleration * dt
```

**Result**: Movement feels immediately responsive while maintaining realistic physics progression.

## Technical Challenges & Solutions

### 1. Frame Rate Independence

**Problem**: Game logic tied to rendering would create inconsistent gameplay on different terminals.

**Solution**: Fixed timestep physics (dt = 1/60) with separate render loop.

### 2. Coordinate System Mapping

**Challenge**: Map abstract game coordinates (80x120) to variable terminal sizes while maintaining aspect ratio.

**Solution**: Dynamic scaling factors with aspect ratio preservation:
```go
scaleX := float64(terminalWidth) / game.GameWidth
scaleY := float64(terminalHeight) / game.GameHeight
```

### 3. Camera Following Logic

**Design Decision**: Camera only moves up, never down (like original DoodleJump).

**Implementation**: 
```go
targetCameraY := g.Player.Position.Y - GameHeight/2
if targetCameraY < g.Camera.Y { // Only move up
    g.Camera.Y = targetCameraY
}
```

**Result**: Creates tension when falling - player disappears below screen, ending game.

## Performance Considerations

### Platform Management
- Generate platforms ahead of camera view
- Remove platforms far below camera to prevent memory growth
- Limit active platforms to ~20-25 for consistent performance

### Collision Optimization
- Skip collision checks for platforms too far from player
- Early exit on first collision found
- Simple circle-rectangle collision for efficiency

## Lessons Learned

### 1. TUI Framework Limitations
Terminal UI frameworks have unique constraints that don't exist in traditional game development. Input handling, in particular, requires creative solutions.

### 2. Physics in Text
Complex physics simulations work well in terminal environments when properly abstracted from rendering.

### 3. User Experience Priorities
Small details like input responsiveness and smooth movement have outsized impact on game feel, even in text-based games.

### 4. Aspect Ratio Matters
Maintaining authentic proportions creates better user experience than maximizing screen usage.

## Potential Talk Themes

1. **"Physics Engines in Unexpected Places"**: How sophisticated physics can enhance any interface
2. **"Working Around Framework Limitations"**: Creative solutions for input handling in TUIs
3. **"Game Feel in Terminal Interfaces"**: What makes terminal games feel responsive and fun
4. **"Architecture for Adaptability"**: Designing systems that work across different constraints
5. **"From Inspiration to Innovation"**: Adapting concepts between different domains (multiplayer stone game → single-player platformer)

## Future Enhancements

- **Power-ups**: Special platforms with different physics properties
- **Particle effects**: Simple ASCII animations for jumps and collisions
- **Sound**: Terminal bell or system sounds for feedback
- **Multiplayer**: Split-screen racing mode
- **Level generation**: Themed platform patterns and challenges