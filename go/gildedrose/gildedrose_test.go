package gildedrose_test

import (
	"testing"

	"github.com/emilybache/gildedrose-refactoring-kata/gildedrose"
)

// 1. I began by identifying the different item types described in the requirements: normal items, Aged Brie, Backstage Passes, Sulfuras, and Conjured items.
// 2. Based on each item's unique behavior, I wrote targeted test cases to cover their respective update logic before and after the sell-by date.
// 3. To streamline validation and reduce repetitive assertions, I implemented a helper function `assertItem`
//    to verify that each item’s `SellIn` and `Quality` values matched the expected outcomes after update.
// 4. I ensured to include edge cases such as:
//    - Quality never going below 0
//    - Quality never exceeding 50
//    - Sulfuras remaining constant regardless of the update
// 5. To increase test coverage and avoid missing critical edge scenarios, I consulted an AI assistant to identify gaps in my test cases.
//    This led to the addition of tests for edge cases like clamping behavior when quality could drop below zero.
// 6. After running the full test suite, all test cases passed successfully, achieving 100% functional test coverage.

// Helper function to assert item properties
func assertItem(t *testing.T, item *gildedrose.Item, expectedSellIn, expectedQuality int) {
	if item.SellIn != expectedSellIn {
		t.Errorf("Expected SellIn %d, got %d", expectedSellIn, item.SellIn)
	}
	if item.Quality != expectedQuality {
		t.Errorf("Expected Quality %d, got %d", expectedQuality, item.Quality)
	}
}

// --- Normal Items ---
func Test_NormalItem_DegradesBy1(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Elixir of the Mongoose", SellIn: 5, Quality: 10},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 4, 9) // quality: 10 - 1 = 9 (normal degradation)
}

func Test_NormalItem_DegradesTwice_AfterSellIn(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Elixir of the Mongoose", SellIn: 0, Quality: 10},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -1, 8) // quality: 10 - 2 = 8 (degrades twice after sell date)
}

func Test_NormalItem_QualityZero_AfterSellIn(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Elixir of the Mongoose", SellIn: 0, Quality: 0},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -1, 0) // quality stays at 0, cannot go negative
}

func Test_QualityNeverGoesNegative_WhenDecreasedByMoreThanQuality(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Conjured Mana Cake", SellIn: 0, Quality: 1},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -1, 0) // quality: 1 - 4 = -3 → clamped to 0 (never negative)
}

// --- Aged Brie ---
func Test_AgedBrie_IncreasesInQuality(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Aged Brie", SellIn: 2, Quality: 0},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 1, 1) // quality: 0 + 1 = 1 (increases as it gets older)
}

func Test_AgedBrie_IncreasesBy2_AfterSellIn(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Aged Brie", SellIn: 0, Quality: 0},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -1, 2) // quality: 0 + 2 = 2 (increases twice as fast after sell date)
}

func Test_AgedBrie_DoesNotExceedQualityLimit(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Aged Brie", SellIn: 2, Quality: 50},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 1, 50) // quality capped at 50 (max limit)
}

func Test_AgedBrie_DoesNotIncreasePast50_AfterSellIn(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Aged Brie", SellIn: -1, Quality: 50},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -2, 50) // quality stays at 50 after sell date
}

// --- Backstage Passes ---
func Test_BackstagePasses_IncreaseBy2_When10OrLess(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Backstage passes to a TAFKAL80ETC concert", SellIn: 10, Quality: 20},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 9, 22) // quality: 20 + 2 = 22 (increases by 2 when 10 days or less)
}

func Test_BackstagePasses_IncreaseBy3_When5OrLess(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Backstage passes to a TAFKAL80ETC concert", SellIn: 5, Quality: 20},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 4, 23) // quality: 20 + 3 = 23 (increases by 3 when 5 days or less)
}

func Test_BackstagePasses_QualityStaysAt50_IfAlreadyMax(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Backstage passes to a TAFKAL80ETC concert", SellIn: 5, Quality: 50},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 4, 50) // quality capped at 50 (max limit)
}

func Test_BackstagePasses_DropsToZero_AfterConcert(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Backstage passes to a TAFKAL80ETC concert", SellIn: 0, Quality: 20},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -1, 0) // quality drops to 0 after concert (sell date passed)
}

func Test_BackstagePasses_DoesNotExceedQualityLimit(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Backstage passes to a TAFKAL80ETC concert", SellIn: 5, Quality: 49},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 4, 50) // quality capped at 50 (max limit)
}

// --- Sulfuras ---
func Test_Sulfuras_DoesNotChange(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Sulfuras, Hand of Ragnaros", SellIn: 0, Quality: 80},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 0, 80) // legendary item, quality never changes
}

func Test_Sulfuras_AlwaysStaysAt80(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Sulfuras, Hand of Ragnaros", SellIn: 5, Quality: 80},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 5, 80) // legendary item, quality stays at 80
}

// --- Conjured Items ---
func Test_ConjuredItem_DegradesBy2(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Conjured Mana Cake", SellIn: 3, Quality: 6},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], 2, 4) // quality: 6 - 2 = 4 (degrades twice as fast)
}

func Test_ConjuredItem_DegradesBy4_AfterSellIn(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Conjured Mana Cake", SellIn: 0, Quality: 6},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -1, 2) // quality: 6 - 4 = 2 (degrades twice as fast after sell date)
}

func Test_ConjuredItem_QualityNeverNegative(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Conjured Mana Cake", SellIn: 0, Quality: 3},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -1, 0) // quality: 3 - 4 = -1 → clamped to 0 (never negative)
}

func Test_ConjuredItem_QualityClamp_ExactDropBeyondZero(t *testing.T) {
	items := []*gildedrose.Item{
		{Name: "Conjured Mana Cake", SellIn: 0, Quality: 2},
	}
	gildedrose.UpdateQuality(items)
	assertItem(t, items[0], -1, 0) // quality: 2 - 4 = -2 → clamped to 0 (never negative)
}
