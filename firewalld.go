package xgfw_ctl

import (
	"errors"
	"fmt"
	"log"
	"strings"
)

const (
	public   = "public"   // 公共区域
	trusted  = "trusted"  // 信任区域
	external = "external" // 外部区域
	home     = "home"     // 家庭区域
	internal = "internal" // 内部区域
	work     = "work"     // 工作区域
	dmz      = "dmz"      // 隔离区域
	block    = "block"    // 限制区域
	drop     = "drop"     // 丢弃区域

	// drop = "drop"
	reject = "reject"

	ipv4 = "ipv4"
	ipv6 = "ipv6"
)

type Option struct {
	opType string

	Permanent    string
	Zone         string
	Family       string
	SourceAddr   string
	PortProtocol string
	Service      string
	Operate      string
}

type Firewalld struct {
	*Runner
	addArgs []string
	Option
}

func (fr *Firewalld) InsertArgs() []string {
	return fr.addArgs
}

type FirewallOption func(*Option)

func WithPermanent() FirewallOption {
	return func(o *Option) {
		o.Permanent = "--permanent"
	}
}

func WithZone(zone string) FirewallOption {
	return func(o *Option) {
		if isExitZone(zone) {
			o.Zone = fmt.Sprintf(`--zone=%s`, zone)
		} else {
			o.Zone = fmt.Sprintf(`--zone=%s`, public)
		}
	}
}

func WithFamily(ipType string) FirewallOption {
	return func(o *Option) {
		var family = ""
		switch ipType {
		case ipv4:
			family = ipv4
		case ipv6:
			family = ipv6
		default:
			family = ipv4
		}
		o.Family = fmt.Sprintf("rule family=%s", family)
	}
}

func WithService(service string) FirewallOption {
	return func(o *Option) {
		o.Service = fmt.Sprintf("service name=%s", service)
	}
}

func WithPort(protoc, port string) FirewallOption {
	return func(o *Option) {
		o.PortProtocol = fmt.Sprintf("port protocol=%s port=%s", protoc, port)
	}
}

func WithSrcAddr(srcAddr string) FirewallOption {
	return func(o *Option) {
		o.PortProtocol = fmt.Sprintf("source address=%s", srcAddr)
	}
}

func WithReject() FirewallOption {
	return func(o *Option) {
		o.Operate = reject
	}
}

func WithDrop() FirewallOption {
	return func(o *Option) {
		o.Operate = drop
	}
}

func ToInert() FirewallOption {
	return func(o *Option) {
		o.opType = "--add-rich-rule="
	}
}

func ToDelete() FirewallOption {
	return func(o *Option) {
		o.opType = "--remove-rich-rule="
	}
}

func NewFirewalld(opts ...FirewallOption) (*Firewalld, error) {
	o := &Option{}
	for _, opt := range opts {
		opt(o)
	}
	if o.Family == "" {
		log.Println("Warning: --family is required")
	}
	var args []string
	arg := strings.Join(append([]string{o.Family, o.SourceAddr, o.PortProtocol, o.Service, o.Operate}), " ")
	words := strings.Fields(arg)

	if o.opType == "" {
		return nil, errors.New("--opType is required")
	}
	args = append(args, o.Permanent, o.Zone, fmt.Sprintf("%s%s", o.opType, strings.Join(words, " ")))

	firewallIns := &Firewalld{
		Runner:  NewRunner(args),
		addArgs: args,
	}

	return firewallIns, nil
}

func DeleteArgsWithInert(opts []string) ([]string, error) {
	var args = opts
	var flag = false
	for index, option := range opts {
		if strings.Contains(option, "--add-rich-rule=") {
			args[index] = strings.ReplaceAll(option, "--add-rich-rule=", "--remove-rich-rule=")
			flag = true
		}
	}
	if !flag {
		return []string{}, errors.New("option flag is required")
	}
	return args, nil
}

func isExitZone(zone string) bool {
	var zones = []string{public, trusted, external, home, internal, work, dmz, block, drop}
	for _, z := range zones {
		if z == zone {
			return true
		}
	}
	return false
}
