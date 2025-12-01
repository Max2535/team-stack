package main

import (
    "fmt"
    "os"
)

func main() {
    if len(os.Args) < 2 {
        fmt.Println("usage: team <command>")
        return
    }

    switch os.Args[1] {
    case "new:module":
        fmt.Println("Generate new module (stub)")
    default:
        fmt.Println("unknown command:", os.Args[1])
    }
}
