package main

import (
	"encoding/xml"
	"fmt"
	"net/url"
	"strings"
)

const (
	pg_xml = `pg_xml%3D%3C%3Fxml%20version%3D%221.0%22%20encoding%3D%22utf-8%22%3F%3E%0A%3Cresponse%3E%3Cpg_payment_id%3E4964603%3C%2Fpg_payment_id%3E%3Cpg_order_id%2F%3E%3Cpg_card_hash%3E4003-03XX-XXXX-5378%3C%2Fpg_card_hash%3E%3Cpg_card_hhash%3Efdf8829f47c1e789538fbf3435087b02%3C%2Fpg_card_hhash%3E%3Cpg_card_month%3E12%3C%2Fpg_card_month%3E%3Cpg_card_year%3E25%3C%2Fpg_card_year%3E%3Cpg_bank%3E%D0%90%D0%9E%20%C2%ABQAZQOM%C2%BB%3C%2Fpg_bank%3E%3Cpg_country%3EKZ%3C%2Fpg_country%3E%3Cpg_card_3ds%3E1%3C%2Fpg_card_3ds%3E%3Cpg_status%3Esuccess%3C%2Fpg_status%3E%3Cpg_card_id%3E1625984%3C%2Fpg_card_id%3E%3Cpg_user_id%3E59e47fb8e368b91c869bb286%3C%2Fpg_user_id%3E%3Cpg_recurring_profile_id%3E15095%3C%2Fpg_recurring_profile_id%3E%3Cpg_type%3Eapprove%3C%2Fpg_type%3E%3Cpg_sig%3E4384ebe74a67df858476f6a73fc5511f%3C%2Fpg_sig%3E%3C%2Fresponse%3E%0A`
)

type ResponseBind struct {
	PgStatus			string		`xml:"pg_status"`				// Статус запроса.
	PgErrorCode			int64		`xml:"pg_error_code"`			// Код ошибки
	PgErrorDescription	string		`xml:"pg_error_description"`	// Текстовое описание ошибки
	PgSalt				string		`xml:"pg_salt"`					// Случайная строка, состоящая из произвольных цифр и латинских букв.
	PgSig				string		`xml:"pg_sig"`					// Подпись запроса
	PgOrderId			int64		`xml:"pg_order_id"`				// ID заказа в системе мерчанта
	PgMerchantId		int64		`xml:"pg_merchant_id"`			// ID мерчанта в системе Paybox
	PgRedirectUrl		string		`xml:"pg_redirect_url"`			// URL который необходимо вызвать в iframe или произвести перенаправление (redirect)
	PgPaymentId			int64		`xml:"pg_payment_id"`			// ID транзакции в системе Paybox
	PgType				string		`xml:"pg_type"`					// Тип транзакции
	PgCard3ds			int64		`xml:"pg_card_3ds"`				// Признак наличия 3ds на карте
	PgCardHash			string		`xml:"pg_card_hash"`			// Маскированный PAN карты
	PgCardId			int64		`xml:"pg_card_id"`				// ID сохраненной карты
	PgUserId			string		`xml:"pg_user_id"`				// ID пользователя в системе мерчанта
	CreatedAt			string		`xml:"created_at"`				// Дата привязки
}

func main() {

	u, err := url.Parse(pg_xml)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if !strings.Contains(u.Path, "pg_xml") {
		fmt.Println("bad response")
		return
	}

	r := &ResponseBind{}
	if err := xml.Unmarshal([]byte(strings.Replace(u.Path, "pg_xml=", "", 1)), r); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Printf("%v\n", r)
}
