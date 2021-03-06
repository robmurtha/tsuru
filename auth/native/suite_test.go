// Copyright 2015 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package native

import (
	"testing"

	"github.com/tsuru/config"
	"github.com/tsuru/tsuru/auth"
	"github.com/tsuru/tsuru/auth/authtest"
	"github.com/tsuru/tsuru/db"
	"github.com/tsuru/tsuru/db/dbtest"
	"golang.org/x/crypto/bcrypt"
	"launchpad.net/gocheck"
)

func Test(t *testing.T) { gocheck.TestingT(t) }

type S struct {
	conn   *db.Storage
	hashed string
	user   *auth.User
	team   *auth.Team
	server *authtest.SMTPServer
	token  auth.Token
}

var _ = gocheck.Suite(&S{})

var nativeScheme = NativeScheme{}

func (s *S) SetUpSuite(c *gocheck.C) {
	config.Set("auth:token-expire-days", 2)
	config.Set("auth:hash-cost", bcrypt.MinCost)
	config.Set("admin-team", "admin")
	config.Set("database:url", "127.0.0.1:27017")
	config.Set("database:name", "tsuru_auth_native_test")
	var err error
	s.server, err = authtest.NewSMTPServer()
	c.Assert(err, gocheck.IsNil)
	config.Set("smtp:server", s.server.Addr())
	config.Set("smtp:user", "root")
	config.Set("smtp:password", "123456")
}

func (s *S) SetUpTest(c *gocheck.C) {
	s.conn, _ = db.Conn()
	s.user = &auth.User{Email: "timeredbull@globo.com", Password: "123456"}
	_, err := nativeScheme.Create(s.user)
	c.Assert(err, gocheck.IsNil)
	s.token, err = nativeScheme.Login(map[string]string{"email": s.user.Email, "password": "123456"})
	c.Assert(err, gocheck.IsNil)
	team := &auth.Team{Name: "cobrateam", Users: []string{s.user.Email}}
	err = s.conn.Teams().Insert(team)
	c.Assert(err, gocheck.IsNil)
	s.team = team
}

func (s *S) TearDownTest(c *gocheck.C) {
	err := dbtest.ClearAllCollections(s.conn.Users().Database)
	c.Assert(err, gocheck.IsNil)
	s.conn.Close()
	cost = 0
	tokenExpire = 0
}

func (s *S) TearDownSuite(c *gocheck.C) {
	s.server.Stop()
}
