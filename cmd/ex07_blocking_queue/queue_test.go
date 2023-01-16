package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCyclicQueue(t *testing.T) {
	// Тест на скорую руку, не судите строго)

	var n int
	var ok bool

	q := NewCyclicQueue[int](3)

	_, ok = q.GetFirst()
	assert.False(t, ok)

	_, ok = q.GetFirst()
	assert.False(t, ok)

	assert.True(t, q.AddLast(1))
	assert.True(t, q.AddLast(2))
	assert.True(t, q.AddLast(3))
	assert.False(t, q.AddLast(4))
	assert.False(t, q.AddLast(5))

	n, ok = q.GetFirst()
	assert.True(t, ok)
	assert.Equal(t, 1, n)

	n, ok = q.GetFirst()
	assert.True(t, ok)
	assert.Equal(t, 2, n)

	n, ok = q.GetFirst()
	assert.True(t, ok)
	assert.Equal(t, 3, n)

	_, ok = q.GetFirst()
	assert.False(t, ok)

	_, ok = q.GetFirst()
	assert.False(t, ok)

	assert.True(t, q.AddLast(6))
	assert.True(t, q.AddLast(7))

	n, ok = q.GetFirst()
	assert.True(t, ok)
	assert.Equal(t, 6, n)

	assert.True(t, q.AddLast(8))
	assert.True(t, q.AddLast(9))

	n, ok = q.GetFirst()
	assert.True(t, ok)
	assert.Equal(t, 7, n)

	n, ok = q.GetFirst()
	assert.True(t, ok)
	assert.Equal(t, 8, n)

	n, ok = q.GetFirst()
	assert.True(t, ok)
	assert.Equal(t, 9, n)

	_, ok = q.GetFirst()
	assert.False(t, ok)

	_, ok = q.GetFirst()
	assert.False(t, ok)
}
