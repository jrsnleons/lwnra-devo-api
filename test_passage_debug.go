package main

import (
	"fmt"
	"lwnra-devo-api/parser"
	"strings"
)

func testPassageDebug() {
	// Sample content: August 9, 2025 devotional
	sampleContent := `DAILY DEVOTIONAL
Read Psalms 71:1-13
August 9, 2025
Psalms 71:1-13 ESV-r
1 In you, O LORD, do I take refuge; let me never be put to shame!
2 In your righteousness deliver me and rescue me; incline your ear to me, and save me!
3 Be to me a rock of refuge, to which I may continually come; you have given the command to save me, for you are my rock and my fortress.
4 Rescue me, O my God, from the hand of the wicked, from the grasp of the unjust and cruel man.
5 For you, O Lord, are my hope, my trust, O LORD, from my youth.
6 Upon you I have leaned from before my birth; you are he who took me from my mother's womb. My praise is continually of you.
7 I have been as a portent to many, but you are my strong refuge.
8 My mouth is filled with your praise, and with your glory all the day.
9 Do not cast me off in the time of old age; forsake me not when my strength is spent.
10 For my enemies speak concerning me; those who watch for my life consult together
11 and say, "God has forsaken him; pursue and seize him, for there is none to deliver him."
12 O God, be not far from me; O my God, make haste to help me!
13 May my accusers be put to shame and consumed; with scorn and disgrace may they be covered who seek my hurt.
REFLECTION QUESTIONS
(Verse 1) What does it look like to take refuge in God instead of distractions or self-reliance?
(Verse 3) How have you seen God's protection in past situations? Are there “fortresses” in your life that are not truly secure?
(Verse 5) Reflect on your journey - how has your relationship with God grown or been tested over time?
(Verses 10-11) Are there negative voices causing fear or doubt? How can God's truth silence them?
(Verse 13) Are you facing situations where you need God to act on your behalf? What would it look like to surrender the outcome to him today? 
A REFUGE IN EVERY SEASON
Reflections in Grace
Whether you're a young adult facing the pressure of proving yourself, a parent feeling overwhelmed, or a senior wondering if your strength is enough for this season, there comes a moment when we all long for a safe place.

Psalm 71 is the prayer of someone who knows what it’s like to be under pressure. The psalmist, likely in the latter years of life, cries out, “In you, O LORD, I have taken refuge; let me never be put to shame” (v.1). This is not a casual statement; it’s an appeal from someone who has walked with God for years and is now facing enemies, fear, and even feelings of abandonment.

What’s striking is the consistent trust in God’s faithfulness: “Be my rock of refuge, to which I can always go” (v.3). The psalmist admits his need for rescue, strength, and justice but more than anything, for the abiding presence of God. 

In a world that often measures worth by strength, success, or youth, this passage reminds us that God's refuge isn't seasonal. It’s for the long haul. When the voices of opposition grow louder, when your past is questioned, or when your future feels fragile, hold fast to this truth: the God who has been your refuge will continue to be, no matter your age, season, or struggle.
PRAYER
Lord, you are my refuge, my unshakable rock. In every season of life, I place my hope in you. Silence every voice that speaks fear and shame, and remind me daily that I am secure in your hands. Amen.`

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

func testTodaysDevotionalBody() {
	sampleContent := `DAILY DEVOTIONAL
Read Psalms 71:1-13
August 9, 2025
Psalms 71:1-13 ESV-r
1 In you, O LORD, do I take refuge; let me never be put to shame!
2 In your righteousness deliver me and rescue me; incline your ear to me, and save me!
3 Be to me a rock of refuge, to which I may continually come; you have given the command to save me, for you are my rock and my fortress.
4 Rescue me, O my God, from the hand of the wicked, from the grasp of the unjust and cruel man.
5 For you, O Lord, are my hope, my trust, O LORD, from my youth.
6 Upon you I have leaned from before my birth; you are he who took me from my mother's womb. My praise is continually of you.
7 I have been as a portent to many, but you are my strong refuge.
8 My mouth is filled with your praise, and with your glory all the day.
9 Do not cast me off in the time of old age; forsake me not when my strength is spent.
10 For my enemies speak concerning me; those who watch for my life consult together
11 and say, "God has forsaken him; pursue and seize him, for there is none to deliver him."
12 O God, be not far from me; O my God, make haste to help me!
13 May my accusers be put to shame and consumed; with scorn and disgrace may they be covered who seek my hurt.
REFLECTION QUESTIONS
(Verse 1) What does it look like to take refuge in God instead of distractions or self-reliance?
(Verse 3) How have you seen God's protection in past situations? Are there “fortresses” in your life that are not truly secure?
(Verse 5) Reflect on your journey - how has your relationship with God grown or been tested over time?
(Verses 10-11) Are there negative voices causing fear or doubt? How can God's truth silence them?
(Verse 13) Are you facing situations where you need God to act on your behalf? What would it look like to surrender the outcome to him today? 
A REFUGE IN EVERY SEASON
Reflections in Grace
Whether you're a young adult facing the pressure of proving yourself, a parent feeling overwhelmed, or a senior wondering if your strength is enough for this season, there comes a moment when we all long for a safe place.

Psalm 71 is the prayer of someone who knows what it’s like to be under pressure. The psalmist, likely in the latter years of life, cries out, “In you, O LORD, I have taken refuge; let me never be put to shame” (v.1). This is not a casual statement; it’s an appeal from someone who has walked with God for years and is now facing enemies, fear, and even feelings of abandonment.

What’s striking is the consistent trust in God’s faithfulness: “Be my rock of refuge, to which I can always go” (v.3). The psalmist admits his need for rescue, strength, and justice but more than anything, for the abiding presence of God. 

In a world that often measures worth by strength, success, or youth, this passage reminds us that God's refuge isn't seasonal. It’s for the long haul. When the voices of opposition grow louder, when your past is questioned, or when your future feels fragile, hold fast to this truth: the God who has been your refuge will continue to be, no matter your age, season, or struggle.
PRAYER
Lord, you are my refuge, my unshakable rock. In every season of life, I place my hope in you. Silence every voice that speaks fear and shame, and remind me daily that I am secure in your hands. Amen.`

	expectedBody := `Whether you're a young adult facing the pressure of proving yourself, a parent feeling overwhelmed, or a senior wondering if your strength is enough for this season, there comes a moment when we all long for a safe place.

Psalm 71 is the prayer of someone who knows what it’s like to be under pressure. The psalmist, likely in the latter years of life, cries out, “In you, O LORD, I have taken refuge; let me never be put to shame” (v.1). This is not a casual statement; it’s an appeal from someone who has walked with God for years and is now facing enemies, fear, and even feelings of abandonment.

What’s striking is the consistent trust in God’s faithfulness: “Be my rock of refuge, to which I can always go” (v.3). The psalmist admits his need for rescue, strength, and justice but more than anything, for the abiding presence of God. 

In a world that often measures worth by strength, success, or youth, this passage reminds us that God's refuge isn't seasonal. It’s for the long haul. When the voices of opposition grow louder, when your past is questioned, or when your future feels fragile, hold fast to this truth: the God who has been your refuge will continue to be, no matter your age, season, or struggle.`

	devo := parser.ParseDevotional(sampleContent)
	fmt.Println("\n=== TEST: Today's Devotional Body Extraction ===")
	fmt.Println("--- Parsed Body ---")
	fmt.Println(devo.Body)
	fmt.Println("--- Expected Body ---")
	fmt.Println(expectedBody)
	if devo.Body == expectedBody {
		fmt.Println("PASS: Body extraction matches expected output.")
	} else {
		fmt.Println("FAIL: Body extraction does not match expected output.")
	}
}

func main() {
	testPassageDebug()
	testTodaysDevotionalBody()
}
