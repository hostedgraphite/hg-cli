package utils

// For testing, we can use a global logger instance to log messages to a file.
// This is useful for debugging and monitoring the application.
// The log file will be created in the root directory of the application.

// import (
// 	"fmt"
// 	"log"
// 	"os"
// )

// // Global logger instance
// var Logger *log.Logger
// var logFile *os.File

// func init() {
// 	var err error
// 	logFile, err = os.OpenFile("app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
// 	if err != nil {
// 		fmt.Println("Error opening log file:", err)
// 		return
// 	}

// 	Logger = log.New(logFile, "", log.Ldate|log.Ltime|log.Lshortfile)
// }

// // CloseLogFile is a function to close the log file when the application exits
// func CloseLogFile() {
// 	if logFile != nil {
// 		logFile.Close()
// 	}
// }
