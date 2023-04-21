package parsers

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v2"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	YAML = "yaml"
	JSON = "json"
)

func ParserConfigurationByFile(format, in string, out interface{}) error {
	data, err := fs.ReadFile(os.DirFS("."), in)

	if err != nil {
		return err
	}

	switch format {
	case YAML:
		return yaml.Unmarshal(data, out)
	case JSON:
		return json.Unmarshal(data, out)
	default:
		return errors.New("invalid file format")
	}
}

func ParseDuration(durationStr string) (time.Duration, error) {

	timeAttr := string(durationStr[len(durationStr)-1])
	realTime := strings.Split(durationStr, timeAttr)

	if len(realTime) != 2 {
		return 0, errors.New("invalid time str")
	}

	var (
		n, _ = strconv.Atoi(realTime[0])
		dur  = time.Duration(n) * time.Second
	)

	switch timeAttr {
	case "m":
		dur *= 60
	case "h":
		dur *= 60 * 60
	case "d":
		dur *= 60 * 60 * 24
	case "y":
		dur *= 60 * 60 * 24 * 365
	}
	return dur, nil
}

func TimeParse(times string) (time.Duration, error) {
	if dur, err := ParseDuration(times); err != nil {
		return 0, err
	} else {
		return dur, nil
	}
}
