package notification

type notiItem interface {
	isNeedSend() bool
}
