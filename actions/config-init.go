package actions

import (
	"encoding/json"
	"fmt"
	"jira-x-toggl/types"
	"os"

	"github.com/urfave/cli/v2"
)

func ConfigInitAction(cCtx *cli.Context) error {
	tmpLog := ""
	isDebug := cCtx.Bool("debug")
	configPath := cCtx.String("config")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		file, _ := os.Create(configPath)
		configData := &types.ConfigData{}
		jsonData, _ := json.MarshalIndent(configData, "", "  ")
		file.Write(jsonData)
		file.Close()

		if isDebug {
			tmpLog = fmt.Sprintf("(i) created empty config file in %s", configPath)
			fmt.Println(tmpLog)
		}
	} else {
		tmpLog = fmt.Sprintf("File already exists in %s", configPath)
		return cli.Exit(tmpLog, 1)
	}

	return nil
}
