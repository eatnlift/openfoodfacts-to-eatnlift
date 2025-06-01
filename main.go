package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

type OpenFoodFactsProduct struct {
	ID              string                 `json:"_id"`
	Code            string                 `json:"code"`
	ProductName     string                 `json:"product_name"`
	ProductNameEn   string                 `json:"product_name_en"`
	ProductNameFr   string                 `json:"product_name_fr"`
	ProductNameEs   string                 `json:"product_name_es"`
	ProductNameDe   string                 `json:"product_name_de"`
	ProductNameIt   string                 `json:"product_name_it"`
	ProductNameNl   string                 `json:"product_name_nl"`
	ProductNamePl   string                 `json:"product_name_pl"`
	ProductNamePt   string                 `json:"product_name_pt"`
	ProductNameUk   string                 `json:"product_name_uk"`
	ProductNameBg   string                 `json:"product_name_bg"`
	ProductNameRo   string                 `json:"product_name_ro"`
	ProductNameEl   string                 `json:"product_name_el"`
	ProductNameRu   string                 `json:"product_name_ru"`
	ProductNameTr   string                 `json:"product_name_tr"`
	ProductNameAr   string                 `json:"product_name_ar"`
	ProductNameHi   string                 `json:"product_name_hi"`
	Lang            string                 `json:"lang"`
	Brands          string                 `json:"brands"`
	ServingSize     string                 `json:"serving_size"`
	Nutriments      map[string]interface{} `json:"nutriments"`
	Allergens       string                 `json:"allergens"`
	AllergensTags   []string               `json:"allergens_tags"`
	IngredientsTags []string               `json:"ingredients_tags"`
}

type FoodItem struct {
	Name                string            `json:"name"`
	OffID               string            `json:"off_id"`
	Brand               string            `json:"brand"`
	Barcode             string            `json:"barcode"`
	ServingSizes        []ServingSize     `json:"serving_sizes"`
	Allergens           []string          `json:"allergens"`
	IngredientAllergens []string          `json:"ingredient_allergens"`
	Translations        map[string]string `json:"translations"`
}

type ServingSize struct {
	MeasurementUnit    string  `json:"measurement_unit"`
	Type               int     `json:"type"`
	Quantity           float64 `json:"quantity"`
	WeightInGrams      float64 `json:"weight_in_grams"`
	Calories           float64 `json:"calories,omitempty"`
	Protein            float64 `json:"protein,omitempty"`
	Fat                float64 `json:"fat,omitempty"`
	Carbs              float64 `json:"carbs,omitempty"`
	Fiber              float64 `json:"fiber,omitempty"`
	Sugar              float64 `json:"sugar,omitempty"`
	Sodium             float64 `json:"sodium,omitempty"`
	Cholesterol        float64 `json:"cholesterol,omitempty"`
	Calcium            float64 `json:"calcium,omitempty"`
	Iron               float64 `json:"iron,omitempty"`
	Potassium          float64 `json:"potassium,omitempty"`
	Magnesium          float64 `json:"magnesium,omitempty"`
	Zinc               float64 `json:"zinc,omitempty"`
	VitaminAIU         float64 `json:"vitamin_a_iu,omitempty"`
	VitaminC           float64 `json:"vitamin_c,omitempty"`
	VitaminD           float64 `json:"vitamin_d,omitempty"`
	VitaminDIU         float64 `json:"vitamin_d_iu,omitempty"`
	VitaminE           float64 `json:"vitamin_e,omitempty"`
	VitaminK           float64 `json:"vitamin_k,omitempty"`
	Thiamin            float64 `json:"thiamin,omitempty"`
	Riboflavin         float64 `json:"riboflavin,omitempty"`
	Niacin             float64 `json:"niacin,omitempty"`
	VitaminB6          float64 `json:"vitamin_b6,omitempty"`
	Folate             float64 `json:"folate,omitempty"`
	VitaminB12         float64 `json:"vitamin_b12,omitempty"`
	Phosphorus         float64 `json:"phosphorus,omitempty"`
	Copper             float64 `json:"copper,omitempty"`
	Manganese          float64 `json:"manganese,omitempty"`
	Selenium           float64 `json:"selenium,omitempty"`
	Water              float64 `json:"water,omitempty"`
	Ash                float64 `json:"ash,omitempty"`
	SaturatedFat       float64 `json:"saturated_fat,omitempty"`
	MonounsaturatedFat float64 `json:"monounsaturated_fat,omitempty"`
	PolyunsaturatedFat float64 `json:"polyunsaturated_fat,omitempty"`
	TransFat           float64 `json:"trans_fat,omitempty"`
}

// AllergenMap maps various allergen strings to standardized values
var AllergenMap = map[string]string{
	// Major allergens (FDA Big 9)
	"MILK":                 "milk",
	"EGGS":                 "eggs",
	"EGG":                  "eggs",
	"FISH":                 "fish",
	"SHELLFISH":            "shellfish",
	"CRUSTACEAN SHELLFISH": "crustacean_shellfish",
	"CRUSTACEAN_SHELLFISH": "crustacean_shellfish",
	"TREE NUTS":            "tree_nuts",
	"TREE_NUTS":            "tree_nuts",
	"PEANUTS":              "peanuts",
	"PEANUT":               "peanuts",
	"WHEAT":                "wheat",
	"SOY":                  "soy",
	"SOYBEAN":              "soy",
	"SOYBEANS":             "soy",
	"SESAME":               "sesame",

	// Specific tree nuts
	"ALMONDS":        "almonds",
	"ALMOND":         "almonds",
	"NUTS":           "nuts",
	"BRAZIL NUT":     "brazil_nuts",
	"BRAZIL_NUT":     "brazil_nuts",
	"BRAZIL NUTS":    "brazil_nuts",
	"BRAZIL_NUTS":    "brazil_nuts",
	"CASHEWS":        "cashews",
	"HAZELNUTS":      "hazelnuts",
	"MACADAMIA NUTS": "macadamia_nuts",
	"MACADAMIA_NUTS": "macadamia_nuts",
	"PECANS":         "pecans",
	"PINE NUTS":      "pine_nuts",
	"PINE_NUTS":      "pine_nuts",
	"PISTACHIOS":     "pistachios",
	"WALNUTS":        "walnuts",

	// Specific fish
	"ANCHOVY":   "anchovy",
	"COD":       "cod",
	"MAHI MAHI": "mahi_mahi",
	"MAHI_MAHI": "mahi_mahi",
	"SALMON":    "salmon",
	"TUNA":      "tuna",

	// Specific shellfish
	"CRAB":     "crab",
	"CRABS":    "crab",
	"LOBSTER":  "lobster",
	"LOBSTERS": "lobster",
	"SHRIMP":   "shrimp",
	"SHRIMPS":  "shrimp",
	"CLAMS":    "clams",
	"CLAM":     "clams",
	"MUSSELS":  "mussels",
	"MUSSEL":   "mussels",
	"OYSTERS":  "oysters",
	"OYSTER":   "oysters",
	"SCALLOPS": "scallops",
	"SCALLOP":  "scallops",

	// Grains containing gluten
	"BARLEY":    "barley",
	"RYE":       "rye",
	"OATS":      "oats",
	"TRITICALE": "triticale",
	"GLUTEN":    "gluten",

	// Other common allergens/sensitivities
	"CELERY":          "celery",
	"MUSTARD":         "mustard",
	"SULFITES":        "sulfites",
	"LUPIN":           "lupin",
	"MOLLUSKS":        "mollusks",
	"CORN":            "corn",
	"GELATIN":         "gelatin",
	"SEEDS":           "seeds",
	"SUNFLOWER SEEDS": "sunflower_seeds",
	"SUNFLOWER_SEEDS": "sunflower_seeds",
	"POPPY SEEDS":     "poppy_seeds",
	"POPPY_SEEDS":     "poppy_seeds",
	"COTTONSEED":      "cottonseed",
	"COCONUT":         "coconut",
	"PALM":            "palm",
	"BUCKWHEAT":       "buckwheat",
	"BEEF":            "beef",
	"PORK":            "pork",
	"CHICKEN":         "chicken",
	"GARLIC":          "garlic",
	"ONION":           "onion",
	"TOMATO":          "tomato",
	"LATEX":           "latex",
	"CARMINE":         "carmine",
	"COCHINEAL":       "cochineal",
	"ANNATTO":         "annatto",
	"MSG":             "msg",
	"SULFUR DIOXIDE":  "sulfur_dioxide",
	"SULFUR_DIOXIDE":  "sulfur_dioxide",
	"BENZOATES":       "benzoates",
	"FOOD COLORS":     "food_colors",
	"FOOD_COLORS":     "food_colors",
	"YELLOW 5":        "yellow_5",
	"YELLOW_5":        "yellow_5",
	"RED 40":          "red_40",
	"RED_40":          "red_40",
}

const INPUT_FILE = "input/openfoodfacts-products.jsonl.gz"
const OUTPUT_DIR = "output"
const CHUNK_SIZE = 50000

func main() {
	// Create output directory if it doesn't exist
	err := os.MkdirAll(OUTPUT_DIR, 0755)
	if err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	inputFile, err := os.Open(INPUT_FILE)
	if err != nil {
		log.Fatalf("Failed to open input file: %v", err)
	}
	defer inputFile.Close()

	gzReader, err := gzip.NewReader(inputFile)
	if err != nil {
		log.Fatalf("Failed to create gzip reader: %v", err)
	}
	defer gzReader.Close()

	decoder := json.NewDecoder(gzReader)

	lineCount := 0
	processedCount := 0
	chunkCount := 0
	var currentFile *os.File
	defer currentFile.Close()

	for {
		if processedCount%CHUNK_SIZE == 0 {
			// Close previous file if it exists
			if currentFile != nil {
				currentFile.Close()
			}

			// Create new file for next chunk
			outputPath := fmt.Sprintf("%s/openfoodfacts_to_eatnlift_%d.jsonl", OUTPUT_DIR, chunkCount)
			currentFile, err = os.Create(outputPath)
			if err != nil {
				log.Fatalf("Failed to create output file: %v", err)
			}
			chunkCount++
		}

		var product OpenFoodFactsProduct
		err := decoder.Decode(&product)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			continue
		}

		lineCount++

		processedProduct, err := ProcessProduct(product)
		if processedProduct == nil {
			log.Printf("Skipping product %s: %v", product.ID, err)
			continue
		}

		// Create a custom encoder that doesn't escape HTML
		buffer := &bytes.Buffer{}
		encoder := json.NewEncoder(buffer)
		encoder.SetEscapeHTML(false)

		err = encoder.Encode(processedProduct)
		if err != nil {
			log.Printf("Error encoding JSON: %v", err)
			continue
		}

		_, err = currentFile.WriteString(buffer.String())
		if err != nil {
			log.Printf("Error writing to output file: %v", err)
			continue
		}

		processedCount++
		if processedCount%10000 == 0 {
			log.Printf("Processed %d products", processedCount)
		}
	}

	log.Printf("Completed processing. Total lines: %d, Products processed: %d, Chunks created: %d",
		lineCount, processedCount, chunkCount)
}

func ProcessProduct(product OpenFoodFactsProduct) (*FoodItem, error) {
	if product.ID == "" || product.Code == "" {
		return nil, fmt.Errorf("product ID or code is empty")
	}

	translations := make(map[string]string)
	if product.ProductNameEn != "" {
		translations["en"] = product.ProductNameEn
	}
	if product.ProductNameFr != "" {
		translations["fr"] = product.ProductNameFr
	}
	if product.ProductNameEs != "" {
		translations["es"] = product.ProductNameEs
	}
	if product.ProductNameDe != "" {
		translations["de"] = product.ProductNameDe
	}
	if product.ProductNameIt != "" {
		translations["it"] = product.ProductNameIt
	}
	if product.ProductNameNl != "" {
		translations["nl"] = product.ProductNameNl
	}
	if product.ProductNamePl != "" {
		translations["pl"] = product.ProductNamePl
	}
	if product.ProductNamePt != "" {
		translations["pt"] = product.ProductNamePt
	}
	if product.ProductNameUk != "" {
		translations["uk"] = product.ProductNameUk
	}
	if product.ProductNameBg != "" {
		translations["bg"] = product.ProductNameBg
	}
	if product.ProductNameRo != "" {
		translations["ro"] = product.ProductNameRo
	}
	if product.ProductNameEl != "" {
		translations["el"] = product.ProductNameEl
	}
	if product.ProductNameRu != "" {
		translations["ru"] = product.ProductNameRu
	}
	if product.ProductNameTr != "" {
		translations["tr"] = product.ProductNameTr
	}
	if product.ProductNameAr != "" {
		translations["ar"] = product.ProductNameAr
	}
	if product.ProductNameHi != "" {
		translations["hi"] = product.ProductNameHi
	}

	name := translations["en"]
	if name == "" {
		name = product.ProductName
	}

	if product.ProductName != "" && len(translations) == 0 && product.Lang != "" {
		lang := strings.ToLower(product.Lang)
		translations[lang] = name
	}

	if name == "" && len(translations) > 0 && product.Lang != "" {
		lang := strings.ToLower(product.Lang)
		name = translations[lang]
	}

	if name == "" && len(translations) > 0 {
		for _, translation := range translations {
			name = translation
			break
		}
	}

	offID := product.ID
	brand := extractBrand(product.Brands)
	barcode := product.Code

	if name == "" && barcode == "" {
		return nil, fmt.Errorf("product name and barcode are empty")
	}

	allergens := product.AllergensTags
	if len(allergens) == 0 && product.Allergens != "" {
		allergens = strings.Split(product.Allergens, ",")
	}
	for i, allergen := range allergens {
		allergens[i] = normalizeAllergen(allergen)
	}

	ingredientAllergens := []string{}
	for _, ingredientTag := range product.IngredientsTags {
		ingredientAllergen := extractIngredientAllergen(ingredientTag)
		if ingredientAllergen != "" {
			ingredientAllergens = append(ingredientAllergens, ingredientAllergen)
		}
	}

	servingSizes := []ServingSize{}

	nutriments := product.Nutriments

	// Always include per-100g serving size
	per100gServing := ServingSize{
		MeasurementUnit: "g",
		Type:            1,
		Quantity:        100,
		WeightInGrams:   100,
	}

	// For per 100g serving size
	mapNutrient(&per100gServing.Calories, nutriments, "", "energy-kcal_100g")
	mapNutrient(&per100gServing.Protein, nutriments, "", "proteins_100g")
	mapNutrient(&per100gServing.Fat, nutriments, "", "fat_100g")
	mapNutrient(&per100gServing.Carbs, nutriments, "", "carbohydrates_100g")
	mapNutrient(&per100gServing.Fiber, nutriments, "", "fiber_100g")
	mapNutrient(&per100gServing.Sugar, nutriments, "", "sugars_100g")
	mapNutrient(&per100gServing.Sodium, nutriments, "g_to_mg", "sodium_100g")
	mapNutrient(&per100gServing.Cholesterol, nutriments, "g_to_mg", "cholesterol_100g")
	mapNutrient(&per100gServing.Calcium, nutriments, "g_to_mg", "calcium_100g")
	mapNutrient(&per100gServing.Iron, nutriments, "g_to_mg", "iron_100g")
	mapNutrient(&per100gServing.Potassium, nutriments, "g_to_mg", "potassium_100g")
	mapNutrient(&per100gServing.Magnesium, nutriments, "g_to_mg", "magnesium_100g")
	mapNutrient(&per100gServing.Zinc, nutriments, "g_to_mg", "zinc_100g")
	mapNutrient(&per100gServing.VitaminAIU, nutriments, "", "vitamin-a_iu_100g")
	mapNutrient(&per100gServing.VitaminC, nutriments, "g_to_mg", "vitamin-c_100g")
	mapNutrient(&per100gServing.VitaminD, nutriments, "g_to_ug", "vitamin-d_100g")
	mapNutrient(&per100gServing.VitaminDIU, nutriments, "", "vitamin-d_iu_100g")
	mapNutrient(&per100gServing.VitaminE, nutriments, "g_to_mg", "vitamin-e_100g")
	mapNutrient(&per100gServing.VitaminK, nutriments, "g_to_ug", "vitamin-k_100g")
	mapNutrient(&per100gServing.Thiamin, nutriments, "g_to_mg", "thiamin_100g")
	mapNutrient(&per100gServing.Riboflavin, nutriments, "g_to_mg", "riboflavin_100g")
	mapNutrient(&per100gServing.Niacin, nutriments, "g_to_mg", "niacin_100g")
	mapNutrient(&per100gServing.VitaminB6, nutriments, "g_to_mg", "vitamin-b6_100g")
	mapNutrient(&per100gServing.Folate, nutriments, "g_to_ug", "folate_100g")
	mapNutrient(&per100gServing.VitaminB12, nutriments, "g_to_ug", "vitamin-b12_100g")
	mapNutrient(&per100gServing.Phosphorus, nutriments, "g_to_mg", "phosphorus_100g")
	mapNutrient(&per100gServing.Copper, nutriments, "g_to_mg", "copper_100g")
	mapNutrient(&per100gServing.Manganese, nutriments, "g_to_mg", "manganese_100g")
	mapNutrient(&per100gServing.Selenium, nutriments, "g_to_ug", "selenium_100g")
	mapNutrient(&per100gServing.Water, nutriments, "", "water_100g")
	mapNutrient(&per100gServing.Ash, nutriments, "", "ash_100g")
	mapNutrient(&per100gServing.SaturatedFat, nutriments, "", "saturated-fat_100g")
	mapNutrient(&per100gServing.MonounsaturatedFat, nutriments, "", "monounsaturated-fat_100g")
	mapNutrient(&per100gServing.PolyunsaturatedFat, nutriments, "", "polyunsaturated-fat_100g")
	mapNutrient(&per100gServing.TransFat, nutriments, "", "trans-fat_100g")

	if per100gServing.Calories == 0 {
		per100gServing.Calories = per100gServing.Protein*4 + per100gServing.Fat*9 + per100gServing.Carbs*4
	}

	// If serving size information is available, include it as an additional serving size
	if product.ServingSize != "" {
		quantity, measurementUnit, weightInGrams, servingType := parseServingSize(product.ServingSize)

		// If weightInGrams is zero, try to calculate it based on the measurementUnit
		if weightInGrams == 0 && quantity > 0 {
			weightInGrams = convertToGrams(quantity, measurementUnit)
		}

		lowerMeasurementUnit := strings.ToLower(measurementUnit)
		skip := lowerMeasurementUnit == "" || lowerMeasurementUnit == "100g" || lowerMeasurementUnit == "100 g" || lowerMeasurementUnit == "100grams" || lowerMeasurementUnit == "100 grams"
		if !skip && quantity > 0 && weightInGrams != 100 {
			if len(measurementUnit) > 1 && isAsciiOnly(measurementUnit) {
				measurementUnit = toTitle(measurementUnit)
			}

			ss := ServingSize{
				MeasurementUnit: measurementUnit,
				Type:            servingType,
				Quantity:        quantity,
				WeightInGrams:   weightInGrams,
			}

			mapNutrient(&ss.Calories, nutriments, "", "energy-kcal_serving")
			mapNutrient(&ss.Protein, nutriments, "", "proteins_serving")
			mapNutrient(&ss.Fat, nutriments, "", "fat_serving")
			mapNutrient(&ss.Carbs, nutriments, "", "carbohydrates_serving")
			mapNutrient(&ss.Fiber, nutriments, "", "fiber_serving")
			mapNutrient(&ss.Sugar, nutriments, "", "sugars_serving")
			mapNutrient(&ss.Sodium, nutriments, "g_to_mg", "sodium_serving")
			mapNutrient(&ss.Cholesterol, nutriments, "g_to_mg", "cholesterol_serving")
			mapNutrient(&ss.Calcium, nutriments, "g_to_mg", "calcium_serving")
			mapNutrient(&ss.Iron, nutriments, "g_to_mg", "iron_serving")
			mapNutrient(&ss.Potassium, nutriments, "g_to_mg", "potassium_serving")
			mapNutrient(&ss.Magnesium, nutriments, "g_to_mg", "magnesium_serving")
			mapNutrient(&ss.Zinc, nutriments, "g_to_mg", "zinc_serving")
			mapNutrient(&ss.VitaminAIU, nutriments, "", "vitamin-a_iu_serving")
			mapNutrient(&ss.VitaminC, nutriments, "g_to_mg", "vitamin-c_serving")
			mapNutrient(&ss.VitaminD, nutriments, "g_to_ug", "vitamin-d_serving")
			mapNutrient(&ss.VitaminDIU, nutriments, "", "vitamin-d_iu_serving")
			mapNutrient(&ss.VitaminE, nutriments, "g_to_mg", "vitamin-e_serving")
			mapNutrient(&ss.VitaminK, nutriments, "g_to_ug", "vitamin-k_serving")
			mapNutrient(&ss.Thiamin, nutriments, "g_to_mg", "thiamin_serving")
			mapNutrient(&ss.Riboflavin, nutriments, "g_to_mg", "riboflavin_serving")
			mapNutrient(&ss.Niacin, nutriments, "g_to_mg", "niacin_serving")
			mapNutrient(&ss.VitaminB6, nutriments, "g_to_mg", "vitamin-b6_serving")
			mapNutrient(&ss.Folate, nutriments, "g_to_ug", "folate_serving")
			mapNutrient(&ss.VitaminB12, nutriments, "g_to_ug", "vitamin-b12_serving")
			mapNutrient(&ss.Phosphorus, nutriments, "g_to_mg", "phosphorus_serving")
			mapNutrient(&ss.Copper, nutriments, "g_to_mg", "copper_serving")
			mapNutrient(&ss.Manganese, nutriments, "g_to_mg", "manganese_serving")
			mapNutrient(&ss.Selenium, nutriments, "g_to_ug", "selenium_serving")
			mapNutrient(&ss.Water, nutriments, "", "water_serving")
			mapNutrient(&ss.Ash, nutriments, "", "ash_serving")
			mapNutrient(&ss.SaturatedFat, nutriments, "", "saturated-fat_serving")
			mapNutrient(&ss.MonounsaturatedFat, nutriments, "", "monounsaturated-fat_serving")
			mapNutrient(&ss.PolyunsaturatedFat, nutriments, "", "polyunsaturated-fat_serving")
			mapNutrient(&ss.TransFat, nutriments, "", "trans-fat_serving")

			if ss.Calories == 0 {
				ss.Calories = ss.Protein*4 + ss.Fat*9 + ss.Carbs*4
			}

			servingSizes = append(servingSizes, ss)
		}

	}

	servingSizes = append(servingSizes, per100gServing)

	foodItem := &FoodItem{
		Name:                name,
		OffID:               offID,
		Brand:               brand,
		Barcode:             barcode,
		ServingSizes:        servingSizes,
		Allergens:           allergens,
		IngredientAllergens: ingredientAllergens,
		Translations:        translations,
	}

	return foodItem, nil
}

func mapNutrient(field *float64, nutriments map[string]interface{}, unitConversion string, keys ...string) {
	for _, key := range keys {
		if value, ok := nutriments[key]; ok {
			if fval, err := toFloat64(value); err == nil {
				switch unitConversion {
				case "g_to_mg":
					*field = fval * 1000
				case "g_to_ug":
					*field = fval * 1e6
				default:
					*field = fval
				}
				return
			}
		}
	}
}

func toFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(strings.TrimSpace(v), 64)
	default:
		return 0, fmt.Errorf("unsupported type")
	}
}

func normalizeAllergen(allergen string) string {
	// Convert to uppercase for case-insensitive matching
	upperAllergen := strings.ToUpper(strings.TrimSpace(allergen))
	upperAllergen = strings.TrimPrefix(upperAllergen, "EN:")

	// Check if we have a mapping for this allergen
	if normalized, ok := AllergenMap[upperAllergen]; ok {
		return normalized
	}

	// Return original (but lowercase) if no mapping exists
	return strings.ToLower(allergen)
}

func extractIngredientAllergen(ingredientTag string) string {
	// Convert to uppercase for case-insensitive matching
	upperIngredientTag := strings.ToUpper(strings.TrimSpace(ingredientTag))
	upperIngredientTag = strings.TrimPrefix(upperIngredientTag, "EN:")

	// Check if we have a mapping for this ingredient tag
	if normalized, ok := AllergenMap[upperIngredientTag]; ok {
		return normalized
	}

	return ""
}
func parseServingSize(servingSizeStr string) (quantity float64, measurementUnit string, weightInGrams float64, servingType int) {
	servingSizeStr = strings.TrimSpace(servingSizeStr)

	// Correct common typos or abbreviations and descriptions
	servingSizeStr = strings.TrimSpace(servingSizeStr)
	servingSizeStr = strings.TrimSuffix(servingSizeStr, "|")
	servingSizeStr = strings.ReplaceAll(servingSizeStr, "OZA", "OZ")
	servingSizeStr = strings.ReplaceAll(servingSizeStr, "OZN", "OZ")
	servingSizeStr = strings.ReplaceAll(servingSizeStr, "ONZ", "OZ")
	servingSizeStr = strings.ReplaceAll(servingSizeStr, "Amount per serving", "Serving")
	servingSizeStr = strings.ReplaceAll(servingSizeStr, "FL.OZ", "FL OZ")

	// Replace commas with periods for decimal numbers
	servingSizeStr = strings.ReplaceAll(servingSizeStr, ",", ".")

	// Remove multiple spaces
	servingSizeStr = regexp.MustCompile(`\s+`).ReplaceAllString(servingSizeStr, " ")

	// Remove extra characters
	servingSizeStr = strings.Trim(servingSizeStr, "|")

	// Handle descriptors like "chips", "slice", "cookie", etc.
	servingSizeStr = handleDescriptors(servingSizeStr)

	// Replace commas with periods for decimal numbers
	servingSizeStr = strings.ReplaceAll(servingSizeStr, ",", ".")

	// Handle fractions (e.g., "1/2")
	servingSizeStr = convertFractions(servingSizeStr)

	// Remove any additional information in parentheses after the main serving size
	servingSizeStr = removeNestedParentheses(servingSizeStr)

	// Pattern to match numeric value followed by "g" or "ml" (e.g., "200g", "200.0ml")
	re := regexp.MustCompile(`^(\d*\.?\d+)\s*(g|gr|grm|gram|kg|mg|ml|l)$`)
	matches := re.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		qtyStr := matches[1]
		unitStr := matches[2]

		// Parse the numeric value
		weightQty, err := strconv.ParseFloat(qtyStr, 64)
		if err != nil {
			// If parsing fails, fallback to default handling
			return 1.0, servingSizeStr, 0.0, 3
		}

		// Convert the unit to grams
		weightInGrams = convertToGrams(weightQty, unitStr)

		// If weightInGrams matches the weightQty (after conversion), adjust the serving size
		if weightInGrams == weightQty || (unitStr == "ml" && weightInGrams == weightQty) {
			quantity = 1.0
			measurementUnit = "Serving"
			servingType = 3
			return quantity, measurementUnit, weightInGrams, servingType
		}
	}

	// Pattern: "quantity measurement_unit (weight_in_grams unit)"
	re = regexp.MustCompile(`(?i)^(?:(\d*\.?\d+)\s+)?([^\(]+?)\s*(?:\(\s*(\d*\.?\d+)\s*(g|gr|grm|gram|kg|mg|ml|l|oz|fl oz|ounces|cup|cups)\s*\))?$`)
	matches = re.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		qtyStr := matches[1]
		unitStr := matches[2]
		weightQtyStr := matches[3]
		weightUnit := strings.ToLower(matches[4])

		// Parse quantity
		var err error
		if qtyStr == "" {
			quantity = 1.0
		} else {
			quantity, err = strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
			if err != nil {
				quantity = 1.0
			}
		}

		measurementUnit = strings.TrimSpace(unitStr)
		measurementUnit = strings.TrimSuffix(measurementUnit, " e")

		// Parse weight
		weightInGrams = 0.0
		if weightQtyStr != "" && weightUnit != "" {
			weightQty, err := strconv.ParseFloat(strings.TrimSpace(weightQtyStr), 64)
			if err == nil {
				weightInGrams = convertToGrams(weightQty, weightUnit)
			}
		}

		// If weightInGrams is still zero, try to extract from measurementUnit
		if weightInGrams == 0.0 {
			weightInGrams = extractWeightFromMeasurementUnit(measurementUnit)
		}

		// If weightInGrams is still zero, estimate based on unit
		if weightInGrams == 0.0 {
			estimatedWeight := estimateWeightFromUnit(quantity, measurementUnit)
			weightInGrams = estimatedWeight
		}

		// Determine servingType
		lowerUnit := strings.ToLower(measurementUnit)
		if isMetricUnit(lowerUnit) {
			servingType = 1
		} else if isImperialUnit(lowerUnit) {
			servingType = 2
		} else {
			servingType = 3
		}

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Pattern: "quantity measurement_unit (weight_in_grams g)"
	re = regexp.MustCompile(`(?i)^(?:(\d*\.?\d+)\s+([^\(]+?))\s*(?:\(\s*(\d*\.?\d+)\s*(g|gr|grm|gram|kg|mg|ml|l|oz|fl oz|ounces)\s*\))?$`)
	matches = re.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		qtyStr := matches[1]
		unitStr := matches[2]
		weightQtyStr := matches[3]
		weightUnit := strings.ToLower(matches[4])

		// Parse quantity
		var err error
		quantity, err = strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
		if err != nil {
			quantity = 1.0
		}

		measurementUnit = strings.TrimSpace(unitStr)

		// Parse weight
		weightInGrams = 0.0
		if weightQtyStr != "" && weightUnit != "" {
			weightQty, err := strconv.ParseFloat(strings.TrimSpace(weightQtyStr), 64)
			if err == nil {
				weightInGrams = convertToGrams(weightQty, weightUnit)
			}
		}

		// If weightInGrams is still zero, try to estimate it
		if weightInGrams == 0.0 {
			estimatedWeight := estimateWeightFromUnit(quantity, measurementUnit)
			weightInGrams = estimatedWeight
		}

		// Determine servingType
		lowerUnit := strings.ToLower(measurementUnit)
		if isMetricUnit(lowerUnit) {
			servingType = 1
		} else if isImperialUnit(lowerUnit) {
			servingType = 2
		} else {
			servingType = 3
		}

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Pattern: "1 BOTTLE (295 ml)" or "1 Tbsp (15 ml)" or "8 OZ (240 ml)"
	re = regexp.MustCompile(`(?i)^(?:(\d*\.?\d+)\s+([^\(]+?))\s*\(\s*(\d*\.?\d+)\s*(g|gr|grm|gram|kg|ml|l|oz|fl oz|ounces)\s*\)$`)
	matches = re.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		qtyStr := matches[1]
		unitStr := matches[2]
		weightQtyStr := matches[3]
		weightUnit := strings.ToLower(matches[4])

		// Parse quantity
		var err error
		quantity, err = strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
		if err != nil {
			quantity = 1.0
		}

		measurementUnit = strings.TrimSpace(unitStr)

		// Parse weight
		weightQty, err := strconv.ParseFloat(strings.TrimSpace(weightQtyStr), 64)
		if err != nil {
			weightQty = 0.0
		}

		// Convert weight to grams
		weightInGrams = convertToGrams(weightQty, weightUnit)

		// Determine servingType
		lowerUnit := strings.ToLower(measurementUnit)
		if isMetricUnit(lowerUnit) {
			servingType = 1
		} else if isImperialUnit(lowerUnit) {
			servingType = 2
		} else {
			servingType = 3
		}

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Pattern: "1 slice 1 oz / 28 g"
	re = regexp.MustCompile(`(?i)^(?:(\d*\.?\d+)\s+([^\(]+?))\s+(\d*\.?\d+)\s*(oz|ounces)\s*/\s*(\d*\.?\d+)\s*(g|gr|grm|gram)$`)
	matches = re.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		qtyStr := matches[1]
		unitStr := matches[2]
		gQtyStr := matches[5]

		// Parse quantity
		var err error
		quantity, err = strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
		if err != nil {
			quantity = 1.0
		}

		measurementUnit = strings.TrimSpace(unitStr)

		// Use grams directly
		weightQty, err := strconv.ParseFloat(strings.TrimSpace(gQtyStr), 64)
		if err != nil {
			weightQty = 0.0
		}

		weightInGrams = weightQty

		// Determine servingType
		servingType = 3 // Non-standard unit

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Regular expression to match patterns like "1.5 g (1 TEA BAG)"
	re = regexp.MustCompile(`(?i)^(?:(\d*\.?\d+)\s*(g|gr|grm|gram|ml|l|oz|fl oz))\s*\(\s*(\d*\.?\d+)?\s*([^\)]+)\s*\)$`)
	matches = re.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		// Extract weight in grams
		weightQtyStr := matches[1]
		weightUnit := strings.ToLower(matches[2])

		var weightQty float64
		var err error
		weightQty, err = strconv.ParseFloat(strings.TrimSpace(weightQtyStr), 64)
		if err != nil {
			weightQty = 0.0
		}

		// Convert weight to grams if necessary
		weightInGrams = convertToGrams(weightQty, weightUnit)

		// Extract quantity and measurement unit from parentheses
		qtyStr := matches[3]
		unitStr := matches[4]

		if qtyStr == "" {
			quantity = 1.0
		} else {
			quantity, err = strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
			if err != nil {
				quantity = 1.0
			}
		}

		measurementUnit = strings.TrimSpace(unitStr)

		// Determine servingType
		lowerUnit := strings.ToLower(measurementUnit)
		if isMetricUnit(lowerUnit) {
			servingType = 1
		} else if isImperialUnit(lowerUnit) {
			servingType = 2
		} else {
			servingType = 3
		}

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Regular expression to match patterns like "2 SLICES (57 g)"
	re = regexp.MustCompile(`(?i)^(?:(\d*\.?\d+)\s+([^\(]+))\s*\(\s*(\d*\.?\d+)\s*(?:g|gr|grm|gram)\s*\)$`)
	matches = re.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		qtyStr := matches[1]
		unitStr := matches[2]
		weightStr := matches[3]

		// Parse quantity
		var err error
		quantity, err = strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
		if err != nil {
			quantity = 1.0
		}

		// Parse measurement unit
		measurementUnit = strings.TrimSpace(unitStr)

		// Parse weight in grams
		weightInGrams, err = strconv.ParseFloat(strings.TrimSpace(weightStr), 64)
		if err != nil {
			weightInGrams = 0.0
		}

		// Determine servingType
		lowerUnit := strings.ToLower(measurementUnit)
		if isMetricUnit(lowerUnit) {
			servingType = 1
		} else if isImperialUnit(lowerUnit) {
			servingType = 2
		} else {
			servingType = 3
		}

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Regular expression to match patterns like "30 g (30 GRM)"
	re = regexp.MustCompile(`(?i)^(?:(\d*\.?\d+)\s*(g|gr|grm|gram|ml|oz|fl oz))\s*\(\s*(\d*\.?\d+)\s*(?:g|gr|grm|gram)\s*\)$`)
	matches = re.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		qtyStr := matches[1]
		unitStr := matches[2]
		weightStr := matches[3]

		// Parse quantity
		var err error
		quantity, err = strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
		if err != nil {
			quantity = 1.0
		}

		// Parse measurement unit
		measurementUnit = strings.TrimSpace(unitStr)

		// Parse weight in grams
		weightInGrams, err = strconv.ParseFloat(strings.TrimSpace(weightStr), 64)
		if err != nil {
			weightInGrams = 0.0
		}

		// Determine servingType
		lowerUnit := strings.ToLower(measurementUnit)
		if isMetricUnit(lowerUnit) {
			servingType = 1
		} else if isImperialUnit(lowerUnit) {
			servingType = 2
		} else {
			servingType = 3
		}

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Pattern for "unit X g"
	reUnit := regexp.MustCompile(`(?i)^([^\s]+)\s+(\d*\.?\d+)\s*(g|gr|grm|gram|ml|oz|fl oz)$`)
	matches = reUnit.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		unitStr := matches[1]
		qtyStr := matches[2]
		weightUnit := matches[3]

		quantity, err := strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
		if err != nil {
			quantity = 1.0
		}

		measurementUnit = strings.TrimSpace(unitStr)
		weightInGrams = convertToGrams(quantity, weightUnit)

		lowerUnit := strings.ToLower(measurementUnit)
		if isMetricUnit(lowerUnit) {
			servingType = 1
		} else if isImperialUnit(lowerUnit) {
			servingType = 2
		} else {
			servingType = 3
		}

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Pattern for "X g"
	reSimple := regexp.MustCompile(`(?i)^(\d*\.?\d+)\s*(g|gr|grm|gram|ml|oz|fl oz)$`)
	matches = reSimple.FindStringSubmatch(servingSizeStr)
	if matches != nil {
		qtyStr := matches[1]
		unitStr := matches[2]

		var err error
		quantity, err = strconv.ParseFloat(strings.TrimSpace(qtyStr), 64)
		if err != nil {
			quantity = 1.0
		}

		measurementUnit = strings.TrimSpace(unitStr)
		weightInGrams = quantity

		lowerUnit := strings.ToLower(measurementUnit)
		if isMetricUnit(lowerUnit) {
			servingType = 1
		} else if isImperialUnit(lowerUnit) {
			servingType = 2
		} else {
			servingType = 3
		}

		return quantity, measurementUnit, weightInGrams, servingType
	}

	// Default handling
	quantity = 1.0
	measurementUnit = servingSizeStr
	servingType = 3
	weightInGrams = 0.0

	return quantity, measurementUnit, weightInGrams, servingType
}

// Convert fractions to decimal numbers
func convertFractions(input string) string {
	re := regexp.MustCompile(`(\d+)\s*/\s*(\d+)`)
	matches := re.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		numerator, _ := strconv.ParseFloat(match[1], 64)
		denominator, _ := strconv.ParseFloat(match[2], 64)
		if denominator != 0 {
			decimalValue := numerator / denominator
			input = strings.Replace(input, match[0], fmt.Sprintf("%.4f", decimalValue), -1)
		}
	}
	return input
}

// Remove nested parentheses beyond the first level
func removeNestedParentheses(input string) string {
	re := regexp.MustCompile(`\([^()]*\)`)
	matches := re.FindAllStringIndex(input, -1)
	if len(matches) > 1 {
		// Keep only the first match
		input = input[:matches[0][1]]
	}
	return input
}

// Handle descriptors like "chips", "slice", "cookie"
func handleDescriptors(input string) string {
	descriptors := []string{"chips", "slice", "slices", "cookie", "cookies", "pouch", "pouches", "can", "cans", "bottle", "bottles", "box", "boxes", "bag", "bags", "piece", "pieces"}
	for _, desc := range descriptors {
		re := regexp.MustCompile(`(?i)^` + desc + `\s+`)
		input = re.ReplaceAllString(input, "")
	}
	return input
}

// Extract weight from measurementUnit if possible
func extractWeightFromMeasurementUnit(measurementUnit string) float64 {
	re := regexp.MustCompile(`(?i)(\d*\.?\d+)\s*(g|gr|grm|gram|kg|mg|ml|l|oz|ounces|fl oz)`)
	matches := re.FindStringSubmatch(measurementUnit)
	if matches != nil {
		weightQtyStr := matches[1]
		weightUnit := matches[2]
		weightQty, err := strconv.ParseFloat(strings.TrimSpace(weightQtyStr), 64)
		if err == nil {
			return convertToGrams(weightQty, weightUnit)
		}
	}
	return 0.0
}

// Estimate weight based on measurement unit
func estimateWeightFromUnit(quantity float64, unit string) float64 {
	lowerUnit := strings.ToLower(unit)
	switch lowerUnit {
	case "cup", "cups":
		return quantity * 240 // 1 cup ~ 240g
	case "tbsp", "tablespoon", "tablespoons":
		return quantity * 15 // 1 tbsp ~ 15g
	case "tsp", "teaspoon", "teaspoons":
		return quantity * 5 // 1 tsp ~ 5g
	case "slice", "slices":
		return quantity * 28 // 1 slice ~ 28g
	case "cookie", "cookies":
		return quantity * 15 // Estimate for a cookie
	// Add more estimations as needed
	default:
		return 0.0
	}
}

// Convert units to grams
func convertToGrams(quantity float64, unit string) float64 {
	switch strings.ToLower(unit) {
	case "g", "gram", "grams", "gr", "grm", "g.", "gr.", "grm.":
		return quantity
	case "kg", "kilogram", "kilograms":
		return quantity * 1000
	case "mg", "milligram", "milligrams":
		return quantity / 1000
	case "oz", "ounces", "oz.", "ounce", "onz", "ozn", "oza":
		return quantity * 28.3495
	case "lb", "pound", "pounds":
		return quantity * 453.592
	case "ml":
		return quantity // Assuming 1 g/ml
	case "l", "liter", "litre", "liters", "litres":
		return quantity * 1000
	default:
		return 0.0
	}
}

// Check if the unit is a metric unit
func isMetricUnit(unit string) bool {
	metricUnits := []string{"g", "gram", "grams", "gr", "grm", "g.", "gr.", "grm.", "kg", "kilogram", "kilograms", "ml", "l", "liter", "litre", "liters", "litres"}
	for _, u := range metricUnits {
		if unit == u {
			return true
		}
	}
	return false
}

// Check if the unit is an imperial unit
func isImperialUnit(unit string) bool {
	imperialUnits := []string{"oz", "ounce", "ounces", "onz", "ozn", "oza", "lb", "pound", "pounds", "fl oz"}
	for _, u := range imperialUnits {
		if unit == u {
			return true
		}
	}
	return false
}

func extractBrand(brands string) string {
	brandsSplit := strings.Split(strings.TrimSpace(brands), ",")
	if len(brandsSplit) == 0 {
		return brands
	}
	return strings.TrimSpace(brandsSplit[0])
}

func toTitle(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func isAsciiOnly(s string) bool {
	return strings.IndexFunc(s, func(r rune) bool {
		return r > unicode.MaxASCII
	}) == -1
}
