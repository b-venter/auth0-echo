#As set in Auth0 API > Settings > Identifier
export AUTH0_AUD="https://some-api"
#Your Auth0 tenant domain
export AUTH0_ISS="https://YOUR-AUTH0-TENANT.REGION.auth0.com/"
#https://YOUR-AUTH0-TENANT-DOMAIN/.well-known/jwks.json
export AUTH0_JWKS="https://YOUR-AUTH0-TENANT.REGION.auth0.com/.well-known/jwks.json"

#Launch API server
go run server.go
