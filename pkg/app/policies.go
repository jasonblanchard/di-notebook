package app

func canReadEntry(p *Principal, entry *Entry) bool {
	switch p.Type {
	case PrincipalUSER:
		return p.ID == entry.CreatorID
	default:
		return false
	}
}

func canDiscardEntry(p *Principal, entry *Entry) bool {
	return p.ID == entry.CreatorID
}

func canChangeEntry(p *Principal, entry *Entry) bool {
	return p.ID == entry.CreatorID
}

func canListEntries(p *Principal, creatorID string) bool {
	return p.ID == creatorID
}

func canUndeleteEntry(p *Principal, entry *Entry) bool {
	return p.ID == entry.CreatorID
}
