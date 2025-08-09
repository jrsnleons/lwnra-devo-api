package main

import (
	"fmt"
	"lwnra-devo-api/parser"
	"strings"
)

func main() {
	testPassageDebug()
}

func testPassageDebug() {
	// Sample problematic content based on the production output
	sampleContent := `August 8, 2025

Revelation 7:9-12 NIV
9 After this I looked and there before me was a great multitude that no one could count, from every nation, tribe, people and language, standing before the throne and in front of the Lamb. They were wearing white robes and were holding palm branches in their hands.
10 And they cried out in a loud voice: "Salvation belongs to our God, who sits on the throne, and to the Lamb."
11 All the angels were standing around the throne and around the elders and the four living creatures. They fell down on their faces before the throne and worshiped God,
12 saying: "Amen! Praise and glory and wisdom and thanks and honor and power and strength be to our God for ever and ever. Amen!"

REFLECTION QUESTIONS:
- How does this picture of heavenly worship challenge the way you view diversity in the church?
- What does it mean for you personally to be "clothed in white" (symbolizing righteousness) before God?
- Do your current worship habits reflect the passionate praise seen in this passage?
- How can you align your life today with the worship of eternity?

A PREVIEW OF FOREVER

Reflections in Grace

Concerts, sporting events, and festivals often unite people who would otherwise never meet. Strangers stand side by side, moved by the same music, cheering for the same team, or celebrating the same cause. For a moment, barriers blur. Age, background, and culture give way to shared passion. What if that kind of unity was just a small echo of something far greater?

Revelation 7 paints a breathtaking picture of what's to come: "a great multitude that no one could count, from every nation, tribe, people and language, standing before the throne and in front of the Lamb" (v.9). There's no division here. Just unified worship, focused on Jesus, the Lamb who saves. These are people from every corner of the world, clothed in white, waving palm branches, crying out in loud voices, "Salvation belongs to our God!"

This vision reminds us that the gospel isn't just for one kind of person. It's for everyone. And the worship of heaven is not quiet, distant, or exclusive. It's vibrant, diverse, and united in adoration of the One who made a way for all.

Let this vision reshape your today. See people not as "us and them," but as potential co-worshipers before the throne. Let worship now be a preview of forever.

{Lord, thank you for a glimpse of heaven's glory. Fill our hearts with the same praise that echoes around your throne. Help us to live with a wide, open love that reflects your heart for every nation and people. Let our worship today honor you as it will for all eternity. Amen.}`

	fmt.Println("=== PARSING TEST ===")
	devo := parser.ParseDevotional(sampleContent)

	fmt.Printf("Date: %s\n", devo.Date)
	fmt.Printf("Reading: %s\n", devo.Reading)
	fmt.Printf("Version: %s\n", devo.Version)
	fmt.Printf("Title: %s\n", devo.Title)
	fmt.Printf("Author: %s\n", devo.Author)
	fmt.Printf("\n=== PASSAGE (PROBLEMATIC) ===\n%s\n", devo.Passage)
	fmt.Printf("\n=== REFLECTION QUESTIONS ===\n")
	for i, q := range devo.ReflectionQs {
		fmt.Printf("%d. %s\n", i+1, q)
	}
	fmt.Printf("\n=== BODY ===\n%s\n", devo.Body)
	fmt.Printf("\n=== PRAYER ===\n%s\n", devo.Prayer)

	// Show what the passage should be
	fmt.Printf("\n=== WHAT PASSAGE SHOULD BE ===\n")
	lines := strings.Split(sampleContent, "\n")
	for i, line := range lines {
		fmt.Printf("%d: %s\n", i, line)
	}
}
