package controllers

import (
	"fmt"
	"strconv"

	"github.com/ndbeals/britannicus-final-project/forms"
	"github.com/ndbeals/britannicus-final-project/models"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

//UserController ...
type UserController struct{}

var userModel = new(models.UserModel)

//getUserID ...
func GetUserID(c *gin.Context) int {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID != nil {
		return models.ConvertToInt64(userID)
	}
	return 0
}

//GetLoggedinUser
func GetLoggedinUser(c *gin.Context) (user models.User, success bool) {
	session := sessions.Default(c)
	userID := session.Get("user_id")

	if userID != nil {
		user.ID = models.ConvertToInt64(userID)
		user.Name = session.Get("user_name").(string)
		user.Email = session.Get("user_email").(string)
	} else {
		return models.User{}, false
	}

	return user, true
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
	var success bool
	// fmt.Println("WHAT IN THE FUCK")
	// fmt.Println("EMAIL FORM", c.DefaultPostForm("email", "default"))
	test, _ := c.GetPostForm("email")
	fmt.Println("ENAIL", test)
	signinForm.Email = test
	signinForm.Password, success = c.GetPostForm("password")

	// if c.BindJSON(&signinForm) != nil {
	// 	c.IndentedJSON(406, gin.H{"message": "Invalid form", "form": signinForm})
	// 	c.Abort()
	// 	return
	// }
	fmt.Println(signinForm)
	if !success {
		// c.(406, "failed to log in")
		fmt.Println("failed")
		c.Abort()
	}

	user, err := userModel.Signin(signinForm)
	if err == nil {
		session := sessions.Default(c)
		fmt.Println("Ifds", session)
		session.Set("user_id", user.ID)
		session.Set("user_email", user.Email)
		session.Set("user_name", user.Name)
		// c.Set("user", user)
		// session.Set("user", user)
		session.Save()

		// c.IndentedJSON(200, gin.H{"message": "User signed in", "user": user})
		c.Redirect(303, "/")
	} else {
		c.IndentedJSON(406, gin.H{"message": "Invalid signin details", "error": err.Error()})
	}

}

//Signup ...
func (ctrl UserController) Signup(c *gin.Context) {
	var signupForm forms.SignupForm

	if c.BindJSON(&signupForm) != nil {
		c.IndentedJSON(406, gin.H{"message": "Invalid form", "form": signupForm})
		c.Abort()
		return
	}

	user, err := userModel.Signup(signupForm)

	if err != nil {
		c.IndentedJSON(406, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	if user.ID > 0 {
		session := sessions.Default(c)
		session.Set("user_id", user.ID)
		session.Set("user_email", user.Email)
		session.Set("user_name", user.Name)
		session.Save()
		c.IndentedJSON(200, gin.H{"message": "Success signup", "user": user})
	} else {
		c.IndentedJSON(406, gin.H{"message": "Could not signup this user", "error": err.Error()})
	}

}

//Signout ...
func (ctrl UserController) Signout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.IndentedJSON(200, gin.H{"message": "Signed out..."})
}

//GetOne ...
func (ctrl UserController) GetOne(c *gin.Context) {
	userid := c.Param("id")

	if userid, err := strconv.ParseInt(userid, 10, 32); err == nil {
		userid := int(userid)

		data, err := userModel.GetOne(userid)
		if err != nil {
			c.IndentedJSON(404, gin.H{"Message": "Article not found", "error": err.Error()})
			c.Abort()
			return
		}
		c.IndentedJSON(200, gin.H{"data": data})
	} else {
		c.IndentedJSON(404, gin.H{"Message": "Invalid parameter"})
	}
}
