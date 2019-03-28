package main

import (
	"os"
)

func main() {
	f, err := os.Open(os.Args[1])
	if err != nil {
		println("Error: 1\n" + err.Error())
	}
	defer f.Close()

	for {
		list, err := f.Readdir(1000)
		if err != nil {
			if err.Error() == "EOF" {
				break;
			}
			println("Error: 2\n" + err.Error())
		}

		for _, fi := range list {
			println(fi.Name() + " @ IsDir? " + btoa(fi.IsDir()))
		}
	}
}

func btoa(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
