package controllers

import (
	"app/models"
	"app/services"
	"app/utils"
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("views/templates/register.html")
		if err != nil {
			http.Error(w, `{"error": "Failed to parse template"}`, http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	//read body request from client
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	//validate user data from client
	registrationErrors := services.ValidateUser(user)
	if len(registrationErrors) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"errors": registrationErrors,
		})
		return
	}

	// check if username or email already exists
	exists, err := services.CheckUserExists(user)
	if err != nil {
		http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		return
	}
	if exists {
		http.Error(w, `{"error": "Username or email already exists"}`, http.StatusBadRequest)
		return
	}
	//hash password
	hashedPassword, err := services.HashPassword(user.Password)
	if err != nil {
		http.Error(w, `{"error": "Failed to hash password"}`, http.StatusInternalServerError)
		return
	}
	//create user with cart_id
	user.Password = string(hashedPassword)
	var cart models.Cart
	err = services.CreateUser(&user, &cart)
	if err != nil {
		http.Error(w, `{"error": "Failed to create user"}`, http.StatusInternalServerError)
		return
	}
	//create cart with user_id
	err = services.CreateCart(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	//send response with message to client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Registration successful! You can now log in on " + user.Username,
	})
}
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("views/templates/login.html")
		if err != nil {
			http.Error(w, `{"error": "Failed to parse template"}`, http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}

	user, err := services.FindUserByUsername(credentials.Username)
	if err != nil {
		if err.Error() == mongo.ErrNoDocuments.Error() {
			http.Error(w, `{"error": "Invalid username"}`, http.StatusUnauthorized)
		} else {
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}
	if user.Locked {
		http.Error(w, `{"error": "user has been locked"}`, http.StatusBadRequest)
		return
	}
	if !services.ComparePassword(credentials.Password, user) {
		http.Error(w, `{"error": "Invalid password"}`, http.StatusUnauthorized)
		return
	}
	if !services.CheckStatusUser(user) {
		http.Error(w, `{"error": "User is blocked"}`, http.StatusUnauthorized)
		return
	}
	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Login successful",
		"token":   token,
	})
}
func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var user models.User

	user, err := services.GetUserByToken(r)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}

	var request struct {
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	if len(request.NewPassword) < 8 {
		http.Error(w, `{"error": "New password must be at least 8 characters"}`, http.StatusBadRequest)
		return
	}
	user, err = services.FindUserByUsername(user.Username)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, `{"error": "User not found"}`, http.StatusUnauthorized)
		} else {
			http.Error(w, `{"error": "Database error"}`, http.StatusInternalServerError)
		}
		return
	}

	if !services.ComparePassword(request.OldPassword, user) {
		http.Error(w, `{"error": "Old password is incorrect"}`, http.StatusUnauthorized)
		return
	}

	if err = services.UpdatePassword(user.Username, request.NewPassword); err != nil {
		http.Error(w, `{"error": "Failed to update password"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Del("Authorization")
	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Authorization", token)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password updated successfully",
	})
}
func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var input struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	user, err := services.GetUserByEmail(input.Email)
	if err != nil {
		http.Error(w, `{"error": "User not found"}`, http.StatusNotFound)
		return
	}
	otp := services.GenerateOTP()
	expiry := services.GenerateExpiry()
	err = services.SaveOTP(otp, expiry, user)
	if err != nil {
		http.Error(w, `{"error": "Save OTP failed"}`, http.StatusInternalServerError)
		return
	}
	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Click here and enter code to reset password http://localhost:8080/verify-otp",
		"otp":     otp,
		"expiry":  expiry,
	})
}
func VerifyOTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var user models.User
	user, err := services.GetUserByToken(r)
	if err != nil {
		http.Error(w, `{"error": "`+err.Error()+`"}`, http.StatusUnauthorized)
		return
	}
	var request struct {
		OTP         string `json:"otp"`
		NewPassword string `json:"new_password"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, `{"error": "Invalid JSON"}`, http.StatusBadRequest)
		return
	}
	otp, expiry, err := services.GetOTP(user)
	if err != nil {
		http.Error(w, `{"error": "Verify OTP failed"}`, http.StatusInternalServerError)
		return
	}
	if request.OTP != otp || time.Now().After(expiry) {
		http.Error(w, `{"error": "Invalid OTP or OTP has expired"}`, http.StatusBadRequest)
		return
	}

	if err = services.UpdatePassword(user.Username, request.NewPassword); err != nil {
		http.Error(w, `{"error": "Failed to update password"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Del("Authorization")
	token, err := utils.GenerateJWT(user)
	if err != nil {
		http.Error(w, `{"error": "Failed to generate token"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Authorization", token)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password reset successfully",
	})
}
func Logout(w http.ResponseWriter, r *http.Request) {
	w.Header().Del("Authorization")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Logout successfully",
	})
}
