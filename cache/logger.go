package cache

import "context"

var _ Logger = (*Discard)(nil)

// Discard is an logger on which all Write calls succeed
// without doing anything.
type Discard struct{}

// NewDiscard a discard logger on which always succeed without doing anything
func NewDiscard() Discard { return Discard{} }

// Errorf implement Logger interface.
func (d Discard) Errorf(context.Context, string, ...any) {}
