package epp

// CheckDomain queries the EPP server for the availability status of one or more domains.
func (c *Conn) CheckDomain(domains ...string) (*DomainCheckResponse, error) {
	req := message{
		Command: &command{
			Check: &check{
				DomainCheck: &domainCheck{
					Domains: domains,
				},
			},
		},
	}
	err := c.writeMessage(&req)
	if err != nil {
		return nil, err
	}
	var res response_
	err = c.readResponse(&res)
	if err != nil {
		return nil, err
	}
	return &res.DomainCheckResponse, nil
}

type DomainCheckResponse struct {
	Checks  []DomainCheck_
	Charges []DomainCharge
}

type DomainCheck_ struct {
	Domain    string
	Reason    string
	Available bool
}

type DomainCharge struct {
	Domain       string
	Category     string
	CategoryName string
}

func init() {
	scanResponse.MustHandleStartElement("epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData", func(c *Context) error {
		c.Value.(*response_).DomainCheckResponse = DomainCheckResponse{}
		return nil
	})
	scanResponse.MustHandleStartElement("epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData > cd", func(c *Context) error {
		dcd := &c.Value.(*response_).DomainCheckResponse
		dcd.Checks = append(dcd.Checks, DomainCheck_{})
		return nil
	})
	scanResponse.MustHandleCharData("epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData > cd > name", func(c *Context) error {
		checks := c.Value.(*response_).DomainCheckResponse.Checks
		check := &checks[len(checks)-1]
		check.Domain = string(c.CharData)
		check.Available = c.AttrBool("", "avail")
		return nil
	})
	scanResponse.MustHandleCharData("epp > response > resData > urn:ietf:params:xml:ns:domain-1.0 chkData > cd > reason", func(c *Context) error {
		checks := c.Value.(*response_).DomainCheckResponse.Checks
		check := &checks[len(checks)-1]
		check.Reason = string(c.CharData)
		return nil
	})
}
