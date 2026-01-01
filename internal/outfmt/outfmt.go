package outfmt

import (
	"context"
	"encoding/json"
	"io"
	"text/tabwriter"
)

type formatKey struct{}
type yesKey struct{}

func WithFormat(ctx context.Context, format string) context.Context {
	return context.WithValue(ctx, formatKey{}, format)
}

func WithYes(ctx context.Context, yes bool) context.Context {
	return context.WithValue(ctx, yesKey{}, yes)
}

func IsJSON(ctx context.Context) bool {
	format, _ := ctx.Value(formatKey{}).(string)
	return format == "json"
}

func GetYes(ctx context.Context) bool {
	yes, _ := ctx.Value(yesKey{}).(bool)
	return yes
}

func WriteJSON(w io.Writer, v interface{}) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func NewTabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 4, 2, ' ', 0)
}
