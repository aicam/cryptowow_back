package server

import (
	"gorm.io/gorm"
	"strconv"
	"strings"
)

func ParseEquippedCache(DB *gorm.DB, equippedCache string) map[int]string {
	splitted := strings.Split(equippedCache, " ")
	var equippedItemsIds []int
	for i := 0; i < len(splitted); i += 2 {
		IDstring, _ := strconv.Atoi(splitted[i])
		equippedItemsIds = append(equippedItemsIds, IDstring)
	}
	for i, itemID := range equippedItemsIds {

	}
}
