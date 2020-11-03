# 修改、改进 RxGo 包

## 概述

本仓库包括了程序包说明文档 `README.md` ，实现程序包的过程以及测试的的文档 `specification.md`，`API.html` 是我生成的线下的 API 文件，方便查看函数的相关用法，`Rxgo` 文件夹中保存了本次我实现的 rxgo 程序包的代码文件和相应的测试文件。`main.go` 用于测试本程序包中的功能实现是否正确。

## 获取包

输入以下的命令即可获取我实现的 `myRxgo` 包

`go get github.com/hupf3/myRxgo/Rxgo`

或者在 src 的相应目录下输入以下命令

`git clone https://github.com/hupf3/myRxgo/Rxgo`

`go build`

`go install`

## 使用说明

**实验环境**：

- 操作系统：mac os
- golang 版本: golang 1.14及以上

为了方便测试过滤操作的设计，我自己写了一个 `main.go` 文件，在 go get 我实现的 rxgo 程序包后，直接运行该文件即可进行测试过滤操作，代码文件如下：

```go
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
```

写好代码文件，运行如下命令 `go run main.go` 即可进行查看输出的结果：

<img src="https://img-blog.csdnimg.cn/20201103224008138.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom: 33%;" />

通过结果可以知道，已经成功完成了filter的各个功能.

至此，示例展示完毕，也是比较简单的进行了实现。也可以通过查看 API 文档进行具体的学习，API 文档生成的过程在 `specification.md` 中有具体的说明！