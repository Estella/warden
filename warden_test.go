package main

import "testing"
import "time"
import "github.com/fzzy/radix/redis"

func cleanup() {
        conn, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10) * time.Second)
        if err != nil {
                panic("Error connecting to Redis")
        }
	conn.Cmd("FLUSHDB")
	conn.Close()
}

func TestSetIPWhitelist(t *testing.T) {
        ip := command{function: "whitelist", actor_type: "ip", actor: "1.1.1.1", reason: "test"}

        set(ip)
        s := status(ip.actor)
        if s != "Whitelisted\n" {
                t.Error("Whitelist ip failed, got", s)
        }

	cleanup()
}

func TestSetCIDRWhitelist(t *testing.T) {
        cidr := command{function: "whitelist", actor_type: "cidr", actor: "1.1.1.0/24", reason: "test"}

        set(cidr)
        s := status(cidr.actor)
        if s != "Whitelisted\n" {
                t.Error("Whitelist cidr failed, got", s)
        }

	cleanup()
}

func TestSetUserWhitelist(t *testing.T) {
        user := command{function: "whitelist", actor_type: "user", actor: "gotest", reason: "test"}

        set(user)
	s := status(user.actor)
        if s != "Whitelisted\n" {
                t.Error("Whitelist user failed, got", s)
        }

	cleanup()
}

func TestSetIPBlacklist(t *testing.T) {
        ip := command{function: "blacklist", actor_type: "ip", actor: "1.1.1.1", reason: "test"}

        set(ip)
        s := status(ip.actor)
        if s != "Blacklisted\n" {
                t.Error("Blacklist ip failed, got", s)
        }

	cleanup()
}

func TestSetCIDRBlacklist(t *testing.T) {
        cidr := command{function: "blacklist", actor_type: "cidr", actor: "1.1.1.0/24", reason: "test"}

        set(cidr)
        s := status(cidr.actor)
        if s != "Blacklisted\n" {
                t.Error("Blacklist cidr failed, got", s)
        }

	cleanup()
}

func TestSetUserBlacklist(t *testing.T) {
        user := command{function: "blacklist", actor_type: "user", actor: "gotest", reason: "test"}

        set(user)
	s := status(user.actor)
        if s != "Blacklisted\n" {
                t.Error("Blacklist user failed, got", s)
        }

	cleanup()
}
