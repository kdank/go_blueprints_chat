package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
	"log"
)

// var ErrNoAvatarUrl = errors.new("chat: Unable to get an avatar URL")
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL.")

type Avatar interface {
	GetAvatarURL(c *client) (string, error)
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (_ AuthAvatar) GetAvatarURL(c *client) (string, error) {
	if url, ok := c.userData["avatar_url"]; ok {
		if urlStr, ok := url.(string); ok {
			return urlStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (_ GravatarAvatar) GetAvatarURL(c *client) (string, error) {
	if email, ok := c.userData["email"]; ok {
		if emailStr, ok := email.(string); ok {
			m := md5.New()
			io.WriteString(m, strings.ToLower(emailStr))
			fmt.Printf("%x", m.Sum(nil))
			log.Println("Final Url", fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)))
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
