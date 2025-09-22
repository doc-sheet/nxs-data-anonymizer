package misc

import (
	"bytes"
	ttemplate "text/template"

	"github.com/Masterminds/sprig/v3"
)

var (
	TemplateNULL = "::NULL::"
	TemplateDrop = "::DROP::"
)

type TemlateRes struct {
	Value   string
	DropRow bool
}

// TemplateExec makes message from given template `tpl` and data `d`
func TemplateExec(tpl string, d any) (TemlateRes, error) {

	var b bytes.Buffer

	// See http://masterminds.github.io/sprig/ for details
	t, err := ttemplate.New("template").Funcs(func() ttemplate.FuncMap {

		// Get current sprig functions
		t := sprig.TxtFuncMap()

		// Add additional functions
		t["null"] = func() string {
			return TemplateNULL
		}
		t["isNull"] = func(v string) bool {
			if v == TemplateNULL {
				return true
			}
			return false
		}
		t["drop"] = func() string {
			return TemplateDrop
		}

		return t
	}()).Parse(tpl)
	if err != nil {
		return TemlateRes{}, err
	}

	err = t.Execute(&b, d)
	if err != nil {
		return TemlateRes{}, err
	}

	// Return empty line if buffer is nil
	if b.Bytes() == nil {
		return TemlateRes{
				Value:   "",
				DropRow: false,
			},
			nil
	}

	// Return `drop` value if buffer is DROP (with special key)
	if bytes.Equal(b.Bytes(), []byte(TemplateDrop)) {
		return TemlateRes{
				Value:   "",
				DropRow: true,
			},
			nil
	}

	// Return buffer content otherwise
	return TemlateRes{
			Value:   b.String(),
			DropRow: false,
		},
		nil
}
