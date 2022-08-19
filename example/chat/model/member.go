package model

type Member string

type Members map[Member]bool

func (m Members) join(n Member) {
	m[n] = true
}

func (m Members) isPresent(n Member) bool {
	_, found := m[n]
	return found
}

func (m Members) isKicked(n Member) bool {
	if active, found := m[n]; !active && found {
		return true
	}

	return false
}

func (m Members) remove(n Member) {
	delete(m, n)
}

func (m Members) kick(n Member) {
	m[n] = false
}

func (m Members) count() int {
	return len(m)
}
