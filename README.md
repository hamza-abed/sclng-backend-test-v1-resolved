# Canvas for Backend Technical Test at Scalingo

## Instructions

* From this canvas, respond to the project which has been communicated to you by our team
* Feel free to change everything

## Execution

```
docker compose up
```

Application will be then running on port `5000`

## Test

1. Health check
```
$ curl localhost:5000/ping
{ "status": "pong" }
```

2. Get/Search Repositories last 100 (or limit) public repositories (by default offset=0 and limit=100)

Accepted parameters : limit, offset, language, licence

Examples : 
1. get last 100 repositories
curl --location 'http://localhost:5000/repos'

2. search all repo using HTML language in the last 100 public repositories
curl --location 'http://localhost:5000/repos?language=ruby'

3. search all repo using MIT Licence in the last 100 public repositories
curl --location 'http://localhost:5000/repos?licence=mit'

4. search all repo using MIT Licence in the last 900 public repositories
curl --location 'http://localhost:5000/repos?licence=mit&limit=900'

## 


## General Architecture

The GitHub API imposes a rate limit of 5000 requests per hour for authenticated users and 60 requests per hour for unauthenticated users.

Therefore, for our API, it is crucial not to rely directly on the GitHub API if we aim to maximize the number of possible requests.

1. **Dedicated Worker:**
    - A dedicated worker is tasked with fetching data from the GitHub API.
    - The worker strictly adheres to the rate limit set by GitHub thanks to a ratelimiter.

2. **Data Processing and Storage:**
    - The worker processes this data and stores it in a PostgreSQL database, utilizing goroutines to handle the maximum amount of data as quickly as possible.

3. **Rate Limit Management:**
    - To avoid spamming the GitHub API, the worker employs a rate limiter and wakes up to make API requests whenever the rate limiter allows.

4. **Scalability:**
    - This API can scale vertically as it depends solely on its PostgreSQL database. 
    - To scale it horizontally (add instances), it is necessary to extract the worker onto a separate instance.

5. **GitHub Token Integration:**
    - If you provide your GitHub token in the configuration, the rate limiter will automatically adjust to make API requests 5000 times per hour.

``````
GithubToken        string `envconfig:"GITHUB_TOKEN" default:""`
``````

# Code Design and Practices

1. **Use of WaitGroup:**
   - The worker utilizes a WaitGroup to process one page of repositories at a time.

2. **Mutex for Database Transactions:**
   - Mutexes are employed to safeguard transactions in the database, preventing deadlocks during the insertion of new records.

3. **Transactional Handling in Save Repo Method:**
   - The `SaveRepository` method executes a single transaction to uphold ACID principles.
   - In case of errors, the entire set of inserts is rolled back, ensuring data integrity.

4. **Public API without Middleware for User Authentication:**
   - Since the API is public, there is no need for middleware to verify user authentication.

5. **Avoidance of ORM for Performance:**
   - No ORM (Object-Relational Mapping) is used to avoid adding unnecessary layers that could potentially slow down API execution.

6. **Data Handling Approach in Tests:**
   - In testing, data is not mocked initially for expedited development.
   - This approach allows testing both the database connection and query execution.

7. **Minimal External Libraries:**
   - The use of the minimum number of external libraries is prioritized to prevent the addition of superfluous layers.


## What I Haven't Done

1. **Usage of Contexts for Tracing:**
   - I did not utilize contexts to add tracing. However, this decision was made to expedite development.

2. **Integration of Swagger/OpenAPI for API Testing:**
   - I did not incorporate Swagger/OpenAPI to simplify API testing.

3. **Creation of a Configuration File:**
   - I did not create a configuration file to ensure that the evaluator doesn't need to make any modifications to launch the project.

4. **Purging of Old Repositories in the Database:**
   - I did not implement the purging of old repositories in the database.


