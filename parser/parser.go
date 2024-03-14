package parser

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type LogMsg struct {
	timestamp  time.Time
	userID     int
	logType    string
	logMessage string
}

// LogParser is a struct that holds the log files and the errors, warns and infos
type LogParser struct {
	logFolder   string
	logFiles    []string
	logMessages map[string][]LogMsg
}

func New(logFolder string) *LogParser {
	logMessages := make(map[string][]LogMsg)
	logMessages["ERROR"] = make([]LogMsg, 0)
	logMessages["WARN"] = make([]LogMsg, 0)
	logMessages["INFO"] = make([]LogMsg, 0)

	return &LogParser{
		logFolder:   logFolder,
		logMessages: logMessages,
	}
}

func (parser *LogParser) openFile(fileName string) (*os.File, error) {
	// format the file path
	filePath := fmt.Sprintf("%s/%s", parser.logFolder, fileName)
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (parser *LogParser) PrintFiles() {
	for _, file := range parser.logFiles {
		println(file)
	}
}

func (parser *LogParser) PrintErrorLog() {
	if errorLogs, ok := parser.logMessages["ERROR"]; ok {
		for _, logMsg := range errorLogs {
			fmt.Println(logMsg)
		}
	} else {
		fmt.Println("no ERROR log messages ")
	}

}

func (parser *LogParser) PrintWarnLog() {
	if warnLogs, ok := parser.logMessages["WARN"]; ok {
		for _, logMsg := range warnLogs {
			fmt.Println(logMsg)
		}
	} else {
		fmt.Println("no WARN log messages ")
	}

}

func (parser *LogParser) PrintInfoLog() {
	if infoLogs, ok := parser.logMessages["INFO"]; ok {
		for _, logMsg := range infoLogs {
			fmt.Println(logMsg)
		}
	} else {
		fmt.Println("no INFO log messages ")
	}

}

func ensureDateString(dateString string) (time.Time, error) {
	dateTime := strings.Split(dateString, "-")
	if len(dateTime) != 2 {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}
	date := dateTime[0]
	timestamp := dateTime[1]

	dateSplit := strings.Split(date, "/")
	if len(dateSplit) != 3 {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}
	year, err := strconv.Atoi(dateSplit[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}
	month, err := strconv.Atoi(dateSplit[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}
	day, err := strconv.Atoi(dateSplit[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}

	timeSplit := strings.Split(timestamp, ":")
	if len(timeSplit) != 3 {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}

	hour, err := strconv.Atoi(timeSplit[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}
	minute, err := strconv.Atoi(timeSplit[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}
	second, err := strconv.Atoi(timeSplit[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date string: %s", dateString)
	}

	return time.Date(year, time.Month(month), day, hour, minute, second, 0, time.UTC), nil
}

func parseLogLine(line string) (LogMsg, error) {
	lineSplit := strings.Split(line, " ")
	if len(lineSplit) <= 3 {
		return LogMsg{}, fmt.Errorf("invalid log line: %s", line)
	}

	dateString := lineSplit[0]
	timestamp, err := ensureDateString(dateString)
	if err != nil {
		return LogMsg{}, err
	}

	userID, err := strconv.Atoi(lineSplit[1])
	if err != nil {
		return LogMsg{}, fmt.Errorf("invalid user id: %s", lineSplit[1])
	}
	//continue

	logMessage := strings.Join(lineSplit[3:], " ") // hello world!
	logType := lineSplit[2]                        // INFO
	if logType != "ERROR" && logType != "INFO" && logType != "WARN" {
		return LogMsg{}, fmt.Errorf("invalid log type: %s", lineSplit[2])
	}

	return LogMsg{
		timestamp:  timestamp,
		userID:     userID,
		logType:    logType,
		logMessage: logMessage}, nil

}

func (parser *LogParser) ParseLines(fileName string) error {

	fp, err := parser.openFile(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()
	parser.logFiles = append(parser.logFiles, fileName)

	var lines []string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}

	// Get []LogMsg of all keys before the for
	errorLogs := parser.logMessages["ERROR"]
	warnLogs := parser.logMessages["WARN"]
	infoLogs := parser.logMessages["WARN"]

	for _, line := range lines {

		logMsg, err := parseLogLine(line)
		if err != nil {
			return err
		}

		if logMsg.logType == "ERROR" {
			errorLogs = append(errorLogs, logMsg)
		}
		if logMsg.logType == "WARN" {
			warnLogs = append(warnLogs, logMsg)
		}
		if logMsg.logType == "INFO" {
			infoLogs = append(infoLogs, logMsg)
		}

	}
	// Update o []LogMsg of each log type of the parser
	parser.logMessages["ERROR"] = errorLogs
	parser.logMessages["WARN"] = warnLogs
	parser.logMessages["INFO"] = infoLogs

	return nil
}

// print the id that appears most in the logs
func (parser *LogParser) MostFrequentID(logType string) {

	var errorsID []int
	errorLogs := parser.logMessages[logType]

	// acess each log massage and then save in slice of "id"
	for _, logError := range errorLogs {
		errorsID = append(errorsID, logError.userID)
	}

	//make a map to put each id as key as the value will be the cont
	listID := make(map[int]int)
	for _, id := range errorsID {
		listID[id]++
	}

	fmt.Println(listID)

	mostFrequentID := errorsID[0]
	maxCount := 0
	for id, count := range listID {
		if count > maxCount {
			mostFrequentID = id
			maxCount = count
		}
	}

	fmt.Println(mostFrequentID)

}

// func countErrors(lines []string) int {
// 	errorCount := 0
// 	for _, line := range lines {
// 		lineSplit := strings.Split(line, " ")
// 		if lineSplit[2] == "ERROR" {
// 			errorCount++
// 		}
// 	}
// 	return errorCount
// }
