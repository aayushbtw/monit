package main

import (
	"errors"
	"net"

	"github.com/aayushbtw/monit/ui"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/bubbletea"
	"github.com/charmbracelet/wish/logging"
)

const (
	host        = "localhost"
	port        = "23234"
	hostKeyPath = ".ssh/host_monit_ed25519"
)

var users = map[string]string{
	"Admin": "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIBL1CPXo/GUe+/wvxxpGrsdVxx8v+W8foAFo7QdM48hv Monit",
	// You can add add your name and public key here :)
}

func keyAuthMiddleware() wish.Middleware {
	return func(next ssh.Handler) ssh.Handler {
		return func(sess ssh.Session) {
			for name, pubkey := range users {
				parsed, _, _, _, err := ssh.ParseAuthorizedKey(
					[]byte(pubkey),
				)
				if err != nil {
					log.Error("Error Parsing user keys", "error", err)
				}

				if ssh.KeysEqual(sess.PublicKey(), parsed) {
					log.Info("Admin connected", "username", name)
					next(sess)
				}
			}
			sess.Exit(1)
		}
	}
}

func main() {
	srv, err := wish.NewServer(
		wish.WithAddress(net.JoinHostPort(host, port)),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithPublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
			return key.Type() == "ssh-ed25519"
		}),
		// Last middleware runs 1st
		wish.WithMiddleware(
			bubbletea.Middleware(ui.Handler),
			activeterm.Middleware(),
			keyAuthMiddleware(),
			logging.Middleware(),
		),
	)

	if err != nil {
		log.Error("Could not start server", "error", err)
	}

	log.Info("Starting SSH server", "host", host, "port", port)
	if err = srv.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		// We ignore ErrServerClosed because it is expected.
		log.Error("Could not start server", "error", err)
	}
}
