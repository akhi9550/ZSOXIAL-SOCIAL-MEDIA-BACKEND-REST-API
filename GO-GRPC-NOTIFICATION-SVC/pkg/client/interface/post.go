package interfaces

type PostClient interface {
	GetUserId(post_id int) (int, error)
}
