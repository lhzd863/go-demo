package jwt

import (
        "bytes"
        "encoding/base64"
        "encoding/json"
        "errors"

        "crypto/hmac"
        "crypto/sha256"
        "crypto/sha512"
        "hash"
)

var separator []byte = []byte{'.'}

var (
        SecretError = errors.New("Signature verification failed")
)

func base64url_encode(b []byte) []byte {
        encoded := []byte(base64.URLEncoding.EncodeToString(b))
        var equalIndex = bytes.Index(encoded, []byte{'='})
        if equalIndex > -1 {
                encoded = encoded[:equalIndex]
        }
        return encoded
}

func base64url_decode(b []byte) ([]byte, error) {
        if len(b)%4 != 0 {
                b = append(b, bytes.Repeat([]byte{'='}, 4-(len(b)%4))...)
        }
        decoded, err := base64.URLEncoding.DecodeString(string(b))
        if err != nil {
                return nil, err
        }
        return decoded, nil
}

func getHash(algorithm string) (func() hash.Hash, error) {
        switch algorithm {
        case "HS256":
                return sha256.New, nil
        case "HS384":
                return sha512.New384, nil
        case "HS512":
                return sha512.New, nil
        }
        return nil, errors.New("Algorithm not supported")
}

func Encode(jwt interface{}, key []byte, algorithm string) ([]byte, error) {
        shaFunc, err := getHash(algorithm)
        if err != nil {
                return []byte{}, err
        }
        sha := hmac.New(shaFunc, key)

        segments := [3][]byte{}

        header, err := json.Marshal(
                map[string]interface{}{
                        "typ": "JWT",
                        "alg": algorithm,
                })
        if err != nil {
                return []byte{}, err
        }
        segments[0] = base64url_encode(header)

        claims, err := json.Marshal(jwt)
        if err != nil {
                return []byte{}, err
        }
        segments[1] = base64url_encode(claims)

        sha.Write(bytes.Join(segments[:2], separator))
        segments[2] = base64url_encode(sha.Sum(nil))

        return bytes.Join(segments[:], separator), nil
}

func Decode(encoded []byte, claims interface{}, key []byte) error {
        segments := bytes.Split(encoded, separator)

        // segments is currently slices make copies so functions like 
        // base64url_decode will not overwrite later portions
        for k, v := range segments {
                newBytes := make([]byte, len(v))
                copy(newBytes, v)
                segments[k] = newBytes
        }

        if len(segments) != 3 {
                return errors.New("Incorrect segment count")
        }

        var header map[string]interface{} = make(map[string]interface{})
        headerBase64, err := base64url_decode(segments[0])
        if err != nil {
                return err
        }
        err = json.Unmarshal(headerBase64, &header)
        if err != nil {
                return err
        }

        algorithm, ok := header["alg"].(string)
        var sha hash.Hash
        if ok {
                shaFunc, err := getHash(algorithm)
                if err != nil {
                        return err
                }
                sha = hmac.New(shaFunc, key)
        } else {
                return errors.New("Algorithm not supported")
        }

        claimsBase64, err := base64url_decode(segments[1])
        if err != nil {
                return err
        }
        err = json.Unmarshal(claimsBase64, claims)
        if err != nil {
                return err
        }

        sha.Write(bytes.Join(segments[:2], separator))
        signature := base64url_encode(sha.Sum(nil))
        if bytes.Compare(signature, segments[2]) != 0 {
                return SecretError
        }
        return nil
}
