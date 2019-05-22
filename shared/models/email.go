package models

// EmailWithList ...
type EmailWithList struct {
	Email string `json:"email"`
	List  string `json:"list"`
}

// SendgridEmail ...
type SendgridEmail struct {
	Email string `json:"email"`
}

// SendgridEmails ...
type SendgridEmails []SendgridEmail

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
