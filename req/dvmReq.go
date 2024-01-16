package main

import (
	"context"
	"log"

	"github.com/nbd-wtf/go-nostr"
	"github.com/sirupsen/logrus"
)

func main() {
	var (
		ctx, cancel = context.WithCancel(context.Background())
	)
	defer cancel()

	sk := "c292c69bbbb9289863005e1beb111807d10cb11418dfa18a5ec19516697c337f"
	pk, _ := nostr.GetPublicKey(sk)

	ev := nostr.Event{
		PubKey:    pk,
		CreatedAt: nostr.Now(),
		Kind:      5500,
		Tags: nostr.Tags{
			{"i", "Hellnostr", "text"},
			//https://raw.githubusercontent.com/ca110us/go-clamav/main/example/test_file/nmap
		},
	}
	ev.Sign(sk)

	relay, err := nostr.RelayConnect(ctx, "wss://nostr-pub.wellorder.net")
	if err != nil {
		log.Fatal(err)
	}

	relay.Publish(ctx, ev)
	Sub, err := relay.Subscribe(
		context.Background(),
		nostr.Filters{
			{
				Kinds: []int{7000, 6500},
				Limit: 10,
				Tags:  nostr.TagMap{"p": []string{pk}},
			},
		},
	)
	for event := range Sub.Events {
		logrus.Info(event)
	}

	relay.Publish(ctx, ev)

	logrus.Info("event sent")
}
