package handlers

import (
	"TMS/ticket-service/internal/models"
	"TMS/ticket-service/internal/repository"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *gin.Context, repo *repository.Repository) {
	var loginReq models.LoginRequest
	err := c.ShouldBindJSON(&loginReq)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	authUser, err := repo.GetAuthUserByUsername(loginReq.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if authUser == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(authUser.PasswordHash), []byte(loginReq.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{ //jwt token oluştu
		"user_id":   authUser.ID,
		"role":      authUser.Role,
		"entity_id": authUser.EntityId,
		"exp":       time.Now().Add(24 * time.Hour).Unix(),
		"iat":       time.Now().Unix(),
	})
	//zaten jwt secret mainde kontrol ediliyor
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate token",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})

}
func Register(c *gin.Context, repo *repository.Repository) {
	tx, error := repo.Db.Begin()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	defer func() {
		_ = tx.Rollback()

	}() //rollback olacak ama zaten commitli olacak

	path := c.FullPath()

	switch path {
	case "/customers":
		var req models.RegisterRequestCustomer
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := repo.GetAuthUserByUsernameTx(tx, req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if user != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "username already exists" + user.Username,
			})
			return
		}
		avab, err := repo.IsMailAvailable(tx, req.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if avab == false {
			c.JSON(http.StatusConflict, gin.H{
				"error": "username already registered" + user.Username,
			})
			return
		}

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		passwordHash := string(hash)
		customer := &models.Customer{
			Name:  req.Name,
			Email: req.Email,
		}
		err = repo.CreateCustomerTx(tx, customer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var insert models.AuthUser
		insert = models.AuthUser{
			Username:     req.Username,
			PasswordHash: passwordHash,
			Role:         "customer",
			EntityId:     customer.ID,
		}
		err = repo.CreateAuthUser(tx, &insert)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "user created",
		})

	case "/representatives":
		var req models.RegisterRequestRepresentative
		role := c.GetString("role")
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"error": "forbidden",
			})
			return
		}
		err := c.ShouldBindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := repo.GetAuthUserByUsernameTx(tx, req.Username)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if user != nil {
			c.JSON(http.StatusConflict, gin.H{
				"error": "username already exists" + user.Username,
			})
			return
		}

		//create user
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(req.Password),
			bcrypt.DefaultCost,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
			return
		}

		passwordHash := string(hash)
		representative := &models.Representative{
			Name: req.Name,
		}
		err = repo.CreateRepresentativeTx(tx, representative)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var insert models.AuthUser
		insert = models.AuthUser{
			Username:     req.Username,
			PasswordHash: passwordHash,
			Role:         "representative",
			EntityId:     representative.ID,
		}
		err = repo.CreateAuthUser(tx, &insert)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		err = tx.Commit()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"message": "user created",
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid role"})
	}
}
func UpdatePassword(c *gin.Context, repo *repository.Repository) {

	var req models.PasswordUpdateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID := int(c.GetFloat64("user_id"))
	user, err := repo.GetAuthUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "password incorrect"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	passwordHash := string(hash)
	err = repo.UpdatePassword(userID, passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "password updated",
	})

}
