package flargs

type Flag[T any] interface {
	Value() T
}

type FlagMeta struct {
	long        string
	short       string
	required    bool
	description string
}

type BoolFlag struct {
	FlagMeta
	value bool
}

type StrFlag struct {
	FlagMeta
	value string
}

type IntFlag struct {
	StrFlag
	value int
}

func newIntFlag(long string, short string, required bool, defaultValue string, description string) IntFlag {
	return IntFlag{
		StrFlag: StrFlag{
			FlagMeta: FlagMeta{
				long:        long,
				short:       short,
				required:    required,
				description: description,
			},
			value: defaultValue,
		},
	}
}
