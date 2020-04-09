package entities

func IsPending(status string) bool {
	pendingStatuses := []string{string(ToDo), string(InProgress), string(Stopped)}
	for _, val := range pendingStatuses {
		if val == status {
			return true
		}
	}
	return false
}
