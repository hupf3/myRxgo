# 修改、改进 RxGo 包

## 简介

ReactiveX是Reactive Extensions的缩写，一般简写为Rx，最初是LINQ的一个扩展，由微软的架构师Erik Meijer领导的团队开发，在2012年11月开源，Rx是一个编程模型，目标是提供一致的编程接口，帮助开发者更方便的处理异步数据流，Rx库支持.NET、JavaScript和C++，Rx近几年越来越流行了，现在已经支持几乎全部的流行编程语言了，Rx的大部分语言库由ReactiveX这个组织负责维护，比较流行的有RxJava/RxJS/Rx.NET，社区网站是 [reactivex.io](http://reactivex.io/documentation/observable.html)。[中文文档](https://mcxiaoke.gitbooks.io/rxdocs/content/Intro.html)

## 课程任务

阅读 ReactiveX 文档。请在 [pmlpml/RxGo](https://github.com/pmlpml/rxgo) 基础上，

1. 修改、改进它的实现
2. 或添加一组新的操作，如 [filtering](http://reactivex.io/documentation/operators.html#filtering)

该库的基本组成：

`rxgo.go` 给出了基础类型、抽象定义、框架实现、Debug工具等

`generators.go` 给出了 sourceOperater 的通用实现和具体函数实现

`transforms.go` 给出了 transOperater 的通用实现和具体函数实现

## 博客地址

[传送门](https://blog.csdn.net/qq_43267773/article/details/109479026)

## 设计说明

### 获取包

输入以下的命令即可获取我实现的 `myRxgo` 包

`go get github.com/hupf3/myRxgo/Rxgo`

或者在 src 的相应目录下输入以下命令

`git clone https://github.com/hupf3/myRxgo/Rxgo`

`go build`

`go install`

### 简单说明

此次作业都是基于老师实现的 [pmlpml/RxGo](https://github.com/pmlpml/rxgo) 基础上增加了 filter 过滤操作，老师原本的代码进行了少量的修改，方便 filter 的实现。

### 包文件结构

<img src="https://img-blog.csdnimg.cn/20201103213308141.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" style="zoom:33%;" />

- `rxgo.go` 给出了基础类型，抽象定义，框架实现，Debug 工具等；该文件基本的内容我没有进行修改，在 `Observable` 结构体中增加了过滤操作所用到的成员变量：

```go
debounce   time.Duration // n秒后发射一个数据
	distinct   bool          // 是否只选择不同的值
	elementAt  int           // 选择指定的索引的元素
	first      bool          // 是否选择第一数据
	last       bool          // 是否选择最后一个数据
	sample     time.Duration // 发射自上次采样以来它最近发射的数据
	skip       int           // 正值跳过前几项，负值跳过后几项
	take       int           // 正值选择前几项，负值选择后几项
	takeOrLast bool          // 判断是哪种take操作
```

- `rxgo_test.go`：该文件是 `rxgo.go` 的测试文件，测试了该函数的功能实现，并没有进行修改

- `transform.go`：给出了 transOperater 的通用实现和具体函数实现，并没有进行修改

- `transform_test.go`：该文件是 `transform.go` 的测试文件，测试了该函数的功能实现，并没有进行修改

- `utility.go`：定义了一些函数，方便文件调用，并没有进行修改

- `generators.go`：给出了 sourceOperater 的通用实现和具体函数实现，并没有进行修改

- `generators_test.go`：该文件是 `generators.go` 的测试文件，测试了该函数的功能实现，并没有进行修改

- `filter.go`：该文件是本次实现 `filter` 操作的代码文件

  - 首先按照老师写的代码文件，我也自己定义了一个过滤操作的结构体 `filterOperator` 

  ```go
  // 过滤操作
  type filterOperator struct {
  	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
  }
  ```

  - 自定义错误类型，方便后面函数的调用

  ```go
  // 自定义错误
  var NoInput = errors.New("No Input!")        // 没有输入
  var OutOfBound = errors.New("Out Of Bound!") // 出界
  ```

  - `newFilterObservable()`：初始化一个过滤 Observable

  ```go
  // 初始化一个过滤Observable
  func (parent *Observable) newFilterObservable(name string) (o *Observable) {
  	//new Observable
  	o = newObservable()
  	o.Name = name
  
  	//chain Observables
  	parent.next = o
  	o.pred = parent
  	o.root = parent.root
  
  	//set options
  	o.buf_len = BufferLen
  	return
  }
  ```

  - `Debounce()`：该函数实现了仅在过了一段指定的时间还没发射数据时才发射一个数据，并且在该函数的后面定义了 filterOperator，后面的操作实现也是如此方法

  ![debounce](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/debounce.png)

  ```go
  // 仅在过了一段指定的时间还没发射数据时才发射一个数据
  func (parent *Observable) Debounce(_debounce time.Duration) (o *Observable) {
  	o = parent.newFilterObservable("debounce")
  	o.first, o.last, o.distinct = false, false, false
  	o.debounce, o.take, o.skip = _debounce, 0, 0
  	o.operator = lastOperator
  	return o
  }
  
  var debounceOperator = filterOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `Distinct()`：该函数实现了抑制（过滤掉）重复的数据项的功能

  ![distinct](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/distinct.png)

  ```go
  // 抑制（过滤掉）重复的数据项
  func (parent *Observable) Distinct() (o *Observable) {
  	o = parent.newFilterObservable("distinct")
  	o.first, o.last, o.distinct = false, false, true
  	o.debounce, o.take, o.skip = 0, 0, 0
  	o.operator = distinctOperator
  	return o
  }
  
  var distinctOperator = filterOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `ElementAt()`：该函数实现了只发射第 n 个数据

  ![elementAt](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/elementAt.png)

  ```go
  // 只发射第N项数据
  func (parent *Observable) ElementAt(index int) (o *Observable) {
  	o = parent.newFilterObservable("elementAt")
  	o.first, o.last, o.distinct, o.takeOrLast = false, false, false, false
  	o.debounce, o.skip, o.take, o.elementAt = 0, 0, 0, index
  	o.operator = elementAtOperator
  	return
  }
  
  var elementAtOperator = filterOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `First()`：该函数实现了只发射第一项（或者满足某个条件的第一项）数据

  ![first](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/first.png)

  ```go
  // 只发射第一项（或者满足某个条件的第一项）数据
  func (parent *Observable) First() (o *Observable) {
  	o = parent.newFilterObservable("first")
  	o.first, o.last, o.distinct = true, false, false
  	o.debounce, o.take, o.skip = 0, 0, 0
  	o.operator = firstOperator
  	return o
  
  }
  
  var firstOperator = filterOperator{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `Last()`：该函数实现了只发射最后一项（或者满足某个条件的最后一项）数据

  ![last](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/last.png)

  ```go
  // 只发射最后一项（或者满足某个条件的最后一项）数据
  func (parent *Observable) Last() (o *Observable) {
  	o = parent.newFilterObservable("last")
  	o.first, o.last, o.distinct = false, true, false
  	o.debounce, o.take, o.skip = 0, 0, 0
  	o.operator = lastOperator
  	return o
  }
  
  var lastOperator = filterOperator{func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `Sample()`：该函数实现了定期发射Observable最近发射的数据项

  ![sample](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/sample.png)

  ```go
  // 定期发射Observable最近发射的数据项
  func (parent *Observable) Sample(_sample time.Duration) (o *Observable) {
  	o = parent.newFilterObservable("sample")
  	o.first, o.last, o.distinct, o.takeOrLast = false, false, false, false
  	o.debounce, o.skip, o.take, o.elementAt, o.sample = 0, 0, 0, 0, _sample
  	o.operator = sampleOperator
  	return o
  }
  
  var sampleOperator = filterOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `Skip()`：该函数实现了抑制Observable发射的前N项数据

  ![skip](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/skip.png)

  ```go
  // 抑制Observable发射的前N项数据
  func (parent *Observable) Skip(num int) (o *Observable) {
  	o = parent.newFilterObservable("skip")
  	o.first, o.last, o.distinct, o.takeOrLast = false, false, false, false
  	o.debounce, o.take, o.skip = 0, 0, num
  	o.operator = skipOperator
  	return o
  }
  
  var skipOperator = filterOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `SkipLast`：该函数实现了抑制Observable发射的后N项数据

  ![skipLast](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/skipLast.png)

  ```go
  // 抑制Observable发射的后N项数据
  func (parent *Observable) SkipLast(num int) (o *Observable) {
  	o = parent.newFilterObservable("skipLast")
  	o.first, o.last, o.distinct, o.takeOrLast = false, false, false, false
  	o.debounce, o.take, o.skip = 0, 0, -num
  	o.operator = skipLastOperator
  	return o
  }
  
  var skipLastOperator = filterOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `Take()`：该函数实现了只发射前面的N项数据

  ![take](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/take.png)

  ```go
  // 只发射前面的N项数据
  func (parent *Observable) Take(num int) (o *Observable) {
  	o = parent.newFilterObservable("Take")
  	o.first, o.last, o.distinct, o.takeOrLast = false, false, false, true
  	o.debounce, o.skip, o.take = 0, 0, num
  	o.operator = takeOperator
  	return o
  }
  
  var takeOperator = filterOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `TakeLast()`：该函数实现了发射Observable发射的最后N项数据

  ![takeLast](https://mcxiaoke.gitbooks.io/rxdocs/content/images/operators/takeLast.n.png)

  ```go
  // 发射Observable发射的最后N项数据
  func (parent *Observable) TakeLast(num int) (o *Observable) {
  	o = parent.newFilterObservable("takeLast")
  	o.first, o.last, o.distinct, o.takeOrLast = false, false, false, true
  	o.debounce, o.skip, o.take = 0, 0, -num
  	o.operator = takeLastOperator
  	return o
  }
  
  var takeLastOperator = filterOperator{opFunc: func(ctx context.Context, o *Observable, x reflect.Value, out chan interface{}) (end bool) {
  	var params = []reflect.Value{x}
  	x = params[0]
  
  	// send data
  	if !end {
  		end = o.sendToFlow(ctx, x.Interface(), out)
  	}
  	return
  },
  }
  ```

  - `takeOrSkip()`：该函数实现了判断是take操作还是skip操作

  ```go
  // 判断是take操作还是skip操作
  func takeOrSkip(_op bool, div int, in []interface{}) ([]interface{}, error) {
  	if (_op && div > 0) || (!_op && div < 0) {
  		if !_op {
  			div = len(in) + div
  		}
  		if div >= len(in) || div <= 0 {
  			return nil, OutOfBound
  		}
  		return in[:div], nil
  	}
  
  	if (_op && div < 0) || (!_op && div > 0) {
  		if _op {
  			div = len(in) + div
  		}
  		if div >= len(in) || div <= 0 {
  			return nil, OutOfBound
  		}
  		return in[div:], nil
  	}
  	return nil, OutOfBound
  }
  ```

  - `op()`：该函数实现了各个操作的初始化和报错信息，以及获取数据流

  ```go
  func (tsop filterOperator) op(ctx context.Context, o *Observable) {
  	// must hold defintion of flow resourcs here, such as chan etc., that is allocated when connected
  	// this resurces may be changed when operation routine is running.
  	in := o.pred.outflow
  	out := o.outflow
  	//fmt.Println(o.name, "operator in/out chan ", in, out)
  	var wg sync.WaitGroup
  
  	// 设置一个时间间隔
  	tspan := o.debounce
  	var _out []interface{}
  
  	go func() {
  		end := false
  		flag := make(map[interface{}]bool)
  
  		// 获取开始的时间
  		start := time.Now()
  		sample_start := time.Now()
  
  		for x := range in {
  			_start := time.Since(start)
  			_sample := time.Since(sample_start)
  			start = time.Now()
  
  			if end {
  				continue
  			}
  
  			if o.sample > 0 && _sample < o.sample {
  				continue
  			}
  
  			if tspan > time.Duration(0) && _start < tspan {
  				continue
  			}
  			// can not pass a interface as parameter (pointer) to gorountion for it may change its value outside!
  			xv := reflect.ValueOf(x)
  			// send an error to stream if the flip not accept error
  			if e, ok := x.(error); ok && !o.flip_accept_error {
  				o.sendToFlow(ctx, e, out)
  				continue
  			}
  
  			o.mu.Lock()
  			_out = append(_out, x)
  			o.mu.Unlock()
  
  			if o.elementAt > 0 {
  				continue
  			}
  
  			if o.take != 0 || o.skip != 0 {
  				continue
  			}
  
  			if o.last {
  				continue
  			}
  
  			if o.distinct && flag[xv.Interface()] {
  				continue
  			}
  			o.mu.Lock()
  			flag[xv.Interface()] = true
  			o.mu.Unlock()
  
  			// scheduler
  			switch threading := o.threading; threading {
  			case ThreadingDefault:
  				if o.sample > 0 {
  					sample_start = sample_start.Add(o.sample)
  				}
  				if tsop.opFunc(ctx, o, xv, out) {
  					end = true
  				}
  			case ThreadingIO:
  				fallthrough
  			case ThreadingComputing:
  				wg.Add(1)
  				if o.sample > 0 {
  					sample_start.Add(o.sample)
  				}
  				go func() {
  					defer wg.Done()
  					if tsop.opFunc(ctx, o, xv, out) {
  						end = true
  					}
  				}()
  			default:
  			}
  			if o.first {
  				break
  			}
  		}
  
  		if o.last && len(_out) > 0 {
  			wg.Add(1)
  			go func() {
  				defer wg.Done()
  				xv := reflect.ValueOf(_out[len(_out)-1])
  				tsop.opFunc(ctx, o, xv, out)
  			}()
  		}
  
  		if o.take != 0 || o.skip != 0 {
  			wg.Add(1)
  			go func() {
  				defer wg.Done()
  				var div int
  				if o.takeOrLast {
  					div = o.take
  				} else {
  					div = o.skip
  				}
  				new_in, err := takeOrSkip(o.takeOrLast, div, _out)
  
  				if err != nil {
  					o.sendToFlow(ctx, err, out)
  				} else {
  					xv := new_in
  					for _, val := range xv {
  						tsop.opFunc(ctx, o, reflect.ValueOf(val), out)
  					}
  				}
  			}()
  		}
  
  		if o.elementAt != 0 {
  			if o.elementAt < 0 || o.elementAt > len(_out) {
  				o.sendToFlow(ctx, OutOfBound, out)
  			} else {
  				xv := reflect.ValueOf(_out[o.elementAt-1])
  				tsop.opFunc(ctx, o, xv, out)
  			}
  		}
  
  		wg.Wait() //waiting all go-routines completed
  		if (o.last || o.first) && len(_out) == 0 && !o.flip_accept_error {
  			o.sendToFlow(ctx, NoInput, out)
  		}
  		o.closeFlow(out)
  	}()
  
  }
  ```

- `filter_test.go`：该文件实现了 `filter.go` 文件中各种过滤操作的测试，会在下面的单元测试中详细说明

## 单元测试

本部分主要讲的是 `filter_test.go` 文件中的函数测试，因为其余的测试文件我并没有进行太多的修改，都是老师原来的代码，本次实现过滤操作的测试代码都写入到了该文件中所以主要讲一下该文件中的测试。

- `func TestDebounce(t *testing.T)`：该函数主要测试了 `Debounce` 操作，测试的方法就是将时间设置很长，然后进行 `Debounce` 操作，得到的数组应为空数组，代码如下：

```go
func TestDebounce(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Debounce(1000000)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{}, res, "Debounce Test Error!")
}
```

- `func TestDistinct(t *testing.T)`：该函数主要测试了 `Distinct` 操作，测试的方法就是在数组中插入一些与原数组相同元素的数字，看看能否过滤出来，代码如下：

```go
func TestDistinct(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50, 30, 40).Map(func(x int) int {
		return x
	}).Distinct()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{10, 20, 30, 40, 50}, res, "Distinct Test Error!")
}
```

- `func TestElementAt(t *testing.T)`：该函数主要测试了 `ElementAt` 操作，测试的方法就是，在函数中传入一个参数，该参数代表想得到的数组中的索引，看看能否正确过滤，代码如下：

```go
func TestElementAt(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).ElementAt(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{40}, res, "SkipLast Test Error!")
}
```

- `func TestFirst(t *testing.T)`：该函数主要测试了 `First` 操作，测试该函数能否返回数组中的第一个元素，代码如下：

```go
func TestFirst(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).First()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{10}, res, "First Test Error!")
}
```

- `func TestLast(t *testing.T)`：该函数主要测试了 `Last` 操作，测试该函数能否返回数组中的最后一个元素，代码如下：

```go
func TestLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Last()
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})

	assert.Equal(t, []int{50}, res, "Last Test Error!")
}
```

- `func TestSkip(t *testing.T)`：该函数主要测试了 `Skip` 操作，测试该函数能否跳过前n个数，读取其他的数据，代码如下：

```go
func TestSkip(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Skip(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{50}, res, "Skip Test Error!")
}
```

- `func TestSkipLast(t *testing.T)`：该函数主要测试了 `SkipLast` 操作，测试该函数能够跳过后n个数，得到其他的数据，代码如下：

```go
func TestSkipLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).SkipLast(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{10}, res, "SkipLast Test Error!")
}
```

- `func TestTake(t *testing.T)`：该函数主要测试了 `Take` 操作，测试该函数能否取得数组中前n个数据，代码如下：

```go
func TestTake(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).Take(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{10, 20, 30, 40}, res, "Take Test Error!")
}
```

- `func TestTakeLast(t *testing.T)`：该函数主要测试了 `TakeLast` 操作，测试该函数能否取得数组中的后n个数据，代码如下：

```go
func TestTakeLast(t *testing.T) {
	res := []int{}
	ob := rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		return x
	}).TakeLast(4)
	ob.Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{20, 30, 40, 50}, res, "TakeLast Test Error!")
}
```

- `func TestSample(t *testing.T)`：该函数主要测试了 `Sample` 操作，测试该函数能否定期采样数据，代码如下：

```go
func TestSample(t *testing.T) {
	res := []int{}
	rxgo.Just(10, 20, 30, 40, 50).Map(func(x int) int {
		time.Sleep(2 * time.Millisecond)
		return x
	}).Sample(3 * time.Millisecond).Subscribe(func(x int) {
		res = append(res, x)
	})
	assert.Equal(t, []int{20, 30, 40, 50}, res, "SkipLast Test Error!")

}
```

测试的最终结果如下图所示：

<img src="https://img-blog.csdnimg.cn/20201103223424593.png#pic_center" alt="在这里插入图片描述" style="zoom:50%;" />

## 功能测试

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

通过结果可以知道，已经成功完成了filter的各个功能，成功完成功能测试

## API文档

生成网页版的 API 文档，输入以下的命令：

`godoc -http=:8080`

然后在浏览器中打开 [http://127.0.0.1:8080](http://127.0.0.1:8080/) ，即可访问网页版的 go doc：

<img src="https://img-blog.csdnimg.cn/20201103224226225.png?x-oss-process=image/watermark,type_ZmFuZ3poZW5naGVpdGk,shadow_10,text_aHR0cHM6Ly9ibG9nLmNzZG4ubmV0L3FxXzQzMjY3Nzcz,size_16,color_FFFFFF,t_70#pic_center" alt="在这里插入图片描述" style="zoom:33%;" />

在目录结构下执行以下命令，即可生成线下的 html 文件

`go doc`

`godoc -url="pkg/github.com/hupf3/myRxgo/Rxgo" > API.html`

我将该文档也保存在了 github 仓库中方便检查

## 总结

通过本次实验了解到了 ReactiveX 异步编程的 API，并且更加熟悉了 chan 等操作，提高了自己的代码能力

