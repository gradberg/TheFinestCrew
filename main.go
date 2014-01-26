package main

import (
    "fmt"
    "os"
    "log"
)

func main() {
    // Initialize Console
    err := Console_Init()
    if err != nil {
        fmt.Printf("Failed to initialize terminal: %v", err)
        return
    }
    defer Console_Close()
    
    err = Console_Validate()
    if err != nil {
        fmt.Printf("Failed to initialize terminal: %v", err)
        return
    }
    
    initializeLogging()

    game := NewGame()
    game.Run()
    
    fmt.Printf("Bye Bye")
}

func initializeLogging() {
    logToFile := true   // This should be read from a config file or something
    
    if (logToFile) {
        s := "tfc_log.txt"
        f, err := os.Create(s)
        if (err == nil) {
            // write to console at least
            log.Printf("Sending log output to %s\n", s)
            log.SetOutput(f)
            return
        }
        
        log.Printf("Failed to set logging to file %s, reason %s\n", s, err)
    }
    
    // If logging to a file is not set, or it fails to create the file, turn the logging off entirely
    log.Printf("Disabling logger\n")
    log.SetOutput(&discardingWriter {})    
}

type discardingWriter struct { }
func (dw *discardingWriter) Write(p []byte) (n int, err error) {
    return len(p), nil
}





