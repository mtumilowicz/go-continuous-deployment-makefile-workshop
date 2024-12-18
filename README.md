[![Build Status](https://app.travis-ci.com/mtumilowicz/go-continuous-deployment-makefile-workshop.svg?branch=main)](https://app.travis-ci.com/mtumilowicz/go-continuous-deployment-workshop)
[![License: GPL v3](https://img.shields.io/badge/License-GPLv3-blue.svg)](https://www.gnu.org/licenses/gpl-3.0)

# go-continuous-deployment-makefile-workshop
 
* references
    * https://chatgpt.com
    * [Better software, faster: principles of Continuous Delivery and DevOps by Bert Jan Schrijver](https://www.youtube.com/watch?v=tcQNK4R1tfs)
    * [WJUG #181 - Continuous Delivery: architektura i praktyka - Łukasz Szydło](https://www.youtube.com/watch?v=kTX45JGCRYU)
    * [GOTO 2019 • Modern Continuous Delivery • Ken Mugrage](https://www.youtube.com/watch?v=wjF4X9t3FMk)
    * https://www.atlassian.com/continuous-delivery/principles/continuous-integration-vs-delivery-vs-deployment
    * https://www.linkedin.com/pulse/new-cloud-norm-finsecdevops-joe-frixon-/
    * [Database Migration with Spring Boot – Pitfalls and Hidden Surprises By Dmitry Belyaev](https://www.youtube.com/watch?v=WuIgPPfgQUU)
    * [Nicolas Frankel - Zero-downtime deployment with Kubernetes, Spring Boot and Flyway](https://www.youtube.com/watch?v=RvCnrBZ0DPY)
    * [Continuous deployment to Kubernetes with the Google Container Tools by David Gageot](https://www.youtube.com/watch?v=3nfNP00Tv1k)
    * [Canary Deploys with Kubernetes and Istio by Jason Yee](https://www.youtube.com/watch?v=VU2ILSrpy_Y)
    * [2018 - Mateusz Dymiński - Zero-Downtime deployments of Java applications with Kubernetes](https://www.youtube.com/watch?v=TVB-sQfJBLc)
    * [Develop and Deploy to Kubernetes like a Googler by David Gageot](https://www.youtube.com/watch?v=YYJ4RZFw4j8)
    * [Better Canary Deploys with Kubernetes and Istio by Jason Yee](https://www.youtube.com/watch?v=R7gUDY_-cFo)
    * [Optimising Kubernetes deployments with Helm by Erwin de Gier](https://www.youtube.com/watch?v=TXZBuBQpm-Q)
    * https://www.optimizely.com/optimization-glossary/feature-toggle/
    * https://www.kameleoon.com/blog/feature-toggles-vs-feature-flags-all-you-need-know
    * https://docs.getunleash.io/topics/feature-flags/feature-flag-best-practices
    * https://launchdarkly.com/blog/is-it-a-feature-flag-or-a-feature-toggle/
    * https://www.exoscale.com/syslog/kubernetes-zero-downtime-deployment/
    * https://www.exoscale.com/syslog/kubernetes-zero-downtime-with-spring-boot/
    * https://andrewlock.net/deploying-asp-net-core-applications-to-kubernetes-part-8-running-database-migrations-using-jobs-and-init-containers/
    * https://tutorialedge.net/golang/makefiles-for-go-developers/
    * https://earthly.dev/blog/golang-makefile/
    * https://www.alexedwards.net/blog/a-time-saving-makefile-for-your-go-projects
    * https://www.reddit.com/r/golang/comments/kgkmgp/how_to_use_makefiles_for_your_golang_development/
    * https://news.ycombinator.com/item?id=21735176
    * https://medium.com/@chaewonkong/simplifying-your-build-process-with-makefiles-in-golang-projects-b125af7a10c4
    * https://dzone.com/articles/what-i-have-understood-about-devops-as-an-engineer
    * [Devoxx Greece 2024 - Kubernetes Resiliency by Chris Ayers](https://www.youtube.com/watch?v=R4u5bvWceBQ)

## preface
* goals of this workshop
    1. understand continuous deployment
        * principles and practices
        * continuous integration vs continuous delivery vs continuous deployment
    1. understand how continuous deployment pipeline may work
        * based on the provided script
    1. able to demonstrate how makefile facilitates golang projects
    1. introduction into feature flags and migrations
        * explore zero-downtime database migrations in a Kubernetes environment
    1. enumerate best practices
* workshop plan
    * demonstration
        1. start k8s locally
            * for example using docker desktop
        1. execute makefile step by step
        1. port forward greeting app
            * k9s -> :namespaces -> default -> /pods -> shift-f -> 8080
        1. verify that it is working: http://localhost:8080/app/greeting
            * should return `greet`
        1. remove helm-workshop dir
        1. execute makefile step by step
        1. port forward greeting app
            * k9s -> :namespaces -> default -> /pods -> shift-f -> 8080
        1. verify that it is working: http://localhost:8080/app/greeting
            * should return `greet2`
    * implementation
        1. add step to run tests before creating docker image

## introduction
* may be worthwhile to take a look at
    * golang workshops
        * https://github.com/mtumilowicz/go-chi-gorilla-wire-workshop
        * https://github.com/mtumilowicz/go-concurrency-goroutine-workshop
    * helm workshop: https://github.com/mtumilowicz/helm-workshop
    * argocd workshop: https://github.com/mtumilowicz/argoCD-workshop

## definitions
* goal of CD?
    * a repeatable, reliable process for releasing software
* continuous integration
    * team members integrate their work frequently
    * commits are verified by automated build and tests
    * pros: avoid integration challenges that can happen when waiting for release day to merge changes
* continuous delivery
    * automated release process on top of automated testing
        * deploy = clicking a button
    * building and testing software in such a way that the software can
    be released to production at any time
    * "ship early, ship often, sacrificing features, never quality" Kule Neath
    * you can decide to release daily, weekly, fortnightly, or whatever suits your business requirements
        * there is no single definition of "continuously"
    * not all applications are suitable for this approach
        * example: highly regulated applications (medical software)
* continuous deployment
    * vs continuous delivery
        ![alt text](img/deployment-vs-delivery.png)
    * every change goes through the build/test pipeline and automatically gets
    put into production
    * clue is not to have multiple releases per day, but to have multiple deployments
        * probably not even doable to have several features shipped to prod every day
        * release != deploy
            * release - make feature accessible to users
            * deploy - put newest code on server
    * during deployment, it's possible that parts of the cluster are running version 1 while other parts are running version 2
        * we have feature toggles, so we can enable the feature once everything is on version 2
            * it is just a deployment, not a release
* devops
    * historical context
        * software teams were usually comprised of development teams (who write code) and operation teams (who deploy & maintain applications)
            * both responsible for software release
        * software releases used to be less frequent
            * rare to saw more than one release per month
            * systems were not designed to change fast after deployment
        * infrastructure was not as readily accessible as it is now
            * operation team was responsible for static capacity planning, provisioning, and maintenance of the same
        * development and operations have contradictory goals
            * development team
                * incentive: develop new features faster
            * operations team
                * incentive: ensure stability of application running in production (time-consuming activity)
            * introduced friction & slows down the entire release process
                * priority is given to completing their incentive, and the rest is not their problem
                    * this stretches the release period from days to weeks to months
                * development team
                    * deployment guide is not well documented
                    * doesn’t consider where the app is getting deployed
                * operations
                    * don’t know how the app works
                    * if something fails, need the help of a developer to figure it out
    * definitions
        * intersection between development and operations
        * is the combination of cultural philosophies, practices, and tools that increases an organization’s ability to deliver applications and services at high velocity
        * in short: take a process that is being done manually and automate it
    * objective: deliver high-quality software faster
        * by removing the roadblocks by introducing automation and streamlining the software delivery process
    * goals
        * identify such roadblocks in the software development lifecycle and try to overcome them by introducing some kind of automation
            * could be via tools and scripts or processes and practices
        *
## makefile
* example
    * `main.go`
        ```
        package main

        import "fmt"

        func main() {
            fmt.Println("Hello")
        }
        ```
    * Makefile
        ```
        build:
        	go build -o bin/main main.go

        run:
        	go run main.go
        ```
    * execution
        1. `make build`
            * creates binary
        1. `make run`
            * outputs `Hello`
* serves two purposes
    * automate common admin tasks
        * streamlines development and automates repetitive tasks with a single command
        * example: running tests, checking for vulnerabilities, pushing changes to a remote repository, and deploying to production
    * provide short aliases for Go commands that are long or difficult to remember
        * example: make commands for running tests with different build tags
            * no one has to remember the go test commands
* introduces a topological sort to build steps
    * allows build steps to run in parallel
    * guarantees order by dependency
    * reason why it is used instead of build shell scripts
* used to run and build most programming languages
    * support incremental builds
        * files that have changed since the last build are compiled
        * significantly reduces build times, especially for large projects
* reduces project build-related errors that can arise with inconsistent manual commands
* typically Makefiles are at the root directory
* `make` uses the Makefile as its source of commands to execute
* commands are defined as a rules in the Makefile
    * single rule defines target, dependencies, and the recipe of the Makefile
        * target
            * example: `build`, `run`, `build_and_clean`
            * main component of a Makefile
            * `make` command executes the recipe by its target name
            * is interface to the commands
        * dependencies
            * example: `build_and_clean` has two dependencies: `build` and `run`
            * target can have dependencies that need to be executed before running the target
        * recipes
            * are the actual commands that will be executed
            * can be a single command or a collection of commands
* variables
    * example
        ```
        BINARY_NAME=hello-world

        run:
         ./${BINARY_NAME}
        ```
    * useful when you want the same configs or outputs to be used for different targets
    * can be substituted by enclosing it ${<variable_name>}
    * can be defined either by using `=` or `:=`
        * `=` will recursively expand the variable
            * replace the value at the point when it is substituted
            * example
                ```
                x = foo
                y = $(x) bar
                x = later

                all:
                 echo $(y) // prints "later bar"
                ```
        * `:=` will expanded variable at the first scan
            * simple expanded variables
            * example
                ```
                x := foo
                y := $(x) bar
                x := later

                all:
                 echo $(y)
                ```
* use cases
    * cross-compile to support every OS and every architecture
        * problem: Golang supports multi-platform builds, needs multiple commands to build the binaries for different platforms
            * example: `GOARCH=amd64 GOOS=darwin go build -o hello-world main.go`
            * didn’t want to manually set the GOOS and GOARCH variables for every command
        * solution
            ```
            Makefile

            compile:
                GOOS=linux GOARCH=386 go build -o bin/main-linux-386 main.go
                GOOS=windows GOARCH=386 go build -o bin/main-windows-386 main.go
                ...
            ```
    * layering commands
        * problem: combining commands
        * solution
            ```
            hello:
            	echo "Hello"

            build:
            	go build -o bin/main main.go

            all: hello build
            ```
    * entrypoint for builds in polyglot environment
        * problem: sometimes it gets confusing to figure out how to build a project
        * solution: standardizing on running `make`
            * example:
    * standardize ci/cd pipelines
        * example
            ```
            # .gitlab-ci.yml

            stages:
              - build
              - test
              - deploy

            variables:
              TARGET: "my_program"

            build:
              stage: build
              script:
                - make build
              artifacts:
                paths:
                  - $TARGET

            test:
              stage: test
              script:
                - make test

            deploy:
              stage: deploy
              script:
                - make deploy
              only:
                - main
            ```
* tips
    * to disable printing the recipe while running the target command, use `@` before the recipe
    * you should label each of your tasks that aren't a file with .PHONY
    * other similar tool: https://github.com/go-task
## feature flag
* is a mechanism that allows code to be turned "on" or "off" remotely without the need for a deploy
    * during runtime, your system will query an outside data source or a service to read the configuration
    * example
        ```
        @GetMapping("/feature-status")
        public String getFeatureStatus() {
            if (featureToggleService.isFeatureEnabled("your-feature-toggle-name")) {
                return "Feature is enabled!";
            } else {
                return "Feature is disabled!";
            }
        }
        ```
* also known as "feature flags", "feature switches", or "release toggles"
* separates feature rollout from code deployment
    * mitigates the risks associated with releasing changes
    * testing changes on small subsets of users
        * example: canary releases
    * enable rapid deployment and rollbacks of new code
        * code changes can be made to the the main trunk instead of having multiple feature branches
            * trunk based development process
* has a lifecycle shorter than an application lifecycle
    * most common use case: protect new functionality
    * roll-out of new functionality is complete => the feature flag should be removed
    * should have expiration dates
        * makes it easier to keep track of old feature flags
    * valid exceptions
        * kill-switches
            * used to gracefully disable part of a system with known weak spots
        * internal flags
            * used to enable additional debugging, tracing, and metrics at runtime
                * too costly to run all the time
* large-scale feature flag system components
    * providers: LaunchDarkly, Unleash
    * Feature Flag Control Service
        * centralized feature flag service that acts as the control plane
        * independent business units or product lines should potentially have their own instances
            * contextual decision based on organization
        * keep the management of the flags as simple as possible
            * avoid the complexity of cross-instance synchronization of feature flag configuration
    * Database or Data Store
        * storing feature flag configurations
    * API Layer
        * allow your application to request feature flag configurations
    * Feature Flag SDK
        * easy-to-use interface for fetching flag configurations and evaluating feature flags at runtime
        * query the local cache and ask the central service for updates in the background
        * continuously updated
            * should handle subscriptions or polling to the feature flag service for updates
* principles
    1. enable run-time control
        * control flags dynamically, not using config files
        * if you need to restart your application to turn on a flag => you are using configuration, not feature flags
    1. never expose PII
        * Personally Identifiable Information (PII)
        * Feature Flag Control Service should only handle the configuration and pass this configuration down to SDKs
            * rationale: feature flags often require contextual data for accurate evaluation
                * example: user IDs, email addresses, or geographical locations
            * example
                ```
                UnleashContext context = UnleashContext.builder()
                        .userId(userId)
                        .build();
                return unleash.isEnabled(featureName, context);
                ```
                and Feature Flag Control Service configuration
                ```
                Feature Toggle Name: your-feature-toggle-name
                Strategy: UserWithId
                Parameters:
                    User IDs: user1,user2,user3
                ```
        * allows offline functionality
        * reduces bandwidth costs
            * local evaluation reduces the amount of data transferred between your application and the feature flag service
    1. evaluate flags as close to the user as possible
        * reduces latency
        * evaluation should always happen server side
    1. decouple reading and writing flags
        * horizontally scale out the read APIs without scaling the write APIs
    1. feature flag payload should be as small as possible
        * problem: feature flag based on individual user IDs
            * categorize these users into logical groupings
    1. favor availability over consistency
        * feature flag system should not be able to take down your main application under any circumstance
        * application's availability should have zero dependencies on the availability of feature flag system
    1. do not confuse flags with application configuration
        * plan to clean up old feature branches
    1. treat feature flags like technical debt
    1. use unique names across all applications
        * enforce naming conventions
        * prevents the reuse of old flag names to protect new features (zombies)
* branching
    * branch per feature
        * cons: usually big and messy changes, often chaotic
        * pros: easier code review (diff branch master)
    * branching through abstraction
        * cons: slower, more thoughtful
        * pros: can't afford to make chaotic changes
            * every commit goes to production
            * initially, you need to make space for your change
        * commiting
            * first commit: make space for the new feature (refactoring)
                * doesn't introduce any changes in functionality, so it can be done safely
                * refactoring doesn’t break anything as long as there is good test coverage
            * second commit: introduce the feature switch (toggle)
            * third commit and subsequent: additional steps for the feature
                * two approaches to handling errors
                    * avoiding errors: through thorough testing and coverage
                    * infrastructure for quickly handling errors
                        * disable the feature with a toggle, and do quick fix
* functional bugs
    * user can't complete some action
        * not serious => a matter of turn-off feature flag
    * user can complete action wrong
        * serious
        * usually leads to data corruption
        * example: something should cost 200 USD but costed 100

## db migrations
* zero-downtime
    * Blue-Green deployment (oldest ideas)
        * two exactly similar environments, one referenced as Green, the other as Blue
        * one of them runs the production app, while the other runs the pre-production app
        * in front of them sits a dispatcher
            * routes requests to the relevant environment: production or pre-production
            * you deploy update to the pre-production environment, test it, and switch the dispatcher
        * migration can take time and can block users in-between environments pending migration
            * not always doable before switch => production data is constantly changing
        * maybe environments should share same database => schema changes must be compatible with the old app version
    * rolling update
        * split the schema update into a series of small side-by-side compatible schema updates
        * application needs to be updated in increments
            * able to cope with the current schema update and the next one
* problem: long running data migration
    * may cause liveness probe to fail what causes multiple restarts of a container and a failed deployment
    * solutions
        1. k8s probes
            * startupProbe
                * prevent premature pod restarts due to long initialization processes
                    * postpones liveness/readiness probes until success
                * used for apps with complex initializations or large data loads
                    * example: migrations 
            * readinessProbe
                * determine if the application is ready to handle traffic
                * prevent routing traffic to pods that are still initializing or are not ready
                * blocks traffic on failure
                    * pod is not restarted
                        * is removed from the list of active service endpoints
            * livenessProbe
                * ensures the pod is alive and running
                * restarts on failure
        1. treat db migration as a long running process
            * use dedicated k8s process
                * jobs
                * init containers
                    * often used for downloading or configuring pre-requisites required by the main container
                    * when Kubernetes deploys a pod, it runs all the init containers first
                        * once all have exited gracefully => main containers be executed
            * separate db migration from service startup


## script description
* script purpose: simulate CI/CD pipeline to deploy by specific commit hash
    * CI needs to retrieve current commit hash and pass it to the script
        * problem: no single source of truth
            * what you pass as a parameter will be deployed
            * if you modify by hand deployment in the cluster there is no reconciliation
        * solution: https://github.com/mtumilowicz/argoCD-workshop
* steps
    1. parse command-line arguments
        * commit hash
        * Docker image name
        * Helm
            * release name
            * chart directory
        * Git repository URL
        * Kubernetes namespace.
    1. validate required inputs
        * required flags check
    1. clone Git repository
    1. checkout commit by provided commit hash within the cloned repository
    1. prepare artifact with gradle
        1. clean & build
        1. run tests
        1. build docker image
            * tagged with the commit hash using Gradle
    1. upgrade helm chart
        * override placeholder `deployment.image.version` with just created docker image
