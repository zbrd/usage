package usage

import (
	"bytes"
	"flag"
	"fmt"
	"testing"
)

func TestPrintDefaults(t *testing.T) {
	tests := []struct {
		tpl  string
		want string
	}{
		{"", ""},
		{"{{.Name}}", "test"},
		{"{{.Exec}}", "/bin/test"},
		{"{{.Author}}", "<author>"},
		{"{{.Usage}}", "  -a\taU\n  -b\tbU\n"},
	}

	for i, tt := range tests {
		name := fmt.Sprintf("test_%d", i)

		t.Run(name, func(t *testing.T) {
			buff := bytes.NewBufferString("")
			fset := flag.NewFlagSet(name, flag.PanicOnError)
			fset.SetOutput(buff)
			_ = fset.Bool("a", false, "aU")
			_ = fset.Bool("b", false, "bU")

			prg := ProgExec(fset, "/bin/test")
			err := prg.PrintUsage(tt.tpl, Data{"Author": "<author>"})

			t.Logf("output = %#v", buff.String())

			if err != nil {
				t.Errorf("err != nil: %v", err)
			}

			if got := buff.String(); got != tt.want {
				t.Errorf("expect %#v, got %#v", tt.want, got)
			}
		})
	}
}
