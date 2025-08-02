package main

import (
	"fmt"
	"lwnra-devo-api/parser"
)

func main() {
	sampleText := `DAILY DEVOTIONAL
Read Matthew 6:16-18
August 2, 2025
Matthew 6:16-18 NIV
16 When you fast, do not look somber as the hypocrites do, for they disfigure their faces to show men they are fasting. I tell you the truth, they have received their reward in full.
17 But when you fast, put oil on your head and wash your face,
18 so that it will not be obvious to men that you are fasting, but only to your Father, who is unseen; and your Father, who sees what is done in secret, will reward you.
REFLECTION QUESTIONS
Jesus warns against fasting to be seen (v. 16). What spiritual habit do you do partly for others to notice?
Why do you think God rewards what's done in secret (v. 18)? What does that reveal about His heart?
When was the last time you practiced a spiritual discipline with no one knowing? How did it feel?
How might living for God's approval change the way you pray, give, or fast today? What step can you take this week to nurture your "secret place" with God?`

	devo := parser.ParseDevotional(sampleText)

	fmt.Printf("Parsed Date: '%s'\n", devo.Date)
	fmt.Printf("Reading: '%s'\n", devo.Reading)
	fmt.Printf("Title: '%s'\n", devo.Title)
	fmt.Printf("Author: '%s'\n", devo.Author)
}
