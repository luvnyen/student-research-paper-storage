package controller

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/luvnyen/student-research-paper-storage/pkg/dto"
	"github.com/luvnyen/student-research-paper-storage/pkg/utils"
	"github.com/luvnyen/student-research-paper-storage/service"
)

type PaperController interface {
	All(ctx *gin.Context)
	FindByID(ctx *gin.Context)
	FindByTitleAuthorAbstract(ctx *gin.Context)
	InsertPaper(ctx *gin.Context)
	DownloadPaper(ctx *gin.Context)
}

type paperController struct {
	paperService service.PaperService
	jwtService   service.JWTService
}

func NewPaperController(paperService service.PaperService, jwtService service.JWTService) PaperController {
	return &paperController{paperService: paperService, jwtService: jwtService}
}

func (db *paperController) ValidateAuthorization(ctx *gin.Context) string {
	authHeader := ctx.GetHeader("Authorization")
	token, err := db.jwtService.ValidateToken(authHeader)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return ""
	}
	claims := token.Claims.(jwt.MapClaims)
	studentID := fmt.Sprintf("%v", claims["user_id"])

	return studentID
}

func (db *paperController) All(ctx *gin.Context) {
	authHeader := ctx.GetHeader("Authorization")
	token, err := db.jwtService.ValidateToken(authHeader)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", "Invalid token", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	studentID := fmt.Sprintf("%v", claims["user_id"])

	papers, err := db.paperService.All(studentID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	response := utils.BuildResponse(true, "Successfully get all papers", papers)
	ctx.JSON(200, response)
}

func (db *paperController) FindByID(ctx *gin.Context) {
	id := ctx.Param("id")

	res, err := db.paperService.FindByID(id)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	response := utils.BuildResponse(true, "Successfully get paper", res)
	ctx.JSON(200, response)
}

func (db *paperController) FindByTitleAuthorAbstract(ctx *gin.Context) {
	var searchDTO dto.SearchDTO

	errDTO := ctx.ShouldBind(&searchDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Failed to process request", errDTO.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, err := db.jwtService.ValidateToken(authHeader)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", "Invalid token", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	studentID := fmt.Sprintf("%v", claims["user_id"])

	res, err := db.paperService.FindByTitleAuthorAbstract(searchDTO, studentID)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	response := utils.BuildResponse(true, "Successfully get papers", res)
	ctx.JSON(200, response)
}

func (db *paperController) InsertPaper(ctx *gin.Context) {
	var paperDTO dto.PaperDTO
	errDTO := ctx.ShouldBind(&paperDTO)
	if errDTO != nil {
		response := utils.BuildErrorResponse("Failed to process request DTO", errDTO.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	authHeader := ctx.GetHeader("Authorization")
	token, err := db.jwtService.ValidateToken(authHeader)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", "Invalid token", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	studentID := fmt.Sprintf("%v", claims["user_id"])

	file, err := ctx.FormFile("file")
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	// fileBytes, err := file.Open()
	// if err != nil {
	// 	response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
	// 	ctx.AbortWithStatusJSON(400, response)
	// 	return
	// }
	// defer fileBytes.Close()
	// fileBase64, err := ioutil.ReadAll(fileBytes)
	// if err != nil {
	// 	response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
	// 	ctx.AbortWithStatusJSON(400, response)
	// 	return
	// }

	// fileBase64String := base64.StdEncoding.EncodeToString(fileBase64)

	extension := filepath.Ext(file.Filename)
	if extension != ".pdf" {
		response := utils.BuildErrorResponse("Failed to process request", "Invalid file extension", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	newFileName := uuid.New().String() + extension

	absPath, err := filepath.Abs("./cdn/paper/" + newFileName)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	if err := ctx.SaveUploadedFile(file, absPath); err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	res, err := db.paperService.InsertPaper(paperDTO, studentID, newFileName)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	response := utils.BuildResponse(true, "Successfully insert paper", res)
	ctx.JSON(200, response)
}

func (db *paperController) DownloadPaper(ctx *gin.Context) {
	id := ctx.Param("id")

	authHeader := ctx.GetHeader("Authorization")
	token, err := db.jwtService.ValidateToken(authHeader)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", "Invalid token", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}
	claims := token.Claims.(jwt.MapClaims)
	studentID := fmt.Sprintf("%v", claims["user_id"])
	studentIDUint64, err := strconv.ParseUint(studentID, 10, 64)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	res, err := db.paperService.FindByID(id)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	log.Println("res.StudentID: ", res.Student.ID)
	log.Println("studentIDUint64: ", studentIDUint64)
	log.Println("studentID: ", studentID)

	if res.Student.ID != studentIDUint64 {
		response := utils.BuildErrorResponse("Failed to process request", "You are not allowed to download this paper because you are not the owner", utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	absPath, err := filepath.Abs("./cdn/paper/" + res.File)
	if err != nil {
		response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
		ctx.AbortWithStatusJSON(400, response)
		return
	}

	ctx.File(absPath)

	// fileBase64, err := base64.RawStdEncoding.DecodeString(res.File)
	// if err != nil {
	// 	response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
	// 	ctx.AbortWithStatusJSON(400, response)
	// 	return
	// }

	// fileName := uuid.New().String() + ".pdf"
	// absPath, err := filepath.Abs("./cdn/paper/" + fileName)
	// if err != nil {
	// 	response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
	// 	ctx.AbortWithStatusJSON(400, response)
	// 	return
	// }

	// file, err := os.Create(absPath)
	// if err != nil {
	// 	response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
	// 	ctx.AbortWithStatusJSON(400, response)
	// 	return
	// }
	// defer file.Close()

	// _, err = file.Write(fileBase64)
	// if err != nil {
	// 	response := utils.BuildErrorResponse("Failed to process request", err.Error(), utils.EmptyObj{})
	// 	ctx.AbortWithStatusJSON(400, response)
	// 	return
	// }

	// ctx.Header("Content-Description", "File Transfer")
	// ctx.Header("Content-Disposition", "attachment; filename="+fileName)
	// ctx.Header("Content-Type", "application/pdf")
	// ctx.Header("Content-Transfer-Encoding", "binary")
	// ctx.File(absPath)
}
