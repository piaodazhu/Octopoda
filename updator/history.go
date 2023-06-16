package main

import "github.com/gin-gonic/gin"

type UpdatorHistory struct {
	Items []UpdatorHistoryItem
}

type UpdatorHistoryItem struct {
	Timestamp uint64
	Message string
}

func GetHistoryHandler(ctx *gin.Context) {

}