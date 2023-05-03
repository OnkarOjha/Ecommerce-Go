package request

type OrderRequest struct{
	CartId string `json:"cartId" validate:"required"`
	ProductId string `json:"productId" validate:"required"`
}