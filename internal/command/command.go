package command

const (
	ADD     = "add"
	REMOVE  = "rm"
	LIST    = "ls"
	GET     = "get"
	CH      = "ch"
	HISTORY = "h"
	START   = "start"
	HELP    = "help"

	HELP_TXT = `
	start - выводит список коман
	help - выводит список команд
	add <url/domain> - добавляет url в список
	rm <id> -  удаляет url/domain из списка 
	ls - список ваших url/domain
	ch <id> изменит статус сервиса на противоположный(отслеживать/не отслеживать)
	get <id> вернёт информацию о сервисе
	h <id> история пингов
	`
)
