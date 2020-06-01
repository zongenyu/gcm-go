package main

import (
	"bitbucket.org/raylios/cloudpost-go/slog"
	"github.com/alexjlockwood/gcm"
	"encoding/json"
	"fmt"
	"strings"
)

// Define Key setting JSON  =========================
//{
//	"gcmKeys":[{"bundle":"com.example.app.android", "apiKey":"aiksidkaqidic"}, ...]
//}

var GcmKeys map[string]interface{}

func send(token, message string) error {

	var key, bundle string

	spls := strings.Split(token, ",")
	if len(spls) != 2 {
		slog.Warning("split token fail: %v", token)
		key = token
		bundle = "com.FFWDcam.cloudservice"
	} else {
		key = spls[0]
		bundle = spls[1]
	}


	apiKey := getGcmKey(bundle)
	if apiKey == "" || len(apiKey) == 0 {
		slog.Err("Error while getting apiKey: %v from GcmKeys: %v", bundle, GcmKeys)
		return nil
	}

	sender := &gcm.Sender{ApiKey: apiKey}

	// Create the message to be sent
	regIds := []string{key}
	data := map[string]interface{}{"message": message}
	msg := gcm.NewMessage(data, regIds...)

	// Send the message and receive the response after at most two retries.
	response, err := sender.Send(msg, 2)
	if err != nil {
		fmt.Println("Failed to send message: " + err.Error())
		return err
	}

	return nil
}

func getGcmKey(bundle string) string {

	var jsonObj map[string]interface{}
	err := json.Unmarshal(GCM_DATA, &jsonObj)
	if err != nil {
		slog.Err("failed parsing json in parseGcmKey with error: %v", err)
		return ""
	}

	GcmKeys, isOK := jsonObj["gcmKeys"].(map[string]interface{})
	if !isOK {
		slog.Err("Failed to assert gcmKeys %v to map!", GcmKeys)
		return ""
	}

	//slog.Debug("gcm keys: %v", GcmKeys)

	key, isOK := GcmKeys[bundle].(string)
	if !isOK || len(key) == 0 {
		slog.Err("Failed to assert apiKey for bundle %v in GcmKeys %v", bundle, GcmKeys)
		return ""
	}
	return key
}
