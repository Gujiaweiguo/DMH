package member

import (
	"errors"
	"strconv"
	"strings"
)

func parseMemberIDFromPath(path string) (int64, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	for i, part := range parts {
		if part == "members" && i+1 < len(parts) {
			memberID, err := strconv.ParseInt(parts[i+1], 10, 64)
			if err != nil {
				return 0, err
			}
			if memberID <= 0 {
				return 0, errors.New("member id must be greater than 0")
			}
			return memberID, nil
		}
	}

	return 0, errors.New("member id is required")
}
