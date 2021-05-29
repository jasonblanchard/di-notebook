package app

func canGetEntry(p *Principal, entry *Entry) bool {
	switch p.Type {
	case PrincipalUSER:
		return p.ID == entry.CreatorID
	default:
		return false
	}
}

func canDeleteEntry(p *Principal, entry *Entry) bool {
	return p.ID == entry.CreatorID
}

func canUpdateEntry(p *Principal, entry *Entry) bool {
	return p.ID == entry.CreatorID
}

func canListEntries(p *Principal, creatorID string) bool {
	return p.ID == creatorID
}

func canUndeleteEntry(p *Principal, entry *Entry) bool {
	return p.ID == entry.CreatorID
}
