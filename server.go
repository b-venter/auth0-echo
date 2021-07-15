package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	auth0 "github.com/b-venter/auth0-go-jwt-middleware"
	jwt "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

/* Sensitive/custom settings */

var aud = os.Getenv("AUTH0_AUD") //As set in Auth0 API > Settings > Identifier

var iss = os.Getenv("AUTH0_ISS") //Your Auth0 tenant domain

var jwksUrl = os.Getenv("AUTH0_JWKS") //https://YOUR-AUTH0-TENANT-DOMAIN/.well-known/jwks.json

/* end */

var jwtClaims []string

type Response struct {
	Message string `json:"message"`
}

type Jwks struct {
	Keys []JSONWebKeys `json:"keys"`
}

type JSONWebKeys struct {
	Kty string   `json:"kty"`
	Kid string   `json:"kid"`
	Use string   `json:"use"`
	N   string   `json:"n"`
	E   string   `json:"e"`
	X5c []string `json:"x5c"`
}

//getPemCert: function to get certificate in PEM format.
//Used by verify() func
func getPemCert(token *jwt.Token) (string, error) {
	cert := ""

	resp, err := http.Get(jwksUrl)

	if err != nil {
		return cert, err
	}
	defer resp.Body.Close()

	var jwks = Jwks{}
	err = json.NewDecoder(resp.Body).Decode(&jwks)

	if err != nil {
		return cert, err
	}

	for k := range jwks.Keys {
		if token.Header["kid"] == jwks.Keys[k].Kid {
			cert = "-----BEGIN CERTIFICATE-----\n" + jwks.Keys[k].X5c[0] + "\n-----END CERTIFICATE-----"
		}
	}

	if cert == "" {
		err := errors.New("unable to find appropriate key")
		return cert, err
	}

	return cert, nil
}

//A jwt.Keyfunc type function, required for JWTMiddleware struct
//It processes the aud, iss validity of the jwt.Token. and gets certificate which CheckJWT will parse and verify.
func verify(token *jwt.Token) (interface{}, error) {

	// Verify 'aud' claim
	checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
	if !checkAud {
		return token, errors.New("invalid audience")
	}

	// Verify 'iss' claim
	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	if !checkIss {
		return token, errors.New("invalid issuer")
	}

	//Store "scope" so long. Use it or lose it later.
	claim := fmt.Sprintf("%v", token.Claims.(jwt.MapClaims)["scope"])
	jwtClaims = strings.Fields(string(claim))

	cert, err := getPemCert(token)
	if err != nil {
		//panic(err.Error())
		return nil, err
	}

	result, err := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))
	if err != nil {
		return nil, err
	}

	return result, nil

}

//Type for unpacking json
type ujson map[string]interface{}

/*Or if you want something more structured:
type ujson struct {
	User string `json:"user"`
	UserV bool `json:"email_verified"`
	FamName string `json:"family_name"`
	GivName string `json:"given_name"`
	Loc string `json:"locale"`
	Name string `json:"name"`
	NicName string `json:"nickname"`
	Pic string `json:"picture"`
	Sub string `json:"sub"`
	Upd string `json:"updated_at"`
}
*/

//Function to call Auth0 /userinfo API endpoint
func getUser(tok string) ujson {
	url := iss + "userinfo"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("Authorization", tok)
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Println(err)
	}

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//Get as json
	var bodyj ujson
	if err := json.Unmarshal(body, &bodyj); err != nil {
		fmt.Println("Error providing json: ", err)
	}

	return bodyj
}

/* MIDDLEWARE */

//Custom middleware

//Middleware that verifies token
func middleJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// c holds Response and Request
		// t has the methods to process the token
		var t *auth0.JWTMiddleware = auth0.New(auth0.Options{
			ValidationKeyGetter: verify,
			SigningMethod:       jwt.SigningMethodRS256,
		})
		vr := t.CheckJWT(c.Response().Writer, c.Request())
		if vr != nil {
			return echo.ErrUnauthorized //Not allowed to proceed. Might be nice to break it down further to "no token", etc
		}

		return next(c) //Proceed to next.
	}
}

func middleUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//Retrieve user info
		token := c.Request().Header.Get("Authorization")
		userInfo := getUser(token)
		fmt.Println("User info: ", userInfo)
		fmt.Println("Possible verification data: ", userInfo["email"])

		//TODO: Replace this with a database lookup
		if userInfo["email"] != "my-email@gmail.com" {
			return echo.ErrForbidden //Note: https://stackoverflow.com/questions/3297048/403-forbidden-vs-401-unauthorized-http-responses
		}

		return next(c) //Proceed to next.
	}
}

/* ENDPOINTS*/

//Simple restricted endpoint
func simple(c echo.Context) error {

	//Process scope from Auth0. Can replace with "if" statement.
	for i := range jwtClaims {
		fmt.Println("scope in jwt Claim: ", jwtClaims[i])
	}

	return c.JSON(http.StatusOK, "A restricted route was successfully accessed.")
}

//Simple unrestricted endpoint
func open(c echo.Context) error {
	return c.JSON(http.StatusOK, "Open route.")
}

func main() {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:4200"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	//Unrestricted API endpoint
	e.GET("/open/", open)

	// Restricted group
	r := e.Group("/simple/")
	r.Use(middleJWT)
	r.GET("", simple, middleUser)

	e.Logger.Fatal(e.Start(":4040"))

}
