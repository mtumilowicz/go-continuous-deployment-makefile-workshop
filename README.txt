# go-continuous-deployment-workshop

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

## preface
* goals of this workshop
    * general understanding
* workshop plan
    1. start k8s locally
        * for example using docker desktop
    1. run script `run/greet`
        * verify that it is working: http://localhost:31234/app/greeting
            * should return `greet`
    1. remove helm-workshop dir
    1. run script `run/greet2`
        * verify that it is working: http://localhost:31234/app/greeting
            * should return `greet2`

## introduction
* may be worthwhile to take a look at
    * golang workshops
        * https://github.com/mtumilowicz/go-chi-gorilla-wire-workshop
        * https://github.com/mtumilowicz/go-concurrency-goroutine-workshop
    * helm workshops: https://github.com/mtumilowicz/helm-workshop

## general
* goal of CD?
    * a repeatable, reliable process for releasing software

## definitions
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
* DevOps
    * development and operations engineers being responsible together
    for the entire product lifecycle
    * usually it is ops team renamed for devops team
    * FinSecDevOps
        * Finance, Security and Development

## feature toggle
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
* pros
    * mitigates the risks associated with releasing changes
    * testing changes on small subsets of users
        * example: canary releases
    * enable rapid deployment and rollbacks of new code
        * code changes can be made to the the main trunk instead of having multiple feature branches
            * trunk based development process
* providers: LaunchDarkly, Unleash
* large-scale feature flag system components
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
* steps
    * first commit: make space for the new feature (refactoring)
        * doesn't introduce any changes in functionality, so it can be done safely
        * refactoring doesn’t break anything as long as there is good test coverage
    * second Commit: introduce the feature switch (toggle)
    * third commit and subsequent: additional steps for the feature
        * two approaches to handling errors
            * avoiding errors: through thorough testing and coverage
            * infrastructure for quickly handling errors
                * disable the feature with a toggle, and do quick fix

## db migrations
* problem
  * long running data migration causes liveness probe to fail what causes multiple restarts of a container and a failed
  deployment
  * solution
    * treat db migration as a long running process
      * use dedicated k8s process
    * separate db migration from service startup
* separate schema changes from data migration and run them in individual transactions
* make schema changes idempotent
* make schema changes backward compatible
  * follow up migration to clean things up
* deployments
    * blue-green deployment
    * rolling update
        * canary deploy

## bugs
* functional bugs
    * user can't complete some action
        * not serious
    * user can complete action wrong
        * serious
        * usually leads to data corruption
        * example: something should cost 200 USD but costed 100


## script description

## argoCD