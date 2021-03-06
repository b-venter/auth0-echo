# Auth0 JWT middleware with Echo framework
Example Go implementation of Auth0 middleware with Labstack's Echo. This is really an adaptation of the Auth0 [tutorial](https://auth0.com/docs/quickstart/backend/golang/01-authorization) by Jim Anderson, with the following changes:
- adapted for [**echo**](https://github.com/labstack/echo) framework (original uses [**negroni**](https://github.com/urfave/negroni))
- it pulls in [my edited version](https://github.com/b-venter/auth0-go-jwt-middleware) of Auth0's Go jwtmiddleware. I think once their v2 is production this will be removed.
- added a function and mddleware to demonstrate getting and verifying user data from the Auth0 `/userinfo` [endpoint](https://auth0.com/docs/api/authentication?shell#get-user-info).


## Giving it a test run
1. You will need to have an [Auth0 account](https://auth0.com/). It is free for basic use, getting your feet wet and possibly even makes great toast.
2. Edit the run_sample.sh to contain your details.
3. You will also need to edit [line 191](https://github.com/b-venter/auth0-echo/blob/9c4945df5ec204f626b73845756a626d5f7aab0b/server.go#L191) of *server.go* with the email address you are using to test with. Or just remove the test.
4. Then a simple `sh run.sh` will set your ENV and launch the server. (I rename *run_sample.sh* to *run.sh*)

## Useful links
* Test tokens: https://jwt.io/
* Human readable iat and exp: https://www.epochconverter.com/
* Auth0 API: https://auth0.com/docs/api/authentication
