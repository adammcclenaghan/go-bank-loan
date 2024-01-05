package shared_models

// *** Model format for the message queue *** //

/*
CreateLoanMessage represents the data passed to the create loan queue. Consumers will pass this data to the bank API
Here, the ApplicationID refers to the ID that we store in our database. Note that it is distinctly different from the
application_id fields returned by the bank API.
*/
type CreateLoanMessage struct {
	ApplicationID string `json:"application_id" binding:"required"`
	FirstName     string `json:"first_name" binding:"required"`
	LastName      string `json:"last_name" binding:"required"`
}

/*
PollLoanMessage represents the data passed to the poll loan queue.
A consumer is expected to use this data to poll the bank jobs API & update our DB

Note that OurApplicationID refers to the applicationID stored in our persistent storage,
which is distinctly different from the application_id field returned by the bank API.
*/
type PollLoanMessage struct {
	OurApplicationID  string `json:"our_id" binding:"required"`
	BankApplicationID string `json:"application_id" binding:"required"`
}
