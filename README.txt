# go-continuous-deployment-workshop

* references
    * https://chatgpt.com
    * [Better software, faster: principles of Continuous Delivery and DevOps by Bert Jan Schrijver](https://www.youtube.com/watch?v=tcQNK4R1tfs)
    * [WJUG #181 - Continuous Delivery: architektura i praktyka - Łukasz Szydło](https://www.youtube.com/watch?v=kTX45JGCRYU)
    * [GOTO 2019 • Modern Continuous Delivery • Ken Mugrage](https://www.youtube.com/watch?v=wjF4X9t3FMk)
    * https://www.atlassian.com/continuous-delivery/principles/continuous-integration-vs-delivery-vs-deployment

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
* continuous delivery
    * building and testing software in such a way that the software can
    be released to production at any time
    * "ship early, ship often, sacrificing features, never quality" kule neath
    * for some, "continuously" means once a day, for others once a week, and for others once a month
        * there is no single definition of "continuously"
    * not all applications are suitable for this approach
        * example: highly regulated applications (medical software)
* continuous deployment
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

## feature toggle
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

## bugs
* functional bugs
    * user can't complete some action
        * nie jest poważny
    * user can complete action wrong (najpierw ustrzegać się przed tymi błędami)
        * jest poważny, miało kosztować 200 zł a naliczyło 100
        * data corruption


## script description