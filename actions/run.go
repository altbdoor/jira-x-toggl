package actions

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"jira-x-toggl/types"

	"github.com/urfave/cli/v2"
)

func RunAction(cCtx *cli.Context) error {
	tmpLog := ""
	isDebug := cCtx.Bool("debug")

	parentIssueId := cCtx.Args().Get(0)
	if parentIssueId == "" {
		return cli.Exit("Please pass in the ID of the parent issue", 1)
	}

	if isDebug {
		fmt.Println("(i) parentIssueId=" + parentIssueId)
	}

	configPath := cCtx.String("config")
	configFile, err := os.Open(configPath)
	if err != nil {
		tmpLog = fmt.Sprintf("Unable to open config file at %s", configPath)
		return cli.Exit(tmpLog, 1)
	}

	var configData types.ConfigData
	configBytes, _ := io.ReadAll(configFile)
	json.Unmarshal(configBytes, &configData)
	configFile.Close()

	client := &http.Client{}

	skipFetchToggl := cCtx.Bool("skip-fetch-toggl")
	if skipFetchToggl {
		if isDebug {
			fmt.Println("(i) skip fetching toggl data")
		}
	} else {
		if isDebug {
			fmt.Println("(i) fetching toggl data")
		}

		today := time.Now().Local()

		startDate := today.AddDate(0, 0, cCtx.Int("start")*-1)
		endDate := today.AddDate(0, 0, cCtx.Int("end")*-1)

		togglPayload := types.TogglPayload{
			StartDate:        startDate.Format("2006-01-02"),
			EndDate:          endDate.Format("2006-01-02"),
			IncludeTimeEntry: false,
			HideRates:        true,
			HideAmounts:      true,
		}
		togglPayloadJson, _ := json.Marshal(togglPayload)

		togglUrl := fmt.Sprint(
			"https://api.track.toggl.com/reports/api/v3/workspace/",
			configData.TogglWorkspaceId,
			"/summary/time_entries",
		)

		req, _ := http.NewRequest("POST", togglUrl, bytes.NewBuffer(togglPayloadJson))
		req.SetBasicAuth(configData.TogglKey, "api_token")
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Accept", "application/json")

		if isDebug {
			dump, _ := httputil.DumpRequest(req, true)
			fmt.Println(string(dump))
		}

		res, _ := client.Do(req)
		file, _ := os.Create("toggl.json")
		io.Copy(file, res.Body)
		file.Close()
		res.Body.Close()
	}

	if isDebug {
		fmt.Println("(i) fetching jira data")
	}

	jiraPayload := types.JiraPayload{
		Jql:        fmt.Sprintf("parentEpic=%s ORDER BY key", parentIssueId),
		MaxResults: 100,
		Fields:     []string{"parent", "assignee", "summary", "timeestimate"},
	}
	jiraPayloadJson, _ := json.Marshal(jiraPayload)
	jiraUrl := fmt.Sprint(configData.JiraUrl, "/rest/api/3/search")

	req, _ := http.NewRequest("POST", jiraUrl, bytes.NewBuffer(jiraPayloadJson))
	req.SetBasicAuth(configData.JiraEmail, configData.JiraKey)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	if isDebug {
		dump, _ := httputil.DumpRequest(req, true)
		fmt.Println(string(dump))
	}

	res, _ := client.Do(req)

	jiraResponse := &types.JiraResponse{}
	json.NewDecoder(res.Body).Decode(jiraResponse)
	res.Body.Close()

	if isDebug {
		tmpLog = fmt.Sprintf("(i) processing %d jira issues", jiraResponse.Total)
		fmt.Println(tmpLog)
	}

	file, _ := os.Open("toggl.json")
	togglData := &types.TogglResponse{}
	json.NewDecoder(file).Decode(togglData)
	file.Close()

	csvFileName := fmt.Sprintf("%s.csv", parentIssueId)
	csvFile, _ := os.Create(csvFileName)
	csvWriter := csv.NewWriter(csvFile)
	csvWriter.Write([]string{"key", "title", "assignee", "estimate_seconds", "total_time_seconds"})
	csvWriter.Flush()

	for _, issue := range jiraResponse.Issues {
		issueSeconds := 0

		for _, group := range togglData.Groups {
			for _, subGroup := range group.SubGroup {
				if strings.HasPrefix(subGroup.Title, issue.Key+" ") {
					issueSeconds += subGroup.Seconds
				}
			}
		}

		csvWriter.Write([]string{
			issue.Key,
			issue.Fields.Summary,
			issue.Fields.Assignee.Email,
			strconv.Itoa(issue.Fields.Estimate),
			strconv.Itoa(issueSeconds),
		})
		csvWriter.Flush()
	}

	csvFile.Close()

	if isDebug {
		tmpLog = fmt.Sprintf("(i) finish process in %s", csvFileName)
		fmt.Println(tmpLog)
	}

	return nil
}
