package generator

// GoGivensOutputGenerator is an interface that can be implemented by anything that can generate file content to be output
// after a test has completed.
type GoGivensOutputGenerator interface{
	Generate(data *PageData) (output string)
	FileExtension() string
}