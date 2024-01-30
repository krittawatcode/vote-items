# VOTE ITEMS

This is the repository for the Vote items application.

You can find my written about the design of this application on my [Medium](https://medium.com/@kritwis/golang-clean-architecture-with-demo-e0938e5be02b).

![App Overview](./application_overview.png)

## Setup

### Traefik Setup

1. Open your hosts file:
   - On Linux or MacOS, the file is `/etc/hosts`.
   - On Windows, the file is `c:\Windows\System32\Drivers\etc\hosts`.

2.  Add the following line to the file: `127.0.0.1 krittawatcode.test`
This will route any requests for `krittawatcode.test` to your local machine.


### Postman Setup
1. Import the VOTE-ITEMS.postman_collection.json file into Postman.
2. Set up your environment variables in Postman as needed.
3. Use the imported collection to test the API endpoints.

### Startup 
- Run your Docker Compose file with docker-compose up.
- Call api login with default user
```
{Email: "admin@mtl.co.th", Password: "adminPassword"},
{Email: "krittawat@mercy.gg", Password: "userPassword"},
```
## NOTE
- for user API: prefix with /user
   ex. krittawatcode.test/api/v1/user/singIn
- for vote item API: prefix with /vote-items
   ex. krittawatcode.test/api/v1/vote_items/1/open
- for open and close vote session API: prefix with vote_session
   ex. krittawatcode.test/api/v1/vote_sessions/1/open
- for user's vote API: prefix with /votes
   ex. krittawatcode.test/api/v1/votes
- for vote's result API: prefix with /vote_results

# API DOCs
- USER prefix with /user
```
1. GET: /me
2. POST: /signUp
3. POST: /signIn
4. POST: /tokens
```

- OPEN / CLOSE VOTE SESSIONS 
```
1. PUT /api/v1/vote_sessions/{id}/open // Open a vote session
2. GET /api/v1/vote_sessions/open  // Get open vote session
3. PUT /api/v1/vote_sessions/{id}/close // Close a vote session
```


- VOTE ITEMS
```
1. GET /api/v1/vote_items // Get all active vote items
2. POST /api/v1/vote_items // Create a new vote item
3. PUT /api/v1/vote_items/{id} // Update item
4. DELETE /api/v1/vote_items/{id} // Delete a vote item by id
5. DELETE /api/v1/vote_items // Clear all vote items
```

- VOTE 
```
1. POST /api/v1/votes // Cast a vote
```

- VOTE RESULT & REPORT 
```
1. GET /api/v1/vote_results/{session_id} // Get vote results by session id
2. GET /api/v1/vote_results/{session_id}?format=csv /// Get vote results by session id in CSV format
```

## END CREDIT
it's such a fun run to do this! hope I get the job I want <3
sadly I cant finish everything I want to do but... hopefully u guys have a mercy on me lol
