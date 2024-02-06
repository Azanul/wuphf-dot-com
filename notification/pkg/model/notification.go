package model

type Notification struct {
	Sender   string `json:"from"`
	Receiver string `json:"to"`
	Msg      string `json:"msg"`
}

func NewNotification(sender, receiver, msg string) (*Notification, error) {
	return &Notification{
		Sender:   sender,
		Receiver: receiver,
		Msg:      msg,
	}, nil
}
