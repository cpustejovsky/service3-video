package testfmt

import (
	"testing"
)

type Marker string

var SuccessMarker Marker = "\u2713"
var FailedMarker Marker = "\u2717"

func Info(t *testing.T, id int, msg string) {
	t.Logf("\tTest %d:\t%v", id, msg)
}

func Error(t *testing.T, id int, msg string, err error) {
	t.Errorf("\t%s\tTest %d:\t%v: %v", FailedMarker, id, msg, err)
}

func Success(t *testing.T, id int, msg string) {
	t.Logf("\t%s\tTest %d:\t%v.", SuccessMarker, id, msg)
}
