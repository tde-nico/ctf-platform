package telegram_bot

type Message struct {
	Text      string
	MediaType string
	Media     string
}

/*
CAACAgQAAxkBAAExXY1nn3i7L9e7bC_Pt5vHj6wpTfL4WwACyRcAAiog-VAHs5wSbsrARDYE - nugget sticker
CAACAgQAAxkBAAExXZFnn3lc9YPvtCjLIVvjTLCA2RLS1wACTRkAAldf-VBMtTrazPGrxDYE - party sticker
CAACAgQAAxkBAAExXZNnn3lsCxmGYHfQzk_LRVlA7MfrlwADFgACxyr4UK7vP5qOleE1NgQ - dramatic sticker
CAACAgQAAxkBAAExXbxnn3_DYC3tM3yGI3yZ2rlN1LWJyAACchwAAkkO-VAiG7AK35yo_zYE - hacker sticker
*/

var messages = []Message{
	{
		Text:      "🏎💨 <user> is smoking everyone 🏎💨\n🩸 FIRST BLOOD! on <chall> 🩸",
		MediaType: "sticker",
		Media:     "CAACAgQAAxkBAAExXY1nn3i7L9e7bC_Pt5vHj6wpTfL4WwACyRcAAiog-VAHs5wSbsrARDYE",
	},
	{
		Text:      "🎉 <user> is throwing a rave in the chat 🎉\n🩸 FIRST BLOOD! on <chall> 🩸",
		MediaType: "sticker",
		Media:     "CAACAgQAAxkBAAExXZFnn3lc9YPvtCjLIVvjTLCA2RLS1wACTRkAAldf-VBMtTrazPGrxDYE",
	},
	{
		Text:      "Everyone panicking after <user>\n🩸 FIRST BLOOD! on <chall> 🩸",
		MediaType: "sticker",
		Media:     "CAACAgQAAxkBAAExXZNnn3lsCxmGYHfQzk_LRVlA7MfrlwADFgACxyr4UK7vP5qOleE1NgQ",
	},
	{
		Text:      "👨‍💻 Live Footage of <user> 👨‍💻\n🩸 FIRST BLOOD! on <chall> 🩸",
		MediaType: "sticker",
		Media:     "CAACAgQAAxkBAAExXbxnn3_DYC3tM3yGI3yZ2rlN1LWJyAACchwAAkkO-VAiG7AK35yo_zYE",
	},
	// {
	// 	Text:      "🖥 STACCA STACCA 🖥\n%s stole some sensitive military data!\n🩸 FIRST BLOOD! on %s 🩸",
	// 	MediaType: "animation",
	// 	Media:     "CgACAgQAAxkBAAEnksZlT_o5YMwWw7_lox819Yyj0jqGfQACQgMAAskZBVNIOKZsmIrdtjME",
	// },
	// {
	// 	Text:      "Grandma better hold tight!\n🏎💨 %s hacking fast! 🏎💨\n🩸 FIRST BLOOD! on %s 🩸",
	// 	MediaType: "sticker",
	// 	Media:     "CAACAgQAAxkBAAEnksJlT_hM7EAhhq0q_0oROm1C_k0LhAACUQkAAtXoWFB-DdtZU0I5hTME",
	// },
	// {
	// 	Text:      "✍️ Polizia Postale should start taking notes ✍️\n'cause %s is making himself noticed!\n🩸 FIRST BLOOD! on %s 🩸",
	// 	MediaType: "animation",
	// 	Media:     "CgACAgQAAxkBAAEnkudlUAkZ8ksOnnhSGO1658u1SaKMpgACPAMAAs--FVOcbXe1qQ6kSDME",
	// },
	// {
	// 	Text:      "💋 %s, bet your love wouldn't want you to be so fast... 💦\n🩸 FIRST BLOOD! on %s 🩸",
	// 	MediaType: "sticker",
	// 	Media:     "CgACAgQAAxkBAAEnkuhlUAkZNvPltQpvxbqbehYr8RKmOwAC5QIAAl40dVP2hs5pV4wQEDME",
	// },
	// {
	// 	Text:      "📆 Sign this date and hour cause you'll want to remember when %s proved everyone to be the fastest! 🚀\n🩸 FIRST BLOOD! on %s 🩸",
	// 	MediaType: "animation",
	// 	Media:     "CgACAgQAAxkBAAEnkullUAkZuFOMYXTIYfppj4npK0YLSAACywIAAp24DFOWwlLw9CpdTTME",
	// },
}
