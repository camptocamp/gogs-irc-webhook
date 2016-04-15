package main

import (
  "encoding/json"
  "github.com/thoj/go-ircevent"
  "github.com/caarlos0/env"
  "os"
  "fmt"
  "strings"
)

type Author struct {
  Name string
}

type Commit struct {
  Message string
  Url string
  Author Author
}

type Repository struct {
  Name string
}

type Pusher struct {
  Name string
}

type GogsPayload struct {
  Secret string `json:"secret"`
  After string `json:"after"`
  Ref string `json:"ref"`
  Commits []Commit `json:"commits"`
  Repository Repository `json:"repository"`
  Pusher Pusher `json:"pusher"`
}

type Config struct {
  Server           string   `env:"IRC_SERVER"`
  Port             string   `env:"IRC_PORT" envDefault:"6667"`
  Room             string   `env:"IRC_ROOM"`
  Nick             string   `env:"IRC_NICK"`
  // Branches         []string `env:"IRC_BRANCHES"` // TODO
  // NickservPassword string   `env:"IRC_NICKSERV_PASSWORD"` // TODO
  // Password         string   `env:"IRC_PASSWORD"` // TODO
  // Ssl              bool     `env:"IRC_SSL" envDefault:"true"` // TODO
  // Join             bool     `env:"IRC_JOIN" envDefault:"true"` // TODO
  Colors           bool     `env:"IRC_COLORS" envDefault:"true"`
  // LongUrl          bool     `env:"IRC_LONG_URL" envDefault:"true"` // TODO
  // Notice           bool     `env:"IRC_NOTICE" envDefault:"false"` // TODO
}

func main() {
  cfg := Config{}
  env.Parse(&cfg)

  payload := GogsPayload{}
  json.Unmarshal([]byte(os.Args[1]), &payload)

  branch := strings.TrimPrefix(payload.Ref, "refs/heads/")

  con := irc.IRC(cfg.Nick, cfg.Nick)
  err := con.Connect(cfg.Server+":"+cfg.Port)
  if err != nil {
    fmt.Println("Failed connecting")
    return
  }
  con.AddCallback("001", func (e *irc.Event) {
    con.Join(cfg.Room)
    str1 := "[%v] %v pushed %v new commit to %v: %v"
    str2 := "%v/%v %v %v: %v"
    if cfg.Colors {
      str1 = "\x0F[\x0313%v\x03] \x0315%v\x03 pushed \x02%v\x02 new commit to \x0306%v\x03: \x1F\x0302%v\x03\x1F"
      str2 = "\x0313%v\x03/\x0306%v\x0306 \x0314%v\x03 \x0315%v\x03: %v"
    }
    con.Privmsgf(cfg.Room, str1, payload.Repository.Name, payload.Pusher.Name, len(payload.Commits), branch, payload.Commits[0].Url)
    con.Privmsgf(cfg.Room, str2, payload.Repository.Name, branch, payload.After[0:7], payload.Commits[0].Author.Name, payload.Commits[0].Message)
    con.Quit()
  })

  con.Loop()
}
