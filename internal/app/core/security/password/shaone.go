package password

import (
	"github.com/ingot-cloud/ingot-go/internal/app/common/digest"
	"strings"
)

// NewSha1Encoder return Encoder
func NewSha1Encoder() Encoder {
	return &Sha1Encoder{}
}

// Sha1Encoder impl Encoder
type Sha1Encoder struct {
}

// Encode the raw password
func (s *Sha1Encoder) Encode(raw string) (string, error) {
	return digest.SHA1String(raw), nil
}

// Matches Verify the encoded password obtained from storage matches the submitted raw
// password after it too is encoded. Returns true if the passwords match, false if
// they do not
func (s *Sha1Encoder) Matches(raw string, encodedPassword string) bool {
	encoderRaw, err := s.Encode(raw)
	if err != nil {
		return false
	}
	return strings.Compare(encoderRaw, encodedPassword) == 0
}
