package log

import (
    "io/ioutil"
    "strings"

    "github.com/Sirupsen/logrus"
)

var logger = logrus.New()


func Debug(
