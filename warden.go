package main

import "fmt"
import "net"
import "log"
import "strings"
import "errors"
import "time"
import "github.com/fzzy/radix/redis"

type command struct {
	function   string
	actor_type string
	actor      string
	reason     string
}

func set(c command) string {
        conn, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10) * time.Second)
        if err != nil {
                panic("Error connecting to Redis")
        }

        q := fmt.Sprintf("%s:repsheet:%s:%s", strings.TrimSpace(c.actor), c.actor_type, c.function)
        r := conn.Cmd("SET", q, strings.TrimSpace(c.reason))
        response, e := r.Str()
        if e != nil {
                fmt.Println("error: ", e)
        }

        conn.Close()

        return response + "\n"
}

func status(actor string) string {
        conn, err := redis.DialTimeout("tcp", "127.0.0.1:6379", time.Duration(10) * time.Second)
        if err != nil {
                panic("Error connecting to Redis")
        }

	q := fmt.Sprintf("%s*", strings.TrimSpace(actor))
	r := conn.Cmd("KEYS", q)

	var whitelisted bool
	var blacklisted bool

	p, _ := r.List()
	for _, a := range p {
		parts := strings.Split(a, ":")
		if parts[len(parts)-1] == "whitelist" {
			whitelisted = true
		}
		if parts[len(parts)-1] == "blacklist" {
			blacklisted = true
		}
	}

	var response string
	if whitelisted {
		response = "Whitelisted"
	} else if blacklisted {
		response = "Blacklisted"
	} else {
		response = "Not Found"
	}

	conn.Close()

	return response + "\n"
}

func dispatch(msg string) ([]byte, error) {
	if msg == "\r\n" {
		return []byte("\n"), nil
	}

        parts := strings.Split(msg, " ")
        if len(parts) == 0 {
                return []byte(""), errors.New("Invalid message")
        }

        var response []byte
        if parts[0] == "blacklist" {
                if len(parts) != 4 {
                        return []byte("Not enough arguments: blacklist <actor_type> <actor> <reason>\n"), nil
                }
		c := command{function: parts[0], actor_type: parts[1], actor: parts[2], reason: parts[3]}
                response = []byte(set(c))
        } else if parts[0] == "whitelist" {
                if len(parts) != 4 {
                        return []byte("Not enough arguments: whitelist <actor_type> <actor> <reason>\n"), nil
                }
		c := command{function: parts[0], actor_type: parts[1], actor: parts[2], reason: parts[3]}
                response = []byte(set(c))
	} else if parts[0] == "status" {
		if len(parts) != 2 {
			return []byte("Not enough arguments: status <actor>"), nil
		}
		response = []byte(status(parts[1]))
        } else {
                response = []byte("Unknown command")
        }

        return response, nil
}

func handleConnection(c net.Conn) {
        for {
                buf := make([]byte, 1024)
                nr, err := c.Read(buf)
                if err != nil {
                        return
                }

                message := string(buf[0:nr])
                response, resperr := dispatch(message)
                if resperr != nil {
                        return
                }
                _, err = c.Write(response)
                if err != nil {
                        log.Fatal("Write Error: ", err)
                }
        }
}

func main() {
        ln, err := net.Listen("tcp", ":3142")

        if err != nil {
                panic("Could not bind to socket")
        }

        for {
                conn, err := ln.Accept()

                if err != nil {
                        fmt.Println("Error receiving message")
                        continue
                }

                go handleConnection(conn)
        }
}
