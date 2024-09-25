package xgfw_ctl

import "testing"

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
