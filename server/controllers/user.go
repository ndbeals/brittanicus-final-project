package controllers

import (
	"strconv"

	"github.com/ndbeals/brittanicus-final-project/forms"
	"github.com/ndbeals/brittanicus-final-project/models"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//UserController ...
type UserController struct{}

var userModel = new(models.UserModel)

//getUserID ...
func getUserID(c *gin.Context) int64 {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID != nil {
		return models.ConvertToInt64(userID)
	}
	return 0
}

//getSessionUserInfo ...
func getSessionUserInfo(c *gin.Context) (userSessionInfo models.UserSessionInfo) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID != nil {
		userSessionInfo.ID = models.ConvertToInt64(userID)
		userSessionInfo.Name = session.Get("user_name").(string)
		userSessionInfo.Email = session.Get("user_email").(string)
	}
	return userSessionInfo
}

//Signin ...
func (ctrl UserController) Signin(c *gin.Context) {
	var signinForm forms.SigninForm

	if c.BindJSON(&signinForm) != nil {
		c.JSON(406, gin.H{"message": "Invalid form", "form": signinForm})
		c.Abort()
		return
	}

	user, err := userModel.Signin(signinForm)
	if err == nil {
		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		session.Set("user_email", user.Email)
		session.Set("user_name", user.Name)
		session.Save()

		c.JSON(200, gin.H{"message": "User signed in", "user": user})
	} else {
		c.JSON(406, gin.H{"message": "Invalid signin details", "error": err.Error()})
	}

}

//Signup ...
func (ctrl UserController) Signup(c *gin.Context) {
	var signupForm forms.SignupForm

	if c.BindJSON(&signupForm) != nil {
		c.JSON(406, gin.H{"message": "Invalid form", "form": signupForm})
		c.Abort()
		return
	}

	user, err := userModel.Signup(signupForm)

	if err != nil {
		c.JSON(406, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	if user.ID > 0 {
		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		session.Set("user_email", user.Email)
		session.Set("user_name", user.Name)
		session.Save()
		c.JSON(200, gin.H{"message": "Success signup", "user": user})
	} else {
		c.JSON(406, gin.H{"message": "Could not signup this user", "error": err.Error()})
	}

}

//Signout ...
func (ctrl UserController) Signout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(200, gin.H{"message": "Signed out..."})
}

//GetOne ...
func (ctrl UserController) GetOne(c *gin.Context) {
	// userID := getUserID(c)

	// if userID == 0 {
	// 	c.JSON(403, gin.H{"message": "Please login first"})
	// 	c.Abort()
	// 	return
	// }

	uid := c.Param("id")

	if uid, err := strconv.ParseInt(uid, 10, 64); err == nil {

		data, err := userModel.GetOne(uid)
		if err != nil {
			c.JSON(404, gin.H{"Message": "Article not found", "error": err.Error()})
			c.Abort()
			return
		}
		c.JSON(200, gin.H{"data": data})
	} else {
		c.JSON(404, gin.H{"Message": "Invalid parameter"})
	}
}