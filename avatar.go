package main

import (
	"errors"
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
