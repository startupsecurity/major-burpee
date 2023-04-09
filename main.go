package main

import (
	"fmt"
	"github.com/slack-go/slack"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	token       = os.Getenv("SLACK_BOT_TOKEN")
	channelName = os.Getenv("SLACK_CHANNEL_NAME")
	counterFile = "counter.txt"
)

func main() {
	api := slack.New(token)

	counter, lastUpdated, err := readCounterFromFile(counterFile)

	if err != nil {
		log.Fatalf("Error reading counter file: %v", err)
	}

	now := time.Now()
	if !isSameDay(now, lastUpdated) {
		counter++
		err = writeCounterToFile(counterFile, counter, now)
		if err != nil {
			log.Fatalf("Error writing counter file: %v", err)
		}
	}

	channelID, err := getChannelID(api, channelName)
	if err != nil {
		log.Fatalf("Error getting channel ID: %v", err)
	}

	_, timestamp, err := api.PostMessage(channelID, slack.MsgOptionText(fmt.Sprintf("BURPEE TIME! Drop and give me #%d", counter), false))
	if err != nil {
		log.Fatalf("Error posting message: %v", err)
	}

	fmt.Printf("Message successfully sent to channel %s at %s", channelName, timestamp)
}

func readCounterFromFile(filename string) (int, time.Time, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		if err != nil {
			if os.IsNotExist(err) {
				counter := 0
				lastUpdated := time.Now()
				err = writeCounterToFile(filename, counter, lastUpdated)
				if err != nil {
					return 0, time.Time{}, err
				}
				return counter, lastUpdated, nil
			}
			return 0, time.Time{}, err
		}
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	counter, err := strconv.Atoi(lines[0])
	if err != nil {
		return 0, time.Time{}, err
	}

	lastUpdated, err := time.Parse(time.RFC3339, lines[1])
	if err != nil {
		return 0, time.Time{}, err
	}

	return counter, lastUpdated, nil
}

func writeCounterToFile(filename string, counter int, lastUpdated time.Time) error {
	data := fmt.Sprintf("%d\n%s", counter, lastUpdated.Format(time.RFC3339))
	return ioutil.WriteFile(filename, []byte(data), 0644)
}

func getChannelID(api *slack.Client, channelName string) (string, error) {
	channels, nextCursor, err := api.GetConversations(&slack.GetConversationsParameters{
		Types:           []string{"public_channel"},
		ExcludeArchived: true,
		Limit:           1000,
	})

	fmt.Printf("%s\n", nextCursor)

	if err != nil {
		return "", err
	}

	for _, channel := range channels {
		fmt.Printf("Channel: %s\n", channel.Name)
		if channel.Name == channelName {
			return channel.ID, nil
		}
	}
	return "", fmt.Errorf("Channel %s not found", channelName)
}

func isSameDay(t1, t2 time.Time) bool {
	y1, m1, d1 := t1.Date()
	y2, m2, d2 := t2.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}
