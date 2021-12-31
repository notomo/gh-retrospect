package outputter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/notomo/gh-retrospect/retrospect"
)

type JSON struct{}

func (p *JSON) Output(w io.Writer, collected *retrospect.Collected) error {
	encorder := json.NewEncoder(w)
	encorder.SetIndent("", "  ")
	if err := encorder.Encode(collected); err != nil {
		return fmt.Errorf("encode json: %w", err)
	}
	return nil
}
