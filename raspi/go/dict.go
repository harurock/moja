package main

type PrefId int

const (
	PrefInvalid PrefId = iota
	PrefHokkaido
	PrefAomori
	PrefIwate
	PrefMiyagi
	PrefAkita
	PrefYamagata
	PrefFukushima
	PrefIbaraki
	PrefTochigi
	PrefGunma
	PrefSaitama
	PrefChiba
	PrefTokyo
	PrefKanagawa
	PrefNiigata
	PrefToyama
	PrefIshikawa
	PrefFukui
	PrefYamanashi
	PrefNagano
	PrefGifu
	PrefShizuoka
	PrefAichi
	PrefMie
	PrefShiga
	PrefKyoto
	PrefOsaka
	PrefHyogo
	PrefNara
	PrefWakayama
	PrefTottori
	PrefShimane
	PrefOkayama
	PrefHiroshima
	PrefYamaguchi
	PrefTokushima
	PrefKagawa
	PrefEhime
	PrefKochi
	PrefFukuoka
	PrefSaga
	PrefNagasaki
	PrefKumamoto
	PrefOita
	PrefMiyazaki
	PrefKagoshima
	PrefOkinawa
)

var (
	PrefDict map[string]PrefId
)

func init() {
	PrefDict = map[string]PrefId{
		"北海道":  PrefHokkaido,
		"青森":   PrefAomori,
		"岩手":   PrefIwate,
		"宮城":   PrefMiyagi,
		"秋田":   PrefAkita,
		"山形":   PrefYamagata,
		"福島":   PrefFukushima,
		"茨城":   PrefIbaraki,
		"栃木":   PrefTochigi,
		"群馬":   PrefGunma,
		"埼玉":   PrefSaitama,
		"千葉":   PrefChiba,
		"東京":   PrefTokyo,
		"秋葉原":  PrefTokyo,
		"秋葉":   PrefTokyo,
		"池袋":   PrefTokyo,
		"練馬区":  PrefTokyo,
		"練馬":   PrefTokyo,
		"板橋区":  PrefTokyo,
		"板橋":   PrefTokyo,
		"北区":   PrefTokyo,
		"足立区":  PrefTokyo,
		"葛飾区":  PrefTokyo,
		"葛飾":   PrefTokyo,
		"江戸川区": PrefTokyo,
		"荒川区":  PrefTokyo,
		"台東区":  PrefTokyo,
		"墨田区":  PrefTokyo,
		"江東区":  PrefTokyo,
		"杉並区":  PrefTokyo,
		"中野区":  PrefTokyo,
		"豊島区":  PrefTokyo,
		"新宿区":  PrefTokyo,
		"新宿":   PrefTokyo,
		"文京区":  PrefTokyo,
		"千代田区": PrefTokyo,
		"中央区":  PrefTokyo,
		"世田谷区": PrefTokyo,
		"世田谷":  PrefTokyo,
		"渋谷区":  PrefTokyo,
		"渋谷":   PrefTokyo,
		"港区":   PrefTokyo,
		"目黒区":  PrefTokyo,
		"品川区":  PrefTokyo,
		"品川":   PrefTokyo,
		"大田区":  PrefTokyo,
		"神奈川":  PrefKanagawa,
		"新潟":   PrefNiigata,
		"富山":   PrefToyama,
		"石川":   PrefIshikawa,
		"福井":   PrefFukui,
		"山梨":   PrefYamanashi,
		"長野":   PrefNagano,
		"岐阜":   PrefGifu,
		"静岡":   PrefShizuoka,
		"愛知":   PrefAichi,
		"三重":   PrefMie,
		"滋賀":   PrefShiga,
		"京都":   PrefKyoto,
		"大阪":   PrefOsaka,
		"枚方":   PrefOsaka,
		"阪大":   PrefOsaka,
		"兵庫":   PrefHyogo,
		"奈良":   PrefNara,
		"和歌山":  PrefWakayama,
		"鳥取":   PrefTottori,
		"島根":   PrefShimane,
		"岡山":   PrefOkayama,
		"倉敷":   PrefOkayama,
		"広島":   PrefHiroshima,
		"山口":   PrefYamaguchi,
		"徳島":   PrefTokushima,
		"香川":   PrefKagawa,
		"愛媛":   PrefEhime,
		"高知":   PrefKochi,
		"福岡":   PrefFukuoka,
		"佐賀":   PrefSaga,
		"長崎":   PrefNagasaki,
		"熊本":   PrefKumamoto,
		"大分":   PrefOita,
		"宮崎":   PrefMiyazaki,
		"鹿児島":  PrefKagoshima,
		"沖縄":   PrefOkinawa,
		"那覇":   PrefOkinawa,
	}
}
