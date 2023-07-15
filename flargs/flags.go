package flargs

import (
	"strconv"
	"time"
)

type flagMeta struct {
	long         string
	short        string
	required     bool
	description  string
	defaultValue string
}

type strFlag struct {
	flagMeta
	strValue string
}

func newStrFlag(long string, short string, required bool, description string, defaultValue string) strFlag {
	return strFlag{
		flagMeta: flagMeta{
			long:         long,
			short:        short,
			required:     required,
			description:  description,
			defaultValue: defaultValue,
		},
		strValue: defaultValue,
	}
}

type boolFlag struct {
	strFlag
}

func newBoolFlag(long string, short string, required bool, description string, defaultValue string) boolFlag {
	return boolFlag{strFlag: newStrFlag(long, short, required, description, defaultValue)}
}

func (f boolFlag) value() (bool, error) {
	val, err := strconv.ParseBool(f.strValue)
	if err != nil {
		return false, err
	}
	return val, nil
}

type intFlag struct {
	strFlag
}

func newIntFlag(long string, short string, required bool, description string, defaultValue string) intFlag {
	return intFlag{strFlag: newStrFlag(long, short, required, description, defaultValue)}
}

func (f intFlag) value() (int, error) {
	val, err := strconv.Atoi(f.strValue)
	if err != nil {
		return 0, err
	}
	return val, nil
}

type uintFlag struct {
	strFlag
}

func newUintFlag(long string, short string, required bool, description string, defaultValue string) uintFlag {
	return uintFlag{strFlag: newStrFlag(long, short, required, description, defaultValue)}
}

func (f uintFlag) value() (uint64, error) {
	val, err := strconv.ParseUint(f.strValue, 10, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

type floatFlag struct {
	strFlag
}

func newFloatFlag(long string, short string, required bool, description string, defaultValue string) floatFlag {
	return floatFlag{strFlag: newStrFlag(long, short, required, description, defaultValue)}
}

func (f floatFlag) value() (float64, error) {
	val, err := strconv.ParseFloat(f.strValue, 64)
	if err != nil {
		return 0, err
	}
	return val, nil
}

type durationFlag struct {
	strFlag
}

func newDurationFlag(long string, short string, required bool, description string, defaultValue string) durationFlag {
	return durationFlag{strFlag: newStrFlag(long, short, required, description, defaultValue)}
}

func (f durationFlag) value() (time.Duration, error) {
	val, err := time.ParseDuration(f.strValue)
	if err != nil {
		return 0, err
	}
	return val, nil
}
