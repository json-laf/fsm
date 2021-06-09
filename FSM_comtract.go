package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type FSMState string            // 状态
type FSMEvent string            // 事件
type FSMHandler func() FSMState // 处理方法，并返回新的状态

//合约
type DB struct {
	contractapi.Contract
}

// 有限状态机
type FSM struct {
	mu       sync.Mutex                           // 排他锁
	state    FSMState                             // 当前状态
	handlers map[FSMState]map[FSMEvent]FSMHandler // 处理地图集，每一个状态都可以出发有限个事件，执行有限个处理
}

// 获取当前状态
func (f *FSM) getState() FSMState {
	return f.state
}

// 设置当前状态
func (f *FSM) setState(newState FSMState) {
	f.state = newState
}

// 某状态添加事件处理方法
func (f *FSM) AddHandler(state FSMState, event FSMEvent, handler FSMHandler) *FSM {
	if _, ok := f.handlers[state]; !ok {
		f.handlers[state] = make(map[FSMEvent]FSMHandler)
	}
	if _, ok := f.handlers[state][event]; ok {
		fmt.Printf("[警告] 状态(%s)事件(%s)已定义过", state, event)
	}
	f.handlers[state][event] = handler
	return f
}

// 事件处理
func (f *FSM) Call(event FSMEvent) FSMState {
	f.mu.Lock()
	defer f.mu.Unlock()
	events := f.handlers[f.getState()]
	if events == nil {
		return f.getState()
	}
	if fn, ok := events[event]; ok {
		oldState := f.getState()
		f.setState(fn())
		newState := f.getState()
		fmt.Println("状态从 [", oldState, "] 变成 [", newState, "]")
	}
	return f.getState()
}

// 实例化FSM
func NewFSM(initState FSMState) *FSM {
	return &FSM{
		state:    initState,
		handlers: make(map[FSMState]map[FSMEvent]FSMHandler),
	}
}

var (
	// 状态(State)：事物的状态，包括初始状态和所有事件触发后的状态
	Poweroff   = FSMState("关闭")
	FirstGear  = FSMState("闪烁")
	SecondGear = FSMState("常亮")
	// 事件(Event)：触发状态变化或者保持原状态的事件
	PowerOffEvent   = FSMEvent("按下关闭按钮")
	FirstGearEvent  = FSMEvent("按下闪烁按钮")
	SecondGearEvent = FSMEvent("按下常亮按钮")
	// 行为或转换(Action/Transition)：执行状态转换的过程
	PowerOffHandler = FSMHandler(func() FSMState {
		err := Leds("blinks")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("灯已关闭")
		return Poweroff
	})
	FirstGearHandler = FSMHandler(func() FSMState {
		err := Leds("blinks")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("灯开始闪烁！")
		return FirstGear
	})
	SecondGearHandler = FSMHandler(func() FSMState {
		err := Leds("light")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("灯开始常亮！")
		return SecondGear
	})
)

// 灯
type Led struct {
	*FSM
}

// 实例化灯
func NewLed(initState FSMState) *Led {
	return &Led{
		FSM: NewFSM(initState),
	}
}

//对等进行控制
func Leds(arg string) (err error) {
	args := []string{"led.py", arg}
	out, err := exec.Command("python3", args...).Output()
	if err != nil {
		return
	}
	result := string(out)
	if strings.Index(result, "success") != 0 {
		err = errors.New(fmt.Sprintf("%s", result))
	}
	return
}

// 初始化
func (d *DB) Init(ctx contractapi.TransactionContextInterface) *Led {
	led := NewLed(Poweroff) // 初始状态是关闭的
	// 关闭状态
	led.AddHandler(Poweroff, PowerOffEvent, PowerOffHandler)
	led.AddHandler(Poweroff, FirstGearEvent, FirstGearHandler)
	led.AddHandler(Poweroff, SecondGearEvent, SecondGearHandler)
	// 闪烁状态
	led.AddHandler(FirstGear, PowerOffEvent, PowerOffHandler)
	led.AddHandler(FirstGear, FirstGearEvent, FirstGearHandler)
	led.AddHandler(FirstGear, SecondGearEvent, SecondGearHandler)
	// 常亮状态
	led.AddHandler(SecondGear, PowerOffEvent, PowerOffHandler)
	led.AddHandler(SecondGear, FirstGearEvent, FirstGearHandler)
	led.AddHandler(SecondGear, SecondGearEvent, SecondGearHandler)
	fmt.Println(led)
	return led
}

// 调用函数
func (d *DB) Call(ctx contractapi.TransactionContextInterface, led *Led, event FSMEvent) error {
	// 开始测试状态变化
	fmt.Println("Called...")
	// led.Call(event)
	return nil
}

//查询
func (d *DB) Query(ctx contractapi.TransactionContextInterface, led *Led) FSMState {
	return led.getState()
}

func main() {
	chaincode, err := contractapi.NewChaincode(new(DB))
	if err != nil {
		fmt.Printf("Error create chaincode: %s", err.Error())
		return
	}
	if err := chaincode.Start(); err != nil {
		fmt.Printf("Error starting chaincode: %s", err.Error())
	}
}
