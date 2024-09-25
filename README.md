# xgfw-ctl

## Description

firewalld and iptables controller of golang

+ firewalld
+ iptables

## exp
### firewalld
```go

	fire := firewalld.BuildOptionOfAdd(
		firewalld.WithPermanent(),
		firewalld.WithFamily("ipv4"),
		firewalld.WithZone("public"),
		firewalld.WithPort("tcp", "23"),
		firewalld.WithReject(),
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

	delArgs, err := firewalld.OptionOfRemoveWithAdd(fire.GetAddArgs())
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
```

### iptables
```go
	testCases := []struct {
		name string
		in   string
		out  string
	}{
		{
			"legacy output",
			"-A foo1 -p tcp -m tcp --dport 1337 -j ACCEPT",
			"-A foo1 -p tcp -m tcp --dport 1337 -j ACCEPT",
		},
		{
			"nft output",
			"[99:42] -A foo1 -p tcp -m tcp --dport 1337 -j ACCEPT",
			"-A foo1 -p tcp -m tcp --dport 1337 -j ACCEPT -c 99 42",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			actual := filterRuleOutput(tt.in)
			if actual != tt.out {
				t.Fatalf("expect %s actual %s", tt.out, actual)
			}
		})
	}
```
