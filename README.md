# 📋 jisub

![ci-build](https://github.com/darkowl91/jisub/actions/workflows/ci-branch.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/darkowl91/jisub)](https://goreportcard.com/report/github.com/darkowl91/jisub)
[![MIT License](https://img.shields.io/github/license/mashape/apistatus.svg?maxAge=2592000)](https://github.com/darkowl91/sys-dia-log/blob/master/LICENSE)

> CLI tool to simplify jira tickets interaction

## Install

Download [latest release](https://github.com/darkowl91/jisub/releases/latest) version. Extract to user home dir.
Add jisub executable to PATH `export PATH=$PATH:~/jisub`

## Config

+ Configure jira instance API path:

```bash
    jisub --config "jira.url https://jira-api.com/jira/rest/api/2"
```

+ Obtain jira token from profile and add it to configuration

```bash
    jisub --config "user.token JIRA_TOKEN"
```

Configuration is stored at `jisub-config.ini` file at jisub home folder.

```ini
[user]
token = 

[jira]
url = 
```

## Usage

+ Create required sub tasks with estimates for the parent ticket:

```bash
    jisub --sub-tasks "BE:3 FE:4 QA:0.5" --fields "storypoints:4 dealsize:3,4,0.5 label:New" JIRA-39106 
```

+ Shorten version:

```bash
    jisub -st "BE:3 FE:4 QA:0.5" -f "storypoints:4 dealsize:3,4,0.5 label:New" JIRA-39106 
```

### Customizing Jira payload

To support various jira fields that may depend on configuration, payload templates are used.
This would allow to have your own jira mapping.

+ `issue.tmpl` - is used for the ticket fields update with the `--fields` flag.
+ `sub-tasks-bulk.tmpl` - is used for the sub tasks creation with the `--sub-tasks` flag.

### Support

[![bmc](https://www.buymeacoffee.com/assets/img/guidelines/download-assets-sm-1.svg)](https://www.buymeacoffee.com/darkowl91)
