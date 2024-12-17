package routes

import (
	"app/controllers"
	mdware "app/middleware"
	"net/http"
)

func AddRoutes() {
	//handler for auth
	http.HandleFunc("/register", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.Register)).ServeHTTP)
	http.HandleFunc("/login", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.Login)).ServeHTTP)
	http.HandleFunc("/change-password", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.ChangePassword)).ServeHTTP)
	http.HandleFunc("/forgot-password", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.ForgotPassword)).ServeHTTP)
	http.HandleFunc("/verify-otp", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.VerifyOTP)).ServeHTTP)

	//handler for categories
	http.HandleFunc("/add-categories", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.CreateCategory)).ServeHTTP)
	//handler for products
	http.HandleFunc("/add-products", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.CreateProduct)).ServeHTTP)
	http.HandleFunc("/get-all-products-by-title", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.GetAllProductsByTitle)).ServeHTTP)
	http.HandleFunc("/get-all-products", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.GetAllProducts)).ServeHTTP)
	http.HandleFunc("/delete-product", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.DeleteProduct)).ServeHTTP)
	http.HandleFunc("/update-product", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.UpdatePriceProduct)).ServeHTTP)
	http.HandleFunc("/get-product-by-category", mdware.LogRequestMiddleware(http.HandlerFunc(controllers.GetProductsByCategory)).ServeHTTP)
}
