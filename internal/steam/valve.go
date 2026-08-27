package steam

import "sync/atomic"

type Valve struct{ open atomic.Bool }

func (v *Valve) Open()        { v.open.Store(true) }
func (v *Valve) Close()       { v.open.Store(false) }
func (v *Valve) IsOpen() bool { return v.open.Load() }
