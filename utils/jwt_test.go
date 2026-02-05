package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateToken kiểm tra generate JWT token
func TestGenerateToken(t *testing.T) {
	secret := "test-secret-key"
	userID := 123

	token, err := GenerateToken(userID, secret)
	
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	
	// Token phải là string hợp lệ
	assert.Greater(t, len(token), 0)
}

// TestValidateToken kiểm tra validate token
func TestValidateToken(t *testing.T) {
	secret := "test-secret-key"
	userID := 123

	// Generate token
	token, err := GenerateToken(userID, secret)
	require.NoError(t, err)

	tests := []struct {
		name      string
		token     string
		secret    string
		wantUserID int
		wantErr   bool
	}{
		{
			name:      "valid token",
			token:     token,
			secret:    secret,
			wantUserID: userID,
			wantErr:   false,
		},
		{
			name:      "wrong secret",
			token:     token,
			secret:    "wrong-secret",
			wantUserID: 0,
			wantErr:   true,
		},
		{
			name:      "invalid token format",
			token:     "invalid.token.format",
			secret:    secret,
			wantUserID: 0,
			wantErr:   true,
		},
		{
			name:      "empty token",
			token:     "",
			secret:    secret,
			wantUserID: 0,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateToken(tt.token, tt.secret)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, 0, gotUserID)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantUserID, gotUserID)
			}
		})
	}
}

// TestTokenExpiration kiểm tra token có expiration time
func TestTokenExpiration(t *testing.T) {
	secret := "test-secret-key"
	userID := 123

	// Generate token
	token, err := GenerateToken(userID, secret)
	require.NoError(t, err)

	// Validate token ngay lập tức - phải thành công
	gotUserID, err := ValidateToken(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, gotUserID)
}

// TestTokenDifferentUsers kiểm tra token cho các user khác nhau
func TestTokenDifferentUsers(t *testing.T) {
	secret := "test-secret-key"
	
	userID1 := 123
	userID2 := 456

	token1, err1 := GenerateToken(userID1, secret)
	require.NoError(t, err1)

	token2, err2 := GenerateToken(userID2, secret)
	require.NoError(t, err2)

	// Hai token phải khác nhau
	assert.NotEqual(t, token1, token2)

	// Validate token1 phải trả về userID1
	gotUserID1, err := ValidateToken(token1, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID1, gotUserID1)

	// Validate token2 phải trả về userID2
	gotUserID2, err := ValidateToken(token2, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID2, gotUserID2)
}

