package httpclient

import "github.com/piaodazhu/Octopoda/protocols/errs"

func InitClients() *errs.OctlError {
	if err := initNsClient(); err != nil {
		return err
	}
	if err := initBrainClient(); err != nil {
		return err
	}
	return nil
}
