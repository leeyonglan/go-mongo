package teaappstat

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/tidwall/gjson"
)

func (c *ConfigDB) getRecipeItemByKeyID(keyField string, key string) Recipe {
	// Implement getItemByKeyID logic based on your requirements
	// This is a placeholder; you may need to fetch data from a database or elsewhere
	keyInt, err := strconv.Atoi(key)
	if err != nil {
		fmt.Println(err)
	}
	return RecipeMaps[keyInt]
}

type Recipe struct {
	RecipeID    int    `json:"recipe_id"`
	Name        string `json:"name"`
	Source      string `json:"source"`
	Pic         string `json:"pic"`
	Rarity      int    `json:"rarity"`
	Flavor      []int  `json:"flavor"`
	Desc        string `json:"desc"`
	Recipe      string `json:"recipe"`
	Hard        int    `json:"hard"`
	FoodType    int    `json:"food_type"`
	FoodRate    int    `json:"food_rate"`
	InPool1     int    `json:"in_pool1"`
	Available   int    `json:"available"`
	CookStep    []int  `json:"cook_step"`
	Step1Var    []int  `json:"step1_var"`
	Step2Var    []int  `json:"step2_var"`
	UnlockValue int    `json:"unlock_value"`
}

type UpgradeExp struct {
	ID            int `json:"id"`
	ToLevel       int `json:"to_level"`
	Normal        int `json:"normal"`
	Rare          int `json:"rare"`
	Excellent     int `json:"excellent"`
	Legendary     int `json:"legendary"`
	NormalGold    int `json:"normal_gold"`
	RareGold      int `json:"rare_gold"`
	ExcellentGold int `json:"excellent_gold"`
	LegendaryGold int `json:"legendary_gold"`
}
type DecorCollection struct {
	ID            int      `json:"id"`
	Sort          int      `json:"sort"`
	Name          string   `json:"name"`
	Type          int      `json:"type"`
	Set           int      `json:"set"`
	Pic           []string `json:"pic"`
	FlowerCoin    int      `json:"flower_coin"`
	Diamond       int      `json:"diamond"`
	DeliveryCD    int      `json:"delivery_CD"`
	Location      int      `json:"location"`
	Sprite        []string `json:"sprite"`
	Function      int      `json:"function"`
	FuncVar       int      `json:"func_var"`
	ShopExp       int      `json:"shop_exp"`
	Init          int      `json:"init"`
	Show          int      `json:"show"`
	Mandatory     int      `json:"mandatory"`
	PicType       int      `json:"picType"`
	ActionType    string   `json:"actionType"`
	ReplaceEffect int      `json:"replace_effect"`
}

func (c *ConfigDB) getUpgradeExp() map[int]UpgradeExp {
	// Implement getUpgradeExp logic based on your requirements
	// This is a placeholder; you may need to fetch data from a database or elsewhere
	return upgradeExpsMaps
}

var (
	upgradeExpsMaps  = map[int]UpgradeExp{}
	RecipeMaps       = map[int]Recipe{}
	DecorMaps        = map[int]DecorCollection{}
	upgradeExps      = []UpgradeExp{}
	DecorCollections = []DecorCollection{}
	Recipes          = []Recipe{}
)

func loadJsonFile(filePath string) (content []byte, err error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Read the file content
	content, err = io.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	return
}

func LoadAndParse() {

	content, _ := loadJsonFile("gameconfig/home.json")
	decorJson := gjson.Get(string(content), "decor_collection").String()
	// Parse JSON data into a struct
	err := json.Unmarshal([]byte(decorJson), &DecorCollections)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	for _, v := range DecorCollections {
		DecorMaps[v.ID] = v
	}

	content, _ = loadJsonFile("gameconfig/recipe.json")
	recipeJson := gjson.Get(string(content), "recipe").String()
	err = json.Unmarshal([]byte(recipeJson), &Recipes)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	for _, v := range Recipes {
		RecipeMaps[v.RecipeID] = v
	}

	upgradexpJson := gjson.Get(string(content), "upgrade_exp").String()
	err = json.Unmarshal([]byte(upgradexpJson), &upgradeExps)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	for _, v := range upgradeExps {
		upgradeExpsMaps[v.ToLevel] = v
	}
}
