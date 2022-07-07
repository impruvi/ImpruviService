# Overview
This package contains the backend code for both the Impruvi 
mobile app and the Impruvi website.

Service code is written in GoLang.
Service is deployed to AWS and infrastructure is managed through AWS CDK.


# Setting up
- Ensure you have golang installed (I'm running version go1.18.3 darwin/amd64)
- If you are using intellij, make sure you have Go Modules integration enabled.
Preferences > Languages & Frameworks > Go > Go Modules


# Deploying
Run the `deploy.sh` script to deploy any changes. 
