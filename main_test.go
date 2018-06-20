package main

import (
	"testing"
	"fmt"
)
var v error

// Основной тест
func TestIncrementNumber1(t *testing.T) {

	fmt.Println("Start main test")

	//Текущее значение
	fmt.Println("Current number value:", getNumber())

	//Пробуем инкремент
	incrementNumber()

	//Значение после инкремента
	fmt.Println("Current number value:", getNumber())
	i := 1

	for i <= 501 {
		incrementNumber()
		i = i + 1
	}
	//Значение после цикла инкрементов
	//Должно быть 501
	fmt.Println("Current number value:", getNumber())

	//Если нет то ошибка
	if getNumber()!=501{
		t.Error("Error mass increment", getNumber())
	}


}

func TestSetMaximumValueNorm(t *testing.T) {
	// number должен быть 501
	// Мы станавливаем  MaximumValue 300
	// После установки должен быть сброс number

	// Number до установки
	fmt.Println("Current number value Before set MaximumValue:", getNumber())

	// Пробуем установить новое значение
	v = setMaximumValue(300)
	if v!=nil{
		t.Error("Error setMaximumValue:", v)
	}

	//Смотрим установленное значение
	fmt.Println("Current number value After set MaximumValue:", getNumber())

	// Должно быть 0 если нет, если нет то ошибка
	//Если при смене максимального значения число резко начинает
	// превышать максимальное значение, то число надо обнулить.

	if getNumber()!=0{
		t.Error("Error reset number for change maximum value", v)
	}

}

func TestSetMaximumValueExcept1(t *testing.T) {
	v = setMaximumValue(-1)

	//Нельзя позволять установить тут число меньше нуля.
	if v!=nil{
		t.Error("Error setMaximumValue:", v)
	}

}

func TestSetMaximumValueExcept2(t *testing.T) {
	v = setMaximumValue(22147483647)
	//Нельзя позволять установить тут число больше максимального значения для int

	if v!=nil{
		t.Error("Error setMaximumValue:", v)
	}
}