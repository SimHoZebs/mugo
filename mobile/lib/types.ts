import { NutritionPayload } from "./api/conversationAPI.schemas";

export type Meal = {
  id: string;
  sessionId: string;
  nutrition: NutritionPayload;
};
