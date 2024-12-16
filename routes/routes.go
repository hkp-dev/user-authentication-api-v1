package routes

import (
	"app/controllers"
	"app/net"
	"net/http"
)

func AddRoutes() {
	//handler for auth
	http.HandleFunc("/register", net.LogRequestMiddleware(http.HandlerFunc(controllers.Register)).ServeHTTP)
	http.HandleFunc("/login", net.LogRequestMiddleware(http.HandlerFunc(controllers.Login)).ServeHTTP)
	http.HandleFunc("/change-password", net.LogRequestMiddleware(http.HandlerFunc(controllers.ChangePassword)).ServeHTTP)
	http.HandleFunc("/forgot-password", net.LogRequestMiddleware(http.HandlerFunc(controllers.ForgotPassword)).ServeHTTP)
	http.HandleFunc("/verify-otp", net.LogRequestMiddleware(http.HandlerFunc(controllers.VerifyOTP)).ServeHTTP)


	//handler for categories
	http.HandleFunc("/add-categories", net.LogRequestMiddleware(http.HandlerFunc(controllers.CreateCategory)).ServeHTTP)

	
	//handler for products
	http.HandleFunc("/add-products", net.LogRequestMiddleware(http.HandlerFunc(controllers.CreateProduct)).ServeHTTP)
}
