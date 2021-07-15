# auth0-echo
Example Go implementation of Auth0 with Labstack's Echo. This is really an adaptation of the Auth0 [tutorial](https://auth0.com/docs/quickstart/backend/golang/01-authorization) by Jim Anderson, with the following changes:
- adapted for **echo** framework (original uses **negroni**)
- it pulls in [my edited version](https://github.com/b-venter/auth0-go-jwt-middleware) of Auth0's Go jwtmiddleware. I think once their v2 is production this will be removed.
- added a function and mddleware to demonstrate getting user data from the Auth0 `/userinfo` [endpoint](https://auth0.com/docs/api/authentication?shell#get-user-info).


