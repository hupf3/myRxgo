package myrxgo

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"time"
)

// 过滤操作
type filterOperator struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

// 自定义错误
var NoInput = errors.New("No Input!")        // 没有输入
var OutOfBound = errors.New("Out Of Bound!") // 出界

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
