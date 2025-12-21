package models

import "strings"

// FilterOptions defines user-specified filtering criteria for displaying rules
type FilterOptions struct {
	AndKeyFilter  string // Filter by AndKey substring (case-insensitive)
	ChannelFilter string // Filter by channel "ONID-TSID-SID" format
	EnabledOnly   bool   // Show only enabled rules (DisableFlag==0)
	DisabledOnly  bool   // Show only disabled rules (DisableFlag==1)
	RegexOnly     bool   // Show only regex rules (RegExpFlag==1)
}

// Matches returns true if the rule passes all active filters
// All filters are AND-ed together
func (f *FilterOptions) Matches(rule *AutoAddRule) bool {
	// Check enabled/disabled
	if f.EnabledOnly && !rule.SearchSettings.IsEnabled() {
		return false
	}
	if f.DisabledOnly && rule.SearchSettings.IsEnabled() {
		return false
	}

	// Check regex
	if f.RegexOnly && !rule.SearchSettings.IsRegex() {
		return false
	}

	// Check AndKey substring (case-insensitive)
	if f.AndKeyFilter != "" {
		andKeyLower := strings.ToLower(rule.SearchSettings.AndKey)
		filterLower := strings.ToLower(f.AndKeyFilter)
		if !strings.Contains(andKeyLower, filterLower) {
			return false
		}
	}

	// Check channel
	if f.ChannelFilter != "" {
		channelMatch := false
		for _, channel := range rule.SearchSettings.ServiceList {
			if channel.String() == f.ChannelFilter {
				channelMatch = true
				break
			}
		}
		if !channelMatch {
			return false
		}
	}

	return true
}

// HasFilters returns true if any filter is active
func (f *FilterOptions) HasFilters() bool {
	return f.AndKeyFilter != "" || f.ChannelFilter != "" ||
		f.EnabledOnly || f.DisabledOnly || f.RegexOnly
}
