package whoami

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/coredns/caddy"

	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
)

func init() { plugin.Register("whoami", setup) }

func setup(c *caddy.Controller) error {
	c.Next() // 'whoami'
	if c.NextArg() {
		return plugin.Error("whoami", c.ArgErr())
	}

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		return Whoami{}
	})
	c.OnRestart(func() error {
		log.Println(`onrestart called`)
		return nil
	})
	c.OnShutdown(func() error {
		log.Println(`onshutdown called`)
		return nil
	})

	go func() {
		log.Println(`setup whoami called`)
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGUSR2)

		for sig := range sigchan {
			switch sig {
			case syscall.SIGUSR2:
				log.Println(`setup whoami got signal`)
			}
		}
	}()
	return nil
}
