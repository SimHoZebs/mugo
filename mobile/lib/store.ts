import { create } from "zustand";
import { devtools, persist } from "zustand/middleware";
import type {} from "@redux-devtools/extension"; // required for devtools typing
import { NutritionResponseBody } from "./api/conversationAPI.schemas";

interface GlobalState {
  meals: NutritionResponseBody["analysis"][];
  setMeals: (meals: NutritionResponseBody["analysis"][]) => void;
}

const useGlobalStore = create<GlobalState>()(
  devtools(
    persist(
      (set) => ({
        meals: [
          {
            name: "Oatmeal with Berries",
            assumptions: [
              {
                id: "1",
                field: "portion_size",
                assumed_value: 1,
                unit: "cup",
                confidence: "high",
                rationale: "Standard serving size for oatmeal",
              },
              {
                id: "8",
                field: "milk_type",
                assumed_value: 1,
                unit: "cup",
                confidence: "medium",
                rationale: "Assumed whole milk for cooking",
              },
            ],
            macros: {
              calories: 320,
              protein: 12,
              carbs: 45,
              fat: 8,
            },
          },
          {
            name: "Grilled Chicken Salad",
            assumptions: [
              {
                id: "2",
                field: "portion_size",
                assumed_value: 150,
                unit: "grams",
                confidence: "medium",
                rationale: "Typical chicken breast portion",
              },
              {
                id: "9",
                field: "cooking_method",
                assumed_value: 1,
                unit: "grilled",
                confidence: "high",
                rationale: "Common preparation method",
              },
              {
                id: "10",
                field: "dressing_amount",
                assumed_value: 2,
                unit: "tbsp",
                confidence: "low",
                rationale: "Light dressing application",
              },
            ],
            macros: {
              calories: 450,
              protein: 35,
              carbs: 20,
              fat: 18,
            },
          },
          {
            name: "Greek Yogurt Parfait",
            assumptions: [
              {
                id: "3",
                field: "portion_size",
                assumed_value: 200,
                unit: "grams",
                confidence: "high",
                rationale: "Standard yogurt container size",
              },
              {
                id: "11",
                field: "toppings",
                assumed_value: 50,
                unit: "grams granola",
                confidence: "medium",
                rationale: "Typical granola addition",
              },
            ],
            macros: {
              calories: 280,
              protein: 20,
              carbs: 30,
              fat: 10,
            },
          },
          {
            name: "Turkey Sandwich",
            assumptions: [
              {
                id: "4",
                field: "portion_size",
                assumed_value: 100,
                unit: "grams",
                confidence: "medium",
                rationale: "Typical deli meat portion",
              },
              {
                id: "12",
                field: "bread_type",
                assumed_value: 2,
                unit: "slices whole wheat",
                confidence: "medium",
                rationale: "Common bread choice",
              },
              {
                id: "13",
                field: "condiments",
                assumed_value: 1,
                unit: "tbsp mayo",
                confidence: "low",
                rationale: "Light spread",
              },
            ],
            macros: {
              calories: 350,
              protein: 25,
              carbs: 35,
              fat: 12,
            },
          },
          {
            name: "Vegetable Stir Fry",
            assumptions: [
              {
                id: "5",
                field: "portion_size",
                assumed_value: 300,
                unit: "grams",
                confidence: "medium",
                rationale: "Standard vegetable serving",
              },
              {
                id: "14",
                field: "cooking_oil",
                assumed_value: 1,
                unit: "tbsp",
                confidence: "medium",
                rationale: "Minimal oil for stir frying",
              },
            ],
            macros: {
              calories: 200,
              protein: 8,
              carbs: 25,
              fat: 10,
            },
          },
          {
            name: "Salmon with Quinoa",
            assumptions: [
              {
                id: "6",
                field: "portion_size",
                assumed_value: 150,
                unit: "grams",
                confidence: "medium",
                rationale: "Typical fish fillet size",
              },
              {
                id: "15",
                field: "cooking_method",
                assumed_value: 1,
                unit: "baked",
                confidence: "high",
                rationale: "Healthy cooking method",
              },
              {
                id: "16",
                field: "seasoning",
                assumed_value: 1,
                unit: "tsp herbs",
                confidence: "medium",
                rationale: "Light seasoning",
              },
            ],
            macros: {
              calories: 400,
              protein: 30,
              carbs: 25,
              fat: 20,
            },
          },
          {
            name: "Banana Smoothie",
            assumptions: [
              {
                id: "7",
                field: "portion_size",
                assumed_value: 1,
                unit: "medium banana",
                confidence: "high",
                rationale: "Average banana size",
              },
              {
                id: "17",
                field: "liquid_base",
                assumed_value: 1,
                unit: "cup milk",
                confidence: "medium",
                rationale: "Common smoothie base",
              },
            ],
            macros: {
              calories: 250,
              protein: 5,
              carbs: 40,
              fat: 3,
            },
          },
        ],
        setMeals: (meals) => set(() => ({ meals })),
      }),
      {
        name: "global-storage",
      },
    ),
  ),
);

export default useGlobalStore;
