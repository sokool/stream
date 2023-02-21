package model

type MemberID string

type Members map[MemberID]bool

func (m Members) join(n MemberID) {
	m[n] = true
}

func (m Members) isPresent(n MemberID) bool {
	_, found := m[n]
	return found
}

func (m Members) isKicked(n MemberID) bool {
	if active, found := m[n]; !active && found {
		return true
	}

	return false
}

func (m Members) remove(n MemberID) {
	delete(m, n)
}

func (m Members) kick(n MemberID) {
	m[n] = false
}

func (m Members) count() int {
	return len(m)
}
