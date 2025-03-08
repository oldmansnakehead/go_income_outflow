package token

import (
	"errors"
	"time"
)

/* const refreshExp = time.Hour * 24 * 30
const refreshNBF = time.Minute * 59
const accessExp = time.Hour * 1 */

const refreshExp = time.Hour * 24 * 30
const refreshNBF = time.Second * 45
const accessExp = time.Second * 60

var ErrInvalidTokenCounter = errors.New("invalid token counter")

var now = time.Now
