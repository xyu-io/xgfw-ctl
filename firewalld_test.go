package xgfwlib

import (
	"log"
	"testing"
)

func TestOption(t *testing.T) {
	fire, err := NewFirewalld(
		WithPermanent(),
		WithFamily("ipv4"),
		WithPort("tcp", "22"),
		WithSrcAddr("193.168.1.1"),
		WithZone("public"),
		WithReject(),
		ToInert(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("firewalld option is : %s", fire.InsertArgs())
}

func TestAddOption(t *testing.T) {
	fire, err := NewFirewalld(
		WithPermanent(),
		WithPort("tcp", "22"),
		WithZone("public"),
		WithReject(),
		ToInert(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("firewalld option is : %s", fire.InsertArgs())
}

func TestAddAndRemoveOption(t *testing.T) {
	fire, err := NewFirewalld(
		WithPermanent(),
		WithFamily("ipv4"),
		WithPort("tcp", "22"),
		WithZone("public"),
		WithReject(),
		ToInert(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("add firewalld option is : %s", fire.InsertArgs())

	fire, err = NewFirewalld(
		WithPermanent(),
		WithFamily("ipv4"),
		WithPort("tcp", "22"),
		WithZone("public"),
		WithReject(),
		ToDelete(),
	)
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("del firewalld option is : %s", fire.InsertArgs())
}

func TestAddAndRemove(t *testing.T) {

	fire, err := NewFirewalld(
		WithPermanent(),
		WithFamily("ipv4"),
		WithZone("public"),
		WithPort("tcp", "23"),
		WithReject(),
		ToInert(),
	)
	if err != nil {
		t.Error(err)
		return
	}
	out, err := fire.Exec()
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
		return
	}
	log.Println("remove", out)

	out, err = fire.ExecArgs([]string{"--list-rich-rules"})
	if err == nil {
		log.Printf("rule list\n%+v", out)
	}
}
