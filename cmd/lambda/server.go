package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/dgrijalva/jwt-go"
	"github.com/jasonblanchard/di-notebook/pkg/app"
	"github.com/jasonblanchard/di-notebook/pkg/openapi"
	"github.com/jasonblanchard/di-notebook/pkg/store/postgres"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

type Server struct {
	*app.App
	Logger *zap.Logger
	Port   string
}

func initConfig(cfgFile string) error {
	viper.SetConfigFile(cfgFile)
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	return nil
}

// NewServer Create a new Server from env
func NewServer() (*Server, error) {
	dbUser := viper.GetString("DB_USER")
	dbPassword := viper.GetString("DB_PASSWORD")
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	database := viper.GetString("DATABASE")
	pretty := viper.GetBool("PRETTY")
	port := viper.GetString("PORT")

	s := &Server{}

	s.Port = port

	db, err := postgres.NewConnection(&postgres.NewConnectionInput{
		User:     dbUser,
		Password: dbPassword,
		Dbname:   database,
		Host:     dbHost,
		Port:     dbPort,
	})

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create database")
	}

	reader := &postgres.Reader{
		Db: db,
	}

	writer := &postgres.Writer{
		Db: db,
	}

	s.App = &app.App{
		StoreReader: reader,
		StoreWriter: writer,
	}

	var logger *zap.Logger

	if pretty == true {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, errors.Wrap(err, "Failed to create logger")
	}

	s.Logger = logger

	return s, nil
}

func (s *Server) HandleMeta(c echo.Context) error {
	apiGwContext, err := echoLambda.GetAPIGatewayContext(c.Request())
	context := fmt.Sprintf("%+v", apiGwContext)

	requestId := apiGwContext.RequestID
	stage := apiGwContext.Stage

	authorizer := apiGwContext.Authorizer

	version := lambdacontext.FunctionVersion

	authorizationHeader := c.Request().Header["Authorization"]

	if err != nil {
		c.JSON(500, err)
		return nil
	}

	c.JSON(200, map[string]interface{}{
		"context":             context,
		"requestId":           requestId,
		"stage":               stage,
		"authorizer":          authorizer,
		"authorizationHeader": authorizationHeader,
		"version":             version,
	})

	return nil
}

func (s *Server) HandleGetEntry(c echo.Context) error {
	authorizationHeader := c.Request().Header["Authorization"]
	sub, err := bearerHeaderToSub(authorizationHeader[0])
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return err
	}
	userId := getUserIdBySub("https://accounts.google.com", sub)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return err
	}

	getEntryInput := &app.GetEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   userId,
		},
		ID: id,
	}

	entry, err := s.App.GetEntry(getEntryInput)
	if err != nil {
		s.Logger.Error(err.Error())
		c.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	response := map[string]interface{}{
		"id":         entry.ID,
		"creator_id": entry.CreatorID,
		"text":       entry.Text,
		"created_at": entry.CreatedAt,
		"updated_at": entry.UpdatedAt,
	}

	c.JSON(200, map[string]interface{}(response))
	return nil
}

func (s *Server) ListEntries(ctx echo.Context, params openapi.ListEntriesParams) error {
	authorizationHeader := ctx.Request().Header["Authorization"]
	sub, err := bearerHeaderToSub(authorizationHeader[0])
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return err
	}
	userId := getUserIdBySub("https://accounts.google.com", sub)

	var after int
	pageToken := params.PageToken
	after, err = strconv.Atoi(*pageToken)
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return err
	}

	input := &app.ListEntriesInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   userId,
		},
		CreatorID: userId,
		First:     int(params.PageSize),
		After:     after,
	}

	output, err := s.App.ListEntries(input)

	entries := []openapi.Entry{}

	for _, entry := range output.Entries {
		id := fmt.Sprintf("%v", entry.ID)
		entry := openapi.Entry{
			Id:        &id,
			CreatorId: &entry.CreatorID,
			CreatedAt: &entry.CreatedAt,
			UpdatedAt: &entry.UpdatedAt,
		}
		entries = append(entries, entry)
	}

	response := map[string]interface{}{
		"next_page_token": fmt.Sprintf("%d", output.Pagination.EndCursor),
		"total_size":      int32(output.Pagination.TotalCount),
		"has_next_page":   output.Pagination.HasNextPage,
		"entries":         entries,
	}

	ctx.JSON(200, response)

	return nil
}

func (s *Server) HandleListEntries(ctx echo.Context) error {
	authorizationHeader := ctx.Request().Header["Authorization"]
	sub, err := bearerHeaderToSub(authorizationHeader[0])
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}
	userId := getUserIdBySub("https://accounts.google.com", sub)

	var after int
	pageToken := ctx.QueryParam("page_token")
	if pageToken == "" {
		after = 0
	} else {
		after, err = strconv.Atoi(pageToken)
	}

	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	var first int
	pageSize := ctx.QueryParam("page_size")
	if pageSize == "" {
		first = 50
	} else {
		first, err = strconv.Atoi(pageSize)
	}

	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	input := &app.ListEntriesInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   userId,
		},
		CreatorID: userId,
		First:     first,
		After:     after,
	}

	output, err := s.App.ListEntries(input)

	entries := []map[string]interface{}{}

	for _, entry := range output.Entries {
		entry := map[string]interface{}{
			"id":         entry.ID,
			"text":       entry.Text,
			"creator_id": entry.CreatorID,
			"created_at": entry.CreatedAt,
		}
		entries = append(entries, entry)
	}

	response := map[string]interface{}{
		"next_page_token": fmt.Sprintf("%d", output.Pagination.EndCursor),
		"total_size":      int32(output.Pagination.TotalCount),
		"has_next_page":   output.Pagination.HasNextPage,
		"entries":         entries,
	}

	ctx.JSON(200, map[string]interface{}(response))
	return nil
}

func (s *Server) HandleCreateEntry(ctx echo.Context) error {
	authorizationHeader := ctx.Request().Header["Authorization"]
	sub, err := bearerHeaderToSub(authorizationHeader[0])
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return err
	}
	userId := getUserIdBySub("https://accounts.google.com", sub)

	type Body struct {
		Text string `json:"text"`
	}

	body := &Body{}

	ctx.Bind(body)

	input := &app.CreateEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   userId,
		},
		CreatorID: userId,
		Text:      body.Text,
	}

	id, err := s.App.CreateEntry(input)
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return err
	}

	response := map[string]interface{}{
		"id":   id,
		"text": body.Text,
	}

	ctx.JSON(200, map[string]interface{}(response))
	return nil
}

func (s *Server) HandleUpdateEntry(ctx echo.Context) error {
	authorizationHeader := ctx.Request().Header["Authorization"]
	sub, err := bearerHeaderToSub(authorizationHeader[0])
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}
	userId := getUserIdBySub("https://accounts.google.com", sub)
	type Body struct {
		Text string `json:"text"`
	}

	body := &Body{}
	ctx.Bind(body)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	input := &app.UpdateEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   userId,
		},
		ID:   id,
		Text: body.Text,
	}

	entry, err := s.App.UpdateEntry(input)
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	response := map[string]interface{}{
		"id":         entry.ID,
		"creator_id": entry.CreatorID,
		"text":       entry.Text,
		"created_at": entry.CreatedAt,
		"updated_at": entry.UpdatedAt,
	}

	ctx.JSON(200, map[string]interface{}(response))
	return nil
}

func (s *Server) HandleDeleteEntry(ctx echo.Context) error {
	authorizationHeader := ctx.Request().Header["Authorization"]
	sub, err := bearerHeaderToSub(authorizationHeader[0])
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}
	userId := getUserIdBySub("https://accounts.google.com", sub)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	input := &app.DeleteEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   userId,
		},
		ID: id,
	}

	entry, err := s.App.DeleteEntry(input)
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	response := map[string]interface{}{
		"id":         entry.ID,
		"creator_id": entry.CreatorID,
		"text":       entry.Text,
		"created_at": entry.CreatedAt,
		"updated_at": entry.UpdatedAt,
		"date_time":  entry.DeleteTime,
	}

	ctx.JSON(200, map[string]interface{}(response))
	return nil
}

func (s *Server) HandleUndeleteEntry(ctx echo.Context) error {
	authorizationHeader := ctx.Request().Header["Authorization"]
	sub, err := bearerHeaderToSub(authorizationHeader[0])
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}
	userId := getUserIdBySub("https://accounts.google.com", sub)
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	input := &app.UndeleteEntryInput{
		Principal: &app.Principal{
			Type: app.PrincipalUSER,
			ID:   userId,
		},
		ID: id,
	}

	entry, err := s.App.UndeleteEntry(input)
	if err != nil {
		s.Logger.Error(err.Error())
		ctx.JSON(500, map[string]interface{}{
			"error": "Something went wrong",
		})
		return nil
	}

	response := map[string]interface{}{
		"id":         entry.ID,
		"creator_id": entry.CreatorID,
		"text":       entry.Text,
		"created_at": entry.CreatedAt,
		"updated_at": entry.UpdatedAt,
		"date_time":  entry.DeleteTime,
	}

	ctx.JSON(200, map[string]interface{}(response))
	return nil
}

func bearerHeaderToSub(header string) (string, error) {
	tokenString := strings.Replace(header, "Bearer ", "", 1)

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return "", nil
	})

	if token == nil {
		return "", errors.New("Invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok != true {
		return "", errors.New("Token does not contain any claims")
	}

	return fmt.Sprintf("%s", claims["sub"]), nil
}

func getUserIdBySub(issuer, sub string) string {
	googleSubs := map[string]string{
		"103156652160725955399": "2b5545ef-3557-4f52-994d-daf89e04c390",
	}

	id, _ := googleSubs[sub]

	return id
}
