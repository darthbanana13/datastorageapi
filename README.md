# Coding challenge - Data API

## Introduction

### Problem Statement

We have a chatbot where we process data from our customers in real time.
A background job pushes the data (text messages) provided by the customers via HTTP to a customer data management API (that we'll call Data API).

Data API is used by data scientists to further improve the chatbot.

In order to be compliant with the data regulation laws, our customers must give an
explicit consent to process the data - Data API also manages that.

### Task

Given the above mentioned problem statement, your goal would be to create
this Data API service that fulfills the requirements and API specification 
mentioned in _Requirements for Data API_.

## Technologies used

### Go programming language

Go was designed by Google to solve web server type problems at scale (Google problems).
But I chose it because it's:
- a special language that has a lot of power compared to the simplicity, meaning there isn't
  a lot of abstraction you can do in the language (think of C vs C++)
- one of the few programming languages who's standard libraries make it trivial to write
  a web server type application
- It's faster than [Java](https://benchmarksgame-team.pages.debian.net/benchmarksgame/fastest/go.html)
  or [Python](https://benchmarksgame-team.pages.debian.net/benchmarksgame/fastest/go-python3.html)
- I also wanted to practice more Go because I love it, I didn't program in it in some time and
  I needed my fix.
- I also didn't program in Java for some time and am not currently very confident in my Python skills

### ArangoDB
ArangoDB is a NoSQL multi-model DB that supports Graph, Key-Value and most importantly for this project
flexible Document storage.
How I chose it was:
- I searched [CNCF](https://landscape.cncf.io/) for a DB project with as many stars as possible
- I needed NoSQL, because it's faster and more scalable generally than SQL (if you don't use exotic non ACID engines)
  more specifically document storage
- It was the [fastest performer](https://www.arangodb.com/2018/02/nosql-performance-benchmark-2018-mongodb-postgresql-orientdb-neo4j-arangodb/) given the points above
  ![Performance summary](https://www.arangodb.com/wp-content/uploads/2018/02/UPDATE-Benchmark-2018.001.jpeg)
- It is open-source
- It wanted to learn something new with this assignment

### Gin framework
I chose it because:
- it is by far the most popular web framework for Go (if we measure popularity by Github stars of course).
- it was the [fastest framework](https://web-frameworks-benchmark.netlify.app/result?asc=0&l=go&order_by=level64) that was not based on [FastHttp](https://github.com/valyala/fasthttp) which is
  faster than [httprouter](https://github.com/julienschmidt/httprouter), what Go Gin is using but,
  because of the hacky way FastHttp uses the Go language that it might break from a Go language update or the
  fact that one can't use a lot of plugins or standard library functions. I think the goal was to go fast without
  loosing sanity and sleepless nights over the implementation, and the use case for this assignment wasn't so niche
  to require such a specialised framework.
- I needed this assignment to be done as fast as reasonably possible, so a framework seemed like it would help me do things quicker
- I never used any web framework with Go so it was a cool learning opportunity

## Requirements for DataAPI

Create an application with an API layer with the following endpoints:

### POST `/data/:customerId/:dialogId`
With the payload
```json
{
    "text": "the text from the customer",
    "language": "EN"
}
```

This is the endpoint used by the background job to push each customer input during their dialogue with our chatbot.


### POST `/consents/:dialogId`
With the payload
`true` or `false`

This endpoint is called AT THE END of the dialogue when the customer is asked if they gives consent for us to store and use their data for further improving the chatbot.

If false it should delete the customer's data

### GET `/data/(?language=:language|customerId=:customerId)`

This endpoint is used by data scientists to retrieve data to improve the chatbot, it should return all the datapoints:

- that match the query params (if any)
- for which we have consent for
- and sorted by most recent data first
- implement pagination for the returned data

## Design decisions
I wanted to reference a `.env` file in my program so docker compose and go would get the configuration from the same file. I used godotenv for that, but then I hit another snag. If I built the binary and put it locally godotenv would reference the `.env` file in the current directory but if I used `go run` then it would reference the `.env` file from a temporary directory. That's why I made the localpath package, so I could always reference the `.env` from the project directory.

Second was the configuration. The problem is that its text and that needs to be validated and configured into magical values for the package. I used initArangoDB & initlogrusfromtext as a generic way to configure the packages from some string values. The concrete initialization happens initApp package.

I tried to use the MVC pattern (it's not really a pattern, I know, just 3 boxes where you put stuff) mostly because that's how I was used from my web app developing days. I know there's no view, because this is a REST API app. But I do have a controller, service, model and a repository for querying the DB.

The controller contains the enpoints and the validation of the endpoints from a client perspective, the services contain the business logic behind the app mostly, the model contains the collection (aka DB table) structure and the repository contains the implementation for communicating with the DB.

Because I could not find a AQL (ArangoDB Query Language) ORM I did my own query builder using the builder pattern. It's not the prettiest thing, it's very generic, it passes around a lot of hashmaps and slices/arrays but I did not have enough time to make it generic then have a concrete class with data structures the app uses.

I've used logrus for logging, because it is compatible with the standard included log interface but also has the possibility to log as json naively. I wanted to extract the logrus class out of every package I wrote and only inject it through a decorator and also an adapter if the new logging interface was not compatible but alas, there was not enough time for the last two. That being said one could easily modify the code to use the standard library logging interface it's just that it isn't as easy as it should be.

I've used a container library to Dependency Inject the DB connection around. I used it because I had the feeling I would need the container for more things but in the end I only ended up using it for the DB connection. It could be removed fairly easily while still fulfilling the container/singleton pattern but as it is I don't think it's that big of a deal.

The filter function from the repository had a ton of parameters so I used the builder pattern again to construct parameters for it. That builder is pretty generic, I wanted to move it to pkg but it wasn't as clean as I would have liked. Still, I had the problem that I could not query in a case insensitive way for the language so I had two options. Extend the aqlBuilder so I could filter it like that but it seemed like a ton of work, or I could extend the builder for the repository filter function with a decorator specific to the chat collection. I chose the latter, it seemed easier to do.

## Things still to do
The code does not abide by the `Clean Code` or `SOLID` principles unfortunately. I didn't have time in the 20ish hours I spent on it to get it to that level. I feel it's understandable and readable though IMO, there aren't very cognitively complex functions in the code with one notable exception being the `FilterAndSortByNewFirst` in the service.

I haven't written a lot of tests unfortunately. I didn't have the time. I wanted to show my design skills for apps first so tests fell at the end of my priorities. It would also be hard to unit test a great deal of the code because some packages aren't generic enough to mock dependencies. Making the packages more generic is something that I wanted to do, like I said above but as it is unit testing is fairly limited. I could however write component or integration tests for many of the packages but that would require some setup for the DB or for the `.env` configuration in some script files before running the tests which I didn't have time to do.

Some performance testing would had also been nice with `siege` or `k6s` but I didn't have the time. As far as I know right now, on my laptop it could easily do about 1K req/sec in debug mode. I don't know the upper limit. Each entry has a nanotimestamp so the data structure can be sorted up to 1Bil req/sec in theory but there would have to be multiple replicas of the DB and the Go app. The framework and the DB are very fast, as shown above, so it should be no problem to scale to 1Mil req/sec in theory, I just haven't put the work to find out how many replicas of each I would need and how much CPU & RAM it would need. Of course the DB & framework/language can be abused to work poorly but I don't believe I've done any major sins in this regard. Also, if the user does not need to wait for the commit to the database the DB insertion could be done in a separate go routine and chunked for every 1K for better performance over the network but data loss would be a concern here. A HA message bus or queue would solve it but again, all of this is theoretical.

Storage of a tens of millions of billions of data points also should not be a problem with Arango DB as long as the index is refreshed fairly often to maintain the performance.

There are also a few cases which I didn't tackle or didn't know what to do:
- inserting the same dialogId & clientId twice
- the user can refuse to consent after they already consented to the data being stored and that would delete previously consented to data 
- there is no cron like job to delete data that was stored for a long time but never consented to
- how many languages would the API serve? Can the API allow Unicode characters?
- I also wanted to implement an endpoint that would delete all data for a specific customer, for GDPR like legislation but all I did was a repository function

The code is also scattered with a lot of `TODO`s which I haven't mentioned in this section.

## Run the project

### Requirements

- Go version 1.17 or greater installed
- A fairly recent version of docker installed, docker engine version 1.13.0+
- Docker-compose with support for version 3.0 file format (haven't tested with podman) 

### Add a .env file
```
cp .env.example .env
```
Then modify the `.env` file to your liking

### Build (optional)

`./build.sh`

### Run

`./run.sh`

### Run the tests (optional)

`./test.sh`

## Example requests

### Insert
```
curl -X POST localhost:8080/data/5/5 \
   -H 'Content-Type: application/json' \
   -d '{"text": "the text from the customer", "language": "en"}'
```

### Consent
```
curl -X POST localhost:8080/consents/5 \
   -H 'Content-Type: application/json' \
   -d '{"consent": false}'
```

### Filter/View
```
curl 'http://localhost:8080/data?language=En&dialogId=5&customerId=5&entriesPerPage=1&page=1'
```
