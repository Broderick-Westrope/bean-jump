package physics

import "math"

type Vector2 struct {
	X, Y float64
}

var zeroVelocity = Vector2{X: 0, Y: 0}

func (v Vector2) isZero() bool {
	return math.Abs(v.X) < 0.01 && math.Abs(v.Y) < 0.01
}

func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

func BlendVector(v1, v2 Vector2, value float64) Vector2 {
	value = math.Min(1.0, math.Max(0.0, value))
	x := v1.X*value + v2.X*(1-value)
	y := v1.Y*value + v2.Y*(1-value)
	return Vector2{X: x, Y: y}
}

func ApplyGravity(velocity Vector2, gravity, dt float64) Vector2 {
	velocity.Y += gravity * dt
	return velocity
}

func ApplyFriction(velocity Vector2, friction, dt float64) Vector2 {
	if friction <= 0 {
		return velocity
	}
	frictionForce := friction * dt
	
	if velocity.X >= 0 {
		velocity.X = math.Max(0, velocity.X-frictionForce)
	} else {
		velocity.X = math.Min(0, velocity.X+frictionForce)
	}
	
	return velocity
}