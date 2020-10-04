package app

func canReadEntry(p *Principal, entry *Entry) bool {
	switch p.Type {
	case PrincipalUSER:
		return p.ID == entry.CreatorID
	default:
		return false
	}
}

func canResetEntries(p *Principal) bool {
	return p.Type == PrincipalTEST
}
