package main

import (
	"math/rand/v2"
)

func shoot(aC int, a int, Chancetohit int, damageRange int) (damage int) {
	if aC*2-a+Chancetohit > rand.IntN(100) {
		return 1 + rand.IntN(damageRange-1)
	}
	return 0
}
