package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const finalWord = "Go!"
const countdownStart = 3
const write = "write"
const sleep = "sleep"

type DefaultSleeper struct{}

type Sleeper interface {
	Sleep()
}

type SpySleeper struct {
	Calls int
}

type SpyCountdownOperations struct {
	Calls []string
}

// two new functions one is sleep that is doing what its meant to do
func (s *SpyCountdownOperations) Sleep() {
	s.Calls = append(s.Calls, sleep)
}

// this is write ops which is basically calling the io.write thingy
// This function satisfies the io.Writer interface.
// Instead of actually printing "3, 2, 1" to a screen, it just writes the word "write" onto the tape.
func (s *SpyCountdownOperations) Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}

func main() {
	sleeper := &DefaultSleeper{}
	Countdown(os.Stdout, sleeper)
}

func (s *SpySleeper) Sleep() {
	s.Calls++
}
func Countdown(out io.Writer, sleeper Sleeper) {
	for i := countdownStart; i > 0; i-- {
		fmt.Fprint(out, i)
		sleeper.Sleep()
	}
	fmt.Fprint(out, finalWord)
}
