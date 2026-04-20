package scanner

type Registry struct {
	modules map[ModuleName]Scanner
}

func NewRegistry() *Registry {
	items := []Scanner{
		PortsScanner{},
		HeadersScanner{},
		TLSScanner{},
		FuzzScanner{},
		XSSScanner{},
		SQLiScanner{},
		CVEScanner{},
	}

	modules := make(map[ModuleName]Scanner, len(items))
	for _, item := range items {
		modules[item.Name()] = item
	}

	return &Registry{modules: modules}
}

func (r *Registry) Resolve(names []ModuleName) []Scanner {
	if len(names) == 0 {
		names = DefaultModules
	}

	scanners := make([]Scanner, 0, len(names))
	for _, name := range names {
		if scanner, ok := r.modules[name]; ok {
			scanners = append(scanners, scanner)
		}
	}
	return scanners
}

