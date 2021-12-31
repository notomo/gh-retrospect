package outputter

import (
	"fmt"
	"io"

	"github.com/notomo/gh-retrospect/retrospect"
)

type Outputter interface {
	Output(w io.Writer, collected *retrospect.Collected) error
}

func Get(typ string) (Outputter, error) {
	if typ == "json" {
		return &JSON{}, nil
	}
	return nil, fmt.Errorf("unexpected outputter type: %s", typ)
}
