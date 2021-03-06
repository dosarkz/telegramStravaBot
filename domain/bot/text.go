package bot

func getStartMessageText() string {
	return "*Вас приветствует беговой бот Friday* 😃🖐" +
		"\n" +
		"/open ➡️ Открыть главное меню\n" +
		"/close ➡️ Закрыть меню\n" +
		"/rating ➡️ Рейтинг Метронома\n" +
		"По техническим вопросам: [Yerzhan](https://t.me/Yerzhan_Ashenov)\n"
}

func getClubMessageText() string {
	return "*Marat#ON Клуб Марафонцев* \n" +
		"Клуб Любителей Бега Marat#ON (Марафон) вновь создан в г. Астана в начале 2017 года" +
		" и объединяет выпускников школы бега Марата Жыланбаева.\n Мастер спорта международного класса," +
		" ультрамарафонец, первый и единственный атлет в истории человечества, в одиночку пробежавший крупнейшие" +
		" пустыни Азии, Африки, Австралии и Америки.\n Установил несколько мировых рекордов, семь из них занесены в" +
		" Книгу рекордов Гиннеса.\n Большая часть мировых рекордов, установленных Жыланбаевым в начале 1990-х годов" +
		" остаются по-прежнему не превзойденными.\n" +
		"\n"
}

func getWorkoutNewMessage() string {
	return "Для того, чтобы создать тренировку, введите текст одним сообщением в следующем виде:\n\n" +
		"`Тренировка в ТП\n" +
		"Легкий бег - 120 ударов 1 час \n" +
		"2022-10-22 06:00```"
}

func amosovMessageText() string {
	return "Уникальная программа разминки Амосова от Марата Толегеновича, делайте ее ежедневно и будете здоровы!\n" +
		"* Каждое упражнение выполняется по 100 раз!* \n" +
		"- Руки на пояс и движение корпусом влево и вправо \n" +
		"- Движение корпусом вниз и вверх \n" +
		"- Сгибание и разгибание рук к центру \n" +
		"- Ноги вместе, сгибание и полуприсед вперед и назад до начала стопы \n" +
		"- Сгибание и полуприсед влево и вправо \n" +
		"- Круговые движения ног по часовой стрелке и против. \n" +
		"- Выпады. \n"
}
