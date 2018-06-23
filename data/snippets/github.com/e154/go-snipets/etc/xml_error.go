package main

import (
	"encoding/xml"
	"fmt"
)

const (
	err_str = `<response><pg_status>error</pg_status><pg_error_code>8016</pg_error_code><pg_error_description>Платежная система не разрешена для магазина, обратитесь к администрации магазина.</pg_error_description></response>`

	card_list = `<?xml version="1.0" encoding="UTF-8"?>
		<response>
		<card>
		<pg_status>success</pg_status>
		<pg_merchant_id>1111</pg_merchant_id>
		<pg_card_id>123456789</pg_card_id>
		<pg_card_hash>414621-XX-XXXX-1219</pg_card_hash>
		<pg_salt>bYFniw1jRJD2pGZg</pg_salt>
		<created_at>2017-08-22 11:17:06</created_at>
		</card>
		<card>
		<pg_status>success</pg_status>
		<pg_merchant_id>1111</pg_merchant_id>
		<pg_card_id>123456789</pg_card_id>
		<pg_card_hash>414621-XX-XXXX-1219</pg_card_hash>
		<pg_salt>bYFniw1jRJD2pGZg</pg_salt>
		<created_at>2017-08-22 11:17:06</created_at>
		</card>
		</response>`
)

type Card struct {
	PgStatus			string		`xml:"pg_status"`				// Статус запроса.
	PgMerchantId		int64		`xml:"pg_merchant_id"`			// ID мерчанта в системе Paybox
	PgCardId			int64		`xml:"pg_card_id"`				// ID сохраненной карты
	PgCardHash			string		`xml:"pg_card_hash"`			// Маскированный PAN карты
	PgSalt				string		`xml:"pg_salt"`					// Случайная строка, состоящая из произвольных цифр и латинских букв.
	CreatedAt			string		`xml:"created_at"`				// Дата привязки
}

type CardList struct {
	Card				[]Card 		`xml:"card"` 					// Список привязанных карт
}

type Error struct {
	PgStatus			string		`xml:"pg_status"`
	PgErrorCode			int			`xml:"pg_error_code"`
	PgErrorDescription	string		`xml:"pg_error_description"`
	PgSalt				string		`xml:"pg_salt"`
	PgSig				string		`xml:"pg_sig"`
}

func main() {

	errObj := &Error{}
	if err := xml.Unmarshal([]byte(err_str), errObj); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", errObj)


	cardlist := &CardList{}
	if err := xml.Unmarshal([]byte(card_list), cardlist); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", cardlist)


	a := []int{1,2,3,4}
	i := 3
	for k,v := range a {
		fmt.Println(k,v)
	}
	a = append(a[:i], a[i+1:]...)
	fmt.Println(a)
}