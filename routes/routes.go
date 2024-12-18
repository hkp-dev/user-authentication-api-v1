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
	http.HandleFunc("/admin/add-products", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.CreateProduct))).ServeHTTP)
	http.HandleFunc("/admin/delete-product", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.DeleteProduct))).ServeHTTP)
	http.HandleFunc("/admin/update-price-product", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.UpdatePriceProduct))).ServeHTTP)
	//category
	http.HandleFunc("/admin/add-categories", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.CreateCategory))).ServeHTTP)
	http.HandleFunc("/admin/get-all-categories", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.GetCategories))).ServeHTTP)
	http.HandleFunc("/admin/delete-category", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.DeleteCategory))).ServeHTTP)

	//user
	http.HandleFunc("/admin/lock-user", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.LockUser))).ServeHTTP)
	http.HandleFunc("/admin/delete-user", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.DeleteUser))).ServeHTTP)
	http.HandleFunc("/admin/unlock-user", mdware.RequireAdminAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.UnlockUser))).ServeHTTP)
}
func RoutesUser() {
	http.HandleFunc("/user/get-all-products-by-title", mdware.RequireUserAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.GetAllProductsByTitle))).ServeHTTP)
	http.HandleFunc("/user/get-all-products", mdware.RequireUserAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.GetAllProducts))).ServeHTTP)
	http.HandleFunc("/user/get-product-by-category", mdware.RequireUserAuth(mdware.LogRequestMiddleware(http.HandlerFunc(controllers.GetProductsByCategory))).ServeHTTP)
}
func RoutesPublic() {
	http.HandleFunc("/register", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.Register)).ServeHTTP)
	http.HandleFunc("/login", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.Login)).ServeHTTP)
	http.HandleFunc("/change-password", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.ChangePassword)).ServeHTTP)
	http.HandleFunc("/forgot-password", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.ForgotPassword)).ServeHTTP)
	http.HandleFunc("/verify-otp", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.VerifyOTP)).ServeHTTP)
}
