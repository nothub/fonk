package flargs

type CmdBuilder struct {
}

type FlagBuilder struct {
}

func New() CmdBuilder {
	return CmdBuilder{}
}
