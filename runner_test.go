package xgfwlib

import (
	"log"
	"testing"
)

func TestFirewalld(t *testing.T) {
	fire, err := NewFirewalld(
		WithPermanent(),
		WithFamily("ipv4"),
		WithPort("tcp", "22"),
		WithZone("public"),
		WithReject(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := fire.Exec()
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(out)
}

func TestFirewalldAll(t *testing.T) {

	fire, err := NewFirewalld(
		WithPermanent(),
		WithFamily("ipv4"),
		WithZone("public"),
		WithPort("tcp", "23"),
		WithReject(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	out, err := fire.Exec()
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("add", out)

	out, err = fire.ExecArgs([]string{"--list-rich-rules"})
	if err == nil {
		log.Printf("rule list\n%+v", out)
	}

	delArgs, err := DeleteArgsWithInert(fire.InsertArgs())
	if err != nil {
		return
	}

	out, err = fire.ExecArgs(delArgs)
	if err != nil {
		t.Error(err)
		return
	}
	log.Println("remove", out)

	out, err = fire.ExecArgs([]string{"--list-rich-rules"})
	if err == nil {
		log.Printf("rule list\n%+v", out)
	}
}
