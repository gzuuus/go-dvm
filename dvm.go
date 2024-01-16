package main

import (
	"context"
	"fmt"

	goNostr "github.com/nbd-wtf/go-nostr"
	"github.com/sebdeveloper6952/go-dvm/domain"
	"github.com/sebdeveloper6952/go-dvm/nostr"
)

type malwareDVM struct {
	sk string
	pk string
}

type res struct {
	Result string
}

func NewMalwareDvm(sk string) (domain.Dvmer, error) {
	dvm := &malwareDVM{}

	return dvm, dvm.SetSk(sk)
}

func (d *malwareDVM) SetSk(sk string) error {
	d.sk = sk
	pk, err := goNostr.GetPublicKey(d.sk)
	if err != nil {
		return err
	}
	d.pk = pk

	return nil
}

func (d *malwareDVM) Pk() string {
	return d.pk
}

func (d *malwareDVM) Sign(e *goNostr.Event) error {
	return e.Sign(d.sk)
}

func (d *malwareDVM) Profile() *nostr.ProfileMetadata {
	return &nostr.ProfileMetadata{
		Name:    "Test DVM",
		About:   "Test DVM about",
		Picture: "https://iconape.com/wp-content/png_logo_vector/virus-2.png",
	}
}

func (d *malwareDVM) KindSupported() int {
	return nostr.KindReqMalwareScan
}

func (d *malwareDVM) AcceptJob(input *nostr.Nip90Input) bool {
	return true
}

func (d *malwareDVM) Run(ctx context.Context, Input *nostr.Nip90Input) (chan *domain.JobUpdate, chan *domain.JobUpdate, chan error) {
	chanToDvm := make(chan *domain.JobUpdate)
	chanToEngine := make(chan *domain.JobUpdate)
	chanErr := make(chan error)

	go func() {
		defer func() {
			close(chanToDvm)
			close(chanToEngine)
			close(chanErr)
		}()

		chanToEngine <- &domain.JobUpdate{
			Status: domain.StatusProcessing,
		}

		chanToEngine <- &domain.JobUpdate{
			Status: domain.StatusSuccess,
			Result: fmt.Sprintf("hello human %s", Input.Input),
		}
	}()

	return chanToDvm, chanToEngine, chanErr
}

type DvmInput struct {
	URL string
}
