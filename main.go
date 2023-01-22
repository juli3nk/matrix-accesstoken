package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/davecgh/go-spew/spew"
	"github.com/juli3nk/go-utils/readinput"
)

type Request struct {
	Type     string `json:"type"`
	User     string `json:"user"`
	Password string `json:"password"`
}

type Response struct {
	AccessToken string `json:"access_token"`
	DeviceID string `json:"device_id"`
	HomeServer string `json:"home_server"`
	UserID string `json:"user_id"`
	WellKnown map[string]map[string]string `json:"well_known"`
}

var (
	flagServerName = flag.String("server-name", "", "Server name")
	flagUsername = flag.String("username", "", "Username")
	flagPassword = flag.String("password", "", "Password")
)

func main() {
	flag.Parse()

	matrixServerName := *flagServerName
	if len(matrixServerName) == 0 {
		for {
			matrixServerName = readinput.ReadInput("Server name")
			if len(matrixServerName)  > 0 {
				break
			}
		}
	}

	matrixUsername := *flagUsername
	if len(matrixUsername) == 0 {
		for {
			matrixUsername = readinput.ReadInput("Username")
			if len(matrixServerName)  > 0 {
				break
			}
		}
	}

	matrixPassword := *flagPassword
	if len(matrixPassword) == 0 {
		for {
			matrixPassword = readinput.ReadPassword("Password")
			if len(matrixPassword)  > 0 {
				break
			}
		}
	}

	homeServerUrl := fmt.Sprintf("https://%s/_matrix/client/r0/login", matrixServerName)

	response, err := getUserLoginInfo(homeServerUrl, matrixUsername, matrixPassword)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n")
	spew.Dump(response)
}

func getUserLoginInfo(url, user, password string) (*Response, error) {
	payload := Request{
		Type:     "m.login.password",
		User:     user,
		Password: password,
	}

	j, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	result := new(Response)

	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}

	return result, nil
}
