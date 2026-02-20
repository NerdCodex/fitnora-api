package services

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type OTPData struct {
	OTP         string
	CreatedTime time.Time
}

var (
	otpStore = make(map[string]OTPData)
	otpMutex sync.RWMutex
)

const otpExpiry = 5 * time.Minute

func GenerateOTP() string {
	return fmt.Sprintf("%06d", rand.Intn(1000000))
}

func SaveOTP(email string, otp string) {
	otpMutex.Lock()
	defer otpMutex.Unlock()

	otpStore[email] = OTPData{
		OTP:         otp,
		CreatedTime: time.Now(),
	}
}

func VerifyOTP(email string, otp string) bool {
	otpMutex.RLock()
	data, exists := otpStore[email]
	otpMutex.RUnlock()

	if !exists {
		return false
	}

	// Check OTP
	if data.OTP != otp {
		return false
	}

	// Check expiry
	if time.Since(data.CreatedTime) > otpExpiry {
		return false
	}

	return true
}

func DeleteOTP(email string) {
	otpMutex.Lock()
	delete(otpStore, email)
	otpMutex.Unlock()
}
