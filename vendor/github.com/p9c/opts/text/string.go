package text

import (
	"encoding/json"
	"fmt"
	"github.com/p9c/opts/meta"
	"github.com/p9c/opts/opt"
	"strings"
	"sync/atomic"
)

// Opt stores a string configuration value
type Opt struct {
	meta.Data
	hook  []func(s string)
	Value *atomic.Value
	Def   string
}

// New creates a new Opt with a given default value set
func New(m meta.Data, def string) *Opt {
	v := &atomic.Value{}
	v.Store([]byte(def))
	return &Opt{Value: v, Data: m, Def: def}
}

// SetName sets the name for the generator
func (x *Opt) SetName(name string) {
	x.Data.Option = strings.ToLower(name)
	x.Data.Name = name
}

// Type returns the receiver wrapped in an interface for identifying its type
func (x *Opt) Type() interface{} {
	return x
}

// GetMetadata returns the metadata of the opt type
func (x *Opt) GetMetadata() *meta.Data {
	return &x.Data
}

// ReadInput sets the value from a string
func (x *Opt) ReadInput(input string) (o opt.Option, e error) {
	if input == "" {
		e = fmt.Errorf("string opt %s %v may not be empty", x.Name(), x.Data.Aliases)
		return
	}
	if strings.HasPrefix(input, "=") {
		// the following removes leading and trailing characters
		input = strings.Join(strings.Split(input, "=")[1:], "=")
	}
	x.Set(input)
	return x, e
}

// LoadInput sets the value from a string
func (x *Opt) LoadInput(input string) (o opt.Option, e error) {
	return x.ReadInput(input)
}

// Name returns the name of the opt
func (x *Opt) Name() string {
	return x.Data.Option
}

// AddHooks appends callback hooks to be run when the value is changed
func (x *Opt) AddHooks(hook ...func(f string)) {
	x.hook = append(x.hook, hook...)
}

// SetHooks sets a new slice of hooks
func (x *Opt) SetHooks(hook ...func(f string)) {
	x.hook = hook
}

// V returns the stored string
func (x *Opt) V() string {
	return string(x.Value.Load().([]byte))
}

// Empty returns true if the string is empty
func (x *Opt) Empty() bool {
	return len(x.Value.Load().(string)) == 0
}

// Bytes returns the raw bytes in the underlying storage
func (x *Opt) Bytes() []byte {
	return x.Value.Load().([]byte)
}

func (x *Opt) runHooks() {
	for i := range x.hook {
		x.hook[i](x.V())
	}
}

// Set the value stored
func (x *Opt) Set(s string) *Opt {
	x.runHooks()
	x.Value.Store([]byte(s))
	return x
}

// SetBytes sets the string from bytes
func (x *Opt) SetBytes(s []byte) *Opt {
	x.runHooks()
	x.Value.Store(s)
	return x
}

// Opt returns a string representation of the value
func (x *Opt) String() string {
	return fmt.Sprintf("%s: '%s'", x.Data.Option, x.V())
}

// MarshalJSON returns the json representation
func (x *Opt) MarshalJSON() (b []byte, e error) {
	v := string(x.Value.Load().([]byte))
	return json.Marshal(&v)
}

// UnmarshalJSON decodes a JSON representation
func (x *Opt) UnmarshalJSON(data []byte) (e error) {
	v := x.Value.Load().([]byte)
	e = json.Unmarshal(data, &v)
	x.Value.Store(v)
	return
}
