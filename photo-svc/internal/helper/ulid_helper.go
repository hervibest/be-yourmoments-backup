package helper

import "github.com/oklog/ulid/v2"

func ParseMultipleULID(ids ...string) error {
	for _, id := range ids {
		if _, err := ulid.Parse(id); err != nil {
			return err
		}
	}
	return nil
}
