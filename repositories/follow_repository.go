package repositories

import (
	"database/sql"
	"news/database"
)

// FollowRepository chứa các method để làm việc với bảng follows
type FollowRepository struct{}

// NewFollowRepository tạo instance mới của FollowRepository
func NewFollowRepository() *FollowRepository {
	return &FollowRepository{}
}

// Follow tạo relationship follow giữa follower và following
// follower_id follow following_id
func (r *FollowRepository) Follow(followerID, followingID int) error {
	query := `INSERT INTO follows (follower_id, following_id) VALUES (?, ?)`
	_, err := database.DB.Exec(query, followerID, followingID)
	return err
}

// Unfollow xóa relationship follow
func (r *FollowRepository) Unfollow(followerID, followingID int) error {
	query := `DELETE FROM follows WHERE follower_id = ? AND following_id = ?`
	_, err := database.DB.Exec(query, followerID, followingID)
	return err
}

// IsFollowing kiểm tra xem follower có đang follow following không
func (r *FollowRepository) IsFollowing(followerID, followingID int) (bool, error) {
	query := `SELECT COUNT(*) FROM follows WHERE follower_id = ? AND following_id = ?`
	var count int
	err := database.DB.QueryRow(query, followerID, followingID).Scan(&count)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return count > 0, nil
}
