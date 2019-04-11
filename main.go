package main

import (
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	if err := CreateConfig(); err != nil {
		panic(err)
	}

	ConfigLogging()

	domains := viper.Get("domains")

	dnsProxy := DNSProxy{
		domains:       domains.(map[string]interface{}),
		defaultServer: viper.GetString("default_dns"),
	}

	dns.HandleFunc(".", func(w dns.ResponseWriter, r *dns.Msg) {
		switch r.Opcode {
		case dns.OpcodeQuery:
			m, err := dnsProxy.getResponse(r)
			if err != nil {
				log.Errorf("Failed lookup for %v with error:\n %v\n", r, err.Error())
				m.SetReply(r)
				w.WriteMsg(m)
				return
			}
			m.SetReply(r)
			w.WriteMsg(m)
		}
	})

	host := viper.GetString("listen_host")
	port := viper.GetString("listen_port")
	server := &dns.Server{Addr: host + ":" + port, Net: "udp"}
	log.Infoln("Starting at", host+":"+port)
	err := server.ListenAndServe()
	if err != nil {
		log.Errorln("Failed to start server:", err.Error())
	}
}
