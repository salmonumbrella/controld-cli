package ui

import (
	"context"
	"fmt"
	"os"

	"golang.org/x/term"
)

type ctxKey struct{}

type UI struct {
	color string
}

func New(color string) *UI {
	return &UI{color: color}
}

func WithUI(ctx context.Context, u *UI) context.Context {
	return context.WithValue(ctx, ctxKey{}, u)
}

func FromContext(ctx context.Context) *UI {
	u, _ := ctx.Value(ctxKey{}).(*UI)
	if u == nil {
		return New("auto")
	}
	return u
}

func (u *UI) useColor() bool {
	switch u.color {
	case "always":
		return true
	case "never":
		return false
	default:
		return term.IsTerminal(int(os.Stdout.Fd()))
	}
}

func (u *UI) Success(msg string) {
	if u.useColor() {
		fmt.Fprintf(os.Stderr, "\033[32m✓\033[0m %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "✓ %s\n", msg)
	}
}

func (u *UI) Error(msg string) {
	if u.useColor() {
		fmt.Fprintf(os.Stderr, "\033[31m✗\033[0m %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "✗ %s\n", msg)
	}
}

func (u *UI) Info(msg string) {
	fmt.Fprintf(os.Stderr, "%s\n", msg)
}

func (u *UI) Warn(msg string) {
	if u.useColor() {
		fmt.Fprintf(os.Stderr, "\033[33m!\033[0m %s\n", msg)
	} else {
		fmt.Fprintf(os.Stderr, "! %s\n", msg)
	}
}
