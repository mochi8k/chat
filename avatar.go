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
	GetAvatarURL(ChatUser) (string, error)
}

type Avatars []Avatar

func (as Avatars) GetAvatarURL(u ChatUser) (string, error) {
	for _, avatar := range as {
		url, err := avatar.GetAvatarURL(u)
		if err == nil {
			return url, nil
		}
	}
	return "", ErrNoAvatarURL
}

type AuthAvatar struct{}

var UseAuthAvatar AuthAvatar

func (AuthAvatar) GetAvatarURL(u ChatUser) (string, error) {
	url := u.AvatarURL()
	if url == "" {
		return "", ErrNoAvatarURL
	}
	return url, nil
}

type GravatarAvatar struct{}

var UseGravatar GravatarAvatar

func (GravatarAvatar) GetAvatarURL(u ChatUser) (string, error) {
	return "//www.gravatar.com/avatar/" + u.getUniqueID(), nil
}

type FileSystemAvatar struct{}

var UseFileSystemAvatar FileSystemAvatar

func (FileSystemAvatar) GetAvatarURL(u ChatUser) (string, error) {
	files, err := ioutil.ReadDir("avatars")
	if err != nil {
		return "", ErrNoAvatarURL
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		match, _ := filepath.Match(u.getUniqueID()+"*", file.Name())
		if match {
			return "/avatars/" + file.Name(), nil
		}
	}

	return "", ErrNoAvatarURL
}
