package token

import (
	"errors"
	"time"
)

/* const refreshExp = time.Hour * 24 * 30
const refreshNBF = time.Minute * 59
const accessExp = time.Hour * 1 */

const refreshExp = time.Hour * 24 * 30 // Refresh Token หมดอายุใน 30 วัน
const refreshNBF = time.Second * 0     // Refresh Token เริ่มใช้งานได้หลังจาก 45 วินาที
const accessExp = time.Second * 3600   // Access Token หมดอายุใน 60 วินาที

var ErrInvalidTokenCounter = errors.New("invalid token counter")

var now = time.Now
