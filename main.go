package main

import (
	"fmt"
	"os"
  "os/exec"
	"strconv"
	"strings"
  "math"
)

const (
  DRIVER = "intel_backlight" // check your driver by running 'ls /sys/class/backlight'
  BASE_PATH = "/sys/class/backlight/" + DRIVER
  LOWER_BOUND = 0.001 // lowest possible percentage of max brightness
)

// Since this program is supposed to run as root with SUID, i decided it would be the
// best to panic at every error, for security reasons.

var MAX_BRIGHTNESS int = getMaxBrightness()
var MIN_BRIGHTNESS int = int(math.Ceil(float64(MAX_BRIGHTNESS) * LOWER_BOUND))

func main() {
	handleArg(os.Args)
}

func getDriver() string {
  cmd := exec.Command("ls", "/sys/class/backlight")
  stdout, err := cmd.Output()
  checkPanic(err)
  return string(stdout[:])
} 

func handleArg(arg []string) {
	if len(arg) != 2 {
		printUsage()
	}

	operator := string(arg[1][0])
	amount, err := strconv.Atoi(string(arg[1][1:]))
	if err != nil {
		printUsage()
	}

	switch operator {
	case "+":
		changeBrightness(amount, intAdd)
	case "-":
		changeBrightness(amount, intSub)
	default:
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Thats the wrong syntax. This program is very simple, the usage is either:")
	fmt.Println("  goblight +[amount]")
	fmt.Println("  goblight -[amount]")
	fmt.Println("")
	fmt.Println("Example: goblight +30")
	os.Exit(0)
}

func changeBrightness(amount int, op intOperation) {
	f, err := os.OpenFile(fmt.Sprintf("%s/brightness", BASE_PATH), os.O_RDWR|os.O_TRUNC, 0755)
	defer f.Close()
	checkPanic(err)

	prev, err := strconv.Atoi(strings.TrimSuffix(string(getCurrentBrightness()), "\n"))
	checkPanic(err)

	diff := op(prev, amount)

	var br []byte

	if diff >= MAX_BRIGHTNESS {
		br = []byte(fmt.Sprint(MAX_BRIGHTNESS))
	} else if diff <= MIN_BRIGHTNESS {
		br = []byte(strconv.Itoa(MIN_BRIGHTNESS))
	} else {
		br = []byte(fmt.Sprint(diff))
	}

	_, err = f.Write(br)
	checkPanic(err)
}

func getCurrentBrightness() []byte {
	data, err := os.ReadFile(fmt.Sprintf("%s/brightness", BASE_PATH))
	checkPanic(err)

	return data
}

func getMaxBrightness() int {
	data, err := os.ReadFile(fmt.Sprintf("%s/max_brightness", BASE_PATH))
	checkPanic(err)

	i, err := strconv.Atoi(strings.TrimSuffix(string(data), "\n"))
	checkPanic(err)

	return i
}

func checkPanic(err error) {
	if err != nil {
		panic(err)
	}
}

// Go requires me to do this crap sometimes...
type intOperation func(int, int) int

func intAdd(x int, y int) int {
	return x + y
}

func intSub(x int, y int) int {
	return x - y
}
