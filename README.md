# go_movie_api

Go-MovieReviews-API EC2 Documentation

Overview
This is a REST API for a movie review application hosted on AWS EC2. It has endpoints for movies, users, comments, and authentication.

Endpoints

Movies
- Get movie: GET /api/v1/movies/{id} 
- Get most popular movies: GET /api/v1/movies/popularity/{limit}

Users 
- Create user: POST /api/v1/users
- Get user: GET /api/v1/users/{id}
- Update user: PUT /api/v1/users/{id}

Comments
- Create comment: POST /api/v1/comments 
- Delete comment: DELETE /api/v1/comments/{id}
- Update comment: PUT /api/v1/comments/{id}

Auth
- Authenticate user: POST /api/v1/auth/login

Returns JSON data. Some endpoints require JWT token in Authorization header.

Hosted on:
http://ec2-18-212-116-192.compute-1.amazonaws.com
