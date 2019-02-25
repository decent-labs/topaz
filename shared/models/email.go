package models

// MarketingSiteEmail ...
type MarketingSiteEmail struct {
	Email string `json:"email"`
}

// MarketingSiteEmails ...
type MarketingSiteEmails []MarketingSiteEmail

// SendgridNewContactRespone ...
type SendgridNewContactRespone struct {
	NewCount            int      `json:"new_count"`
	UpdatedCount        int      `json:"updated_count"`
	ErrorCount          int      `json:"error_count"`
	ErrorIndices        []int    `json:"error_indices"`
	UnmodifiedIndices   []int    `json:"unmodified_indices"`
	PersistedRecipients []string `json:"persisted_recipients"`
	Errors              []string `json:"errors"`
}
