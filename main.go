package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	DRIVER    = "acpi_video0" // check your driver by running 'ls /sys/class/backlight'
	BASE_PATH = "/sys/class/backlight/" + DRIVER
)

// Since this program is supposed to run as root with SUID, i decided it would be the
// best to panic at every error, for security reasons.

var MAX_BRIGHTNESS int = getMaxBrightness()

func main() {
	handleArg(os.Args)
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
	if err != nil {
		panic(err)
	}

	prev, err := strconv.Atoi(strings.TrimSuffix(string(getCurrentBrightness()), "\n"))
	if err != nil {
		panic(err)
	}

	diff := op(prev, amount)

	var br []byte

	if diff >= MAX_BRIGHTNESS {
		br = []byte(fmt.Sprint(MAX_BRIGHTNESS))
	} else if diff <= 0 {
		br = []byte("0")
	} else {
		br = []byte(fmt.Sprint(diff))
	}

	fmt.Printf("diff: %d\n", diff)
	fmt.Printf("br: %s\n", br)

	_, err = f.Write(br)
	if err != nil {
		panic(err)
	}
}

func getCurrentBrightness() []byte {
	data, err := os.ReadFile(fmt.Sprintf("%s/brightness", BASE_PATH))
	if err != nil {
		panic(err)
	}

	return data
}

func getMaxBrightness() int {
	data, err := os.ReadFile(fmt.Sprintf("%s/max_brightness", BASE_PATH))
	if err != nil {
		panic(err)
	}

	i, err := strconv.Atoi(strings.TrimSuffix(string(data), "\n"))
	if err != nil {
		panic(err)
	}

	return i
}

// Go requires me to do this crap sometimes...
type intOperation func(int, int) int

func intAdd(x int, y int) int {
	return x + y
}

func intSub(x int, y int) int {
	return x - y
}
