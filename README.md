# flow-build-server

# Why
The well known tool for implementing CI/CD is Jenkins. But its architecture is 
ugly for current tasks. Jenkins was made for use case of running individual
jobs, but after that everybody understand that job is not enough and it's better
to aggreate jobs to pipelines. So developers started making plugins for Jenkins 
in order to add pipeline functionality to it. But there are a lot of problems
with this plugins because main entity in Jenkins is job.

# Concept
1. The main entity for service is pipeline
2. Pipelines have isolated shared context:
    * File workspace is shared
    * Context variables are shared
    * Network for pipeline is isolated
3. Pipelines for the project are stored within source code in repository.
4. Pipelines are defined in file. 
5. Pipelines could be triggered by difference type of triggers
    * Manually
    * Webhook
    * Schedule
6. Service itself consists from two parts:
    * Backend service with RESTFull API
    * Frontend that consume REST API (web frontend, bot front end)
7. Each step is isolated from others

# Implementation
1. Each step in pipeline is docker compose definition or reference to 
docker compose file + condition for calling next step in pipeline.
2. Workspace is folder in (git+git-lfs) repo that was created for this run.
    * Before each run service fetch "pipeline repo" and mount workspace to container
    * After each run service merge results to "pipeline repo"
    * Variable is also stored in "pipeline repo" but in separate folder
3. Steps should be run on the cluster of machines

![Flow CI Arch](http://i.imgur.com/VxEIY20.jpg)

## Config example
```json
workflows:
    - workflow:
        name: "supper workflow"
        jobs:
            job:
                name: "git"
                action: git-compose.yml
                success:
                    - build
                fail:
                    - email
            job:
                name: build
                action: build-compose.yml
                success:
                    - deploy
                fail:
                    - email
            job:
                name: deploy
                action: qa-deploy-compose.yml
                success:
                    - integration-tests
                fail:
                    - email
            job:
                name: integration-tests
                action: integration-compose.yml
                fail:
                    - email
            job:
                name: email
                action: email-compose.yaml
```
