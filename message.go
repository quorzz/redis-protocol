package protocol

import (
	"errors"
	"strconv"
)

const (
	MessageError  = '-'
	MessageStatus = '+'
	MessageInt    = ':'
	MessageBulk   = '$'
	MessageMutli  = '*'
)

type Message struct {
	Type    byte
	Error   error
	Status  string
	Integer int64
	Bulk    []byte
	Multi   []*Message
}

func (Message *Message) HasError() bool {
	return Message.Type == MessageError
}

func (Message *Message) Bytes() ([]byte, error) {

	switch Message.Type {
	case MessageError:
		return nil, Message.Error
	case MessageStatus:
		return []byte(Message.Status), nil
	case MessageInt:
		return nil, errors.New("Integer Message can not convert to []byte.")
	case MessageBulk:
		return Message.Bulk, nil
	case MessageMutli:
		return nil, errors.New("Mutli Message can not convert to []byte.")
	}
	return nil, errors.New("Invalid Message type.")
}

func (Message *Message) String() (string, error) {

	switch Message.Type {
	case MessageError:
		return "", Message.Error
	case MessageStatus:
		return Message.Status, nil
	case MessageInt:
		return "", errors.New("Integer Message can not convert to string.")
	case MessageBulk:
		return string(Message.Bulk), nil
	case MessageMutli:
		return "", errors.New("Mutli Message can not convert to string.")
	}
	return "", errors.New("Invalid Message type.")
}

func (Message *Message) Int64() (int64, error) {

	switch Message.Type {
	case MessageError:
		return 0, nil
	case MessageStatus:
		if string(Message.Status) == "OK" {
			return 1, nil
		} else {
			return 0, nil
		}
	case MessageInt:
		return Message.Integer, nil
	case MessageBulk:
		return strconv.ParseInt(string(Message.Bulk), 10, 64)
	case MessageMutli:
		return 0, errors.New("Mutli Message can not convert to integer.")
	}
	return 0, errors.New("Invalid Message type.")
}

func (Message *Message) Int() (int, error) {

	i64, err := Message.Int64()
	if err != nil {
		return 0, err
	}
	return int(i64), nil
}

func (Message *Message) Bool() (bool, error) {

	switch Message.Type {
	case MessageError:
		return false, nil
	case MessageStatus:
		if string(Message.Status) == "OK" || string(Message.Status) == "PONG" {
			return true, nil
		} else {
			return false, nil
		}
	case MessageInt:
		return Message.Integer != 0, nil
	case MessageBulk:
		return strconv.ParseBool(string(Message.Bulk))
	case MessageMutli:
		return false, errors.New("Mutli Message can not convert to bool.")
	}
	return false, errors.New("Invalid Message type.")
}

func (Message *Message) StringMap() (map[string]string, error) {
	if Message.Type != MessageMutli {
		return nil, errors.New("Only mutli reponse can convert to [string]string.")
	}

	result := make(map[string]string)
	length := len(Message.Multi)
	if Message.Multi == nil || length <= 0 {
		return result, nil
	}

	for i := 0; i < length/2; i++ {
		key, err := Message.Multi[i*2].String()
		if err != nil {
			return nil, err
		}

		value, err := Message.Multi[i*2+1].String()
		if err != nil {
			return nil, err
		}

		result[key] = value
	}
	return result, nil
}

func (Message *Message) Strings() ([]string, error) {

	if Message.Type != MessageMutli {
		return nil, errors.New("Only mutli reponse can convert to []string.")
	}

	if Message.Multi == nil {
		return nil, nil
	}

	var result []string
	for _, v := range Message.Multi {
		vv, err := v.String()
		if err != nil {
			return nil, err
		}
		result = append(result, vv)
	}

	return result, nil
}
