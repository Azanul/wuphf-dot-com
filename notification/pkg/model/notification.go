package model

type Notification struct {
	Sender    string `json:"sender"`
	Receiver  string `json:"receiver"`
	Msg       string `json:"msg"`
	Reference string `json:"reference"`
}

func NewNotification(sender, receiver, msg, ref string) (*Notification, error) {
	return &Notification{
		Sender:    sender,
		Receiver:  receiver,
		Msg:       msg,
		Reference: ref,
	}, nil
}
