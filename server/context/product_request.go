package context

// delete product request struct
type ProductDeleteRequest struct {
	ProductID string `json:"productId"`
}

// user product rating struct
type UserProductRatingReviewRequest struct {
	ProductID string `json:"productId"`
	Rating    int    `json:"rating"`
	Review    string `json:"review"`
}
