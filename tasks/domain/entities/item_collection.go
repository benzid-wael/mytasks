package entities

import (
	"time"
)

type ItemCollection []Manageable
type ItemPredicate func(item Manageable) bool

// Sort Protocol
func (c ItemCollection) Len() int           { return len(c) }
func (c ItemCollection) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }
func (c ItemCollection) Less(i, j int) bool { return c[i].GetId() < c[j].GetId() }

func (c ItemCollection) Filter(predicate ItemPredicate) ItemCollection {
	result := make(ItemCollection, 0, len(c))
	for _, item := range c {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func (c ItemCollection) Exclude(predicate ItemPredicate) ItemCollection {
	result := make(ItemCollection, 0, len(c))
	for _, item := range c {
		if !predicate(item) {
			result = append(result, item)
		}
	}
	return result
}

func (c ItemCollection) FilterByStatus(status string) ItemCollection {
	return c.Filter(func(item Manageable) bool {
		return item.GetStatus() == status
	})
}

func (c ItemCollection) FilterPending() ItemCollection {
	return c.Filter(func(item Manageable) bool {
		return IsPending(item.GetStatus())
	})
}

func (c ItemCollection) FilterByTag(tag string) ItemCollection {
	return c.Filter(func(item Manageable) bool {
		return item.HasTag(tag)
	})
}

func (c ItemCollection) FilterByTags(tags ...string) ItemCollection {
	return c.Filter(func(item Manageable) bool {
		return item.HasAnyTag(tags...)
	})
}

func (c ItemCollection) FilterByType(kind string) ItemCollection {
	return c.Filter(func(item Manageable) bool {
		return item.GetType() == kind
	})
}

func (c ItemCollection) FilterByCreationDate(fromDate, toDate *time.Time) ItemCollection {
	return c.Filter(func(item Manageable) bool {
		createdAt := item.GetCreationDateTime()
		if fromDate != nil && toDate != nil {
			return fromDate.Before(createdAt) && toDate.After(createdAt)
		} else if fromDate != nil {
			return fromDate.Before(createdAt)
		} else if toDate != nil {
			return toDate.After(createdAt)
		}
		return true
	})
}
