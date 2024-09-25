package xgfw_ctl

import (
	"log"
	"testing"
)

func TestOption(t *testing.T) {
	fire := BuildOptionOfAdd(
		WithPermanent(),
		WithFamily("ipv4"),
		WithPort("tcp", "22"),
		WithSrcAddr("193.168.1.1"),
		WithZone("public"),
		WithReject(),
	)
	t.Logf("firewalld option is : %s", fire.GetAddArgs())
}

func TestAddOption(t *testing.T) {
	fire := BuildOptionOfAdd(
		WithPermanent(),
		WithPort("tcp", "22"),
		WithZone("public"),
		WithReject(),
	)
	t.Logf("firewalld option is : %s", fire.GetAddArgs())
}

func TestAddAndRemoveOption(t *testing.T) {
	fire := BuildOptionOfAdd(
		WithPermanent(),
		WithFamily("ipv4"),
		WithPort("tcp", "22"),
		WithZone("public"),
		WithReject(),
	)
	t.Logf("add firewalld option is : %s", fire.GetAddArgs())

	fire = BuildOptionOfRemove(
		WithPermanent(),
		WithFamily("ipv4"),
		WithPort("tcp", "22"),
		WithZone("public"),
		WithReject(),
	)

	t.Logf("del firewalld option is : %s", fire.GetAddArgs())
}

func TestAddAndRemove(t *testing.T) {

	fire := BuildOptionOfAdd(
		WithPermanent(),
		WithFamily("ipv4"),
		WithZone("public"),
		WithPort("tcp", "23"),
		WithReject(),
	)
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

	delArgs, err := OptionOfRemoveWithAdd(fire.GetAddArgs())
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
