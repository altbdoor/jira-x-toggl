package types

type JiraResponse struct {
	Issues []JiraResponseIssue `json:"issues"`
	Total  int                 `json:"total"`
}

type JiraResponseIssue struct {
	Key    string                  `json:"key"`
	Fields JiraResponseIssueFields `json:"fields"`
}

type JiraResponseIssueFields struct {
	Summary  string                          `json:"summary"`
	Estimate int                             `json:"timeestimate,omitempty"`
	Assignee JiraResponseIssueFieldsAssignee `json:"assignee,omitempty"`
}

type JiraResponseIssueFieldsAssignee struct {
	Email string `json:"emailAddress"`
}
