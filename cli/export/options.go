package export

const Out = "./out.json"

type ExportOptions struct {
	Path   string
	Export bool
}

type ExportOptFunc func(*ExportOptions)

func defaultExportOps() ExportOptions {
	return ExportOptions{
		Path:   Out,
		Export: false,
	}
}

func NewExportOpts(opts ...ExportOptFunc) ExportOptions {
	o := defaultExportOps()

	for _, fn := range opts {
		fn(&o)
	}

	return o
}

func WithExportPath(path string) ExportOptFunc {
	return func(opts *ExportOptions) {
		opts.Path = path
	}
}

func WithExport(opts *ExportOptions) {
	opts.Export = true
}
