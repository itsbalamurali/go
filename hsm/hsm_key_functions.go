package hsm

import (
	"encoding/hex"
	"github.com/rkbalgi/go/crypto"
	"strconv"
)

//generate a check value for the given key
func gen_check_value(key_data []byte) []byte {

	data := make([]byte, 8)
	if len(key_data) == 8 {
		return (crypto.EncryptDes(data, key_data))
	} else {
		return (crypto.EncryptTripleDes(data, key_data))
	}

}

//extract key data (as []byte) from a encoded string
func extract_key_data(key_str string) []byte {

	var key_data []byte
	if key_str[0] == 'U' || key_str[0] == 'Z' || key_str[0] == 'T' || key_str[0] == 'X' || key_str[0] == 'Y' {
		key_data, _ = hex.DecodeString(key_str[1:])
	} else {
		key_data, _ = hex.DecodeString(key_str[:])
	}
	return key_data
}

//encrypt a key under the lmk using thales variants

func encrypt_key(key_str string, key_type string) []byte {

	var key_data []byte
	if key_str[0] == 'U' || key_str[0] == 'Z' || key_str[0] == 'T' {
		key_data, _ = hex.DecodeString(key_str[1:])
	} else {
		key_data, _ = hex.DecodeString(key_str[:])
	}

	lmk_key := key_type_table[key_type[1:]]
	if lmk_key == nil {
		panic("unsupported key type" + key_type)
	}

	return (encrypt_key_thales_v(key_data, lmk_key, key_type))
}

//encrypt a key under the kek using thales variants
func encrypt_key_kek(key_str string, kek []byte, key_type string) []byte {

	key_data := extract_key_data(key_str)
	return (encrypt_key_thales_v(key_data, kek, key_type))
}

func encrypt_key_thales_v(key_data []byte, lmk_key []byte, key_type string) []byte {

	v_key := make([]byte, 16)
	copy(v_key, lmk_key)

	set_variants(v_key, key_type)

	if len(key_data) == 24 {
		//left
		v_key[8] = v_key[8] ^ __variant_triple_len_1
		left := crypto.EncryptTripleDes(key_data[0:8], v_key)

		//now middle
		v_key[8] = v_key[8] ^ __variant_triple_len_1
		v_key[8] = v_key[8] ^ __variant_triple_len_2
		middle := crypto.EncryptTripleDes(key_data[8:16], v_key)

		//right
		v_key[8] = v_key[8] ^ __variant_triple_len_2
		v_key[8] = v_key[8] ^ __variant_triple_len_3
		right := crypto.EncryptTripleDes(key_data[16:], v_key)

		clear_key := make([]byte, 0)
		clear_key = append(clear_key, left...)
		clear_key = append(clear_key, middle...)
		clear_key = append(clear_key, right...)
		return (clear_key)

	} else if len(key_data) == 16 {
		v_key[8] = v_key[8] ^ __variant_dbl_len_1
		left_half := crypto.EncryptTripleDes(key_data[0:8], v_key)

		//now right half
		v_key[8] = v_key[8] ^ __variant_dbl_len_1
		v_key[8] = v_key[8] ^ __variant_dbl_len_2
		right_half := crypto.EncryptTripleDes(key_data[8:], v_key)

		clear_key := make([]byte, 0)
		clear_key = append(clear_key, left_half...)
		clear_key = append(clear_key, right_half...)

		return (clear_key)

	} else if len(key_data) == 8 {
		//single length
		clear_key := crypto.EncryptTripleDes(key_data, v_key)
		return (clear_key)

	} else {
		panic("illegal key size")
	}

}

//decrypt a key encrypted under the lmk using thales variants

func decrypt_key(key_str string, key_type string) []byte {

	key_data := extract_key_data(key_str)

	lmk_key := key_type_table[key_type[1:]]
	if lmk_key == nil {
		panic("unsupported key type" + key_type)
	}

	return (decrypt_key_thales_v(key_data, lmk_key, key_type))

}

//decrypt a key under the kek using thales variants
func decrypt_key_kek(key_str string, kek []byte, key_type string) []byte {
	key_data := extract_key_data(key_str)
	return (decrypt_key_thales_v(key_data, kek, key_type))
}

func decrypt_key_thales_v(key_data []byte, lmk_key []byte, key_type string) []byte {

	v_key := make([]byte, 16)
	copy(v_key, lmk_key)
	set_variants(v_key, key_type)

	if len(key_data) == 24 {
		//left
		v_key[8] = v_key[8] ^ __variant_triple_len_1
		left := crypto.DecryptTripleDes(key_data[0:8], v_key)

		//now middle
		v_key[8] = v_key[8] ^ __variant_triple_len_1
		v_key[8] = v_key[8] ^ __variant_triple_len_2
		middle := crypto.DecryptTripleDes(key_data[8:16], v_key)

		//right
		v_key[8] = v_key[8] ^ __variant_triple_len_2
		v_key[8] = v_key[8] ^ __variant_triple_len_3
		right := crypto.DecryptTripleDes(key_data[16:], v_key)

		clear_key := make([]byte, 0)
		clear_key = append(clear_key, left...)
		clear_key = append(clear_key, middle...)
		clear_key = append(clear_key, right...)
		return (clear_key)

	} else if len(key_data) == 16 {
		v_key[8] = v_key[8] ^ __variant_dbl_len_1
		left_half := crypto.DecryptTripleDes(key_data[0:8], v_key)

		//now right half
		v_key[8] = v_key[8] ^ __variant_dbl_len_1
		v_key[8] = v_key[8] ^ __variant_dbl_len_2
		right_half := crypto.DecryptTripleDes(key_data[8:], v_key)

		clear_key := make([]byte, 0)
		clear_key = append(clear_key, left_half...)
		clear_key = append(clear_key, right_half...)

		return (clear_key)

	} else if len(key_data) == 8 {
		//single length
		clear_key := crypto.DecryptTripleDes(key_data, v_key)
		return (clear_key)

	} else {
		panic("illegal key size")
	}

}

func set_variants(v_key []byte, key_type string) {

	if key_type[0] != '0' {
		//apply variant
		i, _ := strconv.ParseInt(key_type[:1], 10, 32)
		v_key[0] = v_key[0] ^ __variants[i]

	}
}
