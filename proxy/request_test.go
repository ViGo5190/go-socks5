package proxy

import "testing"

type rqstFQDNTestCases struct {
	r   *Rqst
	res string
}

func TestRqst_FQDN(t *testing.T) {
	tc := []rqstFQDNTestCases{
		{
			r: &Rqst{
				addressType: addressTypeIPv4,
				addr:        []byte{0xff, 0x8, 0x8, 0x8},
				port:        8888,
			},
			res: "255.8.8.8:8888",
		},
		{
			r:   &Rqst{},
			res: "<unsupported address type>:0",
		},
		{
			r: &Rqst{
				addressType: addressTypeDomain,
				addr:        []byte{103, 117, 109, 101, 110, 105, 117, 107, 46, 99, 111, 109},
				port:        8888,
			},
			res: "gumeniuk.com:8888",
		},
		{
			r: &Rqst{
				addressType: addressTypeIPv6,
				addr:        []byte{0xff, 0, 0, 0, 0, 0, 0xff, 0, 0, 0, 0, 0, 0, 0, 0, 1},
			},
			res: "ff00:0:0:ff00::1:0",
		},
	}

	for _, tcase := range tc {
		fqdn := tcase.r.FQDN()
		if fqdn != tcase.res {
			t.Errorf("Excpected %v, got %v", tcase.res, fqdn)
		}
	}
}
