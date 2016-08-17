package main

import (
	"errors"
	"io/ioutil"
	"path/filepath"
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
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			return "//www.gravatar.com/avatar/" + useridStr, nil
		}
	}
	return "", ErrNoAvatarURL
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (_ FileSystemAvatar) GetAvatarURL(c *client) (string, error) {
	if userid, ok := c.userData["userid"]; ok {
		if useridStr, ok := userid.(string); ok {
			if files, err := ioutil.ReadDir("avatars"); err == nil {
				for _, file := range files {
					if file.IsDir() {
						continue
					}
					match, _ := filepath.Match(useridStr+"*", file.Name())
					if match {
						return "/avatars/" + file.Name(), nil
					}
				}
			}
		}
	}

	return "", ErrNoAvatarURL
}
