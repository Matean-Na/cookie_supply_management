package logger

import (
	"log"
	"os"
)

var LogFile *os.File

func Init(fileName string) error {
	var err error

	// Открываем файл журнала
	LogFile, err = os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// Устанавливаем файл журнала как вывод для пакета log
	log.SetOutput(LogFile)

	// Возвращаем ошибку
	return nil
}
