package initializers

import "tess/models"

func SyncDatabase() {
	DB.AutoMigrate(&models.User{},&models.Admin{}) 
}
