package models

import "encoding/xml"

// AutoAddRuleResponse represents the XML response from SetAutoAdd API
// Response format:
//   Success: <?xml version="1.0" encoding="UTF-8" ?><entry><success>EPG自動予約を追加しました</success></entry>
//   Error:   <?xml version="1.0" encoding="UTF-8" ?><entry><err>不正値入力</err></entry>
type AutoAddRuleResponse struct {
	XMLName xml.Name `xml:"entry"`
	Success string   `xml:"success"` // Success message
	Error   string   `xml:"err"`     // Error message
}

// IsSuccess returns true if the API call was successful
func (r *AutoAddRuleResponse) IsSuccess() bool {
	return r.Success != "" && r.Error == ""
}

// GetError returns the error message if any
func (r *AutoAddRuleResponse) GetError() string {
	return r.Error
}

// GetMessage returns the success message if any
func (r *AutoAddRuleResponse) GetMessage() string {
	return r.Success
}
