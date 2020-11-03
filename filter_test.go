package myrxgo

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDebounce(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Debounce(1000000)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{}, res, "Debounce Test Error!")
}

func TestDistinct(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Distinct()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{10, 20, 30, 40, 50}, res, "Distinct Test Error!")
}

func TestElementAt(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).ElementAt(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{40}, res, "SkipLast Test Error!")
}

func TestFirst(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).First()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{10}, res, "First Test Error!")
}

func TestLast(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Last()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{50}, res, "Last Test Error!")
}

func TestSample(t *testing.T) {
	res := []int{}
	Just(1, 2, 3, 4, 3, 1, 2, 4, 3).Map(func(x int) int {
		time.Sleep(2 * time.Millisecond)
		return x
	}).Sample(5 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{6, 6, 4, 6}, res, "SkipLast Test Error!")

}

func TestSkip(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Skip(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{50}, res, "Skip Test Error!")
}

func TestSkipLast(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).SkipLast(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{10}, res, "SkipLast Test Error!")
}

func TestTake(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Take(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{10, 20, 30, 40}, res, "Take Test Error!")
}

func TestTakeLast(t *testing.T) {
	res := []int{}
	ob := Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).TakeLast(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{20, 30, 40, 50}, res, "TakeLast Test Error!")
}