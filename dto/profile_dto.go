package dto

// ProfileResponse định dạng response theo RealWorld spec
// {"profile": {"username": "...", "bio": "...", "image": "...", "following": true/false}}
type ProfileResponse struct {
	Profile struct {
		Username  string  `json:"username"`
		Bio       *string `json:"bio"`
		Image     *string `json:"image"`
		Following bool    `json:"following"`
	} `json:"profile"`
}
