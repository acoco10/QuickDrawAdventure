package main

import (
	"testing"
)

func TestShoot(t *testing.T) {
	ac := 3
	a := 2
	cth := 50

	damage := shoot(ac, a, cth, 3)

	if damage < 0 || damage > 3 {
		t.Fatalf(`shoot(ac, a, cth, 3) = %q tdamage to low or too high`, damage)
	}
}
