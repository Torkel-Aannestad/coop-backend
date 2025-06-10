package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Message struct {
	ExternalId string `json:"external_id"`
	Author     string `json:"author"`
	Body       string `json:"body"`
	Platform   string `json:"platform"`
}

type MastodonPost struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Account MastodonAccount
}

type MastodonAccount struct {
	Username string `json:"username"`
}

func (app *application) getMastodonPosts() ([]MastodonPost, error) {
	url := "https://mastodon.social/api/v1/timelines/public?limit=10"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	//todo validate resp struct

	var posts []MastodonPost
	err = json.NewDecoder(resp.Body).Decode(&posts)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (app *application) postMessageToAPI(msg Message) error {
	postBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	//retries
	for i := 1; i < 3; i++ {
		resp, err := http.Post("http://service.api:4000/v1/messages", "application/json", bytes.NewBuffer(postBody))

		if err == nil {
			if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
				app.logger.Error("error posting data to API", "statuscode", resp.StatusCode, "response body", resp.Body)
				return err
			}

			//good case
			return nil
		}
		app.logger.Debug("failed sending message to api", "attemp", i)
		time.Sleep(2 * time.Second)
	}

	return err

}

func (app *application) mastodonBackgroundJob() {
	app.logger.Info("starting mastodonBackgroundJob")

	app.backgroundJob(func() {
		for {
			time.Sleep(app.config.pullingFrequency)

			app.logger.Info("getting Mastodon posts")
			posts, err := app.getMastodonPosts()
			if err != nil {
				app.logger.Error(err.Error())
				continue
			}

			app.mu.Lock()

			//if not seen then send to api
			for _, post := range posts {
				if _, found := app.seenPosts[post.ID]; found {
					app.logger.Info("found skipping", "post.ID", post.ID) //remove
					continue
				}

				message := Message{
					ExternalId: post.ID,
					Author:     post.Account.Username,
					Body:       post.Content,
					Platform:   "mastodon",
				}

				app.logger.Info("new post", "post.ID", post.ID)
				err := app.postMessageToAPI(message)
				if err != nil {
					app.logger.Error("error posting data to API", "error", err.Error())
				}

				app.seenPosts[post.ID] = struct{}{}

			}

			app.mu.Unlock()
		}
	})

	//Todo: cleanup postcache when data becomes stale.
	//need timestamp on seenPosts map value

}

func (app *application) ChecKDependentService() error {
	url := "http://service.api:4000/v1/healthcheck"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	app.logger.Info("api healthcheck", "response", string(body))
	return nil

}
