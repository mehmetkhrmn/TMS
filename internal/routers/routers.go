package routers

import (
	"TMS/internal/models"
	"TMS/internal/repository"
	"database/sql"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func SetupRouter(db *sql.DB) *gin.Engine {

	router := gin.Default()
	gin.SetMode(gin.DebugMode)
	repo := repository.NewRepository(db)
	router.POST("/login", func(context *gin.Context) {
		login(context, repo)
	})
	router.PUT("/answers/:answer_id", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("answer_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "param: answer_id is invalid"})
			return
		}
		var answer models.Answer
		if err := context.ShouldBind(&answer); err != nil { //yine aynı mantıkla var olan verileri bindiliyoruz id yi urlden alcaz
			context.JSON(http.StatusBadRequest, gin.H{"error": "can't bind answer" + err.Error()})
			return
		}
		updateAnswer(id, &answer, context, repo)
	})
	router.POST("/tickets/:ticket_id/answers", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("ticket_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		var answer models.Answer
		if err := context.ShouldBind(&answer); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "can't bind answer" + err.Error()})
			return
		}
		answer.TicketID = id //burada JSON dan alınan değeri değiştiriyoruz JSON daki veri manipüle edilebilir
		createAnswer(&answer, context, repo)
	})

	router.GET("/tickets/:ticket_id/answers", func(context *gin.Context) {
		id, err := strconv.Atoi(context.Param("ticket_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		getAnswers(id, context, repo)
	})
	router.GET("/tickets/:ticket_id/answers/:answer_id", func(context *gin.Context) {
		tid, err := strconv.Atoi(context.Param("ticket_id"))

		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		aid, err := strconv.Atoi(context.Param("answer_id"))
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "answer_id is invalid" + err.Error()})
			return
		}
		getAnswer(tid, aid, context, repo)
	})
	router.PUT("/tickets/:ticket_id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("ticket_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		var ticket models.Ticket
		if err := c.ShouldBindJSON(&ticket); err != nil { //burada bindliyoruz biz gönderilmeyen veriler boş kalcak bu sayede
			c.JSON(http.StatusBadRequest, gin.H{"error": "can't bind ticket" + err.Error()})
			return
		}
		updateTicket(id, &ticket, c, repo)
	})
	//statusa göre ticket döndür
	router.GET("tickets/", func(context *gin.Context) {
		status := context.Query("ticket_status")

		switch status {
		case "open", "closed", "in_progress", "resolved":
			getTicketWith(status, context, repo)

		case "":
			getAllTickets(context, repo)
		default:
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_status is unknown -> " + status})
			return
		}
	})

	//tek bir ticket görüntüleme
	router.GET("/tickets/:ticket_id", func(context *gin.Context) {
		idString := (context.Param("ticket_id")) //url den id yi aldık
		id, err := strconv.Atoi(idString)        //int ye çevir
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		getTicket(id, context, repo)
	})
	//set status
	router.PATCH("tickets/:ticket_id", func(context *gin.Context) {
		idString := context.Param("ticket_id")
		statusString := context.Query("status") //gin iki tane parametre almıyo o yüzden query den alıyoz
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "ticket_id is invalid" + err.Error()})
			return
		}
		status := statusString

		setTicketStatus(id, status, context, repo)

	})
	//bütün ticketleri almak için
	router.GET("/tickets", func(context *gin.Context) {
		getAllTickets(context, repo)
	})
	//oluşturmak için
	router.POST("/tickets", func(context *gin.Context) {
		var ticket models.Ticket
		if err := context.ShouldBind(&ticket); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "cant bind ticket" + err.Error()})
			return
		}
		createTicket(&ticket, context, repo)
	})
	router.GET("/customers", func(context *gin.Context) {
		getCustomers(context, repo)
	})
	router.POST("/customers", func(context *gin.Context) {

		register(context, repo)

	})
	router.GET("/customers/:customer_id", func(context *gin.Context) {
		idString := (context.Param("customer_id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is invalid" + err.Error()})
			return
		}
		getCustomer(id, context, repo)

	})
	router.PUT("/customers/:customers_id", func(context *gin.Context) {
		idString := (context.Param("customers_id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "customer_id is invalid" + err.Error()})
			return
		}
		var customer models.Customer
		if err := context.ShouldBind(&customer); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "cant bind customer " + err.Error()})
			return
		}
		updateCustomer(id, &customer, context, repo)
	})
	router.GET("/representatives", func(context *gin.Context) {
		getAllRepresentatives(context, repo)
	})
	router.POST("/representatives", func(context *gin.Context) {
		register(context, repo)
	})
	router.PUT("/representatives/:representatives_id", func(context *gin.Context) {
		idString := (context.Param("representatives_id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid" + err.Error()})
			return
		}
		var rep models.Representative
		if err := context.ShouldBind(&rep); err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "cant bind representatives" + err.Error()})
			return
		}
		updateRepresentative(id, &rep, context, repo)
	})
	router.GET("/representatives/:representatives_id", func(context *gin.Context) {
		idString := (context.Param("representatives_id"))
		id, err := strconv.Atoi(idString)
		if err != nil {
			context.JSON(http.StatusBadRequest, gin.H{"error": "id is invalid" + err.Error()})
			return
		}
		getRepresentative(id, context, repo)
	})
	return router
}
func createTicket(ticket *models.Ticket, c *gin.Context, repo *repository.Repository) { //repositorydeki structı gönderdik bağlantı kurudk db ile

	//database ekliyoruz
	if err := repo.CreateTicket(ticket); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cant create ticket" + err.Error()})
	}
	c.JSON(http.StatusCreated, ticket) //burada da create ticket fonksyonu içide yazdığımız returning ile güncellenmiş json var

}

//todo:bütün ticketler
func getAllTickets(c *gin.Context, repo *repository.Repository) {
	tickets, error := repo.GetAllTickets()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error})
		return
	}
	c.JSON(http.StatusOK, tickets)
}

//todo:açık olan ticketleri  getir
func getTicketWith(status string, c *gin.Context, repo *repository.Repository) {

	tickets, error := repo.GetAllWithStatus(status)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(http.StatusOK, tickets)

}

//todo:tek bir ticketi görüntüle
func getTicket(id int, c *gin.Context, repo *repository.Repository) {
	ticket, error := repo.GetTicket(id)

	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	if ticket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ticket not found"})
		return
	}
	c.JSON(http.StatusOK, ticket)
}

// todo ticketi editle
func setTicketStatus(id int, status string, c *gin.Context, repo *repository.Repository) {
	err := repo.SetStatus(id, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"id": id, "status": status})
}

func updateTicket(id int, ticket *models.Ticket, c *gin.Context, repo *repository.Repository) {

	err := repo.UpdateTicket(id, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
	}
	c.JSON(200, ticket)
}

// get update create
func createAnswer(answer *models.Answer, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateAnswer(answer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "repo"})
	}
	c.JSON(201, answer)
}
func getAnswer(ticketId int, answerId int, c *gin.Context, repo *repository.Repository) {
	answer, error := repo.GetAnswer(ticketId, answerId)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	if answer == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "answer not found"})
		return
	}
	c.JSON(200, answer)
}
func updateAnswer(id int, ticket *models.Answer, c *gin.Context, repo *repository.Repository) {
	err := repo.UpdateAnswer(id, ticket)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, ticket)

}
func getAnswers(id int, c *gin.Context, repo *repository.Repository) {
	answers, error := repo.GetAnswers(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}

	c.JSON(200, answers)
}
func createCustomer(customer *models.Customer, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateCustomer(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(201, customer)

}
func createCustomerTx(customer *models.Customer, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateCustomer(customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(201, customer)

}
func getCustomers(c *gin.Context, repo *repository.Repository) {
	customers, error := repo.GetCustomers()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, customers)
}
func getCustomer(id int, c *gin.Context, repo *repository.Repository) {
	customer, error := repo.GetCustomer(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, customer)
}
func updateCustomer(id int, customer *models.Customer, c *gin.Context, repo *repository.Repository) {
	_, err := repo.UpdateCustomer(id, customer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, "updated")
}
func createRepresentative(representative *models.Representative, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateRepresentative(representative)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, representative)

}
func createRepresentativeTx(representative *models.Representative, c *gin.Context, repo *repository.Repository) {
	err := repo.CreateRepresentative(representative)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, representative)

}

func getRepresentative(id int, c *gin.Context, repo *repository.Repository) {
	representative, error := repo.GetRepresentative(id)
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	c.JSON(200, representative)
}
func getAllRepresentatives(c *gin.Context, repo *repository.Repository) {
	representatives, err := repo.GetAllRepresentatives()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, representatives)
}
func updateRepresentative(id int, representative *models.Representative, c *gin.Context, repo *repository.Repository) {
	err := repo.UpdateRepresentative(id, representative)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(201, "updated")
}

// -----------------------------------------AUTH
func login(c *gin.Context, repo *repository.Repository) {
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
	secretKey := []byte(os.Getenv("SECRET_KEY"))
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
func register(c *gin.Context, repo *repository.Repository) {
	tx, error := repo.Db.Begin()
	if error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": error.Error()})
		return
	}
	defer tx.Rollback() //rollback oalcak ama zaten commitli olacak

	path := c.FullPath()

	switch path {
	case "/customers":
		var req models.RegisterRequestCustomer
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := repo.GetAuthUserByUsername(req.Username)
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
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		user, err := repo.GetAuthUserByUsername(req.Username)
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
