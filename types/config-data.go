package types

type ConfigData struct {
	JiraUrl          string `json:"jira_url"`
	JiraEmail        string `json:"jira_email"`
	JiraKey          string `json:"jira_key"`
	TogglKey         string `json:"toggl_key"`
	TogglWorkspaceId string `json:"toggl_workspace_id"`
}
