/*
    This has helper methods that define log 'levels' for the logger.
    It also prevents each file from needing to import "log"
*/

package main
import "log"


// General information
func LogInfo(format string, v ...interface{}) {
    log.Printf("_INFO_ " + format + "\n", v...)
}

// Represent potential problems, but not game-breaking ones
func LogWarn(format string, v ...interface{}) {
    log.Printf("_WARN_ " + format + "\n", v...)
}

// Calculation-related logging
func LogCalc(format string, v ...interface{}) {
    log.Printf("_CALC_ " + format + "\n", v...)
}