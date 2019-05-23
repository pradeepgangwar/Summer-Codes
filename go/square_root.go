// This implements the algorithm to calculate the square root of a number
package main

import (
	"fmt"
	"math"
)

func Sqrt(x float64) float64 {
	var change, z float64
	var i int
	z = float64(1)
	change = (z*z -x)/(2*z)
	for i=0; math.Abs(change) > 0.001; i++ {
		z-=change
		change = (z*z -x)/(2*z)
	}
	fmt.Println(i)
	return z
}

func main() {
	fmt.Println(Sqrt(50))
}
