package gildedrose

import "strings"

// Refactored the Gilded Rose inventory system to improve structure, readability, and maintainability.
//
// 1. Started by thoroughly reading the requirements specification document.
// 2. While reading, it became clear that special items would need custom handling logic.
// 3. Designed a system of item-type check functions to encapsulate this logic:
//    - isSulfuras, isAgedBrie, isBackstagePass, and isConjured.
// 4. Analyzed the base code and identified repetition and unclear quality adjustment logic.
// 5. Created two helper functions to centralize and simplify quality updates:
//    - increaseQuality(item *Item, amount int)
//    - decreaseQuality(item *Item, amount int)
// 6. Built an `updateItem` function to manage daily updates for each item.
// 7. Introduced the `updateQuality` function to branch behavior based on item type.
// 8. Implemented `handleExpired` to manage post-sell date behavior using special rules.
// 9. Added quality clamping logic (min 0, max 50) to avoid magic numbers and ensure consistent constraints.
// 10. Refactored the code to be readable, and easily extensible for future item types.
// 11. Additionally, I consulted an AI assistant to cross-check that my implementation aligned with the original requirements, ensuring the logic I applied accurately covered all specified item behaviors.

type Item struct {
	Name            string
	SellIn, Quality int
}

func isSulfuras(item *Item) bool {
	return item.Name == "Sulfuras, Hand of Ragnaros" // legendary item
}

func isAgedBrie(item *Item) bool {
	return item.Name == "Aged Brie" // special item that increases quality
}

func isBackstagePass(item *Item) bool {
	return item.Name == "Backstage passes to a TAFKAL80ETC concert" // special item with increasing quality then drops to 0
}

func isConjured(item *Item) bool {
	return strings.HasPrefix(item.Name, "Conjured") // degrades twice as fast
}

func increaseQuality(item *Item, amount int) {
	item.Quality += amount
	if item.Quality > 50 {
		item.Quality = 50 // quality max capped at 50
	}
}

func decreaseQuality(item *Item, amount int) {
	item.Quality -= amount
	if item.Quality < 0 {
		item.Quality = 0 // quality never negative
	}
}

func UpdateQuality(items []*Item) {
	for _, item := range items {
		updateItem(item) // update each item daily
	}
}

func updateItem(item *Item) {
	if isSulfuras(item) {
		return // Sulfuras never changes
	}

	updateQuality(item) // adjust quality based on item type
	item.SellIn--       // decrease sellIn by 1 day

	if item.SellIn < 0 { // after sell date passed
		handleExpired(item) // special rules apply
	}
}

func updateQuality(item *Item) {
	switch {
	case isAgedBrie(item):
		increaseQuality(item, 1) // Brie increases quality by 1

	case isBackstagePass(item):
		increaseQuality(item, 1) // increase by 1 normally
		if item.SellIn <= 10 {
			increaseQuality(item, 1) // +1 more when 10 days or less
		}
		if item.SellIn <= 5 {
			increaseQuality(item, 1) // +1 more when 5 days or less
		}

	case isConjured(item):
		decreaseQuality(item, 2) // conjured items degrade twice as fast

	default:
		decreaseQuality(item, 1) // normal items degrade by 1
	}
}

func handleExpired(item *Item) {
	switch {
	case isAgedBrie(item):
		increaseQuality(item, 1) // Brie increases quality twice as fast after sell date

	case isBackstagePass(item):
		item.Quality = 0 // backstage passes worthless after concert

	case isConjured(item):
		decreaseQuality(item, 2) // conjured degrade twice as fast after sell date

	default:
		decreaseQuality(item, 1) // normal items degrade twice as fast after sell date
	}
}
