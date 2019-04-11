package main

import (
	"fmt"
	"regexp"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type DNSProxy struct {
	domains       map[string]interface{}
	defaultServer string
}

func (proxy *DNSProxy) getResponse(requestMsg *dns.Msg) (*dns.Msg, error) {
	responseMsg := new(dns.Msg)
	if len(requestMsg.Question) > 0 {
		for _, question := range requestMsg.Question {
			dnsServer := proxy.getIPFromConfigs(question.Name, proxy.domains)
			if dnsServer == "" {
				dnsServer = proxy.defaultServer
			}

			answers, err := proxy.processDNS(dnsServer, &question, requestMsg)
			if err != nil {
				return responseMsg, err
			}
			for _, answer := range answers {
				responseMsg.Answer = append(responseMsg.Answer, answer)
			}
		}
	}

	return responseMsg, nil
}

func (proxy *DNSProxy) processDNS(dnsServer string, q *dns.Question, requestMsg *dns.Msg) ([]dns.RR, error) {
	queryMsg := new(dns.Msg)
	requestMsg.CopyTo(queryMsg)
	queryMsg.Question = []dns.Question{*q}

	msg, err := lookup(dnsServer, queryMsg)
	if err != nil {
		return nil, err
	}
	if len(msg.Answer) > 0 {
		return msg.Answer, nil
	}
	return nil, fmt.Errorf("not found")
}

func (dnsProxy *DNSProxy) getIPFromConfigs(domain string, configs map[string]interface{}) string {

	for k, v := range configs {
		match, _ := regexp.MatchString(k+"\\.", domain)
		if match {
			log.Infoln("Match ip from config", v.(string))
			return v.(string)
		}
	}
	return ""
}

func lookup(server string, m *dns.Msg) (*dns.Msg, error) {
	dnsClient := new(dns.Client)
	dnsClient.Net = "udp"
	response, _, err := dnsClient.Exchange(m, server)
	log.Debugf("Lookup for %v\n"+
		"Request:\n%v\n-----\n"+
		"Response:\n%v\n----\n", server, m, response)
	if err != nil {
		return nil, err
	}
	return response, nil
}
