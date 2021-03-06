package utility

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	//"syscall"
)

type CommandHandler func(args []string) bool

var (
	command                              = make([]byte, 1024)
	reader                               = bufio.NewReader(os.Stdin)
	HandlerMap map[string]CommandHandler = make(map[string]CommandHandler, 20)
)

func StartConsole() {
	go consoleroutine()
}

func StartConsoleWait() {
	consoleroutine()
}

func consoleroutine() {
	for {
		command, _, _ = reader.ReadLine()
		Args := strings.Split(string(command), " ")

		cmdhandler, ok := HandlerMap[Args[0]]
		if ok {
			cmdhandler(Args)
			continue
		}

		switch string(Args[0]) {
		case "cpus":
			fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
		case "routines":
			fmt.Println("Current number of goroutines: ", runtime.NumGoroutine())
		case "setcpus":
			n, _ := strconv.Atoi(Args[1])
			runtime.GOMAXPROCS(n)
			fmt.Println(runtime.NumCPU(), " cpus and ", runtime.GOMAXPROCS(0), " in use")
		case "maxhttpcon":
			fmt.Println("max http conn : ", G_Count)
		case "startgc":
			runtime.GC()
			fmt.Println("gc finished")
		case "lscmd":
			fmt.Println("cmd count:", len(HandlerMap))
			for key, _ := range HandlerMap {
				fmt.Println(key)
			}
		default:
			fmt.Println("Command error, try again.")
		}
	}
}

func HandleFunc(cmd string, mh CommandHandler) {
	if HandlerMap == nil {
		HandlerMap = make(map[string]CommandHandler, 20)
	}

	HandlerMap[cmd] = mh

	return
}

func ColorPrintln(s string) {
	//kernel32 := syscall.NewLazyDLL("kernel32.dll")
	//proc := kernel32.NewProc("SetConsoleTextAttribute")
	//handle, _, _ := proc.Call(uintptr(syscall.Stdout), uintptr(12)) //12 Red light

	//fmt.Println(s)

	//handle, _, _ = proc.Call(uintptr(syscall.Stdout), uintptr(7)) //White dark
	//CloseHandle := kernel32.NewProc("CloseHandle")
	//CloseHandle.Call(handle)
}
