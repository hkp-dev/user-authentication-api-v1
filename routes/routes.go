package routes

import (
	"app/controllers"
	mdware "app/middleware"
	"net/http"
)

func AddRoutes() {
	RoutesAdmin()
	RoutesUser()
	RoutesPublic()
}

func RoutesAdmin() {
	//product
	http.HandleFunc("/admin/add-products", mdware.RequireAdminAuth(controllers.CreateProduct))
	http.HandleFunc("/admin/update-price-product", mdware.RequireAdminAuth(controllers.UpdatePriceProduct))
	http.HandleFunc("/admin/delete-product", mdware.RequireAdminAuth(controllers.DeleteProduct))
	//category
	http.HandleFunc("/admin/add-categories", mdware.RequireAdminAuth(controllers.CreateCategory))
	http.HandleFunc("/admin/delete-categories", mdware.RequireAdminAuth(controllers.DeleteCategory))
	//user
	http.HandleFunc("/admin/lock-user", mdware.RequireAdminAuth(controllers.LockUser))
	http.HandleFunc("/admin/unlock-user", mdware.RequireAdminAuth(controllers.UnlockUser))
	http.HandleFunc("/admin/delete-user", mdware.RequireAdminAuth(controllers.DeleteUser))
}
func RoutesUser() {
	http.HandleFunc("/user/change-password", mdware.RequireUserAuth(controllers.ChangePassword))
	http.HandleFunc("/user/get-all-categories", mdware.RequireUserAuth(controllers.GetAllCategory))
	http.HandleFunc("/user/get-all-products", mdware.RequireUserAuth(controllers.GetAllProducts))
	http.HandleFunc("/user/get-product-by-category", mdware.RequireUserAuth(controllers.GetProductsByCategory))
	http.HandleFunc("/user/get-all-products-by-title", mdware.RequireUserAuth(controllers.GetAllProductsByTitle))
	http.HandleFunc("/user/add-product-to-cart", mdware.RequireUserAuth(controllers.AddProductToCart))
	http.HandleFunc("/user/get-cart", mdware.RequireUserAuth(controllers.GetCart))
	http.HandleFunc("/user/remove-product-from-cart", mdware.RequireUserAuth(controllers.RemoveProductFromCart))
}
func RoutesPublic() {
	http.HandleFunc("/register", mdware.LogRequestMiddleware(controllers.Register))
	http.HandleFunc("/login", mdware.LogRequestMiddleware(controllers.Login))
	http.HandleFunc("/forgot-password", mdware.LogRequestMiddleware(controllers.ForgotPassword))
	http.HandleFunc("/verify-otp", mdware.LogRequestMiddleware(controllers.VerifyOTP))
}
