# proofreader-bot
A Slack bot to proofread messages.
It uses Llama 3.1-B from LM studio that is running locally


## Setup
### LLM configuration
* Download Llama 3.1-B (or any other LLM) and LM studio

### Go installation
```bash
brew install go
```
### Slack configuration:
* The app uses sockets. Enable them in the Slack app configuration - Enable Socket Mode
* Copy SLACK_BOT_TOKEN and SLACK_APP_TOKEN into .env file
* Create a new command "/doublecheck" in Slash Commands section
* Set the request URL to the server URL (eg https://localhost:3000/slack/events)

#### Scopes
```
channels:join
chat:write
commands
groups:write
```


## Run
Run LM studio, select LLM and run server there   

Run the app
```
go run main.go
```
