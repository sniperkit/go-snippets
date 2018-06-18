package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	str := "275 / 45 R19 108 Y Стандартные"
	reg := regexp.MustCompile(`([0-9a-zA-ZСтандарт+])+`)
	tire := reg.FindAllString(str, -1)

	nums := make([]int, len(tire))
	nums[0], _ = strconv.Atoi(tire[0])
	nums[1], _ = strconv.Atoi(tire[1])
	nums[3], _ = strconv.Atoi(tire[3])

	reg = regexp.MustCompile(`([0-9])+`)
	rim := reg.FindAllString(tire[2], 1)
	nums[2], _ = strconv.Atoi(string(rim[0]))


	fmt.Printf("width: %d\n", nums[0])
	fmt.Printf("height: %d\n", nums[1])
	fmt.Printf("rim: %d\n", nums[2])
	fmt.Printf("load index: %d\n", nums[3])
	fmt.Printf("speed index: %s\n", tire[4])
	// standart
	fmt.Printf("standart: %s\n", tire[5])
	if strings.Contains(str, "Стандарт") {
		fmt.Printf("ok\n")
	} else {
		fmt.Printf("not ok\n")
	}

	//
	version := "A6 Quattro 2010->0"
//	version := "A6 Quattro 2010->2012"
//	version := "A6 Quattro"
	verArr := regexp.MustCompile(`([0-9]+)\-\>([0-9]+)`).FindAllString(version, -1)

	if(len(verArr)>0) {
		fmt.Println(verArr[0])
		str := strings.Split(verArr[0], "->")
		start, _ := strconv.Atoi(str[0])
		end, _ := strconv.Atoi(str[1])
		fmt.Printf("start: %d\n", start)
		fmt.Printf("end: %d\n", end)
	}

	//
	str = "275/45R19"
	reg = regexp.MustCompile(`([0-9])+`)
	tire = reg.FindAllString(str, -1)
	fmt.Println(tire)


	//
	// string pattern
	re := regexp.MustCompile(`api.contact.[0-9]+.[0-9]+.name.all`)
//	re2 := regexp.MustCompile(`api/contact/([0-9]+)/([0-9]+)/name/all`)

	fmt.Print("\n")
	fmt.Print(re.FindAllString("api.contact.10.1.name.all", -1))
	fmt.Print("\n")
//	fmt.Print("\n")
//	fmt.Printf("%q\n", re2.FindString("api/contact/10/1/name/all"))

	matched, err := regexp.MatchString("foo.*", "seafood")
	fmt.Println(matched, err)
	fmt.Println(regexp.MatchString("(api.contact.[0-9]+.[0-9]+.name.all):(get)", "api.contact.10.1.name.all:get"))
	fmt.Println(regexp.MatchString("(api.contact.[0-9]+.[0-9]+.name.all):(get)", "api.contact.10.1.name.all:post"))
	fmt.Println(regexp.MatchString("(api.contact.[0-9]+.[0-9]+.name.all):(get)", "api.contact.10.1.name.all"))
	fmt.Println(regexp.MatchString("(api.contact.[0-9]+.[0-9]+.name.all):(get)", "api.contact.10.1.name:get"))
	fmt.Println(regexp.MatchString("api.contact.[0-9]+.[0-9]+.name.all:get", "api.contact.10.1.name.all:get"))
	fmt.Println(regexp.MatchString("/api/post/[0-9//]*:get", "/api/post/:get"))
	fmt.Println(regexp.MatchString("/api/catalog/tire/[0-9]+/[0-9]+/[a-z_-]+:get", "/api/catalog/tire/25/1/-create_time:get"))

	// поиск и вывод
	s := `get published first 10 entries from 'markdown' category`
	fmt.Println(regexp.MatchString("get [a-z]+ [a-z]+ [0-9]+ entries from '[a-z]+' category", s))

	entries := regexp.MustCompile(`([0-9]+) entries`)
	num := entries.FindStringSubmatch(s)
	fmt.Println(num[1])

	category := regexp.MustCompile(`from '(.+)?' category`)
	cat := category.FindStringSubmatch(s)
	fmt.Println(cat[1])

	command := strings.Split(regexp.MustCompile(`(.+)?`).FindAllString(s, -1)[0], " ")
	fmt.Println(command)

	s2 := `get by id 1,2,3,4,5`
	fmt.Println(regexp.MatchString("get by id [0-9]+", s2))
	ids := regexp.MustCompile(`([0-9]+)`).FindAllString(s2, -1)
	fmt.Println(ids[3])

	text := "Edit the Expression & Text to see matches. Roll over matches or the expression for details. Undo mistakes with ctrl-z. Save & Share expressions [qwe:content] with friends or the Community. A full Reference & Help is available in the Library, or watch the [qqqw:content] video Tutorial."
	reg = regexp.MustCompile(`(\[{1}[a-z]{2,64}\:content\]{1})`)
	markers := reg.FindAllString(text, -1)
	fmt.Println(markers)

	s = `100.100.100.100 - - [23/Feb/2015:03:03:56 +0100] "GET /folder/file.mp3 HTTP/1.1" 206 5637064 "-" "AppleCoreMedia/1.0.0.12B466 (iPhone; U; CPU OS 8_1_3 like Mac OS X; da_dk)"`
	reg, _ = regexp.Compile(`([.0-9]+) .*?\[([0-9a-zA-Z:\/+ ]+)\].*?"[A-Z]+ \/([^\/ ]+)\/([a-zA-Z0-9\-.]+).*" ([0-9]{3}) .*"(.*?)"$`)
	a  := reg.FindStringSubmatch(s)
	for index, match := range a {
		fmt.Printf("[%d] %s\n", index, match)
	}

	s = `[body:content]`
	reg, _ = regexp.Compile(`(\[{1})([a-z]{2,64})(\:)([content]+)([\]]{1})`)
	b  := reg.FindStringSubmatch(s)
	for index, match := range b {
		fmt.Printf("[%d] %s\n", index, match)
	}

	// замена в тексте
	background := `http:\/\/wordpress.loc\/wp-content\/uploads\/layerslider\/Previous-demo-slider-with-2D-transitions\/bg1.jpg`
	background_reg, _ := regexp.Compile(`\\/`)
	fmt.Println(background_reg.ReplaceAllString(background, `/`))

	// замена в тексте
	transition := "{\\\"offsetxin\\\":\\\"0\\\",\\\"offsetyin\\\":\\\"bottom\\\",\\\"durationin\\\":\\\"4600\\\",\\\"delayin\\\":\\\"0\\\",\\\"easingin\\\":\\\"easeOutQuad\\\",\\\"fadein\\\":false,\\\"rotatein\\\":\\\"-10\\\",\\\"rotatexin\\\":\\\"0\\\",\\\"rotateyin\\\":\\\"0\\\",\\\"transformoriginin\\\":\\\"50% 50% 0\\\",\\\"skewxin\\\":\\\"0\\\",\\\"skewyin\\\":\\\"0\\\",\\\"scalexin\\\":\\\"1.0\\\",\\\"scaleyin\\\":\\\"1.0\\\",\\\"offsetxout\\\":\\\"0\\\",\\\"offsetyout\\\":\\\"0\\\",\\\"durationout\\\":\\\"1500\\\",\\\"showuntil\\\":\\\"0\\\",\\\"easingout\\\":\\\"easeInOutQuint\\\",\\\"fadeout\\\":true,\\\"rotateout\\\":\\\"0\\\",\\\"rotatexout\\\":\\\"0\\\",\\\"rotateyout\\\":\\\"0\\\",\\\"transformoriginout\\\":\\\"50% 50% 0\\\",\\\"skewxout\\\":\\\"0\\\",\\\"skewyout\\\":\\\"0\\\",\\\"scalexout\\\":\\\"1.0\\\",\\\"scaleyout\\\":\\\"1.0\\\",\\\"parallaxlevel\\\":\\\"0\\\"}"
	offsetxin, _ := regexp.Compile(`\\\"offsetxin\\\":\\\"([0-9]+)\\\"`)
	fmt.Println(offsetxin.ReplaceAllString(transition, `\"offsetxin\":$1`))

	// замена в тексте byte
	styles := `{\\\"width\\\":\\\"100px\\\",\\\"height\\\":\\\"70px\\\",\\\"font-family\\\":\\\"Arial, sans-serif\\\",\\\"font-size\\\":\\\"40px\\\",\\\"line-height\\\":\\\"70px\\\",\\\"color\\\":\\\"white\\\",\\\"background\\\":\\\"#cf431d\\\",\\\"border-radius\\\":\\\"5\\\"}`
	styles_reg, _ := regexp.Compile(`\\\\\\"`)
	fmt.Println(string(styles_reg.ReplaceAll([]byte(styles), []byte(`\"`))))

	// замена в тексте
	post_taxonomy := `{"post_taxonomy": 0}`
	post_taxonomy_reg, _ := regexp.Compile(`"post_taxonomy": ([0-9]+)`)
	fmt.Println(post_taxonomy_reg.ReplaceAllString(post_taxonomy, `"post_taxonomy":"$1"`))

	// поиск в тексте
	sd := `Edit [body:body:cont-ent] the Expression & Text to see matches. [some] f [some:marker] sdf [body:content] Roll over matches or the expression for details`
	reg, _ = regexp.CompilePOSIX(`\[{1}([a-zA-Z\-_:]*)\]{1}`)
	bd  := reg.FindAllStringSubmatch(sd, -1)
	for index, match := range bd {
		fmt.Printf("[%d] %s\n", index, match)
	}
}
