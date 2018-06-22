// +build ignore

package main

import (
	"os"
	"os/exec"
	"regexp"
	"regexp/syntax"
	"strconv"
	"strings"
	"text/template"
)

func main() {
	f, err := os.Create("regexp.go")
	if err != nil {
		panic(err)
	}
	fileTemplate.Execute(f, genRegexps())
	exec.Command("gofmt", "-s", "-l", "-w", "regexp.go").Run()
}

var fileTemplate = template.Must(template.New("file").Parse(`package text

import "regexp"

// Do not modify this file, generate it with 'go run gen_regexp.go'
var (
{{range $name, $data := .}}{{$name}} = regexp.MustCompile({{$data}})
{{end}})`))

func genRegexps() map[string]string {
	charGroups := make(map[string][]string)
	regexen := make(map[string]string)
	res := make(map[string]string)

	add := func(p string, r ...rune) {
		if len(r) == 1 {
			charGroups[p] = append(charGroups[p], string(r[0]))
		} else {
			charGroups[p] = append(charGroups[p], string(r[0])+"-"+string(r[1]))
		}
	}

	interpolatePattern := regexp.MustCompile(`#\{\w+\}`)
	interp := func(s string) string {
		replace := func(g string) string {
			g = g[2 : len(g)-1]
			if regex, ok := regexen[g]; ok {
				return regex
			}
			group, ok := charGroups[g]
			if !ok {
				panic(s + ": unknown group or regexen " + g)
			}
			return strings.Join(group, "")
		}
		return interpolatePattern.ReplaceAllStringFunc(s, replace)
	}

	pattern := func(name string, data string) {
		r, err := syntax.Parse(interp(data), syntax.Perl)
		if err != nil {
			panic(err)
		}
		res[name] = strconv.QuoteToASCII(r.Simplify().String())
	}

	charGroups["spaces"] = []string{
		"\u0020", // White_Space # Zs SPACE
		"\u0085", // White_Space # Cc <control-0085>
		"\u00A0", // White_Space # Zs NO-BREAK SPACE
		"\u1680", // White_Space # Zs OGHAM SPACE MARK
		"\u180E", // White_Space # Zs MONGOLIAN VOWEL SEPARATOR
		"\u2028", // White_Space # Zl LINE SEPARATOR
		"\u2029", // White_Space # Zp PARAGRAPH SEPARATOR
		"\u202F", // White_Space # Zs NARROW NO-BREAK SPACE
		"\u205F", // White_Space # Zs MEDIUM MATHEMATICAL SPACE
		"\u3000", // White_Space # Zs IDEOGRAPHIC SPACE
	}
	add("spaces", '\u0009', '\u000D') // White_Space # Cc [5] <control-0009>..<control-000D>
	add("spaces", '\u2000', '\u200A') // White_Space # Zs [11] EN QUAD..HAIR SPACE

	charGroups["invalid"] = []string{
		"\uFFFE",
		"\uFEFF", // BOM
		"\uFFFF", // Special
	}
	add("invalid", '\u202A', '\u202E') // Directional change

	add("nonLatinHashtag", '\u0400', '\u04ff') // Cyrillic
	add("nonLatinHashtag", '\u0500', '\u0527') // Cyrillic Supplement
	add("nonLatinHashtag", '\u2de0', '\u2dff') // Cyrillic Extended A
	add("nonLatinHashtag", '\ua640', '\ua69f') // Cyrillic Extended B
	add("nonLatinHashtag", '\u0591', '\u05bf') // Hebrew
	add("nonLatinHashtag", '\u05c1', '\u05c2')
	add("nonLatinHashtag", '\u05c4', '\u05c5')
	add("nonLatinHashtag", '\u05c7')
	add("nonLatinHashtag", '\u05d0', '\u05ea')
	add("nonLatinHashtag", '\u05f0', '\u05f4')
	add("nonLatinHashtag", '\ufb12', '\ufb28') // Hebrew Presentation Forms
	add("nonLatinHashtag", '\ufb2a', '\ufb36')
	add("nonLatinHashtag", '\ufb38', '\ufb3c')
	add("nonLatinHashtag", '\ufb3e')
	add("nonLatinHashtag", '\ufb40', '\ufb41')
	add("nonLatinHashtag", '\ufb43', '\ufb44')
	add("nonLatinHashtag", '\ufb46', '\ufb4f')
	add("nonLatinHashtag", '\u0610', '\u061a') // Arabic
	add("nonLatinHashtag", '\u0620', '\u065f')
	add("nonLatinHashtag", '\u066e', '\u06d3')
	add("nonLatinHashtag", '\u06d5', '\u06dc')
	add("nonLatinHashtag", '\u06de', '\u06e8')
	add("nonLatinHashtag", '\u06ea', '\u06ef')
	add("nonLatinHashtag", '\u06fa', '\u06fc')
	add("nonLatinHashtag", '\u06ff')
	add("nonLatinHashtag", '\u0750', '\u077f') // Arabic Supplement
	add("nonLatinHashtag", '\u08a0')           // Arabic Extended A
	add("nonLatinHashtag", '\u08a2', '\u08ac')
	add("nonLatinHashtag", '\u08e4', '\u08fe')
	add("nonLatinHashtag", '\ufb50', '\ufbb1') // Arabic Pres. Forms A
	add("nonLatinHashtag", '\ufbd3', '\ufd3d')
	add("nonLatinHashtag", '\ufd50', '\ufd8f')
	add("nonLatinHashtag", '\ufd92', '\ufdc7')
	add("nonLatinHashtag", '\ufdf0', '\ufdfb')
	add("nonLatinHashtag", '\ufe70', '\ufe74') // Arabic Pres. Forms B
	add("nonLatinHashtag", '\ufe76', '\ufefc')
	add("nonLatinHashtag", '\u200c')           // Zero-Width Non-Joiner
	add("nonLatinHashtag", '\u0e01', '\u0e3a') // Thai
	add("nonLatinHashtag", '\u0e40', '\u0e4e') // Hangul (Korean)
	add("nonLatinHashtag", '\u1100', '\u11ff') // Hangul Jamo
	add("nonLatinHashtag", '\u3130', '\u3185') // Hangul Compatibility Jamo
	add("nonLatinHashtag", '\uA960', '\uA97F') // Hangul Jamo Extended-A
	add("nonLatinHashtag", '\uAC00', '\uD7AF') // Hangul Syllables
	add("nonLatinHashtag", '\uD7B0', '\uD7FF') // Hangul Jamo Extended-B
	add("nonLatinHashtag", '\uFFA1', '\uFFDC') // half-width Hangul
	// Japanese and Chinese
	add("nonLatinHashtag", '\u30A1', '\u30FA')         // Katakana (full-width)
	add("nonLatinHashtag", '\u30FC', '\u30FE')         // Katakana Chouon and iteration marks (full-width)
	add("nonLatinHashtag", '\uFF66', '\uFF9F')         // Katakana (half-width)
	add("nonLatinHashtag", '\uFF70', '\uFF70')         // Katakana Chouon (half-width)
	add("nonLatinHashtag", '\uFF10', '\uFF19')         // \
	add("nonLatinHashtag", '\uFF21', '\uFF3A')         // - Latin (full-width)
	add("nonLatinHashtag", '\uFF41', '\uFF5A')         // /
	add("nonLatinHashtag", '\u3041', '\u3096')         // Hiragana
	add("nonLatinHashtag", '\u3099', '\u309E')         // Hiragana voicing and iteration mark
	add("nonLatinHashtag", '\u3400', '\u4DBF')         // Kanji (CJK Extension A)
	add("nonLatinHashtag", '\u4E00', '\u9FFF')         // Kanji (Unified)
	add("nonLatinHashtag", '\U0002A700', '\U0002B73F') // Kanji (CJK Extension C)
	add("nonLatinHashtag", '\U0002B740', '\U0002B81F') // Kanji (CJK Extension D)
	add("nonLatinHashtag", '\U0002F800', '\U0002FA1F') // Kanji (CJK supplement)
	add("nonLatinHashtag", '\u3003')                   // Kanji iteration mark
	add("nonLatinHashtag", '\u3005')                   // Kanji iteration mark
	add("nonLatinHashtag", '\u303B')                   // Han iteration mark

	add("latinAccent", '\u00c0', '\u00d6')
	add("latinAccent", '\u00d8', '\u00f6')
	add("latinAccent", '\u00f8', '\u00ff')
	// Latin Extended A and B
	add("latinAccent", '\u0100', '\u024f')
	// assorted IPA Extensions
	add("latinAccent", '\u0253', '\u0254')
	add("latinAccent", '\u0256', '\u0257')
	add("latinAccent", '\u0259')
	add("latinAccent", '\u025b')
	add("latinAccent", '\u0263')
	add("latinAccent", '\u0268')
	add("latinAccent", '\u026f')
	add("latinAccent", '\u0272')
	add("latinAccent", '\u0289')
	add("latinAccent", '\u028b')
	// Okina for Hawaiian (it *is* a letter character)
	add("latinAccent", '\u02bb')
	// Combining diacritics
	add("latinAccent", '\u0300', '\u036f')
	// Latin Extended Additional
	add("latinAccent", '\u1e00', '\u1eff')

	// Hashtags
	regexen["punct"] = `\!'#%&'\(\)*\+,\\\-\./:;<=>\?@\[\]\^_{|}~\$`
	regexen["hashSigns"] = "[#＃]"
	regexen["hashtagAlpha"] = interp("[a-z_#{latinAccent}#{nonLatinHashtag}]")
	regexen["hashtagAlphaNumeric"] = interp("[a-z0-9_#{latinAccent}#{nonLatinHashtag}]")
	regexen["hashtagBoundary"] = interp("(?:^|$|[^&a-z0-9_#{latinAccent}#{nonLatinHashtag}])")

	pattern("invalidHashtagEnd", `^(?:#{hashSigns}|://)`)
	pattern("hashtagPattern", "(?im)(?:#{hashtagBoundary})((?:#{hashSigns})(#{hashtagAlphaNumeric}*#{hashtagAlpha}#{hashtagAlphaNumeric}*))")

	// URLs
	regexen["urlPreceding"] = interp("[^A-Za-z0-9@＠$#＃#{invalid}]|^")
	regexen["inUrlWithoutProtocolPreceding"] = `[-_./]$`
	regexen["domainChars"] = interp("[^#{punct}#{spaces}#{invalid}]")
	regexen["subdomain"] = interp(`(?:(?:#{domainChars}(?:[_-]|#{domainChars})*)?#{domainChars}\.)`)
	regexen["domainName"] = interp(`(?:(?:#{domainChars}(?:-|#{domainChars})*)?#{domainChars}\.)`)
	regexen["GTLD"] = "(?:academy|accountants|actor|aero|agency|airforce|archi|arpa|asia|associates|axa|bar|bargains|bayern|berlin|best|bid|bike|biz|black|blackfriday|blue|boutique|build|builders|buzz|cab|camera|camp|capital|cards|care|career|careers|cash|cat|catering|center|ceo|cheap|christmas|citic|claims|cleaning|clinic|clothing|club|codes|coffee|college|cologne|com|community|company|computer|condos|construction|consulting|contractors|cooking|cool|coop|country|credit|creditcard|cruises|dance|dating|democrat|dental|desi|diamonds|digital|directory|discount|dnp|domains|edu|education|email|engineering|enterprises|equipment|estate|eus|events|exchange|expert|exposed|fail|farm|feedback|finance|financial|fish|fishing|fitness|flights|florist|foo|foundation|frogans|fund|furniture|futbol|gal|gallery|gift|glass|globo|gmo|gop|gov|graphics|gratis|gripe|guitars|guru|haus|holdings|holiday|horse|house|immobilien|industries|info|ink|institute|insure|int|international|investments|jetzt|jobs|kaufen|kim|kitchen|kiwi|koeln|kred|land|lease|lighting|limited|limo|link|london|luxury|maison|management|mango|marketing|media|meet|menu|miami|mil|mobi|moda|moe|monash|moscow|museum|nagoya|name|net|neustar|ninja|nyc|okinawa|onl|org|paris|partners|parts|photo|photography|photos|pics|pictures|pink|plumbing|post|pro|productions|properties|pub|qpon|quebec|recipes|red|reisen|ren|rentals|repair|report|rest|reviews|rich|rocks|rodeo|ruhr|ryukyu|saarland|schule|services|sexy|shiksha|shoes|singles|social|sohu|solar|solutions|soy|supplies|supply|support|surgery|systems|tattoo|tax|technology|tel|tienda|tips|today|tokyo|tools|town|toys|trade|training|travel|university|uno|vacations|vegas|ventures|viajes|villas|vision|vodka|vote|voting|voto|voyage|wang|watch|webcam|wed|wien|wiki|works|wtc|wtf|xxx|xyz|yokohama|zone)"
	regexen["CCTLD"] = "(?:ac|ad|ae|af|ag|ai|al|am|an|ao|aq|ar|as|at|au|aw|ax|az|ba|bb|bd|be|bf|bg|bh|bi|bj|bm|bn|bo|br|bs|bt|bv|bw|by|bz|ca|cc|cd|cf|cg|ch|ci|ck|cl|cm|cn|co|cr|cu|cv|cw|cx|cy|cz|de|dj|dk|dm|do|dz|ec|ee|eg|er|es|et|eu|fi|fj|fk|fm|fo|fr|ga|gb|gd|ge|gf|gg|gh|gi|gl|gm|gn|gp|gq|gr|gs|gt|gu|gw|gy|hk|hm|hn|hr|ht|hu|id|ie|il|im|in|io|iq|ir|is|it|je|jm|jo|jp|ke|kg|kh|ki|km|kn|kp|kr|kw|ky|kz|la|lb|lc|li|lk|lr|ls|lt|lu|lv|ly|ma|mc|md|me|mg|mh|mk|ml|mm|mn|mo|mp|mq|mr|ms|mt|mu|mv|mw|mx|my|mz|na|nc|ne|nf|ng|ni|nl|no|np|nr|nu|nz|om|pa|pe|pf|pg|ph|pk|pl|pm|pn|pr|ps|pt|pw|py|qa|re|ro|rs|ru|rw|sa|sb|sc|sd|se|sg|sh|si|sj|sk|sl|sm|sn|so|sr|st|su|sv|sx|sy|sz|tc|td|tf|tg|th|tj|tk|tl|tm|tn|to|tp|tr|tt|tv|tw|tz|ua|ug|uk|us|uy|uz|va|vc|ve|vg|vi|vn|vu|wf|ws|ye|yt|za|zm|zw)"
	regexen["IDNTLD"] = "(?:集团|在线|한국|ভারত|公益|公司|移动|我爱你|москва|қаз|онлайн|сайт|срб|орг|삼성|சிங்கப்பூர்|商城|дети|中文网|中信|中国|中國|భారత్|ලංකා|ભારત|भारत|संगठन|网络|укр|香港|台湾|台灣|мон|الجزائر|عمان|ایران|امارات|بازار|الاردن|بھارت|المغرب|السعودية|مليسيا|شبكة|机构|组织机构|ไทย|سورية|рф|تونس|みんな|世界|ਭਾਰਤ|网址|游戏|مصر|قطر|இலங்கை|இந்தியா|新加坡|فلسطين|政务)"
	regexen["punycode"] = "(?:xn--[0-9a-z]+)"
	regexen["domain"] = interp("(?:#{subdomain}*#{domainName}(?:#{GTLD}|#{CCTLD}|#{IDNTLD}|#{punycode}))")

	pattern("asciiDomain", `(?i)(?:(?:[\-a-z0-9#{latinAccent}]+)\.)+(?:#{GTLD}|#{CCTLD}|#{IDNTLD}|#{punycode})`)
	pattern("invalidBeforeDomain", `[-_./]$`)
	pattern("invalidWithoutPath", "(?i)^#{domainName}#{CCTLD}$")

	regexen["portNumber"] = "[0-9]+"
	regexen["generalUrlPath"] = interp(`[a-z0-9!\*';:=\+,\.\$\/%#\[\]\-_~@|&#{latinAccent}]`)
	regexen["urlBalancedParens"] = interp(`\(#{generalUrlPath}+\)`)
	regexen["urlPathEnding"] = interp(`[\+\-a-z0-9=_#\/#{latinAccent}]|(?:#{urlBalancedParens})`)
	regexen["urlPath"] = interp(`(?:(?:#{generalUrlPath}*(?:#{urlBalancedParens}#{generalUrlPath}*)*#{urlPathEnding})|(?:@#{generalUrlPath}+/))`)
	regexen["urlQuery"] = `[a-z0-9!?\*'@\(\);:&=\+\$/%#\[\]\-_\.,~|]`
	regexen["urlQueryEnding"] = `[a-z0-9_&=#/]`

	pattern("urlPattern", `(?im)(#{urlPreceding})`+ // [1] Preceding character
		`(`+ // [2] URL
		`(https?://)?`+ // [3] Protocol
		`(#{domain})`+ // [4] Domain
		`(?::#{portNumber})?`+ // Port number
		`(/#{urlPath}*)?`+ // [4] URL Path
		`(?:\?#{urlQuery}*#{urlQueryEnding})?)`) // Query String

	return res
}
