# proofreader-bot
A Slack bot to proofread messages


## Setup
### Slack configuration:
* The app uses sockets. Enable them in the Slack app configuration - Enable Socket Mode
* Copy SLACK_BOT_TOKEN and SLACK_APP_TOKEN into .env file
* Create a new command "/typosweep" in Slash Commands section
* Set the request URL to the server URL (eg https://localhost:1234/slack/events)

#### Scopes
```
channels:join
chat:write
commands
groups:write
```
