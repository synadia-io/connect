package cli

import (
	"github.com/nats-io/nkeys"
	log2 "github.com/rs/zerolog/log"
	"os"
	"path/filepath"
)

func GetXKey(xkeyFile string, xkey []byte) (nkeys.KeyPair, error) {
	xkeyFile = os.ExpandEnv(xkeyFile)

	// if no xkey is provided, we will try to read it from the file
	if xkey == nil || len(xkey) == 0 {
		var err error
		xkey, err = readXKeyFileOrCreate(xkeyFile)
		if err != nil {
			return nil, err
		}
	}

	return nkeys.FromCurveSeed(xkey)
}

func readXKeyFileOrCreate(xkeyFile string) ([]byte, error) {
	// if the file exists, read it and return the seed
	seed, err := os.ReadFile(xkeyFile)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	}

	if seed == nil {
		// create the directory to write the xkey file
		if err := os.MkdirAll(filepath.Dir(xkeyFile), 0700); err != nil {
			return nil, err
		}

		// create a new xkey
		xk, err := nkeys.CreateCurveKeys()
		if err != nil {
			return nil, err
		}

		ns, err := xk.Seed()
		if err != nil {
			return nil, err
		}

		err = os.WriteFile(xkeyFile, ns, 0600)
		if err != nil {
			return nil, err
		}
		log2.Info().Msgf("Created new xkey at %s", xkeyFile)

		seed = ns
	}

	return seed, nil
}
