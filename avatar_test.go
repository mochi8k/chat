package main

import "testing"

func TestAuthAvatar(t *testing.T) {
	var authAvatar AuthAvatar
	client := new(client)
	url, err := authAvatar.GetAvatarURL(client)
	if err != ErrNoAvatarURL {
		t.Error("値が存在しない場合、ErrNoAvatarURLを返却すべきです.")
	}

	testUrl := "http://url-to-avatar"
	client.userData = map[string]interface{}{"avatar_url": testUrl}

	url, err = authAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("値が存在する場合、エラーを返却すべきではありません.")
	} else if url != testUrl {
		t.Error("正しいURLを返却すべきです.")
	}

}

func TestGravatarAvatar(t *testing.T) {
	var gravatarAvatar GravatarAvatar
	client := new(client)
	client.userData = map[string]interface{}{"email": "MyEmailAddress@example.com"}
	url, err := gravatarAvatar.GetAvatarURL(client)
	if err != nil {
		t.Error("エラーを返却すべきではありません.")
	}
	if url != "//www.gravatar.com/avatar/0bc83cb571cd1c50ba6f3e8a78ef1346" {
		t.Errorf("%sという誤った値を返しました", url)
	}
}
