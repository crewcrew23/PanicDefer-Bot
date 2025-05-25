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
)

func HelpTXT() string {
	return `📌 *Доступные команды:*

` + "```" + `
/add <url>    - Добавить URL для мониторинга
/rm <id>      - Удалить URL по ID
/ls           - Список всех ваших URL
/ch <id>      - Вкл/Выкл мониторинг
/get <id>     - Информация о сервисе
/h <id>       - История проверок
/help         - Справка по командам
` + "```" + `

*Примеры:*
• /add https://example.com
• /rm 5
• /h 3`
}
