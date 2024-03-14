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
	timestamp time.Time
	userID    int
	message   string
}

// LogParser is a struct that holds the log files and the errors, warns and infos
type LogParser struct {
	logFolder string
	logFiles  []string
	errors    []LogMsg
	warns     []LogMsg
	infos     []LogMsg
}

func New(logFolder string) *LogParser {
	return &LogParser{
		logFolder: logFolder,
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
	for _, logError := range parser.errors {
		fmt.Println(logError.timestamp, logError.userID, logError.message)
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

// ParseLogLine parses a log line and returns a LogMsg
func parseLogLine(line string) (LogMsg, error) {
	lineSplit := strings.Split(line, " ")
	if len(lineSplit) < 3 {
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
	// continue...

	logMessage := strings.Join(lineSplit[2:], " ")
	if len(lineSplit[2:]) < 3 {
		return LogMsg{}, fmt.Errorf("invalid log message: %s", lineSplit[2:])

	}

	return LogMsg{
		timestamp: timestamp,
		userID:    userID,
		message:   logMessage}, nil

}

func (parser *LogParser) ParseLines(fileName string) error {
	// open the file
	fp, err := parser.openFile(fileName)
	if err != nil {
		return err
	}
	defer fp.Close()
	parser.logFiles = append(parser.logFiles, fileName)

	// read the file
	var lines []string
	scanner := bufio.NewScanner(fp)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	// parse the file in []logMsg
	for _, line := range lines {

		lm, err := parseLogLine(line)
		if err != nil {
			return err
		}

		parser.defineLogMsgType(lm)

		// logMsgType := strings.Split(logMsg.message, " ")
		// if logMsgType[0] == "ERROR" {
		// 	parser.errors = append(parser.errors, logMsg)
		// }
		// if logMsgType[0] == "WARN" {
		// 	parser.warns = append(parser.warns, logMsg)
		// }
		// if logMsgType[0] == "INFO" {
		// 	parser.infos = append(parser.infos, logMsg)
		// }
	}

	return nil
}

func (parser *LogParser) defineLogMsgType(lm LogMsg) {
	logMsgType := strings.Split(lm.message, " ")

	if logMsgType[0] == "ERROR" {
		parser.errors = append(parser.errors, lm)
	}
	if logMsgType[0] == "WARN" {
		parser.warns = append(parser.warns, lm)
	}
	if logMsgType[0] == "INFO" {
		parser.infos = append(parser.infos, lm)
	}
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