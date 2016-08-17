package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	gomniauthtest "github.com/stretchr/gomniauth/test"
)

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar

	testUser := &gomniauthtest.TestUser{}
	testUser.On("AvatarURL").Return("", ErrNoAvatarURL)
	testChatUser := &chatUser{User: testUser}
	url, err := authAvatar.GetAvatarURL(testChatUser)

	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、ErrNoAvatarURLを返却すべきです.")
	}

	testUrl := "http://url-to-avatar"
	testUser = &gomniauthtest.TestUser{}
	testChatUser.User = testUser
	testUser.On("AvatarURL").Return(testUrl, nil)
	url, err = authAvatar.GetAvatarURL(testChatUser)

	if err != nil {
		t.Error("値が存在する場合、エラーを返却すべきではありません.")
	} else if url != testUrl {
		t.Error("正しいURLを返却すべきです.")
	}

}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	user := &chatUser{uniqueID: "abc"}

	url, err := gravatarAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("エラーを返却すべきではありません.")
	}
	if url != "//www.gravatar.com/avatar/abc" {
		t.Errorf("%sという誤った値を返しました", url)
	}
}

func TestFileSystemAvatar(t *testing.T) {

	filename := filepath.Join("avatars", "abc.jpg")
	ioutil.WriteFile(filename, []byte{}, 0777)
	defer func() { os.Remove(filename) }()

	var fileSystemAvatar FileSystemAvatar
	user := &chatUser{uniqueID: "abc"}

	url, err := fileSystemAvatar.GetAvatarURL(user)
	if err != nil {
		t.Error("エラーを返却すべきではありません.")
	}
	if url != "/avatars/abc.jpg" {
		t.Errorf("%sという誤った値を返しました", url)
	}
}
