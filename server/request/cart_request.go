package request

type AddToCartRequest struct{
	ProductId    string `json:"productId" validate:"required"`
	UserId string `json:"userId" validate:"required"`
	ProductCount float64 `json:"productCount" validate:"required"`
}

type RemoveFromCart struct{
	CartId string `json:"cartId" validate:"required"`
	ProductId string `json:"productId" validate:"required"`
	UserId string `json:"userId" validate:"required"`
}

type RemoveProduct struct{
	CartId string `json:"cartId" validate:"required"`
	ProductId string `json:"productId" validate:"required"`
	UserId string `json:"userId" validate:"required"`
	ProductCount float64 `json:"productCount" validate:"required"`
}