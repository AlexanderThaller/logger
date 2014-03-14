package logger

import "errors"

func checkPriority(pr Priority) (err error) {
	err = errors.New("priority does not exist")

	for k := range priorities {
		if k == pr {
			err = nil
			break
		}
	}

	return
}
