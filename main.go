package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	// "os"

	"github.com/gin-gonic/gin"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var users []User

func main() {
	r := gin.Default()

	loadUsersFromJSON()

	r.POST("/login", loginHandler)
	r.GET("/logout", logoutHandler)
	r.POST("/register", registerHandler)

	r.Run(":8080")
}

func loadUsersFromJSON() {
	file, err := ioutil.ReadFile("users.json")
	if err != nil {
		fmt.Println("Error reading users.json:", err)
		return
	}

	err = json.Unmarshal(file, &users)
	if err != nil {
		fmt.Println("Error unmarshaling users.json:", err)
		return
	}
}

func saveUsersToJSON() {
	data, err := json.Marshal(users)
	if err != nil {
		fmt.Println("Error marshaling users:", err)
		return
	}

	err = ioutil.WriteFile("users.json", data, 0644)
	if err != nil {
		fmt.Println("Error writing users.json:", err)
		return
	}
}

func loginHandler(c *gin.Context) {
	var input User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	for _, user := range users {
		if user.Username == input.Username && user.Password == input.Password {
			c.JSON(200, gin.H{"message": "Login successful"})
			return
		}
	}

	c.JSON(401, gin.H{"error": "Invalid credentials"})
}

func logoutHandler(c *gin.Context) {
	c.JSON(200, gin.H{"message": "Logout successful"})
}

func registerHandler(c *gin.Context) {
	var newUser User
	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.JSON(400, gin.H{"error": "Invalid JSON"})
		return
	}

	for _, user := range users {
		if user.Username == newUser.Username {
			c.JSON(409, gin.H{"error": "Username already exists"})
			return
		}
	}

	users = append(users, newUser)
	saveUsersToJSON()
	c.JSON(201, gin.H{"message": "User registered successfully"})
}
