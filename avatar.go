package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io"
	"strings"
)

var ErrNoAvatarURL = errors.New("chat: アバターを取得できません.")

type Avatar interface {
	/*
	   指定されたクライアントのアバターのURLを返却する.
	   問題が発生した場合はエラーを返却する.
	   URLを取得できなかった場合はErrNoAvatarURLを返却する.
	*/
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
			return fmt.Sprintf("//www.gravatar.com/avatar/%x", m.Sum(nil)), nil
		}
	}
	return "", ErrNoAvatarURL
}
