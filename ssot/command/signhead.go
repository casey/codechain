package command

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/frankbraun/codechain/archive"
	"github.com/frankbraun/codechain/hashchain"
	"github.com/frankbraun/codechain/internal/def"
	"github.com/frankbraun/codechain/internal/hex"
	"github.com/frankbraun/codechain/secpkg"
	"github.com/frankbraun/codechain/ssot"
	"github.com/frankbraun/codechain/util/file"
	"github.com/frankbraun/codechain/util/homedir"
	"github.com/frankbraun/codechain/util/interrupt"
	"github.com/frankbraun/codechain/util/log"
	"github.com/frankbraun/codechain/util/seckey"
)

func signHead(c *hashchain.HashChain) error {
	// 1. Parse the .secpkg file in the current working directory.
	log.Println("1. parse .secpkg")
	pkg, err := secpkg.Load(secpkg.File)
	if err != nil {
		return err
	}

	// 2. Make sure the project with NAME has been published before.
	//    That is, the directory ~/.config/ssotpub/pkgs/NAME exists.
	log.Println("1. make sure project exists")
	pkgDir := filepath.Join(homedir.SSOTPub(), "pkgs", pkg.Name)
	exists, err := file.Exists(pkgDir)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("package not published yet: '%s' does not exist", pkgDir)
	}

	// 3. Validate the signed head in ~/.config/ssotpub/pkgs/NAME/signed_head
	//    and make sure the corresponding secret key is available.
	log.Println("3. validate signed head")
	signedHeadFile := filepath.Join(pkgDir, "signed_head")
	prevSignedHead, err := ssot.Load(signedHeadFile)
	if err != nil {
		return err
	}

	secKeyFile := filepath.Join(homedir.SSOTPub(), def.SecretsSubDir, prevSignedHead.PubKey())
	secKey, _, _, err := seckey.Read(secKeyFile)
	if err != nil {
		return err
	}

	// 4. Get the HEAD from .codechain/hashchain in the current working directory.
	log.Println("4. get the head")
	head := c.Head()
	fmt.Printf("signing head %x\n", head)

	// 5. Create a new signed head with current HEAD, the counter of the previous
	/*
	   signed head plus 1, and update the saved signed head:

	   - `cp -f ~/.config/ssotpub/pkgs/NAME/signed_head
	            ~/.config/ssotpub/pkgs/NAME/previous_signed_head`
	   - Save new signed head to ~/.config/ssotpub/pkgs/NAME/signed_head (atomic).
	*/
	newSignedHead := ssot.SignHead(head, prevSignedHead.Counter()+1, *secKey)
	if err := newSignedHead.RotateFile(pkgDir); err != nil {
		return err
	}

	// 6. Save the current distribution to:
	//    ~/.config/secpkg/pkgs/NAME/dists/HEAD.tar.gz (`codechain createdist`).
	distDir := filepath.Join(pkgDir, "dists")
	distFile := filepath.Join(distDir, fmt.Sprintf("%x.tar.gz", head))
	if err := archive.CreateDist(c, distFile); err != nil {
		return err
	}

	// 7. Print the distribution name: ~/.config/ssotpkg/pkgs/NAME/dists/HEAD.tar.gz
	fmt.Println("")
	fmt.Printf("Please upload the following distribution file to: %s\n", pkg.URL)
	fmt.Println(distFile)
	fmt.Println("")

	// 8. Print DNS TXT record as defined by the .secpkg and the signed head.
	fmt.Println("Please publish the following DNS TXT record:")
	fmt.Println("")
	newSignedHead.PrintTXT(pkg.DNS)

	// 9. If the HEAD changed, update the .secpkg file accordingly.
	h := hex.Encode(head[:])
	if h != pkg.Head {
		pkg.Head = h
		newSecPkgFile := secpkg.File + "_new"
		err = ioutil.WriteFile(newSecPkgFile, []byte(pkg.Marshal()+"\n"), 0644)
		if err != nil {
			return err
		}
		if err := os.Rename(newSecPkgFile, secpkg.File); err != nil {
			return err
		}
		fmt.Printf("\n%s: updated\n", secpkg.File)
	}

	return nil
}

// SignHead implements the ssotpub 'signhead' command.
func SignHead(argv0 string, args ...string) error {
	fs := flag.NewFlagSet(argv0, flag.ContinueOnError)
	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s -s seckey.bin\n", argv0)
		fmt.Fprintf(os.Stderr, "Sign Codechain head and print it on stdout.\n")
		fs.PrintDefaults()
	}
	//secKey := fs.String("s", "", "Secret key file")
	verbose := fs.Bool("v", false, "Be verbose")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *verbose {
		log.Std = log.NewStd(os.Stdout)
	}
	/*
		if err := seckey.Check(homedir.SSOTPub(), *secKey); err != nil {
			return err
		}
	*/
	if fs.NArg() != 0 {
		fs.Usage()
		return flag.ErrHelp
	}
	c, err := hashchain.ReadFile(def.HashchainFile)
	if err != nil {
		return err
	}
	defer c.Close()
	// add interrupt handler
	interrupt.AddInterruptHandler(func() {
		c.Close()
	})
	// run signHead
	go func() {
		if err := signHead(c); err != nil {
			interrupt.ShutdownChannel <- err
			return
		}
		interrupt.ShutdownChannel <- nil
	}()
	return <-interrupt.ShutdownChannel
}
