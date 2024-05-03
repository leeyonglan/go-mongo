package teaappstat

import "fmt"

type ConfigDB struct {
	// Implement ConfigDB methods if needed
}

type TableEnum struct {
	// Implement TableEnum methods if needed
}

func getHandCardLv(id int, allHandBook map[int]int, configDB *ConfigDB) map[string]interface{} {
	recipeItem := configDB.getRecipeItemByKeyID("recipe_id", fmt.Sprintf("%d", id))
	exp := allHandBook[id]
	if exp == 0 {
		exp = 0
	}
	rare := recipeItem.Rarity
	lv := getCardLvAndExp(rare, exp, configDB)
	return lv
}

func getCardLvAndExp(rare int, exp int, configDB *ConfigDB) map[string]interface{} {
	rareVal := ""
	switch rare {
	case 1:
		rareVal = "normal"
	case 2:
		rareVal = "rare"
	case 3:
		rareVal = "excellent"
	case 4:
		rareVal = "legendary"
	default:
	}

	upgradeExps := configDB.getUpgradeExp() // Replace TableEnum with the actual type
	lv := 0
	curTotal := 0
	nextTotal := 0

	for lv < len(upgradeExps) && exp >= 0 {
		curTotal += getUpgradeExp(upgradeExps[lv], rareVal)

		if lv+1 < len(upgradeExps) {
			nextTotal = curTotal + getUpgradeExp(upgradeExps[lv+1], rareVal)
		} else {
			return map[string]interface{}{
				"lv":   lv + 1,
				"exp":  exp - curTotal,
				"full": true,
			}
		}

		if exp < curTotal {
			if lv == 0 {
				return map[string]interface{}{
					"lv":   lv,
					"exp":  exp,
					"need": getUpgradeExp(upgradeExps[lv], rareVal),
				}
			} else {
				return map[string]interface{}{
					"lv":   lv,
					"exp":  exp - curTotal,
					"need": getUpgradeExp(upgradeExps[lv], rareVal),
				}
			}
		}

		if exp < nextTotal {
			return map[string]interface{}{
				"lv":   lv + 1,
				"exp":  exp - curTotal,
				"need": getUpgradeExp(upgradeExps[lv+1], rareVal),
			}
		}
		lv++
	}

	return nil
}

func getUpgradeExp(upgradeExp UpgradeExp, rareVal string) int {
	switch rareVal {
	case "normal":
		return upgradeExp.Normal
	case "rare":
		return upgradeExp.Rare
	case "excellent":
		return upgradeExp.Excellent
	case "legendary":
		return upgradeExp.Legendary
	default:
		return 0
	}
}

// func main() {
// 	rare := 2
// 	exp := 100

// 	configDB := &ConfigDB{} // Replace with the actual instantiation of ConfigDB
// 	result := getCardLvAndExp(rare, exp, configDB)
// 	fmt.Println(result)
// }
