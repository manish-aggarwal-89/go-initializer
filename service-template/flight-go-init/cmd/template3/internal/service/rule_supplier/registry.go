package rule_supplier

import (
	"{{MODULE_NAME}}/internal/service/rule_supplier/common"
)

type SearchRuleRegistry struct {
	rules map[string]common.SearchRule
}

func NewSearchRuleRegistry() *SearchRuleRegistry {
	return &SearchRuleRegistry{rules: make(map[string]common.SearchRule)}
}

func (r *SearchRuleRegistry) Register(distributionType string, rule common.SearchRule) {
	r.rules[distributionType] = rule
}

func (r *SearchRuleRegistry) Get(distributionType string) (common.SearchRule, bool) {
	rule, ok := r.rules[distributionType]
	return rule, ok
}

type BookingRuleRegistry struct {
	rules map[string]common.BookingRule
}

func NewBookingRuleRegistry() *BookingRuleRegistry {
	return &BookingRuleRegistry{rules: make(map[string]common.BookingRule)}
}

func (r *BookingRuleRegistry) Register(distributionType string, rule common.BookingRule) {
	r.rules[distributionType] = rule
}

func (r *BookingRuleRegistry) Get(distributionType string) (common.BookingRule, bool) {
	rule, ok := r.rules[distributionType]
	return rule, ok
}

type IssuedRuleRegistry struct {
	rules map[string]common.IssuedRule
}

func NewIssuedRuleRegistry() *IssuedRuleRegistry {
	return &IssuedRuleRegistry{rules: make(map[string]common.IssuedRule)}
}

func (r *IssuedRuleRegistry) Register(distributionType string, rule common.IssuedRule) {
	r.rules[distributionType] = rule
}

func (r *IssuedRuleRegistry) Get(distributionType string) (common.IssuedRule, bool) {
	rule, ok := r.rules[distributionType]
	return rule, ok
}

type CancelBookingRuleRegistry struct {
	rules map[string]common.CancelBookingRule
}

func NewCancelBookingRuleRegistry() *CancelBookingRuleRegistry {
	return &CancelBookingRuleRegistry{rules: make(map[string]common.CancelBookingRule)}
}

func (r *CancelBookingRuleRegistry) Register(distributionType string, rule common.CancelBookingRule) {
	r.rules[distributionType] = rule
}

func (r *CancelBookingRuleRegistry) Get(distributionType string) (common.CancelBookingRule, bool) {
	rule, ok := r.rules[distributionType]
	return rule, ok
}
