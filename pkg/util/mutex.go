package util

type Mutex struct {
	locked chan struct{}
}

func NewMutex() *Mutex {
	return &Mutex{
		locked: make(chan struct{}, 1),
	}
}

func (m *Mutex) TryLock() bool {
	select {
	case m.locked <- struct{}{}:
		return true
	default:
		return false
	}
}

func (m *Mutex) Unlock() {
	select {
	case <-m.locked:
		return
	default:
		panic("Releasing unacquired lock")
	}
}

func (m *Mutex) Lock() {
	m.locked <- struct{}{}
}
