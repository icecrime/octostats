package log

import "github.com/Sirupsen/logrus"

var Logger = logrus.New()

func Configure(loglevel string) {
	if level, err := logrus.ParseLevel(loglevel); err != nil {
		Logger.Fatal(err)
	} else {
		Logger.Level = level
	}
}

func LogProgress(item, state string, page int) {
	Logger.WithField("page", page).WithField("page", page).Debugf("loading %s", item)
}
