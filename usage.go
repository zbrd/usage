// Generate a program's usage string from a template and FlagSet.
package usage

import (
	"bytes"
	"io"
	"maps"
	"path"
	"text/template"
)

// Data is the root data context of a usage template.
type Data map[string]any

// FlagSet is an interface matching a subset of [flag.FlagSet].
// By using an interface, caller can provide a different drop-in
// implementation of FlagSet such as [pflag].
//
// [pflag]: https://pkg.go.dev/github.com/spf13/pflag
type FlagSet interface {
	Name() string
	Output() io.Writer
	SetOutput(io.Writer)
	PrintDefaults()
}

// Program wraps a [FlagSet], adding functionality to print its
// usage string using a template (using [Program.PrintUsage])
type Program struct {
	FlagSet
}

// Create a new [Program] that wraps the [FlagSet] fs.
func Prog(fs FlagSet) Program {
	return Program{fs}
}

// Base returns the program's executable base name (without dirs)
func (p Program) Base() string {
	return path.Base(p.Name())
}

// Data returns the root template context [Data], with these
// predefined keys:
//   - Base: the program's executable base name (without dirs)
//   - Name: the program's full executable name (from [os.Args])
//   - Usage: the full usage string from [FlagSet.PrintDefaults]
//
// Additional data can be set by passing other [Data]s.
func (p Program) Data(other ...Data) Data {
	d := Data{
		"Base":  p.Base(),
		"Name":  p.Name(),
		"Usage": getUsage(p.FlagSet),
	}

	for _, mp := range other {
		maps.Copy(d, mp)
	}

	return d
}

// Prints the program's usage using given template tpl and
// additional data d. Writes to the [Program]'s [FlagSet.Output].
func (p Program) PrintUsage(tpl string, d ...Data) error {
	out := p.Output()
	data := p.Data(d...)

	if tpl, err := template.New("usage").Parse(tpl); err != nil {
		return err
	} else if err := tpl.Execute(out, data); err != nil {
		return err
	} else {
		return nil
	}
}

func getUsage(fs FlagSet) string {
	var (
		buf bytes.Buffer
		out = fs.Output()
	)

	fs.SetOutput(&buf)
	fs.PrintDefaults()
	fs.SetOutput(out)

	return buf.String()
}
