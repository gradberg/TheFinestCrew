/*
    Responsible for handling the input of both keboard and timed events
*/

package main

import "time"
import "termbox-go"

type AsyncInputEvent struct {
    Ch rune
    Key termbox.Key
    IsTimerEvent bool
}

var timerChannel chan int64
var asyncInputChannel chan AsyncInputEvent

func AsyncInput_Init() {
    asyncInputChannel = make(chan AsyncInputEvent)
    // making the timerChannel buffered (with at least 1), and emptying it out before
    // ever adding a value prevents deadlocks with turning on and off realtime really fast
    timerChannel = make(chan int64, 1)
    go receiveTerminalInput()
    go runTimerInput()
}

func AsyncInput_Close() {
    close(asyncInputChannel)
    close(timerChannel)
}

func AsyncInput_ReceiveInput() AsyncInputEvent {
    return <- asyncInputChannel
}
func AsyncInput_QueueTimer(milliseconds int64) {
    timerChannel <- milliseconds
}
func AsyncInput_EmptyTimerQueue() {
    select {
        case _ = <- timerChannel:
        default: return
    }
}

func receiveTerminalInput() {
    for {        
        e := Console_ReadKey()
        asyncInputChannel <- AsyncInputEvent {
            Ch: e.Ch,
            Key: e.Key,
        }
    }
}

func runTimerInput() {
    for {
        waitInMilliseconds := <- timerChannel
        time.Sleep(time.Duration(waitInMilliseconds) * time.Millisecond)       
        sendTimerInput()
    }
}
func sendTimerInput() {
    // this has been extracted to its own method so that it can catch a panic when the 
    // timerChannel is closed. As this goroutine can pause for upwards of seconds, I do not
    // know what better mechanism there is for closing it.
    defer func() {
        _ = recover() // discard panic information
    }()
    
    asyncInputChannel <- AsyncInputEvent {
        IsTimerEvent: true,
    }
}
