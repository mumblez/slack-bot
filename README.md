# slack-bot

Bot built specifically to monitor a channel configured for CI/CD notifications (from shippable), listens for new messages and sends it to the github user, which we extract from the message and lookup via a hardcoded map atm (userMappings.json to get their slack username)

Just helps to increase visibility if developers are not directly looking at the ci-notifications channel

# TODO
- Configure the user mapping file as a k8s config map
- Pass in the SLACK_TOKEN as a k8s secret (create a separate namespace?)
- Pass in slack channel to monitor as env (hardcoded atm)
- Create Dockerfile (scratch should be fine)
- Create k8s deployment config
