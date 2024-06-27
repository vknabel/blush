package ast

type ModuleName StaticReference

type ContextModule struct {
	Name    ModuleName
	Symbols *SymbolTable

	Files []*SourceFile
}

func MakeContextModule(name ModuleName) *ContextModule {
	m := &ContextModule{
		Name:  name,
		Files: []*SourceFile{},
	}
	m.Symbols = MakeSymbolTable(m.Symbols.Parent, nil)
	return m
}

func (m *ContextModule) AddSourceFile(sourceFile *SourceFile) {
	m.Files = append(m.Files, sourceFile)
}
