package shared_helpers

import "log"

func FailOnError(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
