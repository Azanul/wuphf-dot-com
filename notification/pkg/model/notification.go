package model

type Notification struct {
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Msg      string `json:"msg"`
}

func NewNotification(sender, receiver, msg string) (*Notification, error) {
	return &Notification{
		Sender:   sender,
		Receiver: receiver,
		Msg:      msg,
	}, nil
}
