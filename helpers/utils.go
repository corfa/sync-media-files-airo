package helpers
import (
    "bytes"
    "fmt"
	"time"
    "io"
    "log"
    "mime/multipart"
    "net/http"
    "os"
	"encoding/base64"
    "encoding/json"
	"strconv"
	"net/url"
	"path"
)

func DownloadFile(url string, filepath string) {
    resp, err := http.Get(url)
    if err != nil {
        log.Printf("Ошибка загрузки файла %s: %v", url, err)
        return
    }
    defer resp.Body.Close()

    outFile, err := os.Create(filepath)
    if err != nil {
        log.Printf("Ошибка создания файла %s: %v", filepath, err)
        return
    }
    defer outFile.Close()

    _, err = io.Copy(outFile, resp.Body)
    if err != nil {
        log.Printf("Ошибка сохранения файла %s: %v", filepath, err)
        return
    }

    fmt.Printf("Файл %s успешно загружен.\n", filepath)
}


func UploadFile(url string, mission int, displayName,  filepath string, username string, password string) error {
	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
    file, err := os.Open(filepath)
    if err != nil {
        return err
    }
    defer file.Close()

    var requestBody bytes.Buffer
    writer := multipart.NewWriter(&requestBody)
	writer.WriteField("mission", strconv.Itoa(mission))
    writer.WriteField("display_name", displayName)

    part, err := writer.CreateFormFile("file", filepath)
    if err != nil {
        return err
    }
    _, err = io.Copy(part, file)
    if err != nil {
        return err
    }

    writer.Close()

    req, err := http.NewRequest("POST", url, &requestBody)
    if err != nil {
        return err
    }
    req.Header.Set("Content-Type", writer.FormDataContentType())
    req.Header.Add("Authorization", "Basic "+auth)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    if resp.StatusCode != http.StatusCreated {
        return fmt.Errorf("ошибка загрузки файла: %s, ответ сервера: %s", resp.Status, string(responseBody))
    }

    log.Printf("Файл %s успешно загружен на %s. Ответ сервера: %s\n", filepath, url, string(responseBody))
    return nil
}


func FetchDataAndProcess(url string, username string, passowrd string) (Response, []byte, error) {
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return Response{}, nil, fmt.Errorf("Ошибка создания запроса: %v", err)
    }

    auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + passowrd))
    req.Header.Add("Authorization", "Basic "+auth)

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return Response{}, nil, fmt.Errorf("Ошибка отправки запроса: %v", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return Response{}, nil, fmt.Errorf("Ошибка чтения ответа: %v", err)
    }

    if resp.StatusCode != http.StatusOK {
        return Response{}, nil, fmt.Errorf("Неправильный статус ответа: %d %s", resp.StatusCode, resp.Status)
    }

    var result Response
    if err := json.Unmarshal(body, &result); err != nil {
        return Response{}, nil, fmt.Errorf("Ошибка декодирования JSON: %v", err)
    }

    return result, body, nil
}


func MapMissions(closeConturMissions, lkMissions []Mission) map[int]int {
    missionMap := make(map[int]int)

    for _, lkMission := range lkMissions {
        lkDate, err := time.Parse(time.RFC3339, lkMission.Date)
        if err != nil {
            log.Printf("Error parsing date for mission %d: %v", lkMission.Id, err)
            continue
        }

        for _, closeConturMission := range closeConturMissions {
            closeConturDate, err := time.Parse(time.RFC3339, closeConturMission.Date)
            if err != nil {
                log.Printf("Error parsing date for mission %d: %v", closeConturMission.Id, err)
                continue
            }

            if closeConturDate.Equal(lkDate) {
                missionMap[lkMission.Id] = closeConturMission.Id
                break
            }
        }
    }

    return missionMap
}

func FetchMissions(url, username, password string) ([]Mission, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error fetching missions: %s", resp.Status)
	}

	var missionsResponse MissionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&missionsResponse); err != nil {
		return nil, err
	}

	return missionsResponse.Results, nil
}

func ExtractFileNameFromURL(fileURL string) (string, error) {
	parsedURL, err := url.Parse(fileURL)
	if err != nil {
		return "", err
	}

	fileName := path.Base(parsedURL.Path)
	return fileName, nil
}