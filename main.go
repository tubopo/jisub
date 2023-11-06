package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

const (
	iniConfig = "/jisub/jisub-config.ini"
)

var (
	version = "dev"

	configFlag   string
	subTasksFlag string
	fieldsFlag   string
	versionFlag  bool
	helpFlag     bool
)

func init() {
	flag.Usage = func() {
		h := "Usage:\n"
		h += "  jisub [OPTIONS] JIRA-39106\n\n"

		h += "Options:\n"
		h += "  -h,  --help 		Print usage\n"
		h += "  -v,  --version 	Print version info\n"
		h += "  -c,  --config 	Create/Update jira configuration\n"
		h += "  -st, --sub-tasks 	Sub tasks to create for provided parent issue\n"
		h += "  -f,  --fields 	Field name, value to update for provided issue\n\n"

		h += "Examples:\n"
		h += "  jisub --config \"jira.url https://jira-api.com/jira/rest/api/2\" \n"
		h += "  jisub --config \"user.token <token>\"\n"
		h += "  jisub --sub-tasks \"QA:2 BE:3 FE:4\" --fields \"storypoints:4 dealsize:2,3,4\" JIRA-39106\n"

		fmt.Fprint(os.Stderr, h)
	}

	// jisub  jira configuration
	flag.StringVar(&configFlag, "config", "", "Create/Update jira configuration")
	flag.StringVar(&configFlag, "c", "", "Create/Update jira configuration")

	// list of issue subtasks
	flag.StringVar(&subTasksFlag, "sub-tasks", "", "Sub tasks to create for provided parent issue")
	flag.StringVar(&subTasksFlag, "st", "", "Sub tasks to create for provided parent issue")

	// fields to update
	flag.StringVar(&fieldsFlag, "fields", "", "Field name, value to update for provided issue")
	flag.StringVar(&fieldsFlag, "f", "", "Field name, value to update for provided issue")

	// version
	flag.BoolVar(&versionFlag, "version", false, "Version information")
	flag.BoolVar(&versionFlag, "v", false, "Version information")

	// help
	flag.BoolVar(&helpFlag, "help", false, "Usage example")
	flag.BoolVar(&helpFlag, "h", false, "Usage example")
}

// > jisub --config "user.token RandomTokenValueStr"
// > jisub --config "jira.url "https://jira-api.com/jira/rest/api/2"

// > jisub --sub-tasks "QA:2 BE:3 FE:4" --fields "storypoints:4 dealsize:2,3,4" JIRA-39106
func main() {
	flag.Parse()

	if versionFlag {
		fmt.Println("jisub version " + version)
		return
	}

	if helpFlag {
		flag.Usage()
		return
	}

	if len(configFlag) > 0 {
		err := updateConfig(configFlag)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(-1)
		}

		return
	}

	err := updateIssue(subTasksFlag, fieldsFlag, flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func updateIssue(subtaskArg, fieldsArg, issueKey string) error {
	if issueKey == "" {
		return fmt.Errorf("missing required issue key")
	}

	jira, err := buildNewJiraFromConfig()
	if err != nil {
		return fmt.Errorf("error creating jira client %w", err)
	}

	issue, err := jira.Issue(issueKey)
	if err != nil {
		return fmt.Errorf("issue not found %v %w", issueKey, err)
	}

	if subtaskArg != "" {
		err = createSubTasks(*jira, *issue, subtaskArg)
		if err != nil {
			return fmt.Errorf("error creating sub tasks %v, %v %w", subtaskArg, issueKey, err)
		}
	}

	if fieldsArg != "" {
		err = updateIssueFields(*jira, *issue, fieldsArg)
		if err != nil {
			return fmt.Errorf("error updating issue fields %v, %v %w", fieldsArg, issueKey, err)
		}
	}

	return nil
}

func createSubTasks(j Jira, parent Issue, subtasksArg string) error {

	subTasksMap := make(map[string]string)
	err := stringToMap(subtasksArg, subTasksMap)
	if err != nil {
		return err
	}

	result, err := j.SubTasks(parent, subTasksMap)
	if err != nil {
		return err
	}

	fmt.Println("sub tasks:")
	for _, issue := range result.Issues {
		fmt.Println(issue.Key)
	}

	return nil
}

func updateIssueFields(j Jira, issue Issue, fieldsArg string) error {

	fieldsUpdatesMap := make(map[string]string)
	err := stringToMap(fieldsArg, fieldsUpdatesMap)
	if err != nil {
		return err
	}

	err = j.IssueUpdate(issue, fieldsUpdatesMap)
	if err != nil {
		return err
	}

	fmt.Printf("issue updated %v\n", issue.Key)

	return nil
}

func updateConfig(arg string) error {
	items := strings.Split(arg, " ")
	// expect key value pair
	if len(items) < 2 {
		return fmt.Errorf("wrong number of arguments provided")
	}

	// load config, in case not loaded create empty
	pwd, _ := os.Getwd()
	cfg, err := ini.Load(pwd + iniConfig)
	if err != nil {
		cfg = ini.Empty()
	}

	// parse config section name, key
	sectionKey := strings.Split(items[0], ".")
	if len(sectionKey) < 2 {
		return fmt.Errorf("incorrect value format, expect: section.key")
	}

	section := cfg.Section(sectionKey[0])
	section.Key(sectionKey[1]).SetValue(items[1])

	err = cfg.SaveTo(pwd + iniConfig)
	if err != nil {
		return err
	}

	return nil
}

func buildNewJiraFromConfig() (*Jira, error) {

	pwd, _ := os.Getwd()
	cfg, err := ini.Load(pwd + iniConfig)
	if err != nil {
		return nil, err
	}

	baseUrl := cfg.Section("jira").Key("url").Value()
	if len(baseUrl) == 0 {
		return nil, fmt.Errorf("missing jira.url value")
	}

	userToken := cfg.Section("user").Key("token").Value()
	if len(userToken) == 0 {
		return nil, fmt.Errorf("missing user.token value")
	}

	return NewJira(baseUrl, BearerAuth(userToken)), nil
}

func stringToMap(str string, resultMap map[string]string) error {
	items := strings.Split(str, " ")
	// no map items provided
	if len(items) == 0 {
		return nil
	}

	for _, v := range items {
		// key:value
		mapEntry := strings.Split(strings.ReplaceAll(v, " ", ""), ":")
		if len(mapEntry) != 2 {
			return fmt.Errorf("incorrect value format, expect: key:value")
		}
		resultMap[mapEntry[0]] = mapEntry[1]
	}

	return nil
}
