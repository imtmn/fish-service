package dto

type (
	WxData struct {
		// 不包括敏感信息的原始数据字符串，用于计算签名。
		RawData string `bson:"rawData" json:"rawData"`
		// 小程序通过 api 得到的加密数据(encryptedData)
		EncryptedData string `bson:"encryptedData" json:"encryptedData"`
		// 使用 sha1( rawData + session_key ) 得到字符串，用于校验用户信息
		Signature string `bson:"signature" json:"signature"`
		// 小程序通过 api 得到的初始向量(iv)
		Iv string `bson:"iv" json:"iv"`
	}
)
