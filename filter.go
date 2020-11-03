package myrxgo

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"sync"
	"time"
)

var (
	NoInputStream      = errors.New("There are no input stream.")
	OutOfBounds        = errors.New("The number of expected flow is beyond the length of the input flow")
	ElementOutOfBounds = errors.New("The element you are looking for is out of the input bound")
)

// 过滤操作
type filterOperator struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

// partion the input stream by giving condition
// isTake ---- true, indicates that it is take mode
// isTake ---- false,  indicates that it is skip mode
func partionFlow(isTake bool, division int, in []interface{}) ([]interface{}, error) {
	// get the first few
	fmt.Println(in)
	if (isTake && division > 0) || (!isTake && division < 0) {
		if !isTake {
			division = len(in) + division
		}
		if division >= len(in) || division <= 0 {
			return nil, OutOfBounds
		}
		return in[:division], nil
	}

	// get the first few
	if (isTake && division < 0) || (!isTake && division > 0) {
		if isTake {
			division = len(in) + division
		}
		if division >= len(in) || division <= 0 {
			return nil, OutOfBounds
		}
		return in[division:], nil
	}
	return nil, OutOfBounds
}

func (tsop filterOperator) op(ctx context.Context, o *Observable) {
	// must hold defintion of flow resourcs here, such as chan etc., that is allocated when connected
	// this resurces may be changed when operation routine is running.
	in := o.pred.outflow
	out := o.outflow
	//fmt.Println(o.name, "operator in/out chan ", in, out)
	var wg sync.WaitGroup

	// 设置一个时间间隔
	interval := o.debounce
	var out_buf []interface{}

	go func() {
		end := false

		is_appear := make(map[interface{}]bool)
		// 获取开始的时间
		start := time.Now()
		sample_start := time.Now()

		for x := range in {
			// receiving time since last receive
			rt := time.Since(start)
			st := time.Since(sample_start)
			start = time.Now()

			if end {
				continue
			}

			// Not satisfy the sample period
			if o.sample > 0 && st < o.sample {
				continue
			}

			// if the receiving time is less than the debounce interval, not send it to the flow
			if interval > time.Duration(0) && rt < interval {
				continue
			}
			// can not pass a interface as parameter (pointer) to gorountion for it may change its value outside!
			xv := reflect.ValueOf(x)
			// send an error to stream if the flip not accept error
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}

			// buffer the input flow
			o.mu.Lock()
			out_buf = append(out_buf, x)
			o.mu.Unlock()

			// if using ElementAt operator, skip the loop
			if o.elementAt > 0 {
				continue
			}

			// if using Take, Takelast, Skip, SkipLast operator, directly skip the loop
			if o.take != 0 || o.skip != 0 {
				continue
			}
			// if using the 'Last' operator, save the input flow val in out_buf and continuing receiving
			if o.last {
				continue
			}
			//_, ok := is_appear[xv]
			// if using the 'Distinct' operator, firstly judge whether the xv has already appeared previously, if so, skip to next loop
			if o.distinct && is_appear[xv.Interface()] {
				continue
			}
			// memorizing the appearance of the element
			o.mu.Lock()
			is_appear[xv.Interface()] = true
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
			// if using 'First' operator, end the loop once the first input is sent to opFunc
			if o.first {
				break
			}
		}

		// if in Last Filter, send the last element in the outbuf to Last
		if o.last && len(out_buf) > 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				xv := reflect.ValueOf(out_buf[len(out_buf)-1])
				tsop.opFunc(ctx, o, xv, out)
			}()
		}

		// Take, TakeLast, Skip, SkipLast operator
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
				new_in, err := partionFlow(o.takeOrLast, div, out_buf)

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

		// ElementAt operator
		if o.elementAt != 0 {
			if o.elementAt < 0 || o.elementAt > len(out_buf) {
				o.sendToFlow(ctx, ElementOutOfBounds, out)
			} else {
				xv := reflect.ValueOf(out_buf[o.elementAt-1])
				tsop.opFunc(ctx, o, xv, out)
			}
		}
		wg.Wait() //waiting all go-routines completed

		if (o.last || o.first) && len(out_buf) == 0 && !o.flip_accept_error {
			o.sendToFlow(ctx, NoInputStream, out)
		}

		o.closeFlow(out)
	}()

}

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
func (parent *Observable) Debounce(t time.Duration) (o *Observable) {
	o = parent.newFilterObservable("debounce")
	o.first, o.last, o.distinct = false, false, false
	o.debounce, o.take, o.skip = t, 0, 0
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
func (parent *Observable) Sample(st time.Duration) (o *Observable) {
	o = parent.newFilterObservable("sample")
	o.first, o.last, o.distinct, o.takeOrLast = false, false, false, false
	o.debounce, o.skip, o.take, o.elementAt, o.sample = 0, 0, 0, 0, st
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
