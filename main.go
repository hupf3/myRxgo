package main

import (
	"fmt"
	"time"

	rxgo "github.com/hupf3/myRxgo/Rxgo"
)

func main() {
	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Debounce(1000000)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Debounce操作后: ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50, 30, 40")
	res = []int{}
	ob = rxgo.Just(10, 20, 30, 40, 50, 30, 40).Map(func(x int) int {
		return x
	}).Distinct()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Distinct操作后: ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res = []int{}
	ob = rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).ElementAt(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Println("ElementAt操作后: ", res[0])
	fmt.Print("\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res = []int{}
	ob = rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).First()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Println("First操作后: ", res[0])
	fmt.Print("\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res = []int{}
	ob = rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Last()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Println("Last操作后: ", res[0])
	fmt.Print("\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res = []int{}
	rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		time.Sleep(2 * time.Millisecond)
		return x
	}).Sample(3 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Sample操作后: ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res = []int{}
	ob = rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Skip(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Skip操作后: ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res = []int{}
	ob = rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).SkipLast(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("SkipLast操作后: ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res = []int{}
	ob = rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Take(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("Take操作后: ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n\n")

	fmt.Println("用于测试的数据为：10, 20, 30, 40, 50")
	res = []int{}
	ob = rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).TakeLast(3)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	fmt.Print("TakeLast操作后: ")
	for _, val := range res {
		fmt.Print(val, "  ")
	}
	fmt.Print("\n")
}
