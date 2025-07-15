package audio

import _ "embed"

//go:embed mp3/greeting/greeting_1.mp3
var greeting_1 []byte

//go:embed mp3/greeting/greeting_2.mp3
var greeting_2 []byte

//go:embed mp3/greeting/greeting_3.mp3
var greeting_3 []byte

//go:embed mp3/greeting/greeting_4.mp3
var greeting_4 []byte

//go:embed mp3/greeting/greeting_5.mp3
var greeting_5 []byte

//go:embed mp3/greeting/greeting_6.mp3
var greeting_6 []byte

//go:embed mp3/greeting/greeting_7.mp3
var greeting_7 []byte

//go:embed mp3/greeting/greeting_8.mp3
var greeting_8 []byte

//go:embed mp3/greeting/greeting_9.mp3
var greeting_9 []byte

//go:embed mp3/greeting/greeting_10.mp3
var greeting_10 []byte

//go:embed mp3/greeting/greeting_11.mp3
var greeting_11 []byte

//go:embed mp3/greeting/greeting_12.mp3
var greeting_12 []byte

//go:embed mp3/greeting/greeting_13.mp3
var greeting_13 []byte

//go:embed mp3/greeting/greeting_14.mp3
var greeting_14 []byte

//go:embed mp3/greeting/greeting_15.mp3
var greeting_15 []byte

//go:embed mp3/greeting/greeting_16.mp3
var greeting_16 []byte

//go:embed mp3/greeting/greeting_17.mp3
var greeting_17 []byte

//go:embed mp3/greeting/greeting_18.mp3
var greeting_18 []byte

//go:embed mp3/greeting/greeting_19.mp3
var greeting_19 []byte

//go:embed mp3/greeting/greeting_20.mp3
var greeting_20 []byte

//go:embed mp3/greeting/greeting_21.mp3
var greeting_21 []byte

//go:embed mp3/greeting/greeting_22.mp3
var greeting_22 []byte

//go:embed mp3/greeting/greeting_23.mp3
var greeting_23 []byte

//go:embed mp3/greeting/greeting_24.mp3
var greeting_24 []byte

//go:embed mp3/greeting/greeting_25.mp3
var greeting_25 []byte

//go:embed mp3/greeting/greeting_26.mp3
var greeting_26 []byte

//go:embed mp3/greeting/greeting_27.mp3
var greeting_27 []byte

//go:embed mp3/greeting/greeting_28.mp3
var greeting_28 []byte

//go:embed mp3/greeting/greeting_29.mp3
var greeting_29 []byte

//go:embed mp3/greeting/greeting_30.mp3
var greeting_30 []byte

//go:embed mp3/greeting/greeting_31.mp3
var greeting_31 []byte

//go:embed mp3/greeting/greeting_32.mp3
var greeting_32 []byte

//go:embed mp3/greeting/greeting_33.mp3
var greeting_33 []byte

//go:embed mp3/greeting/greeting_34.mp3
var greeting_34 []byte

//go:embed mp3/greeting/greeting_35.mp3
var greeting_35 []byte

//go:embed mp3/greeting/greeting_36.mp3
var greeting_36 []byte

//go:embed mp3/greeting/greeting_37.mp3
var greeting_37 []byte

//go:embed mp3/greeting/greeting_38.mp3
var greeting_38 []byte

//go:embed mp3/greeting/greeting_39.mp3
var greeting_39 []byte

//go:embed mp3/greeting/greeting_40.mp3
var greeting_40 []byte

//go:embed mp3/greeting/greeting_41.mp3
var greeting_41 []byte

//go:embed mp3/greeting/greeting_42.mp3
var greeting_42 []byte

//go:embed mp3/greeting/greeting_43.mp3
var greeting_43 []byte

//go:embed mp3/greeting/greeting_44.mp3
var greeting_44 []byte

//go:embed mp3/greeting/greeting_45.mp3
var greeting_45 []byte

func SetupAudioService() *AudioService {
	storage := NewMemoryStorage()
	service := NewAudioService(storage)

	// Register all embedded greeting files
	service.StoreNamedAudio("greeting", "greeting_1", greeting_1)
	service.StoreNamedAudio("greeting", "greeting_2", greeting_2)
	service.StoreNamedAudio("greeting", "greeting_3", greeting_3)
	service.StoreNamedAudio("greeting", "greeting_4", greeting_4)
	service.StoreNamedAudio("greeting", "greeting_5", greeting_5)
	service.StoreNamedAudio("greeting", "greeting_6", greeting_6)
	service.StoreNamedAudio("greeting", "greeting_7", greeting_7)
	service.StoreNamedAudio("greeting", "greeting_8", greeting_8)
	service.StoreNamedAudio("greeting", "greeting_9", greeting_9)
	service.StoreNamedAudio("greeting", "greeting_10", greeting_10)
	service.StoreNamedAudio("greeting", "greeting_11", greeting_11)
	service.StoreNamedAudio("greeting", "greeting_12", greeting_12)
	service.StoreNamedAudio("greeting", "greeting_13", greeting_13)
	service.StoreNamedAudio("greeting", "greeting_14", greeting_14)
	service.StoreNamedAudio("greeting", "greeting_15", greeting_15)
	service.StoreNamedAudio("greeting", "greeting_16", greeting_16)
	service.StoreNamedAudio("greeting", "greeting_17", greeting_17)
	service.StoreNamedAudio("greeting", "greeting_18", greeting_18)
	service.StoreNamedAudio("greeting", "greeting_19", greeting_19)
	service.StoreNamedAudio("greeting", "greeting_20", greeting_20)
	service.StoreNamedAudio("greeting", "greeting_21", greeting_21)
	service.StoreNamedAudio("greeting", "greeting_22", greeting_22)
	service.StoreNamedAudio("greeting", "greeting_23", greeting_23)
	service.StoreNamedAudio("greeting", "greeting_24", greeting_24)
	service.StoreNamedAudio("greeting", "greeting_25", greeting_25)
	service.StoreNamedAudio("greeting", "greeting_26", greeting_26)
	service.StoreNamedAudio("greeting", "greeting_27", greeting_27)
	service.StoreNamedAudio("greeting", "greeting_28", greeting_28)
	service.StoreNamedAudio("greeting", "greeting_29", greeting_29)
	service.StoreNamedAudio("greeting", "greeting_30", greeting_30)
	service.StoreNamedAudio("greeting", "greeting_31", greeting_31)
	service.StoreNamedAudio("greeting", "greeting_32", greeting_32)
	service.StoreNamedAudio("greeting", "greeting_33", greeting_33)
	service.StoreNamedAudio("greeting", "greeting_34", greeting_34)
	service.StoreNamedAudio("greeting", "greeting_35", greeting_35)
	service.StoreNamedAudio("greeting", "greeting_36", greeting_36)
	service.StoreNamedAudio("greeting", "greeting_37", greeting_37)
	service.StoreNamedAudio("greeting", "greeting_38", greeting_38)
	service.StoreNamedAudio("greeting", "greeting_39", greeting_39)
	service.StoreNamedAudio("greeting", "greeting_40", greeting_40)
	service.StoreNamedAudio("greeting", "greeting_41", greeting_41)
	service.StoreNamedAudio("greeting", "greeting_42", greeting_42)
	service.StoreNamedAudio("greeting", "greeting_43", greeting_43)
	service.StoreNamedAudio("greeting", "greeting_44", greeting_44)
	service.StoreNamedAudio("greeting", "greeting_45", greeting_45)

	return service
}
