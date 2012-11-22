package client

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	w             http.ResponseWriter
	req           *http.Request
	requestHeader map[string][]string
	secret        []byte
	expires       int
	where         string
}

func New(w http.ResponseWriter, req *http.Request, requestHeader map[string][]string, secret string) *Client {
	if len(secret) < 16 {
		panic("Secret length must be at least 16.")
	}
	return &Client{
		w,
		req,
		requestHeader,
		[]byte(secret)[0:16],
		3600000,
		"/",
	}
}

func (c *Client) Store(key string, val interface{}) error {
	cookie := &http.Cookie{
		Name:   key,
		Value:  val.(string),
		MaxAge: c.expires,
		Path:   c.where,
	}
	http.SetCookie(c.w, cookie)
	return nil
}

func (c *Client) StoreEncrypted(key string, val interface{}) error {
	enc, err := encryptString(c.secret, val.(string))
	if err != nil {
		return err
	}
	cookie := &http.Cookie{
		Name:   key,
		Value:  enc,
		MaxAge: c.expires,
		Path:   c.where,
	}
	http.SetCookie(c.w, cookie)
	return nil
}

func (c *Client) Get(key string) (interface{}, error) {
	cookie, err := c.req.Cookie("user")
	if err != nil {
		return nil, err
	}
	return cookie.Value, nil
}

func (c *Client) GetDecrypted(key string) (interface{}, error) {
	cookie, err := c.req.Cookie("user")
	if err != nil {
		return nil, err
	}
	dec, err := decryptString(c.secret, cookie.Value)
	if err != nil {
		return nil, err
	}
	return dec, nil
}

func (c *Client) Unstore(key string) error {
	cookie := &http.Cookie{
		Name:   key,
		Value:  "",
		MaxAge: c.expires,
		Path:   c.where,
	}
	http.SetCookie(c.w, cookie)
	return nil
}

func parseAcceptLanguage(l string) []string {
	ret := []string{}
	sl := strings.Split(l, ",")
	c := map[string]struct{}{}
	for _, v := range sl {
		lang := string(strings.Split(v, ";")[0][:2])
		_, has := c[lang]
		if !has {
			c[lang] = struct{}{}
			ret = append(ret, lang)
		}
	}
	return ret
}

// Creates a list of 2 char language abbreviations (for example: []string{"en", "de", "hu"}) out of the value of http header "Accept-Language".
func parseAcceptLanguageSafe(l string) (ret []string) {
	defer func() {
		r := recover()
		if r != nil {
			ret = []string{"en"}
		}
	}()
	ret = parseAcceptLanguage(l)
	return
}

func (c *Client) Languages() []string {
	langs, has := c.req.Header["Accept-Language"]
	if has && len(langs) > 0 {
		return parseAcceptLanguageSafe(langs[0])
	}
	return nil
}

// Encrypts a value and encodes it with base64.
func encryptString(block_key []byte, value string) (string, error) {
	str, err := encDecStr(block_key, value, true)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString([]byte(str)), nil
}

// Decodes a value with base64 and then decrypts it.
func decryptString(block_key []byte, value string) (string, error) {
	decoded_b, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	return encDecStr(block_key, string(decoded_b), false)
}

// Function intended to encrypt the user id before storing it as a cookie.
// encr flag controls
// block_key must be secret.
func encDecStr(block_key []byte, value string, encr bool) (string, error) {
	if len(value) == 0 {
		return "", fmt.Errorf("Nothing to encrypt/decrypt.")
	}
	block, err := aes.NewCipher(block_key)
	if err != nil {
		return "", err
	}
	var bs []byte
	if encr {
		bs, err = encrypt(block, []byte(value))
	} else {
		bs, err = decrypt(block, []byte(value))
	}
	if err != nil {
		return "", err
	}
	// Just in case.
	if bs == nil {
		return "", fmt.Errorf("Somethign went wrong when encoding/decoding.")
	}
	return string(bs), nil
}

// The following functions are taken from securecookie package of the Gorilla web toolkit made by Rodrigo Moraes.
// Only modification was to make the GenerateRandomKey function private.

// encrypt encrypts a value using the given block in counter mode.
//
// A random initialization vector (http://goo.gl/zF67k) with the length of the
// block size is prepended to the resulting ciphertext.
func encrypt(block cipher.Block, value []byte) ([]byte, error) {
	iv := generateRandomKey(block.BlockSize())
	if iv == nil {
		return nil, errors.New("securecookie: failed to generate random iv")
	}
	// Encrypt it.
	stream := cipher.NewCTR(block, iv)
	stream.XORKeyStream(value, value)
	// Return iv + ciphertext.
	return append(iv, value...), nil
}

// decrypt decrypts a value using the given block in counter mode.
//
// The value to be decrypted must be prepended by a initialization vector
// (http://goo.gl/zF67k) with the length of the block size.
func decrypt(block cipher.Block, value []byte) ([]byte, error) {
	size := block.BlockSize()
	if len(value) > size {
		// Extract iv.
		iv := value[:size]
		// Extract ciphertext.
		value = value[size:]
		// Decrypt it.
		stream := cipher.NewCTR(block, iv)
		stream.XORKeyStream(value, value)
		return value, nil
	}
	return nil, errors.New("securecookie: the value could not be decrypted")
}

// GenerateRandomKey creates a random key with the given strength.
func generateRandomKey(strength int) []byte {
	k := make([]byte, strength)
	if _, err := rand.Read(k); err != nil {
		return nil
	}
	return k
}
