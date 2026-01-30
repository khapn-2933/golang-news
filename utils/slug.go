package utils

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// GenerateSlug tạo slug từ title
// Ví dụ: "Hello World!" -> "hello-world"
func GenerateSlug(title string) string {
	// Chuyển về lowercase
	slug := strings.ToLower(title)

	// Loại bỏ các ký tự đặc biệt, chỉ giữ lại chữ cái, số, và khoảng trắng
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")

	// Thay thế khoảng trắng và dấu gạch ngang liên tiếp bằng một dấu gạch ngang
	reg = regexp.MustCompile(`[\s-]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Loại bỏ dấu gạch ngang ở đầu và cuối
	slug = strings.Trim(slug, "-")

	// Nếu slug rỗng, trả về "article"
	if slug == "" {
		return "article"
	}

	return slug
}

// GenerateUniqueSlug tạo slug unique bằng cách thêm số vào cuối nếu cần
// Ví dụ: "hello-world" -> "hello-world-1" nếu "hello-world" đã tồn tại
func GenerateUniqueSlug(baseSlug string, isSlugExists func(string) bool) string {
	slug := baseSlug
	counter := 1

	// Kiểm tra xem slug đã tồn tại chưa
	// Nếu đã tồn tại, thêm số vào cuối cho đến khi tìm được slug chưa tồn tại
	for isSlugExists(slug) {
		slug = baseSlug + "-" + strconv.Itoa(counter)
		counter++
	}

	return slug
}

// isLetterOrDigit kiểm tra xem ký tự có phải là chữ cái hoặc số không
func isLetterOrDigit(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}
