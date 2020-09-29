package app

func canReadEntry(p *Principle, entry *Entry) bool {
	switch p.Type {
	case PrincipleTypeUser:
		return p.ID == entry.CreatorID
	default:
		return false
	}
}

func canResetEntries(p *Principle) bool {
	return p.Type == PrincipleTypeTest
}
