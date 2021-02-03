package fixedtime

import "time"

type NowFunc func() time.Time

var nowFunc NowFunc = time.Now

func Set(f NowFunc) {
	nowFunc = f
}

func Now() time.Time {
	return nowFunc()
}
