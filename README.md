# Jane

Jane is a bot to pull information and conduct operational activities in your chatops scenario - even in a command line way. This bot is written in go and is made to be configuration driven. Contributions are welcome via pull requests. If you want to know why the name 'Jane' was chosen, talk to @kcwinner.



## Getting Started
* This is developed using Go 1.5.3
* Pull the project with 'go get github.com/projectjane/jane'
* Compile with 'go install jane.go'
* Use the samples in the startup folder for different environments


## Configuration
The entire configuration of the site is done via a json config file. The configuration file is expected to be named 'jane.config' and will be looked for in this order:
* -config config.json - Pass in a configuration file location as a command line parameter
* ./jane.json - the location of the jane binary
* ~/jane.json - the home directory of the user
* /etc/jane.json - the global config


## Connectors
Connectors are what Jane uses to pull in information, interpret them and issue out a response. The Routes specify where the results from the input should be written to or * for all. The Target can specify a channel in the case of slack. To add a new connector, Put them in the connectors folder and make an entry in connectors/list.go.

For the connector configuration, when adding routes, you must specify the ID of the connector you want to route response to.

Supported connectors:
* bamboo - Atlassian Bamboo integration
* cli - Command line interface
* email - Email
* exec - Execution of applications
* exec2 - Next generation of the exec/ssh/monitor connector
* imageme - Pull back images or animated gifs
* jira - Atlassian Jira integration
* monitor - Monitor of systems
* response - Text Responses
* rss - RSS Feed
* slack - Slack chat
* ssh - Execute commands on remote systems
* twilio - send SMS alerts
* website - Monitor return code of websites
* webhook - listen at http://.../webhook/ 
* wolfram - Execute queries against Wolfram Alpha

### Exec2 Connector

This connector is the next generation to replace the exec, ssh and monitor connectors. It provides a single means of making local and remote calls to linux systems. You can allow these calls to be made by command, but also mark them with the RunCheck property to set Jane to check them. This combined with the interpreter for output, makes it a very capable monitoring platform.

```
    {"Type": "exec2", "ID": "ExecTwo", "Active": true,
      "Server": "elasticsearch1.somecompany.com", "Port": "22", "Login": "jane", "Pass": "abc123",
      "Commands": [
        {
            "Name": "Apt Check",
            "Match": "jane elasticsearch1 aptcheck",
            "Output": "```%stdout% ```",
            "Cmd": "/usr/lib/nagios/plugins/check_apt",
            "Args": "",
            "HideHelp": false,
            "RunCheck": true,
            "Interval": 1,
            "Remind": 15,
            "Green": "*OK*",
            "Yellow": "*WARNING*",
            "Red": "*CRITICAL*"
        },
      ],
      "Routes": [
        {"Match": "*", "Connectors": "slack", "Target": "#devops"}
      ]
    }
```

*Usage*
* To make local calls to the system, leave out the Server, Port, Login, Pass values.

*Fields*

_Type_ This specifies the type of connector, in this case, 'exec2'

_ID_ This should be a unique identifier for this connector

_Active_ This is a boolean value to set this connector to be activated

_Debug_ This is a boolean value to set if the connector shows debug information in the logs

_Server_ The server address or IP to connect to

_Port_ The port number to connect to (Default: 22)

_Login_ The user to login with

_Pass_ The password to connect with

_Commands_:

    _Name_ Readable name of check
    _Match_ Command (###Matching)[matching]
    _Output_ Formatting for the output of the command, use `%stdout%` as the output
    _Cmd_ The command to execute (do not include arguements)
    _Args_ The arguments, space deliminated
    _HideHelp_ A boolean to show or hide the help when displaying help (Default: false)
    _RunCheck_ A boolean that will have Jane periodically run this (Default: false)
    _Interval_ An integer that is the number of minutes between checks when RunCheck is true (Default: 1)
    _Remind_ An integer which is the number of units of Interval to wait before reminding of a non-Green status, with Zero being no reminders (Default: 0)
    _Green_ A (###Matching)[match] to identify what is in a green state
    _Yellow_ A (###Matching)[match] to identify what is in a yellow state
    _Red_ A (###Matching)[match] to identify what is in a red state


### Routes

Routes can exist for connectors that listen or interpret commands. Routes can have more than one connector if you would like to send messages to more than one place. There is also matching on the message to filter which messages get sent.

```
    "Routes": [
      {"Match": "*", "Connectors": "slack", "Target": "#devops"},
      {"Match": "*DANGER*", "Connectors": "slack", "Target": "@matt"}
    ]
```



### Matching

Matching within Jane has a simple string matching with wild cards.

Examples:
`*failure*`
`*failure`
`failure*`

For anything more complex, you can use Regular Expressions.

Example:
`/fail(.+)/`


### Architecture Notes

The general architecture of Jane is 
