package main

import "testing"

func TestShouldDeliverRequiresFlagAndReadySubscriber(t *testing.T) {
	asset := Asset{ID: "lesson-17", Processed: true}
	subscriber := Subscriber{ID: "sub-42", Active: true}
	if !shouldDeliver(true, asset, subscriber) {
		t.Fatal("expected processed asset for active subscriber")
	}
	if shouldDeliver(false, asset, subscriber) {
		t.Fatal("expected flag-off delivery to be suppressed")
	}
	if shouldDeliver(true, Asset{ID: "lesson-17"}, subscriber) {
		t.Fatal("expected unprocessed asset to be suppressed")
	}
}
