<div align="center">
    <h1>auth-api</h1>
    <p>The authentication API allows users to sign in using their Google accounts and generates a JWT for subsequent authentication.</p>
</div>

## Description
Head to the front page and log in with your Google Account. The login will automatically create a user profile and store it in MongoDB. Use the API to update your profile or to delete it from MongoDB. 

For extra security passwords are not stored in the database, Google provides excellent protection with 2FA options to secure your account.

## Requirements
* Go >=1.17
* MongoDB
* Google Cloud OAuth 2.0 Client ID

## API
**GET /_healthz** - health check

**GET /** - front page
<br><br>
**GET /api/v1/login** - initiates login and redirects to Google

**GET /api/v1/login/callback** - callback endpoint called when user completed log in, creates or updates user profile
<br><br>
**GET /api/v1/users/me** - profile of a current user

**POST /api/v1/users/me** - update current user's profile

**DELETE /api/v1/users/me** - delete current user's profile

## Run
For local testing create a .env file by using the .env.example and insert required credentials, then run the command:
```
make server
```
The API will be accessible at localhost:8080

## Test
```
make test
```
