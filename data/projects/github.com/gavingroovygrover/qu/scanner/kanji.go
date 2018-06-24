// Copyright 2016 Gavin "Groovy" Grover. All rights reserved.
// Use of this source code is governed by the same BSD-style
// license as Go that can be found in the LICENSE file.

package scanner

import(
)

const (
	KeywordKanji = iota + 1
	IdentifierKanji
	PackageKanji
	TentativeKanji
)

type KanjiVal struct {
	Kind           int
	Word           string
	IsSuffixable   bool
	IsScoped       bool
	IsGoReserved   bool
}

// We put all the Kanjis into a single map so we'll always know
// every Kanji has a single unique meaning.
type KanjiMap map[rune]KanjiVal

var SuffixedIdents = map[string]string{
	"整": "int",
	"整8": "int8",
	"整16": "int16",
	"整32": "int32",
	"整64": "int64",

	"绝": "uint",
	"绝8": "uint8",
	"绝16": "uint16",
	"绝32": "uint32",
	"绝64": "uint64",

	"漂32": "float32",
	"漂64": "float64",

	"复": "complex",
	"复64": "complex64",
	"复128": "complex128",
}

var Kanjis = KanjiMap {
	'包': {Kind:KeywordKanji, Word:"package", IsGoReserved:true},
	'入': {Kind:KeywordKanji, Word:"import", IsGoReserved:true},
	'久': {Kind:KeywordKanji, Word:"const", IsGoReserved:true},
	'变': {Kind:KeywordKanji, Word:"var", IsGoReserved:true},
	'种': {Kind:KeywordKanji, Word:"type", IsGoReserved:true},
	'功': {Kind:KeywordKanji, Word:"func", IsGoReserved:true},
	'构': {Kind:KeywordKanji, Word:"struct", IsGoReserved:true},
	'图': {Kind:KeywordKanji, Word:"map", IsGoReserved:true},
	'面': {Kind:KeywordKanji, Word:"interface", IsGoReserved:true},
	'通': {Kind:KeywordKanji, Word:"chan", IsGoReserved:true},

	'如': {Kind:KeywordKanji, Word:"if", IsGoReserved:true},
	'否': {Kind:KeywordKanji, Word:"else", IsGoReserved:true},
	'择': {Kind:KeywordKanji, Word:"switch", IsGoReserved:true},
	'事': {Kind:KeywordKanji, Word:"case", IsGoReserved:true},
	'别': {Kind:KeywordKanji, Word:"default", IsGoReserved:true},
	'掉': {Kind:KeywordKanji, Word:"fallthrough", IsGoReserved:true},
	'选': {Kind:KeywordKanji, Word:"select", IsGoReserved:true},
	'为': {Kind:KeywordKanji, Word:"for", IsGoReserved:true},
	'围': {Kind:KeywordKanji, Word:"range", IsGoReserved:true},
	'终': {Kind:KeywordKanji, Word:"defer", IsGoReserved:true},
	'去': {Kind:KeywordKanji, Word:"go", IsGoReserved:true},
	'回': {Kind:KeywordKanji, Word:"return", IsGoReserved:true},
	'破': {Kind:KeywordKanji, Word:"break", IsGoReserved:true},
	'继': {Kind:KeywordKanji, Word:"continue", IsGoReserved:true},
	'跳': {Kind:KeywordKanji, Word:"goto", IsGoReserved:true},

	'真': {Kind:IdentifierKanji, Word:"true", IsGoReserved:true},
	'假': {Kind:IdentifierKanji, Word:"false", IsGoReserved:true},
	'空': {Kind:IdentifierKanji, Word:"nil", IsGoReserved:true},
	'毫': {Kind:IdentifierKanji, Word:"iota", IsGoReserved:true},

	'能': {Kind:IdentifierKanji, Word:"cap", IsGoReserved:true, IsScoped:true},
	'度': {Kind:IdentifierKanji, Word:"len", IsGoReserved:true, IsScoped:true},
	'实': {Kind:IdentifierKanji, Word:"real", IsGoReserved:true, IsScoped:true},
	'虚': {Kind:IdentifierKanji, Word:"imag", IsGoReserved:true, IsScoped:true},
	'造': {Kind:IdentifierKanji, Word:"make", IsGoReserved:true, IsScoped:true},
	'新': {Kind:IdentifierKanji, Word:"new", IsGoReserved:true, IsScoped:true},
	'关': {Kind:IdentifierKanji, Word:"close", IsGoReserved:true, IsScoped:true},
	'加': {Kind:IdentifierKanji, Word:"append", IsGoReserved:true, IsScoped:true},
	'副': {Kind:IdentifierKanji, Word:"copy", IsGoReserved:true, IsScoped:true},
	'删': {Kind:IdentifierKanji, Word:"delete", IsGoReserved:true, IsScoped:true},
	'丢': {Kind:IdentifierKanji, Word:"panic", IsGoReserved:true, IsScoped:true},
	'抓': {Kind:IdentifierKanji, Word:"recover", IsGoReserved:true, IsScoped:true},
	'写': {Kind:IdentifierKanji, Word:"print", IsGoReserved:true, IsScoped:true},
	'线': {Kind:IdentifierKanji, Word:"println", IsGoReserved:true, IsScoped:true},

	'节': {Kind:IdentifierKanji, Word:"byte", IsGoReserved:true, IsScoped:true},
	'字': {Kind:IdentifierKanji, Word:"rune", IsGoReserved:true, IsScoped:true},
	'串': {Kind:IdentifierKanji, Word:"string", IsGoReserved:true, IsScoped:true},
	'双': {Kind:IdentifierKanji, Word:"bool", IsGoReserved:true, IsScoped:true},
	'错': {Kind:IdentifierKanji, Word:"error", IsGoReserved:true, IsScoped:true},
	'镇': {Kind:IdentifierKanji, Word:"uintptr", IsGoReserved:true, IsScoped:true},

	//suffixable identifiers...
	'整': {Kind:IdentifierKanji, Word:"int", IsGoReserved:true, IsScoped:true, IsSuffixable:true}, //int, int8, int16, int32, int64
	'绝': {Kind:IdentifierKanji, Word:"uint", IsGoReserved:true, IsScoped:true, IsSuffixable:true}, //uint, uint8, uint16, uint32, uint64
	'漂': {Kind:IdentifierKanji, Word:"float", IsGoReserved:true, IsScoped:true, IsSuffixable:true}, //float32, float64
	'复': {Kind:IdentifierKanji, Word:"complex", IsGoReserved:true, IsScoped:true, IsSuffixable:true}, //complex, complex64, complex128

	'让': {Kind:IdentifierKanji}, //verbalize as "let"
	'做': {Kind:IdentifierKanji}, //verbalize as "do"
	'英': {Kind:IdentifierKanji}, //verbalize as "ascii"
	'正': {Kind:IdentifierKanji, Word:"main"},
	'任': {Kind:IdentifierKanji}, //"interface{}" returned; verbalize as "any"

	//packages...
	'形': {Kind:PackageKanji, Word:"fmt"},
	'网': {Kind:PackageKanji, Word:"net"},
	'序': {Kind:PackageKanji, Word:"sort"},
	'数': {Kind:PackageKanji, Word:"math"},
	'大': {Kind:PackageKanji, Word:"math/big"},
	'时': {Kind:PackageKanji, Word:"time"},

	//tentative kanji...
	'这': {Kind:TentativeKanji, Word:"this"},
	'特': {Kind:TentativeKanji, Word:"special"},
	'愿': {Kind:TentativeKanji, Word:"source"},
	'试': {Kind:TentativeKanji, Word:"try"},
	'具': {Kind:TentativeKanji, Word:"util"},
	'动': {Kind:TentativeKanji, Word:"dyn"},
	'指': {Kind:TentativeKanji, Word:"spec"},
	'羔': {Kind:TentativeKanji, Word:"lamb"},
	'程': {Kind:TentativeKanji, Word:"proc"},
	'对': {Kind:TentativeKanji, Word:"assert"},
	'用': {Kind:TentativeKanji, Word:"use"},
	'准': {Kind:TentativeKanji, Word:"prepare"},
	'执': {Kind:TentativeKanji, Word:"execute"},
	'冲': {Kind:TentativeKanji, Word:"flush"},
	'建': {Kind:TentativeKanji, Word:"build"},
	'跑': {Kind:TentativeKanji, Word:"run"},
	'考': {Kind:TentativeKanji, Word:"test"},
	'洗': {Kind:TentativeKanji, Word:"clean"},
	'出': {Kind:TentativeKanji, Word:"exit"},
	'显': {Kind:TentativeKanji, Word:"vars"},
	'后': {Kind:TentativeKanji, Word:"next"},
	'前': {Kind:TentativeKanji, Word:"prev"},
	'学': {Kind:TentativeKanji, Word:"learn"},
	'解': {Kind:TentativeKanji, Word:"parse"},
	'类': {Kind:TentativeKanji, Word:"class"},
	'叫': {Kind:TentativeKanji, Word:"call"},
	'是': {Kind:TentativeKanji, Word:"is"},
	'侯': {Kind:TentativeKanji, Word:"while"},
	'它': {Kind:TentativeKanji, Word:"it"},
	'自': {Kind:TentativeKanji, Word:"self"},
	'滤': {Kind:TentativeKanji, Word:"filter"},
	'减': {Kind:TentativeKanji, Word:"reduce"},
	'组': {Kind:TentativeKanji, Word:"groupby"},
	'颠': {Kind:TentativeKanji, Word:"reverse"},
	'长': {Kind:TentativeKanji, Word:"long"},
	'除': {Kind:TentativeKanji, Word:"exception"},
	'摸': {Kind:TentativeKanji, Word:"pattern"},
}

var KouRadicalChars =
	`㕤㕥㕧㕨㕩㕪㕫㕬㕭㕮㕰㕱㕲㕳㕴㕵㕶㕷㕸㕹㕼㕽㖀㖁㖂㖃㖄㖅㖆㖇㖉㖊㖏㖑㖒㖓㖔㖕㖗㖘㖞㖟㖠㖡㖢㖣㖤㖥㖦㖧㖨㖩㖪㖫㖬㖭㖮㖴㖵㖶`+
	`㖷㖸㖹㖺㖻㖼㖽㖿㗀㗁㗂㗃㗄㗅㗆㗇㗈㗋㗌㗍㗎㗏㗐㗑㗒㗓㗔㗕㗖㗘㗙㗚㗛㗜㗝㗞㗢㗣㗥㗦㗧㗩㗪㗫㗭㗰㗱㗲㗳㗴㗵㗶㗷㗸㗹㗺㗻㗼㗾㗿`+
	`㘀㘁㘂㘃㘄㘅㘆㘇㘈㘉㘊㘋㘌㘍㘎㘐㘑㘓㘔㘕㘖㘗㘙㘚㘛卟叨叩叫叭叮叱叶叹叺叻叼叽叿吀吁吃吅吆吇吋吐吒吓吔吖吗吘吙吚吜吟吠吡吣`+
	`吤吥吧吨吩吪听吭吮吰吱吲吵吶吷吸吹吺吻吼吽呀呁呃呅呋呌呍呎呏呐呒呓呔呕呖呗呚呛呜呝呞呟呠呡呢呣呤呥呦呧呩呪呫呬呭呮呯呱呲`+
	`味呴呵呶呷呸呹呺呻呼呾呿咀咁咂咃咄咆咇咈咉咊咋咍咏咐咑咓咔咕咖咗咘咙咚咛咜咝咞咟咡咣咤咥咦咧咩咪咬咭咮咯咰咱咲咳咴咵咶咷`+
	`咹咺咻咽咾咿哂哃哄哅哆哇哈哊哋哌响哎哏哐哑哒哓哔哕哖哗哘哙哚哜哝哞哟哠哢哣哤哦哧哨哩哪哫哬哮哯哰哱哳哴哵哶哷哸哹哺哻哼哽`+
	`哾唀唁唂唃唄唅唆唈唉唊唋唌唍唎唏唑唒唓唔唕唖唗唙唚唛唝唞唠唡唢唣唤唥唦唧唨唩唪唫唬唭唯唰唱唲唳唴唵唶唷唸唹唺唻唼唽唾唿啀`+
	`啁啂啃啄啅啈啉啊啋啌啍啐啑啒啕啖啗啘啛啜啝啞啡啢啣啤啥啦啧啨啩啪啫啭啮啯啰啱啲啳啴啵啶啷啸啹啺啼啽啾啿喀喁喂喃喅喇喈喉喊`+
	`喋喍喎喏喐喑喒喓喔喕喖喗喘喙喚喛喝喞喟喠喡喢喣喤喥喧喨喩喫喭喯喰喱喲喳喴喵喷喹喺喻喼喽嗁嗂嗃嗄嗅嗆嗈嗉嗊嗋嗌嗍嗎嗏嗐嗑嗒`+
	`嗓嗔嗕嗖嗗嗘嗙嗚嗛嗜嗝嗞嗟嗡嗢嗤嗥嗦嗨嗩嗪嗫嗬嗮嗯嗰嗱嗲嗳嗴嗵嗶嗷嗹嗺嗻嗼嗽嗾嗿嘀嘁嘃嘄嘅嘆嘇嘈嘊嘋嘌嘍嘎嘐嘑嘒嘓嘔嘕嘖`+
	`嘘嘙嘚嘛嘜嘝嘞嘟嘠嘡嘢嘣嘤嘥嘧嘨嘩嘪嘫嘬嘭嘮嘯嘰嘱嘲嘳嘴嘵嘶嘷嘸嘹嘺嘻嘽嘾嘿噀噁噂噃噄噅噆噇噈噉噊噋噌噍噎噏噑噒噓噔噖噗`+
	`噘噙噚噛噜噝噞噠噡噢噣噤噥噦噧噪噫噬噭噮噯噰噱噲噳噴噵噶噷噸噹噺噻噼噾噿嚀嚁嚂嚃嚄嚅嚆嚇嚈嚉嚊嚋嚌嚍嚎嚏嚐嚑嚒嚓嚔嚕嚖嚗`+
	`嚘嚙嚛嚜嚝嚟嚠嚡嚤嚥嚦嚧嚨嚩嚪嚫嚬嚯嚰嚱嚵嚶嚷嚸嚹嚺嚼嚽嚾嚿囀囁囃囄囆囇囈囉囋囌囎囐囑囒囓囔囕囖鳴鸣𠮙𠮜𠮝𠮟𠮤𠮧𠮨𠮩𠮪𠮬`+
	`𠮭𠮱𠮵𠮶𠮹𠮺𠮻𠮼𠮾𠮿𠯀𠯄𠯅𠯆𠯇𠯈𠯋𠯍𠯎𠯏𠯐𠯔𠯖𠯗𠯘𠯙𠯜𠯝𠯞𠯟𠯠𠯡𠯢𠯤𠯥𠯦𠯩𠯪𠯫𠯬𠯯𠯰𠯱𠯲𠯴𠯷𠯸𠯹𠯻𠯼𠯽𠯾𠯿𠰀𠰁𠰂𠰃𠰄𠰆𠰈`+
	`𠰉𠰊𠰋𠰌𠰍𠰏𠰐𠰑𠰒𠰖𠰗𠰘𠰙𠰚𠰜𠰠𠰢𠰧𠰩𠰪𠰭𠰮𠰯𠰱𠰲𠰳𠰴𠰵𠰷𠰸𠰹𠰺𠰻𠰼𠰽𠰾𠰿𠱀𠱁𠱂𠱃𠱅𠱆𠱇𠱈𠱉𠱊𠱋𠱌𠱍𠱎𠱏𠱐𠱓𠱔𠱕𠱖𠱘𠱙𠱚`+
	`𠱜𠱝𠱞𠱟𠱠𠱡𠱢𠱣𠱤𠱥𠱨𠱪𠱱𠱲𠱳𠱴𠱶𠱷𠱸𠱹𠱺𠱻𠱼𠱽𠱾𠱿𠲂𠲃𠲄𠲅𠲇𠲈𠲊𠲋𠲌𠲍𠲎𠲏𠲐𠲓𠲔𠲕𠲖𠲗𠲙𠲚𠲛𠲜𠲝𠲞𠲟𠲠𠲡𠲢𠲣𠲤𠲥𠲦𠲧𠲨`+
	`𠲪𠲫𠲬𠲭𠲮𠲰𠲲𠲳𠲴𠲵𠲶𠲷𠲸𠲺𠲼𠲽𠲾𠲿𠳀𠳁𠳂𠳃𠳈𠳉𠳍𠳎𠳏𠳐𠳑𠳒𠳓𠳔𠳕𠳖𠳗𠳘𠳚𠳜𠳝𠳞𠳟𠳠𠳡𠳣𠳤𠳥𠳦𠳧𠳨𠳩𠳪𠳭𠳰𠳱𠳲𠳳𠳴𠳶𠳷𠳸`+
	`𠳹𠳺𠳻𠳼𠳽𠳾𠳿𠴀𠴁𠴂𠴃𠴄𠴆𠴇𠴈𠴉𠴊𠴋𠴌𠴍𠴎𠴏𠴐𠴑𠴒𠴓𠴔𠴕𠴖𠴗𠴘𠴙𠴚𠴛𠴜𠴝𠴞𠴟𠴠𠴡𠴢𠴣𠴤𠴥𠴧𠴨𠴪𠴫𠴬𠴭𠴮𠴯𠴰𠴱𠴲𠴳𠴴𠴵𠴶𠴷`+
	`𠴹𠴺𠴻𠴼𠴽𠴾𠵃𠵄𠵅𠵆𠵇𠵈𠵉𠵋𠵌𠵎𠵏𠵐𠵑𠵒𠵔𠵕𠵖𠵘𠵙𠵚𠵜𠵝𠵟𠵠𠵡𠵢𠵣𠵨𠵩𠵫𠵭𠵮𠵯𠵰𠵱𠵴𠵷𠵸𠵹𠵺𠵻𠵼𠵽𠵾𠵿𠶀𠶁𠶂𠶃𠶄𠶅𠶆𠶈𠶉`+
	`𠶊𠶋𠶌𠶍𠶎𠶏𠶐𠶑𠶒𠶓𠶔𠶕𠶖𠶗𠶙𠶚𠶛𠶜𠶝𠶞𠶟𠶠𠶡𠶢𠶣𠶤𠶥𠶦𠶧𠶨𠶩𠶪𠶫𠶭𠶯𠶲𠶴𠶸𠶹𠶺𠶻𠶼𠶽𠶾𠶿𠷀𠷁𠷂𠷃𠷄𠷅𠷆𠷇𠷈𠷉𠷊𠷋𠷌𠷍𠷐`+
	`𠷑𠷕𠷖𠷘𠷙𠷚𠷝𠷟𠷢𠷣𠷤𠷥𠷦𠷧𠷨𠷩𠷪𠷬𠷭𠷮𠷯𠷲𠷴𠷵𠷶𠷸𠷹𠷺𠷻𠷼𠷾𠷿𠸀𠸁𠸂𠸃𠸄𠸇𠸉𠸊𠸋𠸌𠸍𠸎𠸏𠸐𠸑𠸒𠸓𠸔𠸕𠸖𠸘𠸚𠸝𠸞𠸟𠸠𠸡𠸢`+
	`𠸣𠸤𠸥𠸦𠸧𠸨𠸩𠸪𠸫𠸬𠸯𠸰𠸳𠸴𠸵𠸷𠸸𠸹𠸺𠸻𠸼𠸽𠸾𠹀𠹁𠹂𠹃𠹄𠹅𠹆𠹇𠹊𠹋𠹌𠹍𠹎𠹏𠹐𠹑𠹒𠹓𠹔𠹕𠹖𠹗𠹘𠹙𠹚𠹛𠹞𠹠𠹡𠹤𠹥𠹦𠹭𠹮𠹯𠹰𠹱`+
	`𠹲𠹳𠹴𠹵𠹶𠹷𠹸𠹹𠹺𠹻𠹼𠹽𠹿𠺀𠺁𠺂𠺄𠺅𠺆𠺈𠺉𠺊𠺋𠺌𠺍𠺏𠺑𠺒𠺓𠺔𠺕𠺖𠺗𠺘𠺙𠺚𠺜𠺝𠺟𠺠𠺡𠺢𠺣𠺦𠺧𠺨𠺩𠺪𠺫𠺬𠺭𠺮𠺰𠺱𠺲𠺳𠺴𠺵𠺶𠺷`+
	`𠺸𠺹𠺺𠺻𠺼𠺽𠺾𠺿𠻀𠻂𠻃𠻄𠻅𠻆𠻈𠻉𠻊𠻋𠻍𠻎𠻏𠻐𠻑𠻒𠻓𠻔𠻕𠻗𠻘𠻙𠻛𠻜𠻞𠻟𠻠𠻢𠻣𠻤𠻥𠻦𠻧𠻨𠻩𠻪𠻫𠻬𠻯𠻱𠻲𠻳𠻴𠻵𠻶𠻷𠻹𠻺𠻻𠻼𠻽𠻾`+
	`𠻿𠼀𠼁𠼂𠼄𠼇𠼈𠼉𠼊𠼋𠼌𠼍𠼎𠼏𠼐𠼒𠼓𠼔𠼕𠼖𠼗𠼘𠼙𠼚𠼜𠼝𠼟𠼠𠼢𠼣𠼤𠼥𠼦𠼩𠼪𠼫𠼬𠼭𠼮𠼯𠼰𠼱𠼲𠼳𠼴𠼵𠼶𠼸𠼹𠼺𠼻𠼼𠼽𠼾𠽀𠽁𠽂𠽃𠽄𠽅`+
	`𠽆𠽇𠽈𠽉𠽊𠽋𠽌𠽍𠽎𠽏𠽐𠽑𠽒𠽓𠽔𠽕𠽖𠽗𠽙𠽛𠽜𠽞𠽟𠽡𠽢𠽣𠽤𠽥𠽦𠽧𠽨𠽩𠽪𠽫𠽬𠽭𠽮𠽯𠽰𠽱𠽲𠽳𠽴𠽵𠽶𠽹𠽻𠽼𠽾𠽿𠾀𠾁𠾆𠾇𠾈𠾊𠾋𠾌𠾍𠾎`+
	`𠾏𠾐𠾑𠾒𠾓𠾔𠾕𠾗𠾘𠾙𠾚𠾛𠾜𠾝𠾞𠾠𠾡𠾢𠾣𠾦𠾨𠾩𠾪𠾫𠾬𠾭𠾮𠾯𠾰𠾱𠾲𠾴𠾵𠾶𠾷𠾸𠾺𠾻𠾼𠾽𠾾𠾿𠿀𠿁𠿂𠿃𠿄𠿅𠿆𠿇𠿈𠿊𠿋𠿌𠿍𠿎𠿏𠿐𠿑𠿒`+
	`𠿓𠿔𠿖𠿗𠿘𠿙𠿚𠿛𠿜𠿝𠿞𠿠𠿢𠿣𠿤𠿥𠿨𠿩𠿪𠿫𠿬𠿭𠿮𠿯𠿰𠿱𠿳𠿴𠿵𠿶𠿷𠿸𠿹𠿺𠿼𠿾𠿿𡀀𡀁𡀂𡀃𡀄𡀅𡀇𡀊𡀌𡀍𡀎𡀏𡀐𡀑𡀔𡀕𡀖𡀗𡀘𡀙𡀚𡀛𡀜`+
	`𡀝𡀞𡀟𡀠𡀡𡀢𡀣𡀥𡀦𡀧𡀨𡀩𡀫𡀬𡀭𡀮𡀯𡀰𡀱𡀲𡀳𡀴𡀵𡀶𡀷𡀹𡀺𡀼𡀽𡀾𡀿𡁀𡁁𡁂𡁃𡁄𡁅𡁆𡁇𡁈𡁊𡁋𡁌𡁍𡁎𡁏𡁐𡁑𡁒𡁓𡁔𡁕𡁖𡁙𡁚𡁛𡁜𡁝𡁞𡁟`+
	`𡁠𡁡𡁣𡁤𡁦𡁧𡁪𡁫𡁬𡁭𡁮𡁯𡁰𡁱𡁲𡁴𡁵𡁶𡁷𡁸𡁹𡁺𡁻𡁼𡁽𡁾𡁿𡂀𡂁𡂂𡂃𡂄𡂅𡂆𡂈𡂉𡂊𡂋𡂌𡂍𡂎𡂏𡂐𡂑𡂒𡂓𡂔𡂕𡂖𡂗𡂘𡂙𡂚𡂛𡂜𡂝𡂞𡂠𡂡𡂢`+
	`𡂣𡂥𡂩𡂪𡂫𡂭𡂮𡂰𡂱𡂳𡂴𡂵𡂷𡂸𡂹𡂺𡂻𡂼𡂿𡃀𡃁𡃂𡃃𡃄𡃅𡃆𡃇𡃈𡃉𡃊𡃌𡃍𡃎𡃏𡃐𡃑𡃒𡃓𡃔𡃕𡃖𡃗𡃘𡃙𡃚𡃛𡃜𡃝𡃞𡃢𡃤𡃥𡃦𡃧𡃨𡃩𡃪𡃮𡃰𡃱`+
	`𡃲𡃳𡃴𡃵𡃶𡃹𡃺𡃻𡃼𡃽𡃾𡃿𡄁𡄃𡄄𡄆𡄇𡄊𡄋𡄍𡄎𡄏𡄐𡄑𡄓𡄔𡄕𡄖𡄗𡄘𡄙𡄟𡄠𡄡𡄢𡄣𡄤𡄥𡄦𡄧𡄨𡄩𡄪𡄫𡄭𡄮𡄯𡄱𡄳𡄴𡄵𡄶𡄷𡄸𡄺𡄼𡄽𡄾𡅁𡅂`+
	`𡅃𡅅𡅆𡅇𡅈𡅉𡅊𡅋𡅌𡅍𡅎𡅏𡅑𡅒𡅓𡅗𡅘𡅙𡅛𡅜𡅞𡅠𡅢𡅣𡅥𡅧𡅨𡅩𡅪𡅫𡅬𡅭𡅯𡅰𡅲𡅳𡅵𡅶𡅷𡅹𡅺𡅼𡅿𡆀𡆁𡆂𡆄𡆅𡆆𡆇𡆈𡆋𡆌𡆍𡆏𡆑𡆓𡆕𡆖𡆗`+
	`𡆘𡆙𡆚𡆜𡆝𡆞𡆟𢒯𧛧𨙫𩐉𩒻𪄨𪄼𪚩𪠳𪠴𪠵𪠶𪠸𪠺𪠻𪠼𪠽𪠾𪠿𪡀𪡁𪡂𪡃𪡄𪡆𪡇𪡈𪡊𪡋𪡏𪡓𪡔𪡕𪡗𪡙𪡚𪡛𪡝𪡞𪡟𪡠𪡡𪡣𪡥𪡦𪡧𪡨𪡩𪡫𪡭𪡮𪡱𪡴`+
	`𪡵𪡶𪡷𪡸𪡺𪡽𪡾𪡿𪢀𪢁𪢂𪢃𪢄𪢅𪢆𪢇𪢉𪢊𪢋𪢌𪢍𪢎𪢐𪢑𪢒𪢔𪢕𪢖𪢗𪢘𪢙𪢚𪢜𪢝𪢟𪢠𪢤𪢥𪢧𫛗𫝘𫝚𫝜𫝞`
	//the last few lhs kou-radical chars are CJK ext-D chars that often display incorrectly in the Windows fonts

