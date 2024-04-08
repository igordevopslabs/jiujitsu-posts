package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/igordevopslabs/jiujitsu-posts/config"
	"github.com/igordevopslabs/jiujitsu-posts/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	// Cria um novo banco de dados mock.
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	// Configura o GORM para usar o banco de dados mock com o driver do PostgreSQL.
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  "host=localhost user=admin password=admin dbname=jiujitsu_posts_test port=5432 sslmode=disable",
		PreferSimpleProtocol: true, // usa o protocolo simples para maior compatibilidade.
		Conn:                 db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return gormDB, mock
}

func TestGetAllPosts(t *testing.T) {
	// Inicializa o Gin e coloca o modo de teste
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	gormDB, mock := SetupMockDB(t)

	// Substitui a config.DB global (ou similar) pelo GORM mock.
	config.DB = gormDB

	// Configura o SQL mock
	rows := sqlmock.NewRows([]string{"id", "title", "content"}).
		AddRow(1, "Post 1", "Content 1").
		AddRow(2, "Post 2", "Content 2")
	mock.ExpectQuery("SELECT \\* FROM \"posts\"").WillReturnRows(rows) // Nota: ajuste na regex para aspas duplas

	// Adiciona a rota ao router do Gin
	router.GET("/posts", GetAllPosts)

	// Cria uma requisição GET para /posts
	req, _ := http.NewRequest("GET", "/posts", nil)
	resp := httptest.NewRecorder()

	// Faz a requisição ao Gin
	router.ServeHTTP(resp, req)

	// Verifica se a resposta é a esperada
	assert.Equal(t, 200, resp.Code)
}

func TestGetPostID(t *testing.T) {
	// Inicializa o Gin e coloca ele no modo de teste
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	gormDB, mock := SetupMockDB(t) //invoca a função

	// Substitui a config.DB global (ou similar) pelo GORM mock.
	config.DB = gormDB

	// Configura o SQL mock
	// Pesquisar substituição correta dos placeholders
	// Pesquiser os parametros suportados pelo WithArgs()

	query := `SELECT \* FROM "posts" WHERE "posts"."id" = \$1 AND "posts"."deleted_at" IS NULL ORDER BY "posts"."id" LIMIT \$2`
	mock.ExpectQuery(query).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "content"}).
			AddRow(1, "Post 1", "Content 1"))

	// Adiciona a rota ao router do Gin em teste mode
	router.GET("/posts/:id", GetPostByID)

	// Cria uma requisição GET para /posts
	req, _ := http.NewRequest("GET", "/posts/1", nil)
	resp := httptest.NewRecorder()

	// Faz a requisição ao Gin
	router.ServeHTTP(resp, req)

	// Verifica se a resposta é a esperada
	assert.Equal(t, 200, resp.Code)

	//tentar fazer got pelo body da req
}

func TestDeletePost(t *testing.T) {
	// Inicializa o Gin e coloca ele no modo de teste
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	gormDB, mock := SetupMockDB(t) //invoca a função

	// Substitui a config.DB global (ou similar) pelo GORM mock.
	config.DB = gormDB

	query := `UPDATE "posts" SET "deleted_at"=\$1 WHERE "posts"."id" = \$2 AND "posts"."deleted_at" IS NULL`
	mock.ExpectBegin()
	mock.ExpectExec(query).WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	// Adiciona a rota ao router do Gin em teste mode
	router.DELETE("/post/:id", middleware.RequireAuth, DeletePost)

	// Cria uma requisição GET para /posts
	req, _ := http.NewRequest("DELETE", "/post/113", nil)

	resp := httptest.NewRecorder()

	// Faz a requisição ao Gin
	router.ServeHTTP(resp, req)

	// Verifica se a resposta é a esperada
	assert.Equal(t, 200, resp.Code)
}
