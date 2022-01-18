package testfmt

import "fmt"

type Marker string
var SuccessMarker Marker = "\u2713"
var FailedMarker Marker = "\u2717"

func Info(id int, msg string) string {
	return fmt.Sprintf("\tTest %d:\t%v", id, msg)
}

func Error(id int, msg string, err error) string {
	return fmt.Sprintf("\t%s\tTest %d:\t%v: %v", FailedMarker, id, msg, err)
}

func Success(id int, msg string) string {
	return fmt.Sprintf("\t%s\tTest %d:\t%v.", SuccessMarker, id, msg)
}
