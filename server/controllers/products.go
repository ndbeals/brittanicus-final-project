package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ndbeals/britannicus-final-project/forms"
	"github.com/ndbeals/britannicus-final-project/models"

	"github.com/gin-gonic/gin"
)

//ProductController ...
type ProductController struct{}

var productModel = models.GetProductModel()

//GetOne ...
func (ctrl ProductController) GetOne(c *gin.Context) {
	productid := c.Param("id")

	if productid, err := strconv.ParseInt(productid, 10, 32); err == nil {
		productid := int(productid)

		data, err := productModel.GetOne(productid)
		if err != nil {
			c.IndentedJSON(404, gin.H{"Message": "Product not found", "error": err.Error()})
			c.Abort()
			return
		}
		c.IndentedJSON(200, gin.H{"data": data})
	} else {
		c.IndentedJSON(404, gin.H{"Message": "Invalid parameter"})
	}
}

//GetList ...
func (ctrl ProductController) GetList(c *gin.Context) {
	page := c.Param("page")
	amount, err := strconv.ParseInt(c.Param("amount"), 10, 32)

	if err != nil {
		amount = 100
	}

	if page, err := strconv.ParseInt(page, 10, 32); err == nil {
		page, amount := int(page), int(amount)
		data, err := productModel.GetList(page, amount)
		if err != nil {
			c.IndentedJSON(404, gin.H{"Message": "Product not found", "error": err.Error()})
			c.Abort()
			return
		}
		c.IndentedJSON(200, data)
	} else {
		c.IndentedJSON(404, gin.H{"Message": "Invalid parameter"})
	}
}

//Update ...
func (ctrl ProductController) Update(c *gin.Context) {
	productid := c.Param("id")
	fmt.Println(productid)

	if productid, err := strconv.ParseInt(productid, 10, 32); err == nil {
		productid := int(productid)

		var updateForm forms.UpdateProductForm
		err := c.BindJSON(&updateForm)
		if err != nil {
			panic(err)
			c.IndentedJSON(404, gin.H{"message": "Invalid form", "form": updateForm})
			c.Abort()
			return
		}

		product, err := productModel.GetOne(productid)
		if err != nil {
			c.IndentedJSON(404, gin.H{"Message": "Product not found", "error": err.Error()})
			c.Abort()
			return
		}

		product.Update(updateForm)

		fmt.Println("updated")

		c.IndentedJSON(200, gin.H{"data": product})
	} else {
		c.IndentedJSON(404, gin.H{"Message": "Invalid parameter"})
	}
}

//Create ...
func (ctrl ProductController) Create(c *gin.Context) {
	// productid := c.Param("id")
	fmt.Println("create")

	// if productid, err := strconv.ParseInt(productid, 10, 32); err == nil {
	// 	productid := int(productid)

	var updateForm forms.UpdateProductForm
	err := c.BindJSON(&updateForm)
	if err != nil {
		panic(err)
		c.IndentedJSON(404, gin.H{"message": "Invalid form", "form": updateForm})
		c.Abort()
		return
	}
	fmt.Println(updateForm.ISBN, updateForm.ProductName, updateForm.Author, updateForm.Genre, updateForm.Publisher, 1, "Soft Cover", updateForm.Description)

	product := models.Product{-1, updateForm.ISBN, updateForm.ProductName, updateForm.Author, updateForm.Genre, updateForm.Publisher, 1, "Soft Cover", updateForm.Description}
	if err != nil {
		c.IndentedJSON(404, gin.H{"Message": "Product not found", "error": err.Error()})
		c.Abort()
		return
	}

	newid, _ := product.Create()

	fmt.Println("Createds", newid)

	c.IndentedJSON(200, gin.H{"data": product, "id": newid})
	// } else {
	// 	c.IndentedJSON(404, gin.H{"Message": "Invalid parameter"})
	// }
}

//Delete ...
func (ctrl ProductController) Delete(c *gin.Context) {
	productid := c.Param("id")
	fmt.Println("Delete")

	if productid, err := strconv.ParseInt(productid, 10, 32); err == nil {
		productid := int(productid)

		product, err := productModel.GetOne(productid)
		if err != nil {
			c.IndentedJSON(404, gin.H{"Message": "Product not found", "error": err.Error()})
			c.Abort()
			return
		}

		_, err = product.Delete()

		fmt.Println("deleted product from api")

		if err != nil {
			c.IndentedJSON(404, gin.H{"Message": "Failed to delete", "error": err.Error()})
			c.Abort()
			return
		}

		c.IndentedJSON(200, gin.H{"Message": "Product Sucessfully deleted"})
	}
}

//ProductListingPage ...
func (ctrl ProductController) ProductListingPage(c *gin.Context) {
	user, _ := GetLoggedinUser(c)

	c.HTML(http.StatusOK, "products.html", gin.H{
		"title": "Products Page",
		"route": "/products",
		"user":  user,
	})
}
