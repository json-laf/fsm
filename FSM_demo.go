package main

import (
	"demo/fsm"
	"errors"
	"fmt"
	"os/exec"
	"strings"
)

var (
	// 状态(State)：事物的状态，包括初始状态和所有事件触发后的状态
	Poweroff   = fsm.FSMState("关闭")
	FirstGear  = fsm.FSMState("闪烁")
	SecondGear = fsm.FSMState("常亮")
	// 事件(Event)：触发状态变化或者保持原状态的事件
	PowerOffEvent   = fsm.FSMEvent("按下关闭按钮")
	FirstGearEvent  = fsm.FSMEvent("按下闪烁按钮")
	SecondGearEvent = fsm.FSMEvent("按下常亮按钮")
	// 行为或转换(Action/Transition)：执行状态转换的过程
	PowerOffHandler = fsm.FSMHandler(func() fsm.FSMState {
		err := Leds("blinks")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("灯已关闭")
		return Poweroff
	})
	FirstGearHandler = fsm.FSMHandler(func() fsm.FSMState {
		err := Leds("blinks")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("灯开始闪烁！")
		return FirstGear
	})
	SecondGearHandler = fsm.FSMHandler(func() fsm.FSMState {
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
	*fsm.FSM
}

// 实例化灯
func NewLed(initState fsm.FSMState) *Led {
	return &Led{
		FSM: fsm.NewFSM(initState),
	}
}

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

// 入口函数
func main() {

	efan := NewLed(Poweroff) // 初始状态是关闭的
	// 关闭状态
	efan.AddHandler(Poweroff, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(Poweroff, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(Poweroff, SecondGearEvent, SecondGearHandler)
	// 闪烁状态
	efan.AddHandler(FirstGear, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(FirstGear, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(FirstGear, SecondGearEvent, SecondGearHandler)
	// 常亮状态
	efan.AddHandler(SecondGear, PowerOffEvent, PowerOffHandler)
	efan.AddHandler(SecondGear, FirstGearEvent, FirstGearHandler)
	efan.AddHandler(SecondGear, SecondGearEvent, SecondGearHandler)

	// 开始测试状态变化
	// efan.Call(FirstGearEvent) // 按下闪烁按钮
	// efan.Call(PowerOffEvent)   // 按下关闭按钮
	efan.Call(SecondGearEvent) // 按下常亮按钮
	// efan.Call(PowerOffEvent)   // 按下关闭按钮
}
