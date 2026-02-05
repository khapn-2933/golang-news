package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestHashPassword kiểm tra hash password hoạt động đúng
func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "empty password",
			password: "",
			wantErr:  false, // bcrypt vẫn hash được empty string
		},
		{
			name:     "long password",
			password: "this-is-a-very-long-password-with-many-characters",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, hash)
				// Hash phải khác password gốc
				assert.NotEqual(t, tt.password, hash)
			}
		})
	}
}

// TestCheckPassword kiểm tra verify password hoạt động đúng
func TestCheckPassword(t *testing.T) {
	// Hash một password
	password := "password123"
	hash, err := HashPassword(password)
	require.NoError(t, err)
	require.NotEmpty(t, hash)

	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "correct password",
			password: password,
			hash:     hash,
			want:     true,
		},
		{
			name:     "wrong password",
			password: "wrongpassword",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty password",
			password: "",
			hash:     hash,
			want:     false,
		},
		{
			name:     "empty hash",
			password: password,
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CheckPassword(tt.password, tt.hash)
			assert.Equal(t, tt.want, result)
		})
	}
}

// TestHashPasswordDifferentHashes kiểm tra mỗi lần hash tạo ra hash khác nhau (do salt)
func TestHashPasswordDifferentHashes(t *testing.T) {
	password := "password123"
	
	hash1, err1 := HashPassword(password)
	require.NoError(t, err1)
	
	hash2, err2 := HashPassword(password)
	require.NoError(t, err2)
	
	// Hai hash phải khác nhau (do salt random)
	assert.NotEqual(t, hash1, hash2)
	
	// Nhưng cả hai đều verify được với password gốc
	assert.True(t, CheckPassword(password, hash1))
	assert.True(t, CheckPassword(password, hash2))
}

