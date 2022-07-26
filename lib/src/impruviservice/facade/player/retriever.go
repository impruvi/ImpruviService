package player

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	playerDao "impruviService/dao/player"
	"impruviService/exceptions"
	"log"
	"time"
)

var jwtKey = []byte("7220a177-278a-48af-83d6-8b575af78ccd")

type Claims struct {
	PlayerId string `json:"playerId"`
	jwt.StandardClaims
}

func DoesPasswordMatch(email, password string) (bool, error) {
	player, err := playerDao.GetPlayerByEmail(email)
	if err != nil {
		return false, err
	}

	if !player.IsActive {
		log.Printf("Account is not active for email: %v\n", email)
		return false, nil
	}

	return player.Password == password, nil
}

func GetPlayerByEmail(email string) (*Player, error) {
	player, err := playerDao.GetPlayerByEmail(email)
	if err != nil {
		return nil, err
	}
	return convert(player)
}

func GetPlayerById(playerId string) (*Player, error) {
	player, err := playerDao.GetPlayerById(playerId)
	if err != nil {
		return nil, err
	}
	return convert(player)
}

func GetPlayerFromToken(token string) (*Player, error) {
	playerId, err := ParseToken(token)
	log.Printf("PlayerId: %v\n", playerId)
	if err != nil {
		return nil, err
	}
	if playerId == "" {
		return nil, exceptions.NotAuthorizedError{Message: "Invalid request token"}
	}

	player, err := GetPlayerById(playerId)
	if err != nil {
		return nil, err
	}
	return player, nil
}

// GenerateToken creates a JWT with the playerId as the claim.
// Token can be passed to ParseToken to retrieve the playerId
func GenerateToken(playerId string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour * 365) // 1 year
	claims := &Claims{
		PlayerId: playerId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func ParseToken(token string) (string, error) {
	claims := &Claims{}

	// Parse the JWT string and store the result in `claims`.
	// Note that we are passing the key in this method as well. This method will return an error
	// if the token is invalid (if it has expired according to the expiry time we set on sign in),
	// or if the signature does not match
	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return "", exceptions.NotAuthorizedError{Message: fmt.Sprintf("Token: %v has invalid signature\n", token)}
		}
		log.Printf("Unexpected error when parsing jwt: %v\n", err)
		return "", err
	}
	if !tkn.Valid {
		return "", exceptions.NotAuthorizedError{Message: fmt.Sprintf("Token: %v is not valid\n", token)}
	}

	return claims.PlayerId, nil
}
