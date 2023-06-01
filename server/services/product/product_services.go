package product

import (
	"main/server/context"
	"main/server/db"
	"main/server/model"
	"main/server/response"
	"main/server/services/order"
	"main/server/utils"

	"github.com/gin-gonic/gin"
)

// Inventory product add
func InventoryProductAddService(ctx *gin.Context, addProduct model.Products) {
	vendorId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
		return
	}
	if vendorId == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Vendor ID provided")
		return
	}
	if db.RecordExist("products", "product_name", addProduct.ProductName) {
		response.ErrorResponse(ctx, 400, "Product already exists")
		return
	}
	if addProduct.ProductCategory == "mobile" {
		addProduct.ProductCategory = "MOBILE"
		var dbconstants model.DbConstant
		dbconstants.ConstantName = "mobile"
		dbconstants.ConstantShortHand = "MOBILE"
		if db.RecordExist("db_constants", "constant_name", addProduct.ProductCategory) {
			return
		}
		err = db.CreateRecord(&dbconstants)
		if err != nil {
			response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Not able to create record")
			return
		}
	}
	addProduct.ProductInventory++
	err = db.CreateRecord(&addProduct)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Not able to create record")
		return
	}

	var vendorInventory model.VendorInventory
	vendorInventory.VendorId = vendorId
	vendorInventory.ProductId = addProduct.ProductId
	err = db.CreateRecord(&vendorInventory)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Not able to create record")
		return
	}

	response.ShowResponse(
		"Success",
		200,
		"Product added successfully",
		addProduct,
		ctx,
	)
}

// Inventory product Update
func InventoryProductUpdateService(ctx *gin.Context, productInventoryEdit model.Products) {
	vendorId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
		return
	}
	if vendorId == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Vendor ID provided")
		return
	}

	if db.BothExists("vendor_inventories", "vendor_id", vendorId, "product_id", productInventoryEdit.ProductId) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "This product does not belong to given vendor")
		return
	}

	db.UpdateRecord(&productInventoryEdit, productInventoryEdit.ProductId, "product_id")

	err = db.FindById(&productInventoryEdit, productInventoryEdit.ProductId, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Not able to find record")
		return
	}
	response.ShowResponse(
		"Success",
		200,
		"Product updated successfully",
		productInventoryEdit,
		ctx,
	)
}

// Inventory product Delete
func InventoryProductDeleteService(ctx *gin.Context, productInventoryDelete context.ProductDeleteRequest) {
	vendorId, err := order.IdFromToken(ctx)
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_UNAUTHORIZED, "Invalid token")
		return
	}
	if vendorId == "" {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "No Vendor ID provided")
		return
	}
	if !db.RecordExist("products", "product_id", productInventoryDelete.ProductID) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Product doesn't exist")
		return
	}
	var vendorInventory model.VendorInventory
	var productInventory model.Products
	if db.BothExists("vendor_inventories", "vendor_id", vendorId, "product_id", productInventoryDelete.ProductID) {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "This product does not belong to given vendor")
		return
	}

	err = db.Delete(&vendorInventory, productInventoryDelete.ProductID, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Unable to delete record")
		return
	}

	err = db.DeleteRecord(&productInventory, productInventoryDelete.ProductID, "product_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Unable to delete record")
		return
	}

	err = db.FindById(&vendorInventory, vendorId, "vendor_id")
	if err != nil {
		response.ErrorResponse(ctx, utils.HTTP_BAD_REQUEST, "Not able to find record")
		return
	}
	response.ShowResponse(
		"Success",
		200,
		"Product deleted successfully",
		vendorInventory,
		ctx,
	)
}
