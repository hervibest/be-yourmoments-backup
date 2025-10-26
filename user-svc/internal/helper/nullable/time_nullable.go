package nullable

import "time"

func TimeToString(t *time.Time, layout string) string {
	if t == nil {
		return "-"
	}
	str := t.Format(layout)
	return str
}
