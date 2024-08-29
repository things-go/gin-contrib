package cache

var _ Logger = (*Discard)(nil)

// Discard is an logger on which all Write calls succeed
// without doing anything.
type Discard struct{}

// NewDiscard a discard logger on which always succeed without doing anything
func NewDiscard() Discard { return Discard{} }

// Debugf implement Logger interface.
func (d Discard) Debugf(string, ...any) {}

// Infof implement Logger interface.
func (d Discard) Infof(string, ...any) {}

// Errorf implement Logger interface.
func (d Discard) Errorf(string, ...any) {}

// Warnf implement Logger interface.
func (d Discard) Warnf(string, ...any) {}

// DPanicf implement Logger interface.
func (d Discard) DPanicf(string, ...any) {}

// Fatalf implement Logger interface.
func (d Discard) Fatalf(string, ...any) {}
