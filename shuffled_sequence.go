package main

import "math/rand"

type ShuffledSequence struct {
	i      int
	values []int
}

func newShuffledSequence(size int) ShuffledSequence {
	values := make([]int, size)
	for i := range values {
		values[i] = i
	}
	rand.Shuffle(size, func(i, j int) {
		values[i], values[j] = values[j], values[i]
	})
	return ShuffledSequence{
		i:      0,
		values: values,
	}
}
