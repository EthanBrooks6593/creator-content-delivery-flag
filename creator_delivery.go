package main

import (
	"errors"
	"fmt"
)

type Asset struct {
	ID        string
	Processed bool
}
type Subscriber struct {
	ID     string
	Active bool
}

func shouldDeliver(flagOn bool, asset Asset, subscriber Subscriber) bool {
	return flagOn && asset.Processed && subscriber.Active
}

func deliver(client *InfraiClient, creator string, asset Asset, subscriber Subscriber) error {
	flagOn, err := client.CreatorFlag("creator-commerce-" + creator)
	if err != nil {
		var statusErr *httpStatusError
		if errors.As(err, &statusErr) && statusErr.status == 404 {
			fmt.Printf("delivery skipped for asset %s and subscriber %s: creator flag is not configured\n", asset.ID, subscriber.ID)
			return nil
		}
		return err
	}
	if !shouldDeliver(flagOn, asset, subscriber) {
		fmt.Printf("delivery skipped for asset %s and subscriber %s\n", asset.ID, subscriber.ID)
		return nil
	}
	fmt.Printf("delivery queued for asset %s and subscriber %s\n", asset.ID, subscriber.ID)
	return nil
}

func main() {
	client, err := NewInfraiClient()
	if err != nil {
		panic(err)
	}
	err = deliver(client, "maya", Asset{ID: "lesson-17", Processed: true}, Subscriber{ID: "sub-42", Active: true})
	if err != nil {
		panic(err)
	}
}
