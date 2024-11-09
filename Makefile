deploy:
	gcloud run deploy doublechecker --source . --region europe-north1
brun-local:
	go build && ./proofreader-bot .env.test 

brun-prod:
	go build && ./proofreader-bot 
