package shortcode

import "math/rand"

// ShortCode .
type ShortCode struct {
	length int
}

// NewShortCode .
func NewShortCode(length int) *ShortCode {
	return &ShortCode{
		length: length,
	}
}

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenerateShortCode .
func (s *ShortCode) GenerateShortCode() string {
	length := len(chars)
	result := make([]byte, s.length)

	for i := 0; i < s.length; i++ {
		result[i] = chars[rand.Intn(length)]
	}
	return string(result)
}

