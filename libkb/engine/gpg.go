package engine

import (
	"fmt"
	"os/exec"

	"github.com/keybase/go/libkb"
	keybase_1 "github.com/keybase/protocol/go"
)

type GPG struct {
	ui       libkb.GPGUI
	secretUI libkb.SecretUI
	last     *libkb.PgpKeyBundle
}

func NewGPG(ui libkb.GPGUI, sui libkb.SecretUI) *GPG {
	return &GPG{ui: ui, secretUI: sui}
}

func (g *GPG) WantsGPG() (bool, error) {
	gpg := G.GetGpgClient()
	if _, err := gpg.Configure(); err != nil {
		if err == exec.ErrNotFound {
			return false, nil
		}
		return false, err
	}

	// they have gpg

	res, err := g.ui.WantToAddGPGKey()
	if err != nil {
		return false, err
	}
	return res, nil
}

func (g *GPG) RunLoadKey(query string) error {
	me, err := libkb.LoadMe(libkb.LoadUserArg{PublicKeyOptional: true})
	if err != nil {
		return err
	}
	sk, err := me.GetDeviceSibkey()
	if err != nil {
		return err
	}
	return g.Run(sk, query)
}

func (g *GPG) Run(signingKey libkb.GenericKey, query string) error {
	gpg := G.GetGpgClient()
	if _, err := gpg.Configure(); err != nil {
		return err
	}
	index, err, warns := gpg.Index(true, query)
	if err != nil {
		return err
	}
	warns.Warn()

	var gks []keybase_1.GPGKey
	for _, key := range index.Keys {
		gk := keybase_1.GPGKey{
			Algorithm:  fmt.Sprintf("%d%s", key.Bits, key.AlgoString()),
			KeyID:      key.GetFingerprint().ToKeyId(),
			Expiration: key.ExpirationString(),
			Identities: key.GetEmails(),
		}
		gks = append(gks, gk)
	}

	res, err := g.ui.SelectKeyAndPushOption(keybase_1.SelectKeyAndPushOptionArg{Keys: gks})
	if err != nil {
		return err
	}
	G.Log.Info("SelectKey result: %+v", res)

	var selected *libkb.GpgPrimaryKey
	for _, key := range index.Keys {
		if key.GetFingerprint().ToKeyId() == res.KeyID {
			selected = key
			break
		}
	}

	if selected == nil {
		return nil
	}

	bundle, err := gpg.ImportKey(true, *(selected.GetFingerprint()))
	if err != nil {
		return fmt.Errorf("ImportKey error: %s", err)
	}

	if err := bundle.Unlock("Import of key into keybase keyring", g.secretUI); err != nil {
		return fmt.Errorf("bundle Unlock error: %s", err)
	}

	G.Log.Info("Bundle unlocked: %s", selected.GetFingerprint().ToKeyId())

	// this seems a little weird to use keygen to post a key, but...
	arg := &libkb.KeyGenArg{
		Pregen:       bundle,
		DoSecretPush: res.DoSecretPush,
		SecretUI:     g.secretUI,
		SigningKey:   signingKey,
	}
	kg := libkb.NewKeyGen(arg)
	if _, err := kg.Run(); err != nil {
		return fmt.Errorf("keygen run error: %s", err)
	}

	G.Log.Info("Key %s imported", selected.GetFingerprint().ToKeyId())

	g.last = bundle

	return nil
}

func (g *GPG) LastKey() *libkb.PgpKeyBundle {
	return g.last
}
