package main

import (
   "fmt"
   "log"
   "os"
   "path/filepath"
   "ecolant/helpers"
)

const (
	targetIdLk = 0                                     // change to yours
    URLLk      = "https://<LK URL>/api/core"          // change to yours
    UsernameLk = "yours username LK"                 // change to yours
    PasswordLk = "yours password LK"                // change to yours 

	targetIdClouseContur = 0
	URLCloseContur      = "https://<CLOUSE CONTUR URL>/api/core"// change to yours
	UsernameCloseContur = "yours username CC"                  // change to yours
	PasswordCloseContur = "yours password CC"                 // change to yours
)

func main() {
	lkUrlMissions := fmt.Sprintf("%s/missions/?target=%d", URLLk, targetIdLk)
	clouseConturUrlMissions := fmt.Sprintf("%s/missions/?target=%d", URLCloseContur, targetIdClouseContur)

	lkMissions, _ := helpers.FetchMissions(lkUrlMissions, UsernameLk, PasswordLk)
	clouseConturMission, _ := helpers.FetchMissions(clouseConturUrlMissions, UsernameCloseContur, PasswordCloseContur)
	missionMap := helpers.MapMissions(clouseConturMission, lkMissions)
	arrayEntity := [3]string{"panoramas", "photos", "videos"}
    for _, value := range arrayEntity {
		for lkId, closeConturId := range missionMap {
			currentUrl := fmt.Sprintf("%s/%s?mission=%d", URLLk, value, lkId)
			
			result, _, _ := helpers.FetchDataAndProcess(currentUrl, UsernameLk, PasswordLk)
			downloadDir := fmt.Sprintf("%s/%s", "downloads", value)
			if err := os.MkdirAll(downloadDir, os.ModePerm); err != nil {
				log.Fatalf("Ошибка создания директории: %v", err)
			}
			for _, item := range result.Results {
					fileName, _ := helpers.ExtractFileNameFromURL(item.File)
					filePath := filepath.Join(downloadDir, fileName)
					helpers.DownloadFile(item.File, filePath)
					uploadUrl := fmt.Sprintf("%s/%s/", URLCloseContur, value)
					err := helpers.UploadFile(uploadUrl, closeConturId, item.DisplayName, filePath, UsernameCloseContur, PasswordCloseContur)
					if err != nil {
						log.Printf("Ошибка загрузки файла на Close Contur: %v", err)
					}
			}
		}
    }

    fmt.Println("Все файлы успешно сохранены.")
}
