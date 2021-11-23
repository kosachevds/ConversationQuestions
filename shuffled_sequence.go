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

func (ss *ShuffledSequence) Next() int {
	result := ss.values[ss.i]
	ss.i = (ss.i + 1) % len(ss.values)
	return result
}

func (ss *ShuffledSequence) NextFrom(values []interface{}) interface{} {
	return values[ss.Next()]
}
