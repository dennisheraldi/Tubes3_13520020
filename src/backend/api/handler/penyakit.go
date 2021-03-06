package handler

import (
	"api/library"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"api/penyakit"
)

type penyakitHandler struct {
	penyakitService penyakit.Service
}

func NewPenyakitHandler(penyakitService penyakit.Service) *penyakitHandler {
	return &penyakitHandler{ penyakitService }
}

func (h *penyakitHandler) GetAllPenyakitHandler(c *gin.Context) {
	penyakits, err := h.penyakitService.FindAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
	}

	var penyakitResponses []penyakit.PenyakitResponse

	for _, p := range penyakits {
		penyakitResponse := convertToPenyakitResponse(p)
		penyakitResponses = append(penyakitResponses, penyakitResponse)
	}

	c.JSON(http.StatusOK, gin.H{
		"data": penyakitResponses,
	})
}

func (h *penyakitHandler) CreatePenyakitHandler(c *gin.Context) {
	var penyakitRequest penyakit.PenyakitRequest

	err := c.ShouldBindJSON(&penyakitRequest)
	if err != nil {
		errorMessages := []string{}
		for _, e := range err.(validator.ValidationErrors) {
			errorMessage := fmt.Sprintf("Error on field %s:, condition: %s", e.Field(), e.ActualTag())
			errorMessages = append(errorMessages, errorMessage)
		}
		c.JSON(http.StatusBadRequest, gin.H{
			"error": errorMessages,
			"status_code": http.StatusBadRequest,
		})
		return
	}

	isDNAValid := library.Sanitasi(penyakitRequest.DNAPenyakit)

	if !isDNAValid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "DNA Penyakit tidak valid",
			"status_code": http.StatusBadRequest,
		})
		return
	}


	penyakit, err := h.penyakitService.FindByName(penyakitRequest.NamaPenyakit)

	if len(penyakit.NamaPenyakit) != 0 {
		penyakitResponse := convertToPenyakitResponse(penyakit)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Nama penyakit sudah ada",
			"status_code": http.StatusBadRequest,
			"data": penyakitResponse,
		})
		return
	}

	penyakit, err = h.penyakitService.Create(penyakitRequest)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
			"status_code": http.StatusBadRequest,
		})
		return
	}

	penyakitResponse := convertToPenyakitResponse(penyakit)

	c.JSON(http.StatusOK, gin.H{
		"data": penyakitResponse,
		"status_code": http.StatusOK,
	})
}

func convertToPenyakitResponse(p penyakit.Penyakit) penyakit.PenyakitResponse {
	return penyakit.PenyakitResponse{
		NamaPenyakit:      p.NamaPenyakit,
		DNAPenyakit: p.DNAPenyakit,
	}
}
