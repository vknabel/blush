package token

type Source struct {
	ModuleName string
	FileName   string
	Offset     int
}

func MakeSource(
	moduleName string,
	fileName string,
	offset int,
) *Source {
	return &Source{
		ModuleName: moduleName,
		FileName:   fileName,
		Offset:     offset,
	}
}
