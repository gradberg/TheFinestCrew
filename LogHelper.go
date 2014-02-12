/*
    This has helper methods that define log 'levels' for the logger.
    It also prevents each file from needing to import "log"
*/

package main
import "log"
import "fmt"

// General information
func LogInfo(format string, v ...interface{}) {
    log.Printf("_INFO__ " + format + "\r\n", v...)
}

// Represents temporary output that should be disabled or removed for normal deployments
func LogDebug(format string, v ...interface{}) {
    log.Printf("_DEBUG_ " + format + "\r\n", v...)
}

func LogAi(cm *CrewMember, format string, v ...interface{}) {
    prefix := fmt.Sprintf("__AI___ [%s:%s] ", cm.GetFullName(), cm.CrewRole.ToString())
    log.Printf(prefix + format + "\r\n", v...)
}

// Represent potential problems, but not game-breaking ones
func LogWarn(format string, v ...interface{}) {
    log.Printf("_WARN__ " + format + "\r\n", v...)
}

// Calculation-related logging
func LogCalc(format string, v ...interface{}) {
    log.Printf("_CALC__ " + format + "\r\n", v...)
}